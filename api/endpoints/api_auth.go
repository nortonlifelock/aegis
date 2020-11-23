package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nortonlifelock/aegis/pkg/crypto"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/pkg/errors"
)

func logout(w http.ResponseWriter, r *http.Request) {
	var bearerTokens []string
	var exists bool
	var bearerToken string
	var err error
	var wrapper errorWrapper

	bearerTokens, exists = r.Header["Authorization"]
	if exists && len(bearerTokens) >= 1 {
		bearerToken = bearerTokens[0]
		bearerToken = strings.TrimPrefix(bearerToken, "Bearer ")
		_, _, err = Ms.DeleteSessionByToken(bearerToken)
		(&wrapper).addError(err, databaseError)

	} else {
		err = errors.Errorf("logged out without providing a bearer token to terminate")
		(&wrapper).addError(err, authenticationError)
	}

	if err == nil {
		respondToUser(nil, newResponse(nil, 0), w, wrapper, logoutEndpoint)
	} else {
		respondToUserWithStatusCode(nil, newResponse(nil, 0), w, wrapper, logoutEndpoint, http.StatusUnauthorized)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	var obj interface{}

	var wrapper errorWrapper
	var req = &loginReq{}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, mbSize))
	if err == nil {
		err = r.Body.Close()
		if err == nil {
			err = json.Unmarshal(body, req)
			if err == nil {
				var groups []string
				obj, groups, err = authenticateUser(req, &wrapper)
				_ = groups
			} else {
				(&wrapper).addError(err, requestFormatError)
			}
		} else {
			(&wrapper).addError(err, requestFormatError)
		}
	} else {
		(&wrapper).addError(err, requestFormatError)
	}

	if err == nil {
		respondToUser(user, newResponse(obj, 0), w, wrapper, loginEndpoint)
	} else {
		respondToUserWithStatusCode(user, newResponse(obj, 0), w, wrapper, loginEndpoint, http.StatusUnauthorized)
	}
}

func authenticateUser(loginReq *loginReq, wrapper *errorWrapper) (obj interface{}, groups []string, err error) {
	var user domain.User
	user, err = Ms.GetUserByUsername(loginReq.Username)
	if err == nil {
		if user != nil {
			if !user.IsDisabled() {

				user, err = encryptOrDecryptUser(user, crypto.DecryptMode)
				if err == nil {
					var adConfig ADConfig
					adConfig, err = grabADConfig(loginReq.OrgID)

					if err == nil {
						var token string
						token, groups, err = authenticate(adConfig, user, loginReq.Password, loginReq.OrgID)
						if err == nil {
							obj, err = createSession(user, token, wrapper)
						} else {
							wrapper.addError(err, authenticationError)
						}
					} else {
						wrapper.addError(err, backendError)
					}
				} else {
					wrapper.addError(fmt.Errorf("error while decrypting user for authentication - %s", err.Error()), backendError)
				}

			} else {
				err = errors.New("tried to create a session for a disabled user")
				wrapper.addError(err, authenticationError)
			}
		} else {
			err = errors.Errorf("could not find user with the username %s in the Aegis database", loginReq.Username)
			wrapper.addError(err, authenticationError)
		}
	} else {
		wrapper.addError(err, databaseError)
	}

	return obj, groups, err
}

func grabADConfig(orgID string) (adConfig ADConfig, err error) {
	if len(orgID) > 0 {
		if OrgADConfigs[orgID] != nil {
			adConfig = *OrgADConfigs[orgID].Con
		} else {
			err = fmt.Errorf("couldn't find AD configuration for organization [%s]", orgID)
		}
	} else {
		// TODO currently we are grabbing a random config when first logging in - do we want to have users login to an organization?
		for _, val := range OrgADConfigs {
			adConfig = *val.Con
			break
		}
	}

	return adConfig, err
}

func createSession(user domain.User, token string, wrapper *errorWrapper) (obj interface{}, err error) {
	var defaultPermissions domain.Permission
	defaultPermissions, err = Ms.GetPermissionOfLeafOrgByUserID(user.ID())
	if err == nil {
		if defaultPermissions != nil {
			_, _, err = Ms.CreateUserSession(user.ID(), defaultPermissions.OrgID(), token)
			if err == nil {
				obj = token
			} else {
				err = errors.New("could not create user session in database")
				wrapper.addError(err, databaseError)
			}
		} else {
			wrapper.addError(errors.Errorf("no permissions found user %s", sord(user.Username())), authorizationError)
		}
	} else {
		wrapper.addError(err, databaseError)
	}

	return obj, err
}

func updateSessionPermissions(w http.ResponseWriter, r *http.Request) {
	var status int
	var err error
	var user domain.User
	var token string
	var obj interface{}
	var wrapper errorWrapper

	params := mux.Vars(r)

	token, _, _, err = handleRequest(w, r, updateUserEndpoint)
	if err == nil {
		user, _, err = validateToken(token)
		if err == nil {
			if status, obj, err = updateUserSession(user, params, &wrapper); err != nil {
				(&wrapper).addError(err, authenticationError)
			}
		} else {
			status = http.StatusUnauthorized
			(&wrapper).addError(err, authenticationError)
		}
	} else {
		status = http.StatusUnauthorized
		(&wrapper).addError(err, authenticationError)
	}

	respondToUserWithStatusCode(user, newResponse(obj, 0), w, wrapper, updateUserEndpoint, status)
}

func updateUserSession(user domain.User, params map[string]string, wrapper *errorWrapper) (status int, obj interface{}, err error) {
	if user.Username() != nil {
		var orgIDString = params["org"]
		if len(orgIDString) > 0 {
			status, obj, err = recreateSession(user, orgIDString, wrapper)
		} else {
			wrapper.addError(errors.Errorf("org id not included in apiRequest"), requestFormatError)
		}
	} else {
		wrapper.addError(errors.Errorf("username could not be found from bearer token"), requestFormatError)
	}

	return status, obj, err
}

func recreateSession(user domain.User, orgID string, wrapper *errorWrapper) (status int, obj interface{}, err error) {
	var token string
	// generateJWT returns an error if user doesn't have permission for org
	token, err = generateJWT(sord(user.Username()), orgID)
	if err == nil {
		if !user.IsDisabled() {

			status, obj, err = createSessionForOrg(user, orgID, token, wrapper)
		} else {
			err = errors.New("tried to create a session for a disabled user")
			wrapper.addError(err, authenticationError)
		}
	} else {
		status = http.StatusUnauthorized
		wrapper.addError(err, authenticationError)
	}

	return status, obj, err
}

func createSessionForOrg(user domain.User, orgID string, token string, wrapper *errorWrapper) (status int, obj interface{}, err error) {
	var orgSpecificPermissions domain.Permission
	orgSpecificPermissions, err = Ms.GetPermissionByUserOrgID(user.ID(), orgID)
	if err == nil {
		if orgSpecificPermissions != nil {
			_, _, err = Ms.CreateUserSession(user.ID(), orgSpecificPermissions.OrgID(), token)
			if err == nil {
				obj = token
				status = http.StatusOK
			} else {
				err = errors.New("could not create user session in database")
				wrapper.addError(err, databaseError)
			}
		} else {
			wrapper.addError(errors.Errorf("no permissions found for org %s and user %s", orgID, sord(user.Username())), authorizationError)
		}
	} else {
		wrapper.addError(err, databaseError)
	}

	return status, obj, err
}
