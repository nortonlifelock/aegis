package domain

import "time"

type ImageFinding interface {
	CVSS2() *float32
	CVSS3() *float32
	ImageName() string
	ImageVersion() string
	ImageTag() string
	Registry() string
	LastFound() *time.Time
	FirstFound() *time.Time
	LastUpdated() *time.Time
	Patchable() *string
	Solution() *string
	Summary() *string
	VendorReference() string
	VulnerabilityID() string
}
