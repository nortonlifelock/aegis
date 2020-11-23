package nexpose

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nortonlifelock/log"
	"github.com/pkg/errors"
)

// GetScans loads scans from the nexpose based on the active flag that is passed in.
// The returned data can be sorted by passing a sort string that uses the formatting available in the nexpose API.
func (a *Session) GetScans(ctx context.Context, active bool, sort string) (<-chan *Scan, error) {
	var summaries = make(chan *Scan)
	var err error

	go func(summaries chan<- *Scan) {
		defer handleRoutinePanic(a.lstream)
		defer close(summaries)

		fields := map[string]string{"sort": sort}
		if active {
			fields["active"] = "true"
		}

		var data <-chan interface{}
		if data, err = a.getPagedData(ctx, &PageOfScan{}, apiGetScans, fields); err == nil {
			for {
				select {
				case <-ctx.Done():
					return
				case d, ok := <-data:
					if ok {
						var summary Scan
						if summary, ok = d.(Scan); ok {
							summaries <- &summary
						} else {
							a.lstream.Send(log.Error("unable to cast paged return as detection", err))
						}
					} else {
						return
					}
				}
			}
		} else {
			a.lstream.Send(log.Error("error executing paged call against nexpose", err))
		}
	}(summaries)

	return summaries, err
}

// GetScan loads the scan data for the scan ID passed in
func (a *Session) GetScan(id string) (scan *Scan, err error) {

	scan = &Scan{}
	err = a.execute(http.MethodGet, fmt.Sprintf(apiGetScan, id), nil, nil, scan)

	return scan, err
}

// CreateScan creates a new scan at the specified site with corresponding parameters as passed through engineID, templateID,
// hostIPs and name
func (a *Session) CreateScan(siteID, engineID, templateID, name string, hostIPs []string) (scanID string, err error) {

	if len(siteID) > 0 {
		if len(engineID) > 0 {
			if len(templateID) > 0 {
				if len(hostIPs) > 0 {
					if len(name) > 0 {

						var postbody = &AdhocScan{
							EngineID:   engineID,
							TemplateID: templateID,
							Hosts:      hostIPs,
							Name:       name,
						}

						var body []byte
						if body, err = json.Marshal(postbody); err == nil {

							ref := &Reference{}
							if err = a.execute(http.MethodPost, fmt.Sprintf(apiPostSiteScan, encode(siteID)), nil, bytes.NewBuffer(body), ref); err == nil {

								idInterface := ref.ID

								// Return the ID of the newly created scan
								switch idInterface.(type) {
								case float32, float64:
									scanID = fmt.Sprintf("%f", idInterface)
								case string:
									scanID = fmt.Sprintf("%s", idInterface)
								case int, int32, int64:
									scanID = fmt.Sprintf("%d", idInterface)
								default:
									scanID = fmt.Sprintf("%v", idInterface)
								}
							}
						}
					} else {
						err = errors.New("the scan name cannot be empty")
					}
				} else {
					err = errors.New("the list of hosts to scan cannot be empty")
				}
			} else {
				err = errors.New("the template id cannot be empty")
			}
		} else {
			err = errors.New("the scan engine id must be greater than 0")
		}
	} else {
		err = errors.New("the site id for a scan must be greater than 0")
	}

	return scanID, err
}

// GetSites returns a list of sites from Nexpose
func (a *Session) GetSites() (err error) {
	err = a.execute(http.MethodGet, apiGetSites, nil, nil, nil)
	return err
}

// GetEngineForSite loads the scan engine for the site ID that is passed in
func (a *Session) GetEngineForSite(siteID string) (engine *Engine, err error) {
	engine = &Engine{}
	err = a.execute(http.MethodGet, fmt.Sprintf(apiGetSiteScanEngine, encode(siteID)), nil, nil, engine)
	return engine, err
}

// GetScanTemplates loads the scan templates from the nexpose API
func (a *Session) GetScanTemplates() (templates []*ScanTemplate, err error) {

	templateResult := &ScanTemplates{}
	if err = a.execute(http.MethodGet, apiGetScanTemplates, nil, nil, templateResult); err == nil {
		templates = templateResult.Templates
	}

	return templates, err
}

