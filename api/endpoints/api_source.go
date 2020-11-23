package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/crypto"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/jira"
	"github.com/pkg/errors"
	"net/http"
	//nexpose "github.com/nortonlifelock/drivers/nexpose/connector"
	"github.com/nortonlifelock/aegis/internal/integrations"
	qualys "github.com/nortonlifelock/aegis/pkg/qualys/connector"
	"strconv"
	"strings"
)

func getPayloadSkeletons(trans *transaction) (sourceToSkeleton map[string]string) {
	sourceToSkeleton = make(map[string]string)

	for _, id := range []string{integrations.AWS, integrations.Azure, integrations.JIRA, integrations.Qualys, integrations.Nexpose} {
		if len(id) > 0 {
			var toMarshall interface{}

			switch id {
			case integrations.AWS:
				break
			case integrations.Azure:
				break
			case integrations.JIRA:
				toMarshall = jira.PayloadJira{}
				break
			case integrations.Qualys:
				toMarshall = qualys.QSPayload{}
				break
			case integrations.Nexpose:
				//toMarshall = nexpose.NXPayload{}
				break
			default:
				trans.err = fmt.Errorf("could not find payload skeleton for %s", id)
			}

			if trans.err == nil {
				if toMarshall != nil {
					var retValByte []byte
					retValByte, trans.err = json.Marshal(toMarshall)
					if trans.err == nil {
						sourceToSkeleton[id] = string(retValByte)
					} else {
						break
					}
				} else {
					sourceToSkeleton[id] = "{}"
				}
			} else {
				break
			}
		}
	}

	return sourceToSkeleton
}

func getAllSourceConfigs(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getSourcesEndpoint, admin|manager|reporter, func(trans *transaction) {
		var sourceConfigs []domain.SourceConfig
		sourceConfigs, trans.err = Ms.GetSourceConfigByOrgID(trans.permission.OrgID())
		if trans.err == nil {

			var sourceDtos sourceDtoContainer
			sourceDtos = toSourceConfigDtoSlice(sourceConfigs)
			sourceDtos.SourceToSkeleton = getPayloadSkeletons(trans)
			if trans.err == nil {
				trans.status = http.StatusOK
				trans.obj = sourceDtos
			} else {
				(&trans.wrapper).addError(trans.err, backendError)
			}

		} else {
			(&trans.wrapper).addError(trans.err, databaseError)
		}
	})
}

