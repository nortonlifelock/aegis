package endpoints

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/jira"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

// PieChartContents is a member of pieChartInfo and must be exported in order to be marshaled
type PieChartContents struct {
	Labels []string `json:"labels"`
	Data   []int    `json:"data"`
}

// pieChartInfo is used to populate the pie chart on the UI dashboard which shows the percentage of tickets in each status
type pieChartInfo struct {
	Contents *PieChartContents
	Updated  time.Time
}

// BurndownPackage holds the information of the timelines in which tickets were closed in a particular project
// it populates the burndown pie char
type BurndownPackage struct {
	MediumStats *BurndownStats `json:"medium"`
	HighStats   *BurndownStats `json:"high"`
	CritStats   *BurndownStats `json:"crit"`
}

// BurndownStats is a member of BurndownPackage and must be exported
type BurndownStats struct {
	ZeroToThirty  int  `json:"zero_thirty"`
	ThirtyToSixty int  `json:"thirty_sixty"`
	SixtyToNinety int  `json:"sixty_ninety"`
	NinetyPlus    int  `json:"ninety_up"`
	Overdue       int  `json:"overdue"`
	Done          bool `json:"done"`
}

// these are the statuses specifically requested for the UI pie chart
type statusRegistry struct {
	open                   string
	reopened               string
	inProgress             string
	resolvedRemediated     string
	resolvedFalsePositive  string
	resolvedDecommissioned string
	resolvedException      string
	closedRemediated       string
	closedFalsePositive    string
	closedDecommission     string
	closedException        string
}

func newStatusRegistry() statusRegistry {
	return statusRegistry{
		"open",
		"reopened",
		"in-progress",
		"resolved-remediated",
		"resolved-falsepositive",
		"resolved-decommissioned",
		"resolved-exception",
		"closed-remediated",
		"closed-false-positive",
		"closed-decommission",
		"closed-exception",
	}
}

func getStatuses() []string {
	var statuses = make([]string, 0)
	statusReflection := reflect.ValueOf(newStatusRegistry())
	for i := 0; i < (&statusReflection).NumField(); i++ {
		valueField := statusReflection.Field(i)
		statuses = append(statuses, valueField.String())

	}

	return statuses
}

type backendStatusAndCustomStatus struct {
	BackendStatus []string `json:"backend_status"`
	CustomStatus  []string `json:"custom_status"`
}

type backendFieldsAndCustomFields struct {
	BackendFields []string `json:"backend_fields"`
	CustomFields  []string `json:"custom_fields"`
}

type exceptionTicket struct {
	domain.Ticket
	cerf   string
	status string
}

func (t *exceptionTicket) CERF() string {
	return t.cerf
}

func (t *exceptionTicket) Status() *string {
	return &t.status
}

func attachExceptionToTicket(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, attachCERFToTicketEndpoint, admin|manager, func(trans *transaction) {
		params := mux.Vars(r)
		var ticket = params[ticketParam]
		var cerf = params[idParam]
		if len(ticket) > 0 && len(cerf) > 0 {
			var engine integrations.TicketingEngine
			engine, trans.err = getTicketingConnection(integrations.JIRA, trans.permission.OrgID())
			if trans.err == nil {
				var jiraTicket domain.Ticket
				if jiraTicket, trans.err = engine.GetTicket(ticket); trans.err == nil {
					_, _, trans.err = engine.UpdateTicket(&exceptionTicket{
						Ticket: jiraTicket,
						cerf:   cerf,
						status: engine.GetStatusMap(jira.StatusClosedException),
					}, fmt.Sprintf("%s added to ticket at request of %s", cerf, sord(trans.user.Username())))

					if trans.err == nil {
						trans.status = http.StatusOK
						trans.obj = "Ticket updated! Please allow a couple minutes for the changes to reflect in the Aegis UI"
					} else {
						(&trans.wrapper).addError(trans.err, backendError)
					}
				} else {
					(&trans.wrapper).addError(trans.err, backendError)
				}
			} else {
				(&trans.wrapper).addError(trans.err, backendError)
			}
		} else {
			(&trans.wrapper).addError(fmt.Errorf("empty CERF or ticket field"), requestFormatError)
		}
	})
}

