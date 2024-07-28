package authorizer

import "github.com/mamaart/viewmodel"

type ErrorBox struct {
	Title  string
	Errors []struct {
		Message string
	}
}

type vm struct {
	ClientID string
	ErrorBox ErrorBox
}

func newVM(clientID string) *vm {
	return &vm{ClientID: clientID}
}

func (vm *vm) Data() viewmodel.VM { return nil }
func (vm *vm) Templ() []string {
	return []string{
		"./templates/login.html",
		"./templates/error.html",
	}
}
