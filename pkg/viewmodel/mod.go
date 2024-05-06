package viewmodel

import (
	"html/template"
	"net/http"
)

type VM interface {
	Templ() []string
	Data() VM
}

func New[T VM](title string, data T) *rootModel {
	return &rootModel{
		Title: title,
		data:  data,
	}
}

func (vm *rootModel) Execute(w http.ResponseWriter) {
	if err := mustV(template.
		New("index").
		Funcs(template.FuncMap{
			"safeHTML": func(s string) template.HTML { return template.HTML(s) },
		}).
		ParseFiles(allPaths(vm)...)).Execute(w, &vm); err != nil {
		panic(err)
	}
}

func allPaths(vm VM) []string {
	if vm != nil {
		return append(vm.Templ(), allPaths(vm.Data())...)
	}
	return []string{}
}

func mustV[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
