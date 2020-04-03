package endpoints

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nortonlifelock/aegis/internal/implementations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/files"
	"github.com/pkg/errors"
)

const (
	// allowedUploadsPerUser defines the amount of uploads a user is allowed to do within a 10 minute windows
	allowedUploadsPerUser = 64
)

var (
	userConnectionMapLock   = &sync.Mutex{}
	mapUsernameToConnection = make(map[string]*websocket.Conn)

	userFileUploadCountLock      = &sync.Mutex{}
	mapUsernameToFileUploadCount = make(map[string]int)
)

func bulkUpdateUpload(w http.ResponseWriter, r *http.Request) {

	executeTransaction(w, r, bulkUpdateUploadEndpoint, admin|reporter|manager, func(trans *transaction) {
		var req = &apiRequest{}
		trans.err = json.Unmarshal([]byte(trans.originalBody), req)
		if trans.err == nil {
			trans.wrapper, trans.status, trans.obj = createBulkUpdateUpload(req, trans.user)

			if trans.wrapper.detailedError == nil {

				// Track how many files the user uploads every 10 minutes for limiting
				go func(username string) {
					defer handleRoutinePanic(trans, w, bulkUpdateUploadEndpoint)
					addFileCountForUser(username)
				}(*trans.user.Username())

			}
		} else {
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func addFileCountForUser(username string) {
	userFileUploadCountLock.Lock()
	mapUsernameToFileUploadCount[username]++
	userFileUploadCountLock.Unlock()

	time.Sleep(time.Minute * 10)

	userFileUploadCountLock.Lock()
	mapUsernameToFileUploadCount[username]--
	userFileUploadCountLock.Unlock()
}

// TODO could include slack token in apiRequest for slack updates - would make part of job history payload
func createBulkUpdateJobHistory(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, bulkUpdateEndpoint, admin|reporter|manager, func(trans *transaction) {
		var req = &apiRequest{}
		trans.err = json.Unmarshal([]byte(trans.originalBody), req)
		if trans.err == nil {
			trans.wrapper, trans.status, trans.obj = createBulkUpdateJobIfPermitted(req, trans.user, trans.permission)
		} else {
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func createBulkUpdateUpload(request *apiRequest, user domain.User) (wrapper errorWrapper, status int, obj interface{}) {
	var err error
	status = http.StatusBadRequest

	if request != nil && request.BulkUpdateFile != nil {

		var count int
		userFileUploadCountLock.Lock()
		count = mapUsernameToFileUploadCount[*user.Username()]
		userFileUploadCountLock.Unlock()

		if count < allowedUploadsPerUser {
			if len(request.BulkUpdateFile.Filename) > 0 && len(request.BulkUpdateFile.Contents) > 0 {

				if len(request.BulkUpdateFile.Contents) < int(math.Pow(2, 20)) {
					if filepath.Ext(request.BulkUpdateFile.Filename) == ".csv" {

						var updateContents = request.BulkUpdateFile.Contents
						updateContents = strings.Replace(updateContents, "\\r", "\r", -1)
						updateContents = strings.Replace(updateContents, "\\n", "\n", -1)

						_, err = csv.NewReader(strings.NewReader(updateContents)).ReadAll()

						if err == nil {

							var path = fmt.Sprintf("%s/bulk_updates", WorkingDir)
							var finishedPath = fmt.Sprintf("%s/bulk_updates/finished", WorkingDir)

							if _, err = os.Stat(path); os.IsNotExist(err) {
								err = os.Mkdir(path, os.ModeDir|os.ModePerm)
								if err != nil {
									fmt.Printf("could not initialize %s", path)
								}
							}

							if _, err = os.Stat(finishedPath); os.IsNotExist(err) {
								err = os.Mkdir(finishedPath, os.ModeDir|os.ModePerm)
								if err != nil {
									fmt.Printf("could not initialize %s", finishedPath)
								}
							}

							err = files.WriteFile(fmt.Sprintf("%s/%s", path, request.BulkUpdateFile.Filename), updateContents)

							if err == nil {
								status = http.StatusOK
								obj = "Upload success"
							} else {
								(&wrapper).addError(err, backendError)
							}

						} else {
							(&wrapper).addError(err, requestFormatError)
						}

					} else {
						(&wrapper).addError(errors.Errorf("file [%s] did not have a .csv extension", request.BulkUpdateFile.Filename), requestFormatError)
					}
				} else {
					(&wrapper).addError(errors.Errorf("uploaded file was over 500kb"), requestFormatError)
				}
			} else {
				err = errors.Errorf("the apiRequest must contain both the filename and its contents")
				(&wrapper).addError(err, requestFormatError)
			}
		} else {
			(&wrapper).addError(errors.Errorf("only %d files are allowed to be uploaded by each user every 10 minutes", allowedUploadsPerUser), requestFormatError)
		}
	} else {
		err = errors.Errorf("apiRequest did not appear to be for a bulk update")
		(&wrapper).addError(err, requestFormatError)
	}

	return wrapper, status, obj
}

// first checks the amount of recent uploads before creating the job history
func createBulkUpdateJobIfPermitted(request *apiRequest, user domain.User, permission domain.Permission) (wrapper errorWrapper, status int, obj interface{}) {
	status = http.StatusBadRequest

	if len(request.BulkUpdateJob.Filenames) <= allowedUploadsPerUser {
		if request != nil && request.BulkUpdateJob != nil {
			if len(request.BulkUpdateJob.Filenames) > 0 && len(request.BulkUpdateJob.ServiceURL) > 0 {

				var filesExist bool
				var err error

				filesExist, err = verifyFilesExist(request.BulkUpdateJob.Filenames)

				if filesExist {
					if err == nil {
						status, obj = createBulkUpdateJobInDb(request, user, permission, status, wrapper)
					} else {
						wrapper.addError(err, requestFormatError)
					}
				} else {
					(&wrapper).addError(errors.Errorf("one or more of the files in the apiRequest were not found on the server"), requestFormatError)
				}
			} else {
				(&wrapper).addError(errors.Errorf("bulk update job apiRequest did not provide all required components"), requestFormatError)
			}
		} else {
			(&wrapper).addError(errors.Errorf("apiRequest did not appear to be a bulk update job apiRequest"), requestFormatError)
		}
	} else {
		(&wrapper).addError(errors.Errorf("user uploaded too many files in a 10 minute timeframe"), requestFormatError)
	}

	return wrapper, status, obj
}

func createBulkUpdateJobInDb(request *apiRequest, user domain.User, permission domain.Permission, status int, wrapper errorWrapper) (int, interface{}) {
	var err error
	var obj interface{}

	var payloadObj = &implementations.BulkUpdatePayload{}
	payloadObj.Filenames = request.BulkUpdateJob.Filenames
	payloadObj.UsernameOfRequester = sord(user.Username())
	payloadObj.ServiceURL = request.BulkUpdateJob.ServiceURL
	payloadObj.OrgID = permission.OrgID()
	var payloadByte []byte
	if payloadByte, err = json.Marshal(payloadObj); err == nil {

		var baseJob domain.JobRegistration
		if baseJob, err = Ms.GetJobsByStruct("BulkUpdateJob"); err == nil {
			if baseJob != nil {

				var configs []domain.JobConfig
				configs, err = Ms.GetJobConfigByOrgIDAndJobID(permission.OrgID(), baseJob.ID())

				if err == nil {
					if configs != nil && len(configs) > 0 {

						_, _, err = Ms.CreateJobHistory(
							configs[0].JobID(),
							configs[0].ID(),
							domain.JobStatusPending,
							iord(configs[0].PriorityOverride()),
							"",
							0,
							string(payloadByte),
							"",
							time.Now().UTC(),
							sord(user.Username()),
						)
						if err == nil {
							status = http.StatusOK
							var resp = &GeneralResp{}
							resp.Message = "File uploaded and job history created"
							obj = resp
						} else {
							(&wrapper).addError(err, databaseError)
						}

					} else {
						err = errors.Errorf("could not find config with org id [%s] and job id [%d]", permission.OrgID(), baseJob.ID())
						(&wrapper).addError(err, requestFormatError)
					}
				} else {
					err = errors.Errorf("error while retreiving job config from database")
					(&wrapper).addError(err, databaseError)
				}
			} else {
				err = errors.Errorf("multiple bulk update jobs found in database")
				(&wrapper).addError(err, databaseError)
			}
		} else {
			err = errors.Errorf("could not find bulk update job in database")
			(&wrapper).addError(err, databaseError)
		}
	}
	return status, obj
}

func verifyFilesExist(files []string) (exist bool, err error) {

	if len(files) > 0 {
		exist = true

		for _, file := range files {
			var path = fmt.Sprintf("%s/bulk_updates/%s", WorkingDir, file)

			if _, err = os.Stat(path); os.IsNotExist(err) {
				exist = false
			}
		}
	}

	return exist, err
}

// maybe have a method that creates the connection which also takes a function as an argument for what to do with the connection
func createWebsocketConnectionBulkUpdate(w http.ResponseWriter, r *http.Request) {
	executeWebsocketTransaction(w, r, burnDownWebSocketEndpoint, func(trans *websocketTransaction) {
		type msg struct {
			Message string `json:"message"`
		}

		var message = &msg{}

		if trans.user.Username() != nil {
			message.Message = "server established connection with client"
			userConnectionMapLock.Lock()
			mapUsernameToConnection[*trans.user.Username()] = trans.connection
			userConnectionMapLock.Unlock()
		} else {
			message.Message = "you do not have permissions to perform that action"
		}

		trans.err = trans.connection.WriteJSON(message)
	})
}

func listenForMessagesFromJobToSendToClient(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{}

	c, err := upgrader.Upgrade(w, r, nil)
	if err == nil {
		defer func() {
			_ = c.Close()
		}()

		for {
			// receive message from bulkUpdateJob and parse it
			var messageByte []byte
			_, messageByte, err = c.ReadMessage()
			if err == nil {
				var message = &implementations.BulkUpdateMessage{}
				err = json.Unmarshal(messageByte, message)
				if err == nil {
					// find the destination user from the message and grab their connection
					if len(message.User) > 0 {
						userConnectionMapLock.Lock()
						connection := mapUsernameToConnection[message.User]
						userConnectionMapLock.Unlock()
						message.User = ""
						// forward the message to use user
						if connection != nil {
							_ = connection.WriteJSON(message)
						}
					} else {
						fmt.Println("Request with empty user came from Aegis job")
					}

				}
			} else {
				fmt.Printf("could not parse message: %s\n", err.Error())
				break
			}
		}

	} else {
		fmt.Println("upgrade:", err)
	}
}
