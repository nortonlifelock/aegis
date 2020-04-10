package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/subscriptions"
	network2 "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/nortonlifelock/aws"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"os"
	"strings"
	"sync"
)

// ConnectionAzure holds the authorization information and URL for Azure. tagNames acts as a cache holding the unique list of tag names, as
// the API call to gather them all is expensive
type ConnectionAzure struct {
	authorizer      autorest.Authorizer
	url             string
	subscriptionIDs []string
	lstream         logger

	tagNames []string // TODO should I do a ttl?
	ips      []domain.CloudIP
}

type logger interface {
	Send(log log.Log)
}

// Send implements the logger interface and lets the Azure connection be used in place of the logger
func (connection *ConnectionAzure) Send(log log.Log) {
	connection.lstream.Send(log)
}

type authStruct struct {
	TenantID     string `json:"TenantID"`
	ClientID     string `json:"ClientID"`
	ClientSecret string `json:"Password"`
}

const (
	mac = "mac"
)

// Azure authenticates using environment variables, so we use a mutex to prevent multiple threads from accessing those environment variables simultaneously
var azureAuthLock sync.Mutex

// CreateConnection establishes a connection to Azure and collects a list of subscriptions that the auth information has access to
func CreateConnection(authInfo, url string, lstream logger) (connection *ConnectionAzure, err error) {
	azureAuthLock.Lock()

	defer func() {
		_ = os.Setenv(auth.TenantID, "")
		_ = os.Setenv(auth.ClientID, "")
		_ = os.Setenv(auth.ClientSecret, "")
		azureAuthLock.Unlock()
	}()

	var parseAuth = &authStruct{}
	if err = json.Unmarshal([]byte(authInfo), parseAuth); err == nil {

		err = os.Setenv(auth.TenantID, parseAuth.TenantID) // directory id
		if err == nil {
			err = os.Setenv(auth.ClientID, parseAuth.ClientID) // application id
			if err == nil {
				err = os.Setenv(auth.ClientSecret, parseAuth.ClientSecret)
			}
		}

		if err == nil {
			var authorizer autorest.Authorizer
			if authorizer, err = auth.NewAuthorizerFromEnvironment(); err == nil {
				connection = &ConnectionAzure{
					authorizer: authorizer,
					url:        url, // e.g. https://management.azure.com
					lstream:    lstream,
				}

				connection.subscriptionIDs, err = connection.getSubscriptionIDs()
			} else {
				err = fmt.Errorf("error while authorizing Azure client - %s", err.Error())
			}
		}
	}

	return connection, err
}

// GetIPTagMapping returns a map where the first key is an IP that returns a map containing a key->value pair containing tagNames->tagValues
// this method also creates a list of unique tag names and caches them for later use by GetAllTagNames
func (connection *ConnectionAzure) GetIPTagMapping() (ipToKeyToValue map[domain.CloudIP]map[string]string, err error) {
	ipToKeyToValue = make(map[domain.CloudIP]map[string]string)

	wg := &sync.WaitGroup{}
	lock := &sync.Mutex{}

	for _, subID := range connection.subscriptionIDs {
		wg.Add(1)
		go func(subID string) {
			defer wg.Done()
			threadErr := connection.getIPTagsForSub(subID, lock, ipToKeyToValue)
			if threadErr != nil {
				err = threadErr
			}
		}(subID)
	}

	wg.Wait()

	connection.tagNames = make([]string, 0)

	var seenTag = make(map[string]bool)
	for ip := range ipToKeyToValue {
		for key := range ipToKeyToValue[ip] {
			if !seenTag[key] {
				seenTag[key] = true
				connection.tagNames = append(connection.tagNames, key)
			}
		}
	}

	connection.ips = make([]domain.CloudIP, 0)

	var seenIP = make(map[domain.CloudIP]bool)
	for ip := range ipToKeyToValue {

		if !seenIP[ip] {
			seenIP[ip] = true
			connection.ips = append(connection.ips, ip)
		}
	}

	return ipToKeyToValue, err
}

