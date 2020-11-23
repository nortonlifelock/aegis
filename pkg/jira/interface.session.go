package domain

// Session defines the interface
type Session interface {
	IsDisabled() (param bool)
	OrgID() (param string)
	SessionKey() (param string)
	UserID() (param string)
}