// GetScanTemplate loads a specific scan template from the nexpose api identified by the template ID passed in
func (a *Session) GetScanTemplate(templateID string) (template *ScanTemplate, err error) {

	template = &ScanTemplate{}
	err = a.execute(http.MethodGet, fmt.Sprintf(apiGetScanTemplate, encode(templateID)), nil, nil, template)

	return template, err
}

// DuplicateScanTemplate creates a duplicate scan template of the scan templateID passed in with a unique identifier
func (a *Session) DuplicateScanTemplate(templateID string) (newTemplateID string, err error) {
	if len(templateID) > 0 {
		var template *ScanTemplate
		if template, err = a.GetScanTemplate(encode(templateID)); err == nil {
			newID := uuid.New()

			// Set the new name of the template
			template.Name = fmt.Sprintf("%d_%s", time.Now().Year(), newID.String())
			newTemplateID, err = a.CreateScanTemplate(template)
		} else {
			err = fmt.Errorf("error while gathering scan template - %v", err.Error())
		}
	} else {
		err = fmt.Errorf("empty template ID")
	}

	return newTemplateID, err
}

// DuplicateScanTemplateWVulns creates a duplicate scan template of the scan templateID passed in with a unique identifier
// and returns that identifier to the requestor after a successful creation
func (a *Session) DuplicateScanTemplateWVulns(ctx context.Context, templateID string, vulnerabilityIDs []string) (newTemplateID string, err error) {

	if len(templateID) > 0 {
		var template *ScanTemplate
		if template, err = a.GetScanTemplate(encode(templateID)); err == nil {

			var checks <-chan string
			if checks, err = a.GetChecksIDsForVulnerabilities(ctx, vulnerabilityIDs); err == nil {

				var newChecks []string
				var seen = make(map[string]bool)
				var seenLock = sync.Mutex{}

				func() {
					for {
						select {
						case <-ctx.Done():
							return
						case check, ok := <-checks:
							if ok {
								seenLock.Lock()
								isNew := !seen[check]
								seenLock.Unlock()

								if isNew {
									// Append the incomming checks to this template
									newChecks = append(newChecks, check)
								}
							} else {

								// Channel is closed so all of the checks have been loaded
								if len(newChecks) > 0 {

									newID := uuid.New()

									// Set the new name of the template
									template.Name = fmt.Sprintf("%s_%s", time.Now().Format("2006-01-02"), newID.String())
									template.Checks.Individual.Enabled = newChecks

									newTemplateID, err = a.CreateScanTemplate(template)
								} else {
									err = errors.New("no valid checks were loaded for the vulnerabilities")
								}

								return
							}
						}
					}
				}()
			}
		}
	} else {
		err = fmt.Errorf("empty template ID")
	}

	return newTemplateID, err
}

// CreateScanTemplate creates a scan template using the template passed in and returns the templateID of the
// new template back to the user
func (a *Session) CreateScanTemplate(template *ScanTemplate) (templateID string, err error) {

	if template != nil {

		var body []byte
		if body, err = json.Marshal(template); err == nil {

			ref := &Reference{}
			if err = a.execute(http.MethodPost, apiPostScanTemplates, nil, bytes.NewBuffer(body), ref); err == nil {

				// Return the ID of the newly created scan
				templateID = fmt.Sprintf("%v", ref.ID)
			}
		}
	} else {
		err = errors.New("nil template passed to create scan template")
	}

	return templateID, err
}

// DeleteScanTemplate deletes the scan template based on the templateID that is passed in
func (a *Session) DeleteScanTemplate(templateID string) (err error) {

	if len(templateID) > 0 {
		err = a.execute(http.MethodDelete, fmt.Sprintf(apiDeleteScanTemplate, encode(templateID)), nil, nil, nil)
	} else {
		err = errors.New("template ID for deletion cannot be empty")
	}

	return err
}
