package awsclient

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/nortonlifelock/aegis/pkg/domain"
)

// GetIPTagMapping returns a map, where the first key is an IP address that returns another map
// the inner map is a series of key-value pairs for tags and their values
func (connection *CloudConnection) GetIPTagMapping() (ipToKeyToValue map[domain.CloudIP]map[string]string, err error) {
	ipToKeyToValue = make(map[domain.CloudIP]map[string]string)

	var instances []*instanceInfo
	instances, err = connection.getAllInstances()

	for index := range instances {
		instance := instances[index]

		for nwIndex := range instance.instance.NetworkInterfaces {
			networkInterface := instance.instance.NetworkInterfaces[nwIndex]

			if networkInterface.PrivateIpAddress != nil {
				var currentIP = *networkInterface.PrivateIpAddress

				if len(currentIP) > 0 {
					var keyValues = make(map[string]string)

					for _, tag := range instance.instance.Tags {

						if tag.Key != nil && tag.Value != nil {

							if len(*tag.Key) > 0 && len(*tag.Value) > 0 {
								keyValues[*tag.Key] = *tag.Value
							}

						}
					}

					if instance.instance.State != nil {
						if instance.instance.State.Name != nil {
							keyValues[State] = *instance.instance.State.Name
						}
					}

					if instance.instance.InstanceId != nil {
						keyValues[InstanceID] = *instance.instance.InstanceId
					}

					ipToKeyToValue[instance] = keyValues
				}
			}
		}
	}

	return ipToKeyToValue, err
}

// GetAllTagNames returns a unique list of the names of tags found from the API
func (connection *CloudConnection) GetAllTagNames() (tagNames []string, err error) {
	var instances []*instanceInfo
	instances, err = connection.getAllInstances()

	if err == nil {
		var seen = make(map[string]bool)
		tagNames = make([]string, 0)

		for iIndex := range instances {
			instance := instances[iIndex]

			for tIndex := range instance.instance.Tags {
				tag := instance.instance.Tags[tIndex]

				if tag.Key != nil {

					if !seen[*tag.Key] {
						tagNames = append(tagNames, *tag.Key)
						seen[*tag.Key] = true
					}

				}

			}
		}

		// while these are not explicitly tags, they are used as tags
		tagNames = append(tagNames, State, InstanceID)
	}

	return tagNames, err
}

// IPAddresses returns the IP Address associated with the AWS connection
func (connection *CloudConnection) IPAddresses() (ips []domain.CloudIP, err error) {
	var seen = make(map[string]bool)
	ips = make([]domain.CloudIP, 0)
	var instances []*instanceInfo
	if instances, err = connection.getAllInstances(); err == nil {
		for _, instance := range instances {
			if len(instance.IP()) > 0 {
				if !seen[instance.IP()] {
					seen[instance.IP()] = true
					ips = append(ips, instance)
				}
			}
		}
	}

	return ips, err
}

// GetAllRegions retrieves a list of regions from AWS
func (connection *CloudConnection) GetAllRegions() (regions []string, err error) {
	regions = make([]string, 0)
	//req := &ec2.DescribeRegionsInput{}
	//
	//var resp *ec2.DescribeRegionsOutput
	//resp, err = connection.ec2Service.DescribeRegions(req)
	//if err == nil {
	//	for index := range resp.Regions {
	//		if resp.Regions[index].RegionName != nil {
	//			regions = append(regions, *resp.Regions[index].RegionName)
	//		}
	//	}
	//}

	return regions, err
}

type instanceInfo struct {
	instance *ec2.Instance
	region   string
}

func (i *instanceInfo) IP() string {
	var ip string
	if i.instance != nil {
		if i.instance.PrivateIpAddress != nil {
			ip = *i.instance.PrivateIpAddress
		}
	}
	return ip
}

func (i instanceInfo) MAC() string {
	var mac string
	if i.instance != nil {
		for _, networkInterface := range i.instance.NetworkInterfaces {
			if networkInterface.MacAddress != nil {
				if len(*networkInterface.MacAddress) > 0 {
					mac = *networkInterface.MacAddress // TODO what if there's more than one MAC?
					break
				}
			}
		}
	}
	return mac
}

func (i *instanceInfo) Region() string {
	return i.region
}

func (i *instanceInfo) State() string {
	var state string
	if i.instance.State != nil {
		if i.instance.State.Name != nil {
			state = *i.instance.State.Name
		}
	}

	if state == "terminated" {
		state = domain.DeviceDecommed
	}
	return state
}

func (i *instanceInfo) InstanceID() string {
	var id string
	if i.instance.InstanceId != nil {
		id = *i.instance.InstanceId
	}
	return id
}

func (connection *CloudConnection) getAllInstances() (instances []*instanceInfo, err error) {
	instances = make([]*instanceInfo, 0)
	threadErrors := make([]error, 0)

	lock := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for _, region := range connection.regions {

		wg.Add(1)
		go func(region string) {
			defer wg.Done()

			instancesForRegion, err := connection.getInstancesForRegion(region)
			lock.Lock()
			defer lock.Unlock()
			if err == nil {
				for _, instance := range instancesForRegion {
					instances = append(instances, &instanceInfo{
						instance: instance,
						region:   region,
					})
				}
			} else {
				err = fmt.Errorf("error while grabbing instances for region [%v] - %v", region, err.Error())
				threadErrors = append(threadErrors, err)
			}
		}(region)
	}

	wg.Wait()

	for _, subError := range threadErrors {
		err = fmt.Errorf("[%v]-[%v]", err, subError)
	}

	return instances, err
}

func (connection *CloudConnection) getInstancesForRegion(region string) (instances []*ec2.Instance, err error) {
	var regionConnection *CloudConnection
	if regionConnection, err = createCloudConnectionForRegion(connection.config.Credentials, region); err == nil {

		instances = make([]*ec2.Instance, 0)
		req := &ec2.DescribeInstancesInput{}

		for {
			var resp *ec2.DescribeInstancesOutput
			resp, err = regionConnection.ec2Service.DescribeInstances(req)
			if err == nil {
				for index := range resp.Reservations {
					res := resp.Reservations[index]
					instances = append(instances, res.Instances...)
				}

				if resp != nil {
					if resp.NextToken != nil {
						req.SetNextToken(*resp.NextToken)
						continue
					}
				}
			}

			break
		}
	} else {
		err = fmt.Errorf("error while creating regional connection for [%v] - %v", region, err.Error())
	}

	return instances, err
}
