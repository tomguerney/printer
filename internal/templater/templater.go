package templater

import (
	"fmt"
	"html/template"
	"os"
	"reflect"
)

type Templater struct {
	templates []Template
}

type Template struct {
	Id       string
	Template string
	Colors   map[string]string
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
