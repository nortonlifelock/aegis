package domain

import "time"

type CodeFinding interface {
	ProjectID() string
	ProjectName() string
	ProjectVersion() string
	ProjectOwner() string
	ComponentName() string
	ComponentVersion() string
	ViolatedPolicyName() string
	ViolatedPolicySeverity() string
	CVSS() float32
	Description() string
	Summary() string
	Updated() time.Time
	VulnerabilityID() string
}
