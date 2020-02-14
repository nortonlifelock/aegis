package domain

// LogType defines the interface
type LogType interface {
	ID() (param int)
	LogType() (param string)
	Name() (param string)
}
