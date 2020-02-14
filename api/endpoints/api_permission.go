package endpoints

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nortonlifelock/aegis/internal/database/dal"
	"github.com/nortonlifelock/domain"
	"net/http"
	"strings"
)

func createPermissionsForUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var orgID = params[orgParam]

	executeTransaction(w, r, createUserPermEndpoint, baseAdmin, func(trans *transaction) {
		var ID = params[userParam]
		var userUpdating domain.User
		userUpdating, trans.err = Ms.GetUserByID(ID, trans.permission.OrgID())
		if trans.err == nil {
			_, _, trans.err = Ms.CreateUserPermissions(userUpdating.ID(), orgID)
			if trans.err == nil {
				trans.err = updateUserPermissions(trans.originalBody, userUpdating.ID(), orgID)
				if trans.err == nil {
					trans.status = http.StatusOK
				} else {
					(&trans.wrapper).addError(trans.err, processError)
				}
			} else {
				(&trans.wrapper).addError(trans.err, databaseError)
			}
		} else {
			(&trans.wrapper).addError(trans.err, databaseError)
		}
	})
}

func getPermissionList(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getPermListEndpoint, allAllowed, func(trans *transaction) {
		var generalPermissions = &dal.Permission{}

		var permissionArray []*Permission
		permissionArray, trans.err = dalPermissionToPermissionArray(generalPermissions)
		if trans.err == nil {
			trans.obj = permissionArray
			trans.status = http.StatusOK
		} else {
			(&trans.wrapper).addError(trans.err, processError)
		}
	})
}

func getUserPermissionsByUserOrgID(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getUserPermEndpoint, admin|manager, func(trans *transaction) {
		params := mux.Vars(r)
		var orgID = params[orgParam]
		var userID = params[userParam]

		var user domain.User
		user, trans.err = Ms.GetUserByID(userID, trans.permission.OrgID())
		if trans.err == nil {
			var userPermissions domain.Permission
			userPermissions, trans.err = Ms.GetPermissionByUserOrgID(user.ID(), orgID)
			if trans.err == nil {
				var permissionArray []*Permission
				permissionArray, trans.err = dalPermissionToPermissionArray(userPermissions)
				if trans.err == nil {
					trans.obj = permissionArray
					trans.status = http.StatusOK
				} else {
					(&trans.wrapper).addError(trans.err, processError)
				}
			} else {
				(&trans.wrapper).addError(trans.err, databaseError)
			}
		} else {
			(&trans.wrapper).addError(trans.err, databaseError)
		}

	})
}

func updateUserPermissionsByUserOrgID(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, updateUserPermEndpoint, baseAdmin, func(trans *transaction) {
		params := mux.Vars(r)
		var orgID = params[orgParam]
		var ID = params[userParam]

		var user domain.User
		user, trans.err = Ms.GetUserByID(ID, trans.permission.OrgID())
		if trans.err == nil {
			trans.err = updateUserPermissions(trans.originalBody, user.ID(), orgID)
			if trans.err == nil {
				trans.status = http.StatusOK
			} else {
				(&trans.wrapper).addError(trans.err, processError)
			}
		} else {
			(&trans.wrapper).addError(trans.err, databaseError)
		}

	})
}

// TODO this function must be manually updated as permissions are added to the permissions table
// TODO
func updateUserPermissions(originalBody string, userID string, orgID string) (err error) {
	// original body contains an array of strings, the elements of this array contains the permissions that should be set to true for the user
	// if a permission is not present in this array, the permission should be set to false for this user

	type permissionStruct struct {
		Permissions []string `json:"permissions"`
	}

	// load the permissions from the original body into a struct
	var inPermission = &permissionStruct{}
	err = json.Unmarshal([]byte(originalBody), inPermission)
	if err == nil {
		var truePermissions = strings.Join(inPermission.Permissions, ",")
		truePermissions = strings.ToLower(truePermissions)

		// TODO I don't like this solution. it will require maintenance as more permissions are added. an alternative is not evident to me
		var (
			isAdmin    = strings.Index(truePermissions, "admin") >= 0
			isManager  = strings.Index(truePermissions, "manager") >= 0
			isReader   = strings.Index(truePermissions, "reader") >= 0
			isReporter = strings.Index(truePermissions, "reporter") >= 0
		)

		_, _, err = Ms.UpdatePermissionsByUserOrgID(
			userID, orgID,

			isAdmin,
			isManager,
			isReader,
			isReporter)
	}

	return err
}
