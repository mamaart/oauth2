package viewmodel

import (
	"html/template"
	"net/http"
)

type VM interface {
	Templ() string
	Data() VM
}

func New[T VM](title string, data T) *rootModel {
	return &rootModel{
		Title: title,
		data:  data,
		path:  "./templates/index.html",
	}
}

type rootModel struct {
	Title string
	data  VM
	path  string
}

func (vm *rootModel) Data() VM {
	return vm.data
}

func (vm *rootModel) Templ() string {
	return vm.path
}

func (vm *rootModel) Execute(w http.ResponseWriter) {
	if err := template.
		Must(template.ParseFiles(allPaths(vm)...)).
		ExecuteTemplate(w, "index", &vm); err != nil {
		panic(err)
	}
}

func allPaths(vm VM) []string {
	if vm != nil {
		return append(allPaths(vm.Data()), vm.Templ())
	}
	return []string{}
}