func (connection *ConnectionAzure) getIPTagsForSub(subID string, lock *sync.Mutex, ipToKeyToValue map[domain.CloudIP]map[string]string) (err error) {
	subWg := &sync.WaitGroup{}
	subWg.Add(3)

	var errChan = make(chan error, 3)

	go func() {
		defer close(errChan)

		go func() {
			defer subWg.Done()
			connection.getIPTagsForVirtualMachines(subID, lock, ipToKeyToValue, errChan)
		}()

		go func() {
			defer subWg.Done()
			connection.getIPTagsForApplicationGateways(subID, lock, ipToKeyToValue, errChan)
		}()

		go func() {
			defer subWg.Done()
			connection.getIPTagsForLoadBalancers(subID, lock, ipToKeyToValue, errChan)
		}()
	}()

	subWg.Wait()

	for {
		if subErr, ok := <-errChan; ok {
			if err == nil {
				err = subErr
			} else {
				err = fmt.Errorf("%v,%v", err, subErr)
			}
		} else {
			break
		}
	}

	return err
}

func (connection *ConnectionAzure) getIPTagsForLoadBalancers(subID string, lock *sync.Mutex, ipToKeyToValue map[domain.CloudIP]map[string]string, errChan chan error) {
	lbTags, err := connection.loadBalancerTags(subID)
	if err == nil {
		for ip := range lbTags {
			lock.Lock()

			cip := &cloudIP{
				ip:         ip,
				subID:      subID,
				state:      "",
				instanceID: lbTags[ip][awsclient.InstanceID],
			}

			if ipToKeyToValue[cip] == nil {
				ipToKeyToValue[cip] = make(map[string]string)
			}

			for key, value := range lbTags[ip] {
				ipToKeyToValue[cip][key] = value
			}
			lock.Unlock()
		}
	} else {
		connection.Send(log.Errorf(err, "error while gathering load balancer tags for subscription id [%s]", subID))
		errChan <- err
	}
}

func (connection *ConnectionAzure) getIPTagsForApplicationGateways(subID string, lock *sync.Mutex, ipToKeyToValue map[domain.CloudIP]map[string]string, errChan chan error) {
	appTags, err := connection.applicationGatewayTags(subID)
	if err == nil {
		for ip := range appTags {
			lock.Lock()

			cip := &cloudIP{
				ip:         ip,
				subID:      subID,
				state:      "",
				instanceID: appTags[ip][awsclient.InstanceID],
			}

			if ipToKeyToValue[cip] == nil {
				ipToKeyToValue[cip] = make(map[string]string)
			}

			for key, value := range appTags[ip] {
				ipToKeyToValue[cip][key] = value
			}
			lock.Unlock()
		}
	} else {
		connection.Send(log.Errorf(err, "error while gathering application gateway tags for subscription id [%s]", subID))
		errChan <- err
	}
}

func (connection *ConnectionAzure) getIPTagsForVirtualMachines(subID string, lock *sync.Mutex, ipToKeyToValue map[domain.CloudIP]map[string]string, errChan chan error) {
	vmTags, err := connection.virtualMachineNetworkInterfaceTags(subID)
	if err == nil {
		for ip := range vmTags {
			lock.Lock()

			cip := &cloudIP{
				ip:         ip,
				subID:      subID,
				state:      "",
				mac:        vmTags[ip][mac],
				instanceID: vmTags[ip][awsclient.InstanceID],
			}

			if ipToKeyToValue[cip] == nil {
				ipToKeyToValue[cip] = make(map[string]string)
			}

			for key, value := range vmTags[ip] {
				ipToKeyToValue[cip][key] = value
			}
			lock.Unlock()
		}
	} else {
		connection.Send(log.Errorf(err, "error while gathering VM network tags for subscription id [%s]", subID))
		errChan <- err
	}
}

// GetAllTagNames grabs a unique list of tag names
// due to the expensive nature of gathering all tags, we will cache the tag names as we create the ip mapping
func (connection *ConnectionAzure) GetAllTagNames() (tagNames []string, err error) {
	if connection.tagNames == nil {
		_, err = connection.GetIPTagMapping() // this method populates the tag names
	}

	return connection.tagNames, err
}

// IPAddresses retrieves a list of IP Addresses for the azure subscription
func (connection *ConnectionAzure) IPAddresses() (ips []domain.CloudIP, err error) {
	_, err = connection.GetIPTagMapping()
	return connection.ips, err
}

