package aqua

import (
	"fmt"
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
	val := "No"
	if vr.VPatchStatus == "patch_available" {
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
	return vr.NvdURL
}

func (vr *VulnerabilityResult) VulnerabilityID() string {
	if len(vr.Resource.Path) > 0 {
		return fmt.Sprintf("%s;%s", vr.Name, vr.Resource.Path) // for files
	} else {
		return fmt.Sprintf("%s;%s", vr.Name, vr.Resource.Cpe) // for packages
	}
}
