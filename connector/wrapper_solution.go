package connector

type solution struct {
	text string
}

func (s *solution) Summary() string {
	return s.text
}

// TODO does Qualys provide remediation steps?
func (s *solution) Steps() string {
	return s.text
}

func (s *solution) String() string {
	return s.text
}