func (connection *ConnectionAzure) getSubscriptionIDs() (subscriptionIDs []string, err error) {
	var subscriptionClient = subscriptions.NewClient()
	subscriptionClient.Authorizer = connection.authorizer
	subscriptionIDs = make([]string, 0)

	var resultPage subscriptions.ListResultPage
	resultPage, err = subscriptionClient.List(context.Background())
	if err == nil {
		subscriptionContracts := resultPage.Values()

		for _, subscription := range subscriptionContracts {
			if subscription.ID != nil {
				//fmt.Println(*subscription.ID, *subscription.DisplayName)
				fields := strings.Split(*subscription.ID, "/")
				if len(fields) == 3 {
					subscriptionIDs = append(subscriptionIDs, fields[2])
				} else {
					err = fmt.Errorf("invalid id formatting for subscription [%s]", *subscription.ID)
				}
			}
		}
	} else {
		err = fmt.Errorf("error while gathering subscription ids - %s", err.Error())
	}

	if err == nil && len(subscriptionIDs) == 0 {
		err = fmt.Errorf("account does not appear to have access to any subscriptions")
	}

	return subscriptionIDs, err
}

func (connection *ConnectionAzure) virtualMachineNetworkInterfaceTags(subscriptionID string) (ipToKeyToValue map[string]map[string]string, err error) {
	ipToKeyToValue = make(map[string]map[string]string)

	interfacesClient := network.NewInterfacesClient(subscriptionID)
	interfacesClient.Authorizer = connection.authorizer
	resultPage, err := interfacesClient.ListAll(context.Background())
	if err == nil {
		interfaces := resultPage.Values()

		for _, networkInterface := range interfaces {
			if networkInterface.InterfacePropertiesFormat != nil {
				if networkInterface.IPConfigurations != nil {
					for _, ipConfig := range *networkInterface.IPConfigurations {
						if ipConfig.InterfaceIPConfigurationPropertiesFormat != nil {
							if ipConfig.InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress != nil {
								ip := *ipConfig.InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress
								if ipToKeyToValue[ip] == nil {
									ipToKeyToValue[ip] = make(map[string]string)
								}

								if networkInterface.Name != nil {
									instanceID := *networkInterface.Name
									// TODO should we share the const?
									ipToKeyToValue[ip][awsclient.InstanceID] = instanceID
								}

								if networkInterface.MacAddress != nil {
									ipToKeyToValue[ip][mac] = *networkInterface.MacAddress
								}

								for key, value := range networkInterface.Tags {
									if value != nil {
										ipToKeyToValue[ip][key] = *value
									}
								}
							}
						}
					}
				}
			}
		}
	} else {
		err = fmt.Errorf("error while grabbing network interface information - %s", err.Error())
	}

	return ipToKeyToValue, err
}

func (connection *ConnectionAzure) applicationGatewayTags(subscriptionID string) (ipToKeyToValue map[string]map[string]string, err error) {
	ipToKeyToValue = make(map[string]map[string]string)

	appGatewayClient := network.NewApplicationGatewaysClient(subscriptionID)
	appGatewayClient.Authorizer = connection.authorizer
	var applicationGatewayListResultPage network.ApplicationGatewayListResultPage
	applicationGatewayListResultPage, err = appGatewayClient.ListAll(context.Background())

	if err == nil {
		gateways := applicationGatewayListResultPage.Values()
		for _, gateway := range gateways {

			if gateway.ApplicationGatewayPropertiesFormat != nil && gateway.Tags != nil {

				connection.applicationGatewayFrontendIP(gateway, ipToKeyToValue)

				connection.applicationGatewayBackendIP(gateway, ipToKeyToValue)
			}

		}
	} else {
		err = fmt.Errorf("error while grabbing gateway information - %s", err.Error())
	}

	return ipToKeyToValue, err
}

