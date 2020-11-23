package qualys

import (
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/log"
	"time"
)

// LoadVulnerabilities downloads the ENTIRE qualys knowledge base on vulnerabilities in a single API call. There is currently
// no Qualys support for paging on this API call, so this method can be quite expensive (> 10 minutes)
func (session *Session) LoadVulnerabilities(since *time.Time) (output *QKnowledgeBaseVulnOutput, err error) {
	output = &QKnowledgeBaseVulnOutput{}

	// Set the status of the KB load as started
	session.lstream.Send(log.Info("Starting load of vulnerabilities from knowledge base"))

	var fields = make(map[string]string)
	fields["action"] = "list"
	fields["details"] = "All" // Pull ALL details for the vulnerabilities
	if since != nil {
		// there is also published_after field, but it matches the last_modified_after field if the vulnerability has been created but not updated
		fields["last_modified_after"] = since.Format("2006-01-02")
	}

	// Execute the post call against the Qualys API
	if err = session.post(session.Config.Address()+qsVulnerabilities, fields, &output); err != nil {
		session.lstream.Send(log.Errorf(err, "Vulnerability Information failed to load [%s]", err.Error()))
	}

	return output, err
}

// LoadVulnerability loads a single vulnerability from the Qualys knowledge base
func (session *Session) LoadVulnerability(id string) (vuln *QVulnerability, err error) {
	var output = &QKnowledgeBaseVulnOutput{}

	var fields = make(map[string]string)
	fields["action"] = "list"
	fields["details"] = "All" // Pull ALL details for the vulnerabilities
	fields["ids"] = id

	// Execute the post call against the Qualys API
	if err = session.post(session.Config.Address()+qsVulnerabilities, fields, &output); err == nil {
		if len(output.Vulnerabilities) == 1 {
			vuln = &output.Vulnerabilities[0]
		} else {
			err = fmt.Errorf("expected 1 vulnerability when searching for QID [%s] but got %d", id, len(output.Vulnerabilities))
		}
	} else {
		session.lstream.Send(log.Errorf(err, "Vulnerability Information failed to load [%s]", err.Error()))
	}

	return vuln, err
}
