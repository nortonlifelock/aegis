package dal

// Solution contains the information needed to remediate a vulnerability
type Solution struct {
	text string
}

func (s Solution) String() string {
	return s.text
}

// Summary returns a summary of the solution
func (s Solution) Summary() string {
	return s.text
}

// Steps returns a step-by-step guide to remediate the vulnerability
// TODO
func (s Solution) Steps() string {
	return s.text
}
