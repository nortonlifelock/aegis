package nexpose

import (
	"fmt"
	"github.com/benjivesterby/validator"
)

// ResourcesSolution is the json return for a list of solutions returned from nexpose
type ResourcesSolution struct {
	Page Page `json:"Page"`

	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`

	// The resources returned.
	Resources []Solution `json:"resources,omitempty"`
}

// Solution is a vulnerability solution json representation for the nexpose api
type Solution struct {

	// Additional information or resources that can assist in applying the remediation.
	AdditionalInformation *struct {

		// Hypertext Markup Language (HTML) representation of the content.
		HTML string `json:"html,omitempty"`

		// Textual representation of the content.
		Text string `json:"text,omitempty"`
	} `json:"additionalInformation,omitempty"`

	// The systems or software the solution applies to.
	AppliesTo string `json:"appliesTo,omitempty"`

	// The estimated duration to apply the solution, in ISO 8601 format.
	// For example: PT5M.
	Estimate string `json:"estimate,omitempty"`

	// The identifier of the solution.
	ID string `json:"id,omitempty"`

	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`

	// The steps required to remediate the vulnerability.
	SolutionSteps *struct {

		// Textual representation of the content.
		HTML string `json:"html,omitempty"`

		// Textual representation of the content.
		Text string `json:"text,omitempty"`
	} `json:"steps,omitempty"`

	// The summary of the solution.
	SolutionSummary *struct {

		// Textual representation of the content.
		HTML string `json:"html,omitempty"`

		// Textual representation of the content.
		Text string `json:"text,omitempty"`
	} `json:"summary,omitempty"`

	// The type of the solution.
	// One of:
	// - Configuration
	// - Rollup patch
	// - Patch
	Type string `json:"type,omitempty"`
}

// Summary returns the summary text of the solution if the summary exists
func (s *Solution) Summary() (summary string) {
	if validator.IsValid(s) {

		if s.SolutionSummary != nil {
			summary = s.SolutionSummary.Text
		}
	}

	return summary
}

// Steps returns the steps to correct the vulnerability if they exist
func (s *Solution) Steps() (steps string) {

	if validator.IsValid(s) {

		if s.SolutionSteps != nil {
			steps = s.SolutionSteps.Text
		}
	}

	return steps
}

func (s *Solution) String() (solution string) {
	return fmt.Sprintf("Summary:\n%s\nSteps\n%s\n", s.Summary(), s.Steps())
}
