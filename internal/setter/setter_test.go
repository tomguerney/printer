package setter

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockWriter is a mock io.Writer for testing
type MockWriter struct {
	mock.Mock
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	args := m.Called(string(p))
	return args.Int(0), args.Error(1)
}

// MockFormatter is a mock formatter for testing
type MockFormatter struct {
	mock.Mock
}

func (m *MockFormatter) Msg(text string, a ...interface{}) string {
	args := m.Called(text, a)
	return args.String(0)
}

func (m *MockFormatter) Error(text string, a ...interface{}) string {
	args := m.Called(text, a)
	return args.String(0)
}

func (m *MockFormatter) Tabulate(headers []string, table [][]string) []string {
	args := m.Called(headers, table)
	return args.Get(0).([]string)
}

// MockTemplater is a mock templater for testing
type MockStenciller struct {
	mock.Mock
}

func (m *MockStenciller) AddTmplStencil(
	id,
	template string,
	colors map[string]string,
) error {
	args := m.Called(id, template, colors)
	return args.Error(0)
}

func (m *MockStenciller) AddTableStencil(
	id string,
	headers, columnOrder []string,
	colors map[string]string,
) error {
	args := m.Called(id, headers, columnOrder, colors)
	return args.Error(0)
}

func (m *MockStenciller) TmplStencil(
	id string,
	data map[string]string,
) (string, error) {
	args := m.Called(id, data)
	return args.String(0), args.Error(1)
}

func (m *MockStenciller) TableStencil(
	id string,
	rows []map[string]string,
) ([][]string, error) {
	args := m.Called(id, rows)
	return args.Get(0).([][]string), args.Error(1)
}

func (m *MockStenciller) TableStencilHeaders(
	id string,
	rows []map[string]string,
) ([][]string, bool) {
	args := m.Called(id, rows)
	return args.Get(0).([][]string), args.Bool(1)
}

func (m *MockStenciller) PrefixHeaders(rows [][]string) [][]string {
	args := m.Called(rows)
	return args.Get(0).([][]string)
}

type SetterSuite struct {
	suite.Suite
	Setter     *Setter
	Writer     *MockWriter
	ErrWriter  *MockWriter
	Formatter  *MockFormatter
	Stenciller *MockStenciller
}

func (suite *SetterSuite) SetupTest() {
	suite.Writer = new(MockWriter)
	suite.Writer.On("Write", mock.Anything).Return(0, nil)
	suite.Formatter = new(MockFormatter)
	suite.Stenciller = new(MockStenciller)
	suite.Setter = &Setter{suite.Writer, suite.ErrWriter, suite.Formatter, suite.Stenciller}
}

func (suite *SetterSuite) TestMessage() {
	msg := "test message"
	expected := "formatted string"
	suite.Formatter.On("Msg", msg, mock.Anything).Return(expected)
	suite.Setter.Msg(msg)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *SetterSuite) TestMessageWithArgument() {
	msg := "test message"
	args := "test args"
	expected := "formatted string"
	suite.Formatter.On("Msg", mock.Anything, mock.Anything).Return(expected)
	suite.Setter.Msg(msg, args)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *SetterSuite) TestSMessage() {
	msg := "test message"
	expected := "formatted string"
	suite.Formatter.On("Msg", msg, mock.Anything).Return(expected)
	actual := suite.Setter.SMsg(msg)
	suite.Equal(expected, actual)
}

func (suite *SetterSuite) TestSMessageWithArgument() {
	msg := "test message"
	args := "test args"
	expected := "formatted string"
	suite.Formatter.On("Msg", mock.Anything, mock.Anything).Return(expected)
	actual := suite.Setter.SMsg(msg, args)
	suite.Equal(expected, actual)
}

func (suite *SetterSuite) TestError() {
	errMsg := "test error message"
	expected := "formatted error string"
	suite.Formatter.On("Error", errMsg, mock.Anything).Return(expected)
	suite.Setter.Error(errMsg)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *SetterSuite) TestErrorWithArgument() {
	errMsg := "test message %v"
	args := "test arg"
	expected := "formatted error string"
	suite.Formatter.On("Error", errMsg, mock.Anything).Return(expected)
	suite.Setter.Error(errMsg, args)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *SetterSuite) TestSError() {
	errMsg := "test error message"
	expected := "formatted error string"
	suite.Formatter.On("Error", errMsg, mock.Anything).Return(expected)
	actual := suite.Setter.SError(errMsg)
	suite.Equal(expected, actual)
}

func (suite *SetterSuite) TestSErrorWithArgument() {
	errMsg := "test message %v"
	args := "test arg"
	expected := "formatted error string"
	suite.Formatter.On("Error", errMsg, mock.Anything).Return(expected)
	actual := suite.Setter.SError(errMsg, args)
	suite.Equal(expected, actual)
}

