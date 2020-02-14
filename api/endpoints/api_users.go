package endpoints

import (
	"fmt"
	"github.com/nortonlifelock/aegis/internal/database/dal"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nortonlifelock/crypto"
	"github.com/nortonlifelock/domain"
	"github.com/pkg/errors"
)

// because a user may belong to many organizations, we do not encrypt the user fields using an organization-specific key
// instead, we encrypt using the application-level KMS key
func encryptOrDecryptUserFields(first, last, email string, encryptOrDecrypt int) (outFirst, outLast, outEmail string, err error) {
	var client crypto.Client
	client, err = crypto.NewEncryptionClientWithDirectKey(crypto.KMS, EncryptionKey, KMSRegion())

	if err == nil {
		var encryptOrDecryptMethod func(string) (string, error)
		if encryptOrDecrypt == crypto.EncryptMode {
			encryptOrDecryptMethod = client.Encrypt
		} else {
			encryptOrDecryptMethod = client.Decrypt
		}

		outFirst, err = encryptOrDecryptMethod(first)
		if err == nil {
			outLast, err = encryptOrDecryptMethod(last)
			if err == nil {
				outEmail, err = encryptOrDecryptMethod(email)
			}
		}
	}

	return outFirst, outLast, outEmail, err
}

// Encrypted fields: First, Last, Email
func encryptOrDecryptUser(input domain.User, encryptOrDecrypt int) (output domain.User, err error) {
	var outFirst, outLast, outEmail string
	outFirst, outLast, outEmail, err = encryptOrDecryptUserFields(input.FirstName(), input.LastName(), input.Email(), encryptOrDecrypt)
	if err == nil {
		output = &dal.User{
			IDvar:         input.ID(),
			Usernamevar:   input.Username(),
			IsDisabledvar: input.IsDisabled(),
			FirstNamevar:  outFirst,
			LastNamevar:   outLast,
			Emailvar:      outEmail,
		}
	}

	return output, err
}

func getMyName(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getUsersNameEndpoint, allAllowed, func(trans *transaction) {
		trans.obj = fmt.Sprintf("%s %s", trans.user.FirstName(), trans.user.LastName())
		trans.status = http.StatusOK
	})
}

