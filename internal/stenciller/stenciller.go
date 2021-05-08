package stenciller

import (
	"fmt"
	"html/template"
	"os"
	"reflect"
)

// Stenciller takes structs and formats their field values to predefined templates and colors
type Stenciller struct {
	stencils []stencil
}

type stencil struct {
	ID       string
	Template string
	Colors   map[string]string
}

// New returns a pointer to a new Stenciller struct
func New() *Stenciller {
	return &Stenciller{}
}

// AddStencil adds a new stencil
func (t *Stenciller) AddStencil(id, template string, colors map[string]string) {
	t.stencils = append(t.stencils, stencil{id, template, colors})
}

// UseStencil applies an interface to the stencil with the passed ID
func (t *Stenciller) UseStencil(id string, s interface{}) (string, error) {
	return "nil", nil
}

func modifyInterface(ifce interface{}) {

	v := reflect.ValueOf(ifce)

	structFields := []reflect.StructField{}

	for i := 0; i < v.NumField(); i++ {
		structFields = append(structFields, reflect.StructField{
			Name: v.Type().Field(i).Name,
			Type: reflect.TypeOf(string("")),
		})
	}

	typ := reflect.StructOf(structFields)

	x := reflect.New(typ).Elem()

	for i := 0; i < v.NumField(); i++ {
		x.FieldByName(v.Type().Field(i).Name).SetString(fmt.Sprintf("modified: %v", v.Field(i)))
	}

	s := x.Addr().Interface()

	actuallyDoTemplate(s)
}

func actuallyDoTemplate(i interface{}) {
	tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, i)
	if err != nil {
		panic(err)
	}
}

type Inventory struct {
	Material string
	Count    uint
}

func main() {
	sweaters := Inventory{"wool", 17}
	modifyInterface(sweaters)
}
