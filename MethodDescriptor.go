package scaffold

type methodDescriptor struct {
	Name   string
	Method string
}

func newMethod(name string, method string) *methodDescriptor {
	return &methodDescriptor{name, method}
}