func updateSourceConfig(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, updateSourceEndpoint, admin|manager, func(trans *transaction) {
		if trans.endpoint.verify() {
			if source, isSource := trans.endpoint.(*Source); isSource {
				trans.obj, trans.status, trans.err = source.update(trans.user, trans.permission, trans.originalBody)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = fmt.Errorf("tried to pass a non-source as a source")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = fmt.Errorf("invalid length of format of apiRequest fields")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func deleteSourceConfig(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, deleteSourceEndpoint, admin, func(trans *transaction) {
		if trans.endpoint.verify() {
			if source, isSource := trans.endpoint.(*Source); isSource {
				trans.obj, trans.status, trans.err = source.delete(trans.user, trans.permission)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = fmt.Errorf("tried to pass a non-source as a source")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = fmt.Errorf("invalid length of format of apiRequest fields")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func createSourceConfig(w http.ResponseWriter, r *http.Request) {

	executeTransaction(w, r, createSourceEndpoint, admin|manager, func(trans *transaction) {
		//TODO source config does not have a created by column
		if trans.endpoint.verify() {
			if source, isSource := trans.endpoint.(*Source); isSource {
				trans.obj, trans.status, trans.err = source.create(trans.user, trans.permission)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = fmt.Errorf("tried to pass a non-source as a source")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = fmt.Errorf("invalid length of format of apiRequest fields")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func (source *Source) create(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest

	var client crypto.Client
	client, err = crypto.NewEncryptionClient(crypto.AES256, Ms, EncryptionKey, permission.OrgID(), AppConfig.KMSProfile(), AppConfig.KMSRegion())
	if err == nil {
		var encryptedPassword string
		var encryptedPrivateKey string
		var encryptedConsumerKey string
		var encryptedToken string

		encryptedPassword, encryptedPrivateKey, encryptedConsumerKey, encryptedToken, err = source.encryptSourceInformation(client)
		if err == nil {
			var port = ""

			if source.Port > 0 {
				port = strconv.Itoa(source.Port)
			}

			var sourceInDb domain.Source
			sourceInDb, err = Ms.GetSourceByName(source.Source)
			if err == nil {
				if sourceInDb != nil {
					_, _, err = Ms.CreateSourceConfig(
						source.Source,
						sourceInDb.ID(),
						permission.OrgID(),
						source.Address,
						port,
						source.Username,
						encryptedPassword,
						encryptedPrivateKey,
						encryptedConsumerKey,
						encryptedToken,
						source.Payload)
					if err == nil {

						// return the recently created source
						var sourceConfigsForOrg []domain.SourceConfig
						sourceConfigsForOrg, err = Ms.GetSourceConfigByOrgID(permission.OrgID())
						if err == nil {

							status = http.StatusOK
							generalResp.Response = fmt.Sprintf("%s", sourceConfigsForOrg[len(sourceConfigsForOrg)-1].ID())

						} else {
							err = errors.Errorf("error while gathering created source config - %s", err.Error())
						}
					} else {
						err = errors.Errorf("error while creating source config - %s", err.Error())
					}
				} else {
					err = fmt.Errorf("could not find source in database with name [%s]", source.Source)
				}
			} else {
				err = fmt.Errorf("error while gathering source from database with name [%s] - %s", source.Source, err.Error())
			}

		} else {
			err = errors.New("error while encrypting provided password")
		}
	}

	return generalResp, status, err
}

func (source *Source) encryptSourceInformation(client crypto.Client) (encryptedPassword string, encryptedPrivateKey string, encryptedConsumerKey string, encryptedToken string, err error) {
	if len(source.Password) > 0 {
		encryptedPassword, err = client.Encrypt(source.Password)
	}
	if err == nil { // prevent error overwrite
		if len(source.PrivateKey) > 0 {
			encryptedPrivateKey, err = client.Encrypt(source.PrivateKey)
		}
	}
	if err == nil { // prevent error overwrite
		if len(source.ConsumerKey) > 0 {
			encryptedConsumerKey, err = client.Encrypt(source.ConsumerKey)
		}
	}
	if err == nil { // prevent error overwrite
		if len(source.Token) > 0 {
			encryptedToken, err = client.Encrypt(source.Token)
		}
	}

	return encryptedPassword, encryptedPrivateKey, encryptedConsumerKey, encryptedToken, err
}

func (source *Source) delete(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest

	_, _, err = Ms.DisableSource(
		source.ID,
		permission.OrgID(),
		sord(user.Username()))
	if err == nil {
		status = http.StatusOK
		generalResp.Response = fmt.Sprintf("source config deleted")
	} else {
		err = errors.New("error while deleting source config")
	}

	return generalResp, status, err
}

func (source *Source) update(user domain.User, permission domain.Permission, originalBody string) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest

	var sourceConfigs domain.SourceConfig
	sourceConfigs, err = Ms.GetSourceConfigByID(source.ID)
	if err == nil {
		if sourceConfigs != nil {
			var bodyString = originalBody
			var address = source.Address
			var username = source.Username
			var password = source.Password
			var privateKey = source.PrivateKey
			var consumerKey = source.ConsumerKey
			var token = source.Token
			var payload = source.Payload
			var port = strconv.Itoa(source.Port)

			address, username, port, privateKey, consumerKey, token, password, payload, err = source.extractFieldsToUpdate(bodyString, address, sourceConfigs, username, port, privateKey, consumerKey, token, password, payload, permission)
			if err == nil {
				_, _, err = Ms.UpdateSourceConfig(
					source.ID,
					permission.OrgID(),
					address,
					username,
					password,
					privateKey,
					consumerKey,
					token,
					port,
					payload,
					sord(user.Username()))
			}

			if err == nil {
				status = http.StatusOK
				generalResp.Response = fmt.Sprintf("source config updated")
			} else {
				err = errors.Errorf("error while updating source config - %s", err.Error())
			}
		} else {
			err = errors.Errorf("could not find source config with id [%s]", source.ID)
		}
	} else {
		err = errors.Errorf("error while retreiving existing source config from database - %s", err.Error())
	}

	return generalResp, status, err
}

func (source *Source) extractFieldsToUpdate(bodyString string, address string, sourceConfigs domain.SourceConfig, username string, port string, privateKey string, consumerKey string, token string, password string, payload string, permission domain.Permission) (string, string, string, string, string, string, string, string, error) {
	var err error
	var allAuth domain.AllAuth
	if err = json.Unmarshal([]byte(sourceConfigs.AuthInfo()), &allAuth); err == nil {
		if strings.Index(bodyString, "address") < 0 {
			address = sourceConfigs.Address()
		}
		if strings.Index(bodyString, "username") < 0 {
			username = allAuth.Username
		}
		if strings.Index(bodyString, "port") < 0 {
			port = sourceConfigs.Port()
		}

		privateKey, consumerKey, token, password, err = source.extractFieldsToUpdateThatRequireDecryption(bodyString, privateKey, allAuth, consumerKey, token, password, permission)

		if err == nil {
			if strings.Index(bodyString, "payload") < 0 {
				payload = sord(sourceConfigs.Payload())
			}
		}
	}

	return address, username, port, privateKey, consumerKey, token, password, payload, err
}

func (source *Source) extractFieldsToUpdateThatRequireDecryption(bodyString string, privateKey string, allAuth domain.AllAuth, consumerKey string, token string, password string, permission domain.Permission) (string, string, string, string, error) {
	var client crypto.Client
	var err error

	client, err = crypto.NewEncryptionClient(crypto.AES256, Ms, EncryptionKey, permission.OrgID(), AppConfig.KMSProfile(), AppConfig.KMSRegion())
	if err == nil {
		if strings.Index(bodyString, "private_key") < 0 {
			privateKey = allAuth.PrivateKey
		} else {
			privateKey, err = client.Encrypt(privateKey)
			if err != nil {
				privateKey = allAuth.PrivateKey
				err = errors.Errorf("error while encrypting private key: %s\n", err.Error())
			}
		}
		if err == nil { // prevent error overwrite
			if strings.Index(bodyString, "consumer_key") < 0 {
				consumerKey = allAuth.ConsumerKey
			} else {
				consumerKey, err = client.Encrypt(consumerKey)
				if err != nil {
					consumerKey = allAuth.ConsumerKey
					err = errors.Errorf("error while encrypting consumer key: %s\n", err.Error())
				}
			}
		}
		if err == nil { // prevent error overwrite
			if strings.Index(bodyString, "token") < 0 {
				token = allAuth.Token
			} else {
				token, err = client.Encrypt(token)
				if err != nil {
					token = allAuth.Token
					err = errors.Errorf("error while encrypting token: %s\n", err.Error())
				}
			}
		}
		if err == nil { // prevent error overwrite
			if strings.Index(bodyString, "password") < 0 {
				password = allAuth.Password
			} else {
				password, err = client.Encrypt(password)
				if err != nil {
					password = allAuth.Password
					err = errors.Errorf("error while encrypting password: %s\n", err.Error())
				}
			}
		}
	}

	return privateKey, consumerKey, token, password, err
}

func (source *Source) verify() (verify bool) {
	// TODO what verifications can be done on the payload?

	checkAddress := source.Address
	checkAddress = strings.Replace(checkAddress, "https://", "", 1)
	checkAddress = strings.Replace(checkAddress, "http://", "", 1)

	verify = source.verifyFieldLengths(checkAddress)

	if verify {
		if check(usernameRegex, source.Username) {
			if check(lettersRegex, source.Source) {
				if check(passRegex, source.Password) {
					if check(urlRegex, checkAddress) || len(checkAddress) == 0 {
						verify = true
					}
				}
			}
		}
	}

	return verify
}

func (source *Source) verifyFieldLengths(checkAddress string) (verify bool) {

	if len(source.ID) >= 0 { // can equal zero when creating source
		if len(source.Source) > 0 && len(source.Source) < upperBoundNameLen {
			if len(source.Username) >= 0 && len(source.Username) <= upperBoundNameLen {
				if len(source.Password) >= 0 && len(source.Password) <= upperBoundNameLen {
					if len(source.PrivateKey) < oauthBoundLen && len(source.ConsumerKey) < oauthBoundLen && len(source.Token) < oauthBoundLen {
						if source.Port >= -1 {
							if len(checkAddress) >= 0 && len(checkAddress) < 100 {
								verify = true
							}
						}
					}
				}
			}
		}
	}
	return verify
}

func getAllSources(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getAllSourceConfigsEndpoint, allAllowed, func(trans *transaction) {
		trans.obj, trans.status, trans.totalRecords, trans.err = readSources(trans.permission)
		if trans.err != nil {
			trans.err = fmt.Errorf("error while retrieving Sources [%s]", trans.err.Error())
			(&trans.wrapper).addError(trans.err, databaseError)
		}
	})
}

func readSources(permission domain.Permission) (sourceDtos []*Src, status int, totalRecords int, err error) {
	status = http.StatusBadRequest
	if permission != nil && len(permission.OrgID()) > 0 {

		var sources []domain.Source

		sources, err = Ms.GetSources()
		if err == nil {
			status = http.StatusOK
			if sources != nil {
				sourceDtos = toSourceDtoSlice(sources)
			}
		} else {
			err = fmt.Errorf("error while getting sources [%s]", err.Error())
		}
	} else {
		err = fmt.Errorf("error: Permission is nil or ogranization is not provided")
	}

	return sourceDtos, status, totalRecords, err
}