func (suite *SetterSuite) TestTabulate() {
	table := [][]string{
		{"The", "first", "row"},
		{"This", "is", "another", "row"},
		{"The", "tertiary", "row"},
	}
	expected := []string{
		"row1",
		"row2",
		"row3",
	}
	suite.Formatter.On("Tabulate", mock.Anything, table).Return(expected)
	suite.Setter.Tabulate(table)
	suite.Writer.AssertCalled(suite.T(), "Write", fmt.Sprintln(expected[0]))
	suite.Writer.AssertCalled(suite.T(), "Write", fmt.Sprintln(expected[1]))
	suite.Writer.AssertCalled(suite.T(), "Write", fmt.Sprintln(expected[2]))
	suite.Writer.AssertNumberOfCalls(suite.T(), "Write", 3)
}
func (suite *SetterSuite) TestSTabulate() {
	table := [][]string{
		{"The", "first", "row"},
		{"This", "is", "another", "row"},
		{"The", "tertiary", "row"},
	}
	expected := []string{
		"row1",
		"row2",
		"row3",
	}
	suite.Formatter.On("Tabulate", mock.Anything, table).Return(expected)
	actual := suite.Setter.STabulate(table)
	suite.Equal(expected, actual)
}

func (suite *SetterSuite) TestTmplStencil() {
	id := "test id"
	data := map[string]string{"key": "value"}
	expected := "stencilled string"
	suite.Stenciller.On("TmplStencil", id, data).Return(expected, nil)
	err := suite.Setter.TmplStencil(id, data)
	suite.NoError(err)
	suite.Writer.AssertCalled(suite.T(), "Write", fmt.Sprintln(expected))
}

func (suite *SetterSuite) TestTmplStencilWithError() {
	id := "test id"
	data := map[string]string{"key": "value"}
	suite.Stenciller.On("TmplStencil", id, data).
		Return("", errors.New("error"))
	err := suite.Setter.TmplStencil(id, data)
	suite.Error(err)
	suite.Writer.AssertNotCalled(suite.T(), "Write", mock.Anything)
}

func (suite *SetterSuite) TestSTmplStencil() {
	id := "test id"
	data := map[string]string{"key": "value"}
	expected := "stencilled string"
	suite.Stenciller.On("TmplStencil", id, data).Return(expected, nil)
	actual, err := suite.Setter.STmplStencil(id, data)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *SetterSuite) TestSTmplStencilWithError() {
	id := "test id"
	data := map[string]string{"key": "value"}
	suite.Stenciller.On("TmplStencil", id, data).
		Return("", errors.New("error"))
	actual, err := suite.Setter.STmplStencil(id, data)
	suite.Error(err)
	suite.Equal("", actual)
}

func (suite *SetterSuite) TestTableStencil() {
	id := "test id"
	rows := []map[string]string{{"key": "value"}}
	tableStencilResult := [][]string{{"row1a", "row1b"}, {"row2a", "row2b"}}
	tabulateResult := []string{"row1", "row2"}
	suite.Stenciller.On("TableStencil", id, rows).
		Return(tableStencilResult, nil)
	suite.Formatter.On("Tabulate", mock.Anything, mock.Anything).Return(tabulateResult)
	err := suite.Setter.TableStencil(id, rows)
	suite.NoError(err)
	suite.Writer.AssertCalled(suite.T(), "Write", fmt.Sprintln(tabulateResult[0]))
}

func (suite *SetterSuite) TestTableStencilWithError() {
	id := "test id"
	rows := []map[string]string{{"key": "value"}}
	suite.Stenciller.On("TableStencil", id, rows).
		Return([][]string{}, errors.New("error"))
	err := suite.Setter.TableStencil(id, rows)
	suite.Error(err)
	suite.Formatter.AssertNotCalled(suite.T(), "Tabulate", mock.Anything)
	suite.Writer.AssertNotCalled(suite.T(), "Write", mock.Anything)
}

func (suite *SetterSuite) TestSTableStencil() {
	id := "test id"
	rows := []map[string]string{{"key": "value"}}
	tableStencilResult := [][]string{{"row1a", "row1b"}, {"row2a", "row2b"}}
	expected := []string{"row1", "row2"}
	suite.Stenciller.On("TableStencil", id, rows).
		Return(tableStencilResult, nil)
	suite.Formatter.On("Tabulate", mock.Anything, mock.Anything).Return(expected)
	actual, err := suite.Setter.STableStencil(id, rows)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *SetterSuite) TestSTableStencilWithError() {
	id := "test id"
	rows := []map[string]string{{"key": "value"}}
	suite.Stenciller.On("TableStencil", id, rows).
		Return([][]string{}, errors.New("error"))
	result, err := suite.Setter.STableStencil(id, rows)
	suite.Error(err)
	suite.Len(result, 0)
	suite.Formatter.AssertNotCalled(suite.T(), "Tabulate", mock.Anything)
}

func TestSetterSuite(t *testing.T) {
	suite.Run(t, new(SetterSuite))
}