func getFieldMaps(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getFieldMapsEndpoint, admin|manager, func(trans *transaction) {

		params := mux.Vars(r)
		var id = params[idParam]

		if len(id) > 0 {

			var jiraSourceConfig domain.SourceConfig
			jiraSourceConfig, trans.err = Ms.GetSourceConfigByID(id)
			if trans.err == nil {
				var eng integrations.TicketingEngine
				eng, trans.err = integrations.GetEngine(context.Background(), integrations.JIRA, Ms, logger{}, AppConfig, jiraSourceConfig)
				if trans.err == nil {

					if jiraEng, ok := eng.(*jira.ConnectorJira); ok {
						var customFields []string
						customFields, trans.err = jiraEng.GetFieldsForProject(jiraEng.GetProject(), true)
						if trans.err == nil {

							var backendFieldsAndCustomFields backendFieldsAndCustomFields
							backendFieldsAndCustomFields.CustomFields = customFields
							backendFieldsAndCustomFields.BackendFields = jira.MappableFields

							// TODO CERFEXPIRATION is not returned from the API because it is not a part of the Aegis project
							// TODO we grab the expiration date from the CERF project and use it to populate the field
							const backendCERFExpiration = "Actual Expiration Date"
							backendFieldsAndCustomFields.CustomFields = append(backendFieldsAndCustomFields.CustomFields, backendCERFExpiration)

							trans.status = http.StatusOK
							trans.obj = backendFieldsAndCustomFields

						} else {
							(&trans.wrapper).addError(trans.err, backendError)
						}
					} else {
						(&trans.wrapper).addError(fmt.Errorf("did not load a jira connector"), backendError)
					}
				} else {
					(&trans.wrapper).addError(trans.err, backendError)
				}
			} else {
				(&trans.wrapper).addError(trans.err, databaseError)
			}

		} else {
			(&trans.wrapper).addError(fmt.Errorf("must include an id in the apiRequest"), requestFormatError)
		}

	})
}

// TODO there are status checks in JQLs that aren't here
func getStatusMaps(w http.ResponseWriter, r *http.Request) {

	// TODO change name of endpoint
	executeTransaction(w, r, getStatusMapsEndpoint, admin|manager, func(trans *transaction) {
		params := mux.Vars(r)
		var id = params[idParam]
		if len(id) > 0 {

			var jiraSourceConfig domain.SourceConfig
			jiraSourceConfig, trans.err = Ms.GetSourceConfigByID(id)
			if trans.err == nil {
				var eng integrations.TicketingEngine
				eng, trans.err = integrations.GetEngine(context.Background(), integrations.JIRA, Ms, logger{}, AppConfig, jiraSourceConfig)
				if trans.err == nil {

					if jiraEng, ok := eng.(*jira.ConnectorJira); ok {
						var projectStatuses = make([]string, 0)

						for key := range jiraEng.Statuses {
							projectStatuses = append(projectStatuses, key)
						}

						var retVal backendStatusAndCustomStatus
						retVal.BackendStatus = getStatuses()
						retVal.CustomStatus = projectStatuses

						trans.obj = retVal
						trans.status = http.StatusOK
					} else {
						(&trans.wrapper).addError(fmt.Errorf("did not load a jira connector"), backendError)
					}
				} else {
					(&trans.wrapper).addError(trans.err, backendError)
				}
			} else {
				(&trans.wrapper).addError(trans.err, databaseError)
			}

		} else {
			(&trans.wrapper).addError(fmt.Errorf("must include an id in the apiRequest"), requestFormatError)
		}
	})
}

var pieChartLock = &sync.Mutex{}
var pieChartCache = make(map[string]*pieChartInfo)

var burnDownLock = &sync.Mutex{}
var burnDownCache = make(map[string]*BurndownPackage)

func getJiraUrls(w http.ResponseWriter, r *http.Request) {

	executeTransaction(w, r, getJIRAURLSEndpoint, admin|manager|reporter, func(trans *transaction) {
		var configs []domain.SourceConfig
		configs, trans.err = Ms.GetSourceConfigByNameOrg(integrations.JIRA, trans.permission.OrgID())
		if trans.err == nil {

			var urls = make([]string, 0)
			var seen = make(map[string]bool)

			for _, config := range configs {
				if !seen[config.Address()] {
					urls = append(urls, config.Address())
					seen[config.Address()] = true
				}
			}

			trans.obj = urls
			trans.status = http.StatusOK
		} else {
			(&trans.wrapper).addError(trans.err, databaseError)
		}
	})
}

func getFieldsForJiraProject(w http.ResponseWriter, r *http.Request) {

	executeTransaction(w, r, getFieldsForProjectEndpoint, admin|manager|reporter, func(trans *transaction) {

		var engine integrations.TicketingEngine
		engine, trans.err = getTicketingConnection(integrations.JIRA, trans.permission.OrgID())
		if trans.err == nil {
			if connector, ok := engine.(*jira.ConnectorJira); ok {
				trans.obj, trans.err = connector.GetFieldsForProject(connector.GetProject(), false)
				if stringSlice, isStringSlice := trans.obj.([]string); isStringSlice && trans.err == nil {
					stringSlice = append(stringSlice, "Labels")
					trans.obj = stringSlice
					trans.status = http.StatusOK
				}
			} else {
				trans.err = errors.New("ticketing connection did not appear to be jira connection")
				(&trans.wrapper).addError(trans.err, trans.err.Error())
			}

		} else {
			(&trans.wrapper).addError(trans.err, "error while gathering ticketing connection")
		}

	})
}

func establishBurnDownWebsocket(w http.ResponseWriter, r *http.Request) {
	executeWebsocketTransaction(w, r, burnDownWebSocketEndpoint, func(trans *websocketTransaction) {
		var burnDown = burnDownCache[trans.permission.OrgID()]

		for {

			if burnDown != nil {
				trans.err = trans.connection.WriteJSON(burnDown)
				if trans.err == nil {
					if burnDown.MediumStats.Done && burnDown.HighStats.Done && burnDown.CritStats.Done {
						break
					}
				} else {
					break
				}
			} else {
				burnDown = burnDownCache[trans.permission.OrgID()]
			}

			time.Sleep(time.Second)
		}

		_ = trans.connection.Close()
	})
}

