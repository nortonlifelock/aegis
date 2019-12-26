package connector

type vendorReference struct {
	source int
	name   string
}

func (vr *vendorReference) Reference() string {
	return vr.name
}

func (vr *vendorReference) String() string {
	return vr.name
}
