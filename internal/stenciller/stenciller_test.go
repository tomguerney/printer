package stenciller

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type StencillerSuite struct {
	suite.Suite
	Stenciller *Stenciller
}

func (suite *StencillerSuite) SetupTest() {
	suite.Stenciller = &Stenciller{}
}

func (suite *StencillerSuite) TestAddStencil() {
	suite.Empty(suite.Stenciller.stencils)
	id := "test-id"
	template := "{{ .Test }} template"
	colors := map[string]string{
		"test": "red",
	}
	suite.Stenciller.AddStencil(id, template, colors)
	suite.Len(suite.Stenciller.stencils, 1)
}

func (suite *StencillerSuite) TestGetValue() {
	type TestType struct {
		Field1 string
		Field2 int
		Field3 string
	}
	data := TestType{"one", 2, "three"}
	orig := reflect.ValueOf(data)
	value, _ := suite.Stenciller.copyValue(orig)
	suite.Equal(3, value.NumField())
}

func (suite *StencillerSuite) TestGetValueAllFieldsAreStrings() {
	type TestType struct {
		Field1 string
		Field2 int
		Field3 byte
		Field4 bool
	}
	data := TestType{"one", 2, byte(12), true}
	orig := reflect.ValueOf(data)
	value, _ := suite.Stenciller.copyValue(orig)
	for i := 0; i < value.NumField(); i++ {
		suite.IsType("", value.Field(i).Interface())
	}
}

func (suite *StencillerSuite) TestGetValueIgnoresUnexportedfields() {
	type TestType struct {
		field1 string
		Field2 int
		Field3 string
	}
	data := TestType{"one", 2, "three"}
	orig := reflect.ValueOf(data)
	value, _ := suite.Stenciller.copyValue(orig)
	suite.Equal(2, value.NumField())
}

func (suite *StencillerSuite) TestFindStencil() {
	stencil1 := &stencil{ID: "1"}
	stencil2 := &stencil{ID: "2"}
	suite.Stenciller.stencils = []*stencil{stencil1, stencil2}
	actual, err := suite.Stenciller.findStencil("1")
	suite.NoError(err)
	suite.Equal(stencil1, actual)
	suite.NotEqual(stencil2, actual)
}

func (suite *StencillerSuite) TestNotFindStencil() {
	stencil1 := &stencil{ID: "1"}
	stencil2 := &stencil{ID: "2"}
	suite.Stenciller.stencils = []*stencil{stencil1, stencil2}
	actual, err := suite.Stenciller.findStencil("3")
	suite.Errorf(err, "Unable to find stencil with id of 3")
	suite.Nil(actual)
}

func (suite *StencillerSuite) TestColorFields() {
	stencil := &stencil{ID: "1"}
	type TestOrig struct {
		Field1 string
		Field2 int
		Field3 string
	}
	origData := TestOrig{"one", 2, "three"}
	orig := reflect.ValueOf(origData)
	type TestCopy struct {
		Field1 string
		Field2 string
		Field3 string
	}
	copyData := TestCopy{}
	copy := reflect.ValueOf(copyData)
	suite.Stenciller.colorFields(stencil, orig, copy)

}

func TestStencillerSuite(t *testing.T) {
	suite.Run(t, new(StencillerSuite))
}
