package implementations

import (
	"context"
	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
)

func loadTickets(lstream log.Logger, ticketing integrations.TicketingEngine, ticketSlice []string) (tickets []domain.Ticket, err error) {
	tickets = make([]domain.Ticket, 0)

	for _, t := range ticketSlice {

		// Get ticket to be processes
		var ticket domain.Ticket
		if ticket, err = ticketing.GetTicket(t); err == nil {

			if ticket != nil {
				tickets = append(tickets, ticket)
			} else {
				lstream.Send(log.Warningf(err, "Unable to Load Ticket [%s]", t))
			}

		} else {
			lstream.Send(log.Errorf(err, "Error while loading ticket [%s]: %s", t, err.Error()))
		}
	}

	return tickets, err
}

// validInputsMultipleSources is used for jobs that may use multiple in/out source configs
func validInputsMultipleSources(ctx context.Context, id string, appConfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (context.Context, string, domain.Config, domain.DatabaseConnection, log.Logger, string, domain.JobConfig, []domain.SourceConfig, []domain.SourceConfig, bool) {
	var valid bool

	if len(id) > 0 {
		if appConfig != nil {
			if db != nil {
				if lstream != nil {
					if jobConfig != nil {
						if inSource != nil && outSource != nil {

							if len(inSource) >= 1 && len(outSource) >= 1 {

								if len(payload) > 0 {

									if ctx == nil {
										ctx = context.Background()
									}

									valid = true
								}
							}
						}
					}
				}
			}
		}
	}

	return ctx, id, appConfig, db, lstream, payload, jobConfig, inSource, outSource, valid
}

// validInputs is used for jobs that only use a single in/out source config
func validInputs(ctx context.Context, id string, appConfig domain.Config, db domain.DatabaseConnection, lstream log.Logger, payload string, jobConfig domain.JobConfig, inSource []domain.SourceConfig, outSource []domain.SourceConfig) (context.Context, string, domain.Config, domain.DatabaseConnection, log.Logger, string, domain.JobConfig, domain.SourceConfig, domain.SourceConfig, bool) {
	var valid bool
	var firstSourceIn domain.SourceConfig
	var firstSourceOut domain.SourceConfig
	ctx, id, appConfig, db, lstream, payload, jobConfig, inSource, outSource, valid = validInputsMultipleSources(ctx, id, appConfig, db, lstream, payload, jobConfig, inSource, outSource)
	if valid {
		valid = len(inSource) == 1 && len(outSource) == 1
		if valid {
			firstSourceIn = inSource[0]
			firstSourceOut = outSource[0]
		}
	}

	return ctx, id, appConfig, db, lstream, payload, jobConfig, firstSourceIn, firstSourceOut, valid
}

func getPermitThread(simultaneousCount int) (permit chan bool) {
	permit = make(chan bool, simultaneousCount)
	for i := 0; i < simultaneousCount; i++ {
		permit <- true
	}

	return permit
}

func getCategoryBasedOnRule(rules []domain.CategoryRule, vulnTitle, vulnCategory, vulnType string) (category string) {
	for _, rule := range rules {
		var invalid bool

		if rule.VulnerabilityTitle() != nil {
			if sord(rule.VulnerabilityTitle()) != vulnTitle {
				invalid = true
			}
		}

		if rule.VulnerabilityCategory() != nil {
			if sord(rule.VulnerabilityCategory()) != vulnCategory {
				invalid = true
			}
		}

		if rule.VulnerabilityType() != nil {
			if sord(rule.VulnerabilityType()) != vulnType {
				invalid = true
			}
		}

		if !invalid {
			category = rule.Category()
			break
		}
	}

	return category
}
