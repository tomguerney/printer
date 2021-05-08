package stenciller

import (
	"fmt"
	"html/template"
	"os"
	"reflect"
)

// Stenciller takes structs and formats their field values to predefined templates and colors
type Stenciller struct {
	stencils []*stencil
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
func (t *Stenciller) AddStencil(id, template string, colors map[string]string) error {
	//TODO: validate input: id is not duplicated, colors exist, template is valid
	t.stencils = append(t.stencils, &stencil{id, template, colors})
	return nil
}

// UseStencil applies an interface to the stencil with the passed ID
func (t *Stenciller) UseStencil(id string, data interface{}) (string, error) {
	stencil, err := t.findStencil(id)
	if err != nil {
		return "", err
	}
	orig := reflect.ValueOf(data)
	copy, err := t.copyValue(orig)
	if err != nil {
		return "", err
	}
	_, err = t.colorFields(stencil, orig, copy)
	return "nil", nil
}

func (t *Stenciller) findStencil(id string) (*stencil, error) {
	for _, stencil := range t.stencils {
		if stencil.ID == id {
			return stencil, nil
		}
	}
	return nil, fmt.Errorf("Unable to find stencil with id of %v", id)
}

func (t *Stenciller) copyValue(orig reflect.Value) (copy reflect.Value, err error) {
	fields := []reflect.StructField{}
	// For each field in value...
	for i := 0; i < orig.NumField(); i++ {
		field := orig.Type().Field(i)
		//Skip field if it's unexported
		if len(field.PkgPath) != 0 {
			continue
		}
		//Convert field to a string
		fields = append(fields, reflect.StructField{
			Name: field.Name,
			Type: reflect.TypeOf(string("")),
		})
	}
	// Create new type with same field names as value but all fields are strings
	typ := reflect.StructOf(fields)
	// Return new value of type
	return reflect.Indirect(reflect.New(typ)), nil
}

func (t *Stenciller) colorFields(stencil *stencil, orig, copy reflect.Value) (reflect.Value, error) {
	for i := 0; i < orig.NumField(); i++ {
		fmt.Println(orig.Field(i))
	}

	return reflect.Value{}, nil
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
