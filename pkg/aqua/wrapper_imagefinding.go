package aqua

import (
	"fmt"
	"strings"
	"time"
)

const timeFormat = "2006-01-02"

func (vr *VulnerabilityResult) CVSS2() *float32 {
	val := float32(vr.NvdCvss2Score)
	return &val
}

func (vr *VulnerabilityResult) CVSS3() *float32 {
	val := float32(vr.NvdCvss3Score)
	return &val
}

func (vr *VulnerabilityResult) ImageTag() string {
	return strings.Replace(vr.ImageNameVar, fmt.Sprintf("%s:", vr.ImageRepositoryName), "", 1)
}

func (vr *VulnerabilityResult) ImageName() string {
	return vr.ImageRepositoryName // DeviceID
}

func (vr *VulnerabilityResult) ImageVersion() string {
	return vr.ImageDigest // Hostname
}

func (vr *VulnerabilityResult) Registry() string {
	return vr.RegistryVar // GroupID
}

func (vr *VulnerabilityResult) LastFound() *time.Time {
	timeVal, _ := time.Parse(timeFormat, vr.LastFoundDate)
	return &timeVal
}

func (vr *VulnerabilityResult) FirstFound() *time.Time {
	timeVal, _ := time.Parse(timeFormat, vr.FirstFoundDate)
	return &timeVal
}

func (vr *VulnerabilityResult) LastUpdated() *time.Time {
	timeVal, _ := time.Parse(timeFormat, vr.ModificationDate)
	return &timeVal
}

func (vr *VulnerabilityResult) Patchable() *string {
	//val := "No"
	//if vr.VPatchStatus == "patch_available" {
	//	val = "Yes"
	//}

	val := "No"
	if len(vr.FixVersion) > 0 {
		val = "Yes"
	}

	return &val
}

func (vr *VulnerabilityResult) Solution() *string {
	return &vr.SolutionVar
}

func (vr *VulnerabilityResult) Summary() *string {
	return &vr.DescriptionVar
}

func (vr *VulnerabilityResult) VendorReference() string {
	var ref = vr.NvdURL
	if len(vr.VendorURL) > 0 {
		ref = fmt.Sprintf("%s\n%s", ref, vr.VendorURL)
	}
	return ref
}

func (vr *VulnerabilityResult) VulnerabilityLocation() string {
	if len(vr.Resource.Path) > 0 {
		return vr.Resource.Path // for files
	} else {
		return vr.Resource.Cpe // for packages
	}
}

func (vr *VulnerabilityResult) VulnerabilityID() string {
	return vr.Name
}

func (vr *VulnerabilityResult) Exception() bool {
	return false // we will override this method for exceptions
}
