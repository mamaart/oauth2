package viewmodel

type rootModel struct {
	Title string
	data  VM
}

func (vm *rootModel) Data() VM        { return vm.data }
func (vm *rootModel) Templ() []string { return []string{"./templates/index.html"} }
