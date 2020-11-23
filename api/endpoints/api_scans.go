package endpoints

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/pkg/errors"
)

func getAGForTicket(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getAgsEndpoint, allAllowed, func(trans *transaction) {
		params := mux.Vars(r)
		var ip = params[ticketParam]

		if len(ip) > 0 {
			var matched bool
			matched, trans.err = regexp.MatchString(ipRegex, ip)
			if trans.err == nil {
				if matched {
					var ags []domain.AssignmentGroup
					ags, trans.err = Ms.GetAssignmentGroupByOrgIP(trans.permission.OrgID(), ip)
					if trans.err == nil {
						if ags != nil && len(ags) > 0 {
							trans.status = http.StatusOK
							trans.obj = ags[0].GroupName()
						} else {
							(&trans.wrapper).addError(errors.Errorf("could not find assignment group entry for ip [%s] and organization [%s]", ip, trans.permission.OrgID()), databaseError)
						}
					} else {
						(&trans.wrapper).addError(trans.err, databaseError)
					}
				} else {
					(&trans.wrapper).addError(errors.Errorf("ip [%s] did not pass regex validation", ip), requestFormatError)
				}
			} else {
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			(&trans.wrapper).addError(errors.Errorf("empty ip address passed to GetAgForTicket"), requestFormatError)
		}
	})
}

func getScansForScanner(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getScansEndpoint, manager|admin|reporter, func(trans *transaction) {
		// TODO permission
		params := mux.Vars(r)

		var scanSummaries []domain.ScanSummary
		var scanner = strings.ToLower(params[scannerParam])
		scanSummaries, trans.err = Ms.GetScanSummariesBySourceName(trans.permission.OrgID(), scanner)

		if trans.err == nil {
			if scanSummaries != nil {
				trans.obj = toScanSummaryDtoSlice(scanSummaries, trans.permission.OrgID())
			}
			trans.status = http.StatusOK
		} else {
			(&trans.wrapper).addError(trans.err, databaseError)
		}
	})
}

func createWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	executeWebsocketTransaction(w, r, createScanWebsocketEndpoint, func(trans *websocketTransaction) {
		for {
			var summaries []domain.ScanSummary
			summaries, trans.err = Ms.GetRecentlyUpdatedScanSummaries(trans.permission.OrgID())
			if trans.err == nil {
				var summariesDto = toScanSummaryDtoSlice(summaries, trans.permission.OrgID())
				trans.err = trans.connection.WriteJSON(summariesDto)
				if trans.err != nil {
					fmt.Printf("Error while sending websocket connection, closing connection...\n")
					break
				}
			}

			time.Sleep(time.Second * 30)
		}

		_ = trans.connection.Close()
	})
}
