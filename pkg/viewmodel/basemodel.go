package viewmodel

type baseModel struct {
	paths []string
}

// If the viewmodel has values it is not basic
func Basic(paths ...string) *baseModel { return &baseModel{paths: paths} }
func (vm *baseModel) Templ() []string  { return vm.paths }
func (vm *baseModel) Data() VM         { return nil }
