package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {

	var err error

	apiPath := flag.String("url", "", "The output file for the resulting file")
	user := flag.String("user", "", "The output file for the resulting file")
	deleteTemp := flag.Bool("d", false, "Flag to delete the templates instead of just printing them")
	templatePrefix := flag.String("prefix", strconv.Itoa(time.Now().Year())+"-", "Delete templates with an ID that starts with this")

	flag.Parse()

	if *deleteTemp {
		fmt.Println("Delete flag set, will delete returned templates")
	} else {
		fmt.Println("Delete flag [-d] not included, will only print returned templates")
	}

	fmt.Printf("Using template prefix [%s]\n", *templatePrefix)

	if apiPath != nil && len(*apiPath) > 0 {
		if user != nil && len(*user) > 0 {

			var password []byte
			fmt.Println("Enter Nexpose Password: ")

			if password, err = terminal.ReadPassword(0); true {

				fmt.Println("Pulling templates from Nexpose...")

				var req *http.Request
				if req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/3/scan_templates", *apiPath), nil); err == nil {

					req.SetBasicAuth(*user, string(password))

					tr := &http.Transport{
						TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
						TLSHandshakeTimeout:   1 * time.Minute,
						ResponseHeaderTimeout: 2 * time.Minute,
						ExpectContinueTimeout: 1 * time.Second,
					}

					client := &http.Client{Transport: tr, Timeout: 5 * time.Minute}

					var resp *http.Response
					if resp, err = client.Do(req); err == nil {
						if resp != nil {
							defer resp.Body.Close()
						}

						var body []byte
						var records = records{}
						if body, err = ioutil.ReadAll(resp.Body); err == nil {

							if err = json.Unmarshal(body, &records); err == nil {
								deletedCount, failedDelete := deleteTemplates(records, templatePrefix, deleteTemp, apiPath, user, password, client)
								fmt.Printf("TOTAL Templates: %v\nDELETED Templates: %d FAILED: %d\n", len(records.Resources), deletedCount, failedDelete)
							} else {
								fmt.Println(fmt.Sprintf("Failed to load templates"))
							}
						} else {
							fmt.Println(fmt.Sprintf("Unable to parse response body: %s\n", err))
						}
					} else {
						fmt.Println(err.Error())
					}
				} else {
					fmt.Println(err.Error())
				}
			} else {
				fmt.Println("Error while reading password from terminal")
			}
		} else {
			fmt.Println("Username (-user) is Required")
		}
	} else {
		fmt.Println("API Path (-url) is Required")
	}

	fmt.Println("Deletion complete")
}

func deleteTemplates(r records, templatePrefix *string, deleteTemp *bool, apiPath *string, user *string, password []byte, client *http.Client) (int, int) {
	wg := &sync.WaitGroup{}
	lock := &sync.Mutex{}
	deletedCount := 0
	failedDelete := 0
	index := 0
	for key := range r.Resources {

		if strings.Index(r.Resources[key].Name, *templatePrefix) == 0 {

			if *deleteTemp {

				wg.Add(1)
				go func(tempID string) {
					defer wg.Done()
					if err := deleteTemplate(apiPath, tempID, user, password, client); err == nil {
						lock.Lock()
						deletedCount++
						lock.Unlock()
					} else {
						lock.Lock()
						failedDelete++
						lock.Unlock()
					}
				}(r.Resources[key].ID)

				index++
				if index%100 == 0 {
					wg.Wait()
				}

			} else {
				fmt.Printf("%s\n", r.Resources[key].ID)
			}
		}
	}
	wg.Wait()
	return deletedCount, failedDelete
}

func deleteTemplate(apiPath *string, tempID string, user *string, password []byte, client *http.Client) (err error) {
	if req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/3/scan_templates/%s", *apiPath, tempID), nil); err == nil {
		req.SetBasicAuth(*user, string(password))

		resp, err := client.Do(req)
		if resp != nil {
			resp.Body.Close()
		}

		if err == nil {
			fmt.Printf("Deleted %s\n", tempID)
		} else {
			fmt.Printf("[+] Error while deleting template - %s\n", err.Error())
			return err
		}
	} else {
		fmt.Printf("[+] Error while creating request - %s\n", err.Error())
		return err
	}

	return err
}

type records struct {
	Resources []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"resources"`
}