func updateUser(w http.ResponseWriter, r *http.Request) {

	executeTransaction(w, r, updateUserEndpoint, admin|manager, func(trans *transaction) {
		if trans.endpoint.verify() {
			if this, isUser := trans.endpoint.(*User); isUser {
				trans.obj, trans.status, trans.err = this.update(trans.user, trans.permission, trans.originalBody)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = fmt.Errorf("tried to pass a non-user as an user")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = fmt.Errorf("invalid length of format of apiRequest fields")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func deleteUser(w http.ResponseWriter, r *http.Request) {

	executeTransaction(w, r, deleteUserEndpoint, admin, func(trans *transaction) {
		if trans.endpoint.verify() {
			if this, isUser := trans.endpoint.(*User); isUser {
				trans.obj, trans.status, trans.err = this.delete(trans.user, trans.permission)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = fmt.Errorf("tried to pass a non-user as a user")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = fmt.Errorf("invalid length of format of apiRequest fields")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

// TODO need to make sure that a person in the same organization does not have the same username/email
func createUser(w http.ResponseWriter, r *http.Request) {

	executeTransaction(w, r, createUserEndpoint, admin|manager, func(trans *transaction) {
		if trans.endpoint.verify() {
			if this, isUser := trans.endpoint.(*User); isUser {
				trans.obj, trans.status, trans.err = this.create(trans.user, trans.permission)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = fmt.Errorf("tried to pass a non-user as a user")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = fmt.Errorf("invalid length of format of apiRequest fields")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getAllUsersEndpoint, admin|manager, func(trans *transaction) {
		trans.obj, trans.status, trans.totalRecords, trans.err = readUsers(trans.permission)
		(&trans.wrapper).addError(trans.err, processError)
	})
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getUserByIDEndpoint, admin|manager, func(trans *transaction) {
		params := mux.Vars(r)
		var ID = params[idParam]
		var user domain.User
		user, trans.err = Ms.GetUserByID(ID, trans.permission.OrgID())
		if trans.err == nil {
			var decryptedUser *User
			decryptedUser, trans.err = toUserDto(user)
			if trans.err == nil {
				trans.obj = decryptedUser
				trans.status = http.StatusOK
			} else {
				(&trans.wrapper).addError(trans.err, backendError)
			}
		} else {
			(&trans.wrapper).addError(trans.err, databaseError)
		}
	})
}

func readUsers(permission domain.Permission) (usersDto []*User, status int, totalRecords int, err error) {
	status = http.StatusBadRequest

	var users []domain.User
	users, err = Ms.GetUsersByOrg(permission.OrgID())
	if err == nil {
		status = http.StatusOK
		usersDto = toUserDtoSlice(users)
		totalRecords = len(usersDto)
	} else {
		err = fmt.Errorf("error while loading users from database - %s", err.Error())
	}

	return usersDto, status, totalRecords, err
}

func (u *User) create(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest

	if len(u.Username) > 0 {
		if len(u.Email) > 0 {
			if len(u.FirstName) > 0 {
				if len(u.LastName) > 0 {

					var encryptFirst, encryptLast, encryptEmail string
					encryptFirst, encryptLast, encryptEmail, err = encryptOrDecryptUserFields(strings.TrimSpace(u.FirstName), strings.TrimSpace(u.LastName), u.Email, crypto.EncryptMode)
					if err == nil {
						_, _, err = Ms.CreateUser(u.Username, encryptFirst, encryptLast, encryptEmail)
						if err == nil {

							var createdUser domain.User
							createdUser, err = Ms.GetUserByUsername(u.Username)
							if err == nil {
								if createdUser != nil {
									status = http.StatusOK
									generalResp.Response = createdUser.ID()
								} else {
									err = fmt.Errorf("error while grabbing newly created user")
									status = http.StatusBadRequest
								}
							} else {
								err = fmt.Errorf("error while creating new user")
								status = http.StatusBadRequest
							}

						} else {
							err = fmt.Errorf("error while creating user - %s", err.Error())
							status = http.StatusBadRequest
						}
					} else {
						err = fmt.Errorf("error while encrypting user fields - %s", err.Error())
						status = http.StatusBadRequest
					}
				} else {
					err = fmt.Errorf("must include a lastName in your apiRequest")
					status = http.StatusBadRequest
				}
			} else {
				err = fmt.Errorf("must include a firstName in your apiRequest")
				status = http.StatusBadRequest
			}
		} else {
			err = fmt.Errorf("must include an email in your apiRequest")
			status = http.StatusBadRequest
		}
	} else {
		err = fmt.Errorf("must include the username in your apiRequest")
		status = http.StatusBadRequest
	}

	return generalResp, status, err
}

func (u *User) update(user domain.User, permission domain.Permission, originalBody string) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest

	if len(u.ID) > 0 {
		var existingUser domain.User
		existingUser, err = Ms.GetUserByID(u.ID, permission.OrgID())
		if err == nil {

			var disabled = u.IsDisabled

			// We can't pass an "empty" boolean to keep the existing value so we must pass the existing user's value
			if strings.Index(originalBody, "isdisabled") < 0 {
				disabled = existingUser.IsDisabled()
			}

			var encryptFirst, encryptLast, encryptEmail string
			encryptFirst, encryptLast, encryptEmail, err = encryptOrDecryptUserFields(strings.TrimSpace(u.FirstName), strings.TrimSpace(u.LastName), u.Email, crypto.EncryptMode)
			if err == nil {
				// If an empty string is passed to u stored procedure that field is not updated
				_, _, err = Ms.UpdateUserByID(u.ID, encryptFirst, encryptLast, encryptEmail, disabled)

				if err == nil {
					status = http.StatusOK
					generalResp.Response = "user updated"
				} else {
					err = fmt.Errorf("error while updating user")
				}
			} else {
				err = fmt.Errorf("error while encrypting fields - %s", err.Error())
			}

		} else {
			err = fmt.Errorf("error while gathering details of existing user")
		}
	} else {
		err = fmt.Errorf("must include the id of the user you'd like to update")
	}

	return generalResp, status, err
}

// TODO should supply UUID for deletion
func (u *User) delete(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest

	_, _, err = Ms.DeleteUserByUsername(u.Username)
	if err == nil {
		status = http.StatusOK
		generalResp.Response = "user deleted"
	} else {
		err = errors.New("error while deleting user")
	}

	return generalResp, status, err
}

func (u *User) verify() (valid bool) {
	valid = false
	if len(u.ID) >= 0 {
		if len(u.Username) >= lowerBoundNameLen && len(u.Username) <= upperBoundNameLen {
			if len(u.FirstName) >= lowerBoundNameLen && len(u.FirstName) <= upperBoundNameLen {
				if len(u.LastName) >= lowerBoundNameLen && len(u.LastName) <= upperBoundNameLen {
					if len(u.Email) >= lowerBoundEmailLen && len(u.Email) <= upperBoundNameLen {
						if check(lettersRegex, u.FirstName) {
							if check(lettersRegex, u.LastName) {
								if check(emailRegex, u.Email) {
									if check(usernameRegex, u.Username) {
										valid = true
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return valid
}