func getCountOfJiraTicketsInStatus(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getTicketStatusCountEndpoint, admin|manager|reporter|reader, func(trans *transaction) {
		var engine integrations.TicketingEngine
		engine, trans.err = getTicketingConnection(integrations.JIRA, trans.permission.OrgID())
		if trans.err == nil {
			if connector, ok := engine.(*jira.ConnectorJira); ok {

				var needsCacheUpdate = false
				pieChartLock.Lock()
				if pieChartCache[trans.permission.OrgID()] == nil {
					needsCacheUpdate = true
				} else if pieChartCache[trans.permission.OrgID()].Contents == nil {
					needsCacheUpdate = true
				} else if time.Since(pieChartCache[trans.permission.OrgID()].Updated) > time.Hour {
					needsCacheUpdate = true
				}
				pieChartLock.Unlock()

				var returnVal *PieChartContents
				if needsCacheUpdate {
					var organization domain.Organization
					var err error
					organization, err = Ms.GetOrganizationByID(trans.permission.OrgID())

					var orgCode string
					if err == nil {
						orgCode = organization.Code()
					}

					returnVal = getPieChartContents(connector, orgCode)
					pieChartLock.Lock()
					pieChartCache[trans.permission.OrgID()] = &pieChartInfo{
						Contents: returnVal,
						Updated:  time.Now(),
					}
					pieChartLock.Unlock()
				} else {
					returnVal = pieChartCache[trans.permission.OrgID()].Contents
				}

				trans.obj = returnVal
				trans.status = http.StatusOK
			} else {
				trans.err = errors.New("ticketing connection did not appear to be jira connection")
				(&trans.wrapper).addError(trans.err, trans.err.Error())
			}

		} else {
			(&trans.wrapper).addError(trans.err, "error while gathering ticketing connection")
		}
	})
}

func getPieChartContents(connector *jira.ConnectorJira, orgCode string) *PieChartContents {
	var wg = &sync.WaitGroup{}
	var labelLock = &sync.Mutex{}
	var statusLabels = make([]string, 0)
	var countLabels = make([]int, 0)

	var statuses = getStatuses()

	for _, status := range statuses {
		wg.Add(1)
		go func(status string) {
			defer wg.Done()
			setCountStatusLabelForStatus(connector, status, labelLock, &statusLabels, &countLabels, orgCode)
		}(status)

	}

	wg.Wait()
	results := &PieChartContents{}
	results.Labels = statusLabels
	results.Data = countLabels
	return results
}

func setCountStatusLabelForStatus(connector *jira.ConnectorJira, status string, labelLock *sync.Mutex, statusLabels *[]string, countLabels *[]int, orgCode string) {
	var count int
	var err error
	count, err = connector.GetCountOfTicketsInStatus(status, orgCode)
	if err == nil {
		labelLock.Lock()
		*statusLabels = append(*statusLabels, fmt.Sprintf("%s - (%d)", status, count))
		*countLabels = append(*countLabels, count)
		labelLock.Unlock()
	} else {
		fmt.Println(err.Error())
	}
}

func calculateBurndownForPriority(sess *jira.ConnectorJira, priority string, SLA time.Duration, stats *BurndownStats, orgCode string) {
	var date = time.Now().AddDate(0, 0, -90).Format("2006/01/02")

	var baseJQL = fmt.Sprintf("project = %s AND \"VRR Priority\" = \"%s\" AND created > \"%s\" AND Org ~ \"%s\" ORDER BY created asc", sess.GetProject(), priority, date, orgCode)
	var ticketChan = sess.GetByCustomJQLChan(baseJQL)

	for {
		tic, ok := <-ticketChan
		if ok {
			if tic != nil {
				if tic.CreatedDate() != nil {
					var daysOpen time.Duration
					if tic.ResolutionDate() != nil && !tic.ResolutionDate().IsZero() {
						daysOpen = tic.ResolutionDate().Sub(*tic.CreatedDate())
					} else {
						daysOpen = time.Now().Sub(*tic.CreatedDate())
					}

					var daysOpenFloat = daysOpen.Hours() / 24

					var overDue = daysOpen > (24 * time.Hour * SLA)

					if overDue {
						stats.Overdue++
					}

					if tic.Status() != nil {
						makeCountForBurnDownChart(tic, daysOpenFloat, stats)
					}

				}
			}
		} else {
			break
		}
	}

	stats.Done = true
}

func makeCountForBurnDownChart(tic domain.Ticket, daysOpenFloat float64, stats *BurndownStats) {
	if strings.Contains(strings.ToLower(*tic.Status()), "closed") {
		if daysOpenFloat < 30 {
			stats.ZeroToThirty++
		} else if daysOpenFloat < 60 {
			stats.ThirtyToSixty++
		} else if daysOpenFloat < 90 {
			stats.SixtyToNinety++
		} else if daysOpenFloat >= 90 {
			stats.NinetyPlus++
		}
	}
}
