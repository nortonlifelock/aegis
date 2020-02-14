package domain

// Permission defines the interface
type Permission interface {
	Admin() (param bool)
	Manager() (param bool)
	OrgID() (param string)
	ParentOrgPermission() (param Permission)
	Reader() (param bool)
	Reporter() (param bool)
	UserID() (param string)
}