func (connection *ConnectionAzure) applicationGatewayBackendIP(gateway network2.ApplicationGateway, ipToKeyToValue map[string]map[string]string) {
	if gateway.BackendAddressPools != nil {
		for _, ipConfig := range *gateway.BackendAddressPools {

			if ipConfig.BackendIPConfigurations != nil {
				for _, configVal := range *ipConfig.BackendIPConfigurations {

					if configVal.InterfaceIPConfigurationPropertiesFormat != nil {
						if configVal.PrivateIPAddress != nil {
							if ipToKeyToValue[*configVal.PrivateIPAddress] == nil {
								ipToKeyToValue[*configVal.PrivateIPAddress] = make(map[string]string)
							}

							if gateway.Name != nil {
								instanceID := *gateway.Name
								// TODO should we share the const?
								ipToKeyToValue[*configVal.PrivateIPAddress][awsclient.InstanceID] = instanceID
							}

							for key, value := range gateway.Tags {
								if value != nil {
									ipToKeyToValue[*configVal.PrivateIPAddress][key] = *value
								}
							}
						}
					}
				}
			}
		}
	}
}

func (connection *ConnectionAzure) applicationGatewayFrontendIP(gateway network2.ApplicationGateway, ipToKeyToValue map[string]map[string]string) {
	if gateway.FrontendIPConfigurations != nil {
		for _, ipConfig := range *gateway.FrontendIPConfigurations {
			if ipConfig.PrivateIPAddress != nil {
				if ipToKeyToValue[*ipConfig.PrivateIPAddress] == nil {
					ipToKeyToValue[*ipConfig.PrivateIPAddress] = make(map[string]string)
				}

				if gateway.Name != nil {
					// TODO should we share the const?
					instanceID := *gateway.Name
					ipToKeyToValue[*ipConfig.PrivateIPAddress][awsclient.InstanceID] = instanceID
				}

				for key, value := range gateway.Tags {
					if value != nil {
						ipToKeyToValue[*ipConfig.PrivateIPAddress][key] = *value
					}
				}
			}

		}
	}
}

func (connection *ConnectionAzure) loadBalancerTags(subscriptionID string) (ipToKeyToValue map[string]map[string]string, err error) {
	ipToKeyToValue = make(map[string]map[string]string)

	loadBalancerClient := network.NewLoadBalancersClient(subscriptionID)
	loadBalancerClient.Authorizer = connection.authorizer

	loadBalancerInfo, err := loadBalancerClient.ListAll(context.Background())
	if err == nil {

		loadBalancerInfoValues := loadBalancerInfo.Values()

		loadBalancerIPClient := network.NewLoadBalancerFrontendIPConfigurationsClient(subscriptionID)
		loadBalancerIPClient.Authorizer = connection.authorizer

		for _, lb := range loadBalancerInfoValues {

			if lb.ID != nil {

				fields := strings.Split(*lb.ID, "/")

				if len(fields) == 9 {

					var ipConfig network.LoadBalancerFrontendIPConfigurationListResultPage
					ipConfig, err = loadBalancerIPClient.List(context.Background(), fields[4], fields[8])
					if err == nil {
						ipsInfo := ipConfig.Values()

						for _, ipsInfo := range ipsInfo {
							if ipsInfo.PrivateIPAddress != nil {
								if ipToKeyToValue[*ipsInfo.PrivateIPAddress] == nil {
									ipToKeyToValue[*ipsInfo.PrivateIPAddress] = make(map[string]string)
									fmt.Println("lb", *ipsInfo.PrivateIPAddress)
								}

								if lb.Name != nil {
									instanceID := *lb.Name
									// TODO should we share the const?
									ipToKeyToValue[*ipsInfo.PrivateIPAddress][awsclient.InstanceID] = instanceID
								}

								for key, value := range lb.Tags {
									if value != nil {
										ipToKeyToValue[*ipsInfo.PrivateIPAddress][key] = *value
									}
								}
							}
						}
					} else {
						err = fmt.Errorf("error while grabbing ip information for load balancer [%s]", *lb.ID)
					}

				} else {
					err = fmt.Errorf("could not extract resource group and load balancer name from %s", *lb.ID)
					break
				}

			}
		}

	} else {
		err = fmt.Errorf("error while gathering load balancer information - %s", err.Error())
	}

	return ipToKeyToValue, err
}

type cloudIP struct {
	ip         string
	subID      string
	state      string
	mac        string
	instanceID string
}

func (c *cloudIP) IP() string {
	return c.ip
}

func (c *cloudIP) Region() string {
	return c.subID
}

// TODO
func (c *cloudIP) State() string {
	return domain.DeviceRunning
}

// TODO mac not populated
func (c *cloudIP) MAC() string {
	return c.mac
}

func (c *cloudIP) InstanceID() string {
	return c.instanceID
}
