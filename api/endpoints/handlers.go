package endpoints

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/nortonlifelock/domain"
	"math"
	"net/http"
	"strings"
)

//TODO include the jwt contents (user, expiration, permissions) in the response because decoding on the front end isn't good

type errorWrapper struct {
	detailedError error
	friendlyError error
}

func (wrapper *errorWrapper) addError(err error, message string) {
	if err != nil {
		wrapper.detailedError = err
		wrapper.friendlyError = fmt.Errorf(message)
	}
}

type websocketTransaction struct {
	connection *websocket.Conn
	err        error
	permission domain.Permission
	user       domain.User
}

type transaction struct {
	status       int
	wrapper      errorWrapper
	err          error
	user         domain.User
	permission   domain.Permission
	endpoint     endpoint
	originalBody string
	token        string
	obj          interface{}
	totalRecords int
}

func executeWebsocketTransaction(w http.ResponseWriter, r *http.Request, endpointName string, process func(trans *websocketTransaction)) {
	trans := &websocketTransaction{}

	trans.connection, trans.err = createSocketConnection(w, r)
	if trans.err == nil {
		var token string
		token, _, _, trans.err = handleRequest(w, r, endpointName)
		if trans.err == nil {
			trans.user, trans.permission, trans.err = validateToken(token)
			if trans.err == nil {
				process(trans)
			}
		}
	}

	if trans.err != nil {
		fmt.Println(trans.err.Error())
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
}

func executeTransaction(w http.ResponseWriter, r *http.Request, endpointName string, permissionMode int, process func(trans *transaction)) {
	var trans = &transaction{
		status: http.StatusBadRequest,
	}

	defer handleRoutinePanic(trans, w, endpointName)

	if verifyUIOrigin(r) {
		trans.token, trans.endpoint, trans.originalBody, trans.err = handleRequest(w, r, endpointName)
		if trans.err == nil {
			trans.user, trans.permission, trans.err = validateToken(trans.token)
			if trans.err == nil {
				if checkPermission(permissionMode, trans.permission) {
					process(trans)
				} else {
					(&trans.wrapper).addError(fmt.Errorf("user does not have permissions for this action"), permissionError)
					trans.status = http.StatusForbidden
				}
			} else {
				(&trans.wrapper).addError(trans.err, authorizationError)
				trans.status = http.StatusUnauthorized
			}
		} else {
			(&trans.wrapper).addError(trans.err, authorizationError)
			trans.status = http.StatusUnauthorized
		}
	} else {
		(&trans.wrapper).addError(fmt.Errorf("apiRequest did not originate from UI"), authorizationError)
	}

	respondToUserWithStatusCode(trans.user, newResponse(trans.obj, trans.totalRecords), w, trans.wrapper, endpointName, trans.status)
}

func checkPermission(permissionMode int, basePermission domain.Permission) bool {
	var isAllowed bool
	var currentPermission = basePermission
	if currentPermission != nil {
		for currentPermission != nil {
			isAllowed = (currentPermission.Admin() && currentPermission.ParentOrgPermission() == nil) || // admin at the root of an organization hierarchy can do anything
				(permissionMode&admin == admin && currentPermission.Admin()) ||
				(permissionMode&manager == manager && currentPermission.Manager()) ||
				(permissionMode&reporter == reporter && currentPermission.Reporter()) ||
				(permissionMode&reader == reader && currentPermission.Reader())

			if !isAllowed {
				currentPermission = currentPermission.ParentOrgPermission()
			} else {
				break
			}
		}
	}

	return isAllowed
}

var (
	// EncryptionKey is used to encrypt/decrypt fields to/from the database
	EncryptionKey string

	// mbSize is the size of a megabyte to prevent huge requests
	mbSize = int64(math.Pow(2, 20))

	// Ms is the database connection
	Ms domain.DatabaseConnection

	// AppConfig holds relevant fields from the app config
	AppConfig listenerConfig

	// WorkingDir defines the path to the Aegis directory, this is used over os.WorkingDir() as it is inconsistent
	WorkingDir string

	// OrgADConfigs holds the AD configurations for each organization. The AD configuration from the root of an
	// organizational hierarchy is used for the leaf organizations. The key is the organization ID
	OrgADConfigs map[string]*OrgConfigWrapper
)

type listenerConfig interface {
	TransportProtocol() string
	UILocation() string
	EncryptionKey() string
	KMSProfile() string
	KMSRegion() string
}

// This method ensures our server does not make a websocket connection when requests come from
// anywhere besides our UI
func verifyUIOrigin(r *http.Request) (validOrigin bool) {
	validOrigin = false
	var originUI = fmt.Sprintf("%s://%s", AppConfig.TransportProtocol(), AppConfig.UILocation())

	origin := r.Header["Origin"]
	if len(origin) == 1 {
		// Will there ever be more than 1 origin?
		if origin[0] == originUI {
			validOrigin = true
		}
	}

	return validOrigin
}

func createSocketConnection(w http.ResponseWriter, r *http.Request) (conn *websocket.Conn, err error) {
	var upgrader = websocket.Upgrader{}
	upgrader.CheckOrigin = verifyUIOrigin

	conn, err = upgrader.Upgrade(w, r, nil)
	if err == nil {
		var bearerToken string
		var inputMessage []byte
		if _, inputMessage, err = conn.ReadMessage(); err == nil {
			inputString := string(inputMessage)
			var lookingFor = "Bearer="

			bearerIndex := strings.Index(inputString, lookingFor)
			if bearerIndex >= 0 {
				bearerToken = inputString[bearerIndex+len(lookingFor):]
			}
		}

		if len(bearerToken) > 0 {
			r.Header.Add("Authorization", "Bearer "+bearerToken)
		} else {
			err = errors.New("could not find bearer token in document cookie")
		}

	}

	return conn, err
}
