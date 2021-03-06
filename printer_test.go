package printer

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/tomguerney/printer/internal/formatter"
	"github.com/tomguerney/printer/internal/stenciller"
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

func (m *MockFormatter) Text(i interface{}, a ...interface{}) string {
	args := m.Called(i, a)
	return args.String(0)
}

func (m *MockFormatter) Tabulate(rows [][]string, headers ...string) []string {
	args := m.Called(rows, headers)
	return args.Get(0).([]string)
}

func (m *MockFormatter) SetTabwriterOptions(twOptions *formatter.TabwriterOptions) {
	m.Called(twOptions)
}

// MockStenciller is a mock templater for testing
type MockStenciller struct {
	mock.Mock
}

func (m *MockStenciller) AddTemplateStencil(stencil *stenciller.TemplateStencil) error {
	args := m.Called(stencil)
	return args.Error(0)
}

func (m *MockStenciller) AddTableStencil(stencil *stenciller.TableStencil) error {
	args := m.Called(stencil)
	return args.Error(0)
}

func (m *MockStenciller) UseTemplateStencil(id string, data map[string]string) (string, error) {
	args := m.Called(id, data)
	return args.String(0), args.Error(1)
}

func (m *MockStenciller) UseTableStencil(id string, rows []map[string]string) ([][]string, error) {
	args := m.Called(id, rows)
	return args.Get(0).([][]string), args.Error(1)
}

func (m *MockStenciller) Color(text, color string) (string, bool) {
	args := m.Called(text, color)
	return args.String(0), args.Bool(1)
}

// MockPrompter is a mock prompter for testing
type MockPrompter struct {
	mock.Mock
}

func (m *MockPrompter) Select(label string, table []string) (i int, err error) {
	args := m.Called(label, table)
	return args.Int(0), args.Error(1)
}

type PrinterSuite struct {
	suite.Suite
	OutWriter  *MockWriter
	ErrWriter  *MockWriter
	Formatter  *MockFormatter
	Stenciller *MockStenciller
	Prompter   *MockPrompter
}

func (suite *PrinterSuite) SetupTest() {
	suite.OutWriter = new(MockWriter)
	suite.OutWriter.On("Write", mock.Anything).Return(0, nil)
	suite.ErrWriter = new(MockWriter)
	suite.ErrWriter.On("Write", mock.Anything).Return(0, nil)
	suite.Formatter = new(MockFormatter)
	suite.Stenciller = new(MockStenciller)
	suite.Prompter = new(MockPrompter)
	singleton = &Printer{suite.OutWriter, suite.ErrWriter, suite.Formatter, suite.Stenciller, suite.Prompter}
}

func (suite *PrinterSuite) TestOut() {
	text := "test message"
	expected := "formatted string"
	suite.Formatter.On("Text", text, mock.Anything).Return(expected)
	Out(text)
	suite.OutWriter.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestOutWithArgs() {
	text := "test message"
	args := "test args"
	expected := "formatted string"
	suite.Formatter.On("Text", mock.Anything, mock.Anything).Return(expected)
	Out(text, args)
	suite.OutWriter.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestErr() {
	text := "test error message"
	expected := "formatted error string"
	suite.Formatter.On("Text", text, mock.Anything).Return(expected)
	Err(text)
	suite.ErrWriter.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestErrorWithArgument() {
	text := "test message %v"
	args := "test arg"
	expected := "formatted error string"
	suite.Formatter.On("Text", text, mock.Anything).Return(expected)
	Err(text, args)
	suite.ErrWriter.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestTabulate() {
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
	suite.Formatter.On("Tabulate", table, mock.Anything).Return(expected)
	Tabulate(table)
	suite.OutWriter.AssertCalled(suite.T(), "Write", fmt.Sprintln(expected[0]))
	suite.OutWriter.AssertCalled(suite.T(), "Write", fmt.Sprintln(expected[1]))
	suite.OutWriter.AssertCalled(suite.T(), "Write", fmt.Sprintln(expected[2]))
	suite.OutWriter.AssertNumberOfCalls(suite.T(), "Write", 3)
}

func (suite *PrinterSuite) TestTabulateWithHeaders() {
	table := [][]string{
		{"The", "first", "row"},
		{"This", "is", "another", "row"},
		{"The", "tertiary", "row"},
	}
	headers := []string{"header1", "header2"}
	expected := []string{
		"headers",
		"row1",
		"row2",
		"row3",
	}
	suite.Formatter.On("Tabulate", table, headers).Return(expected)
	Tabulate(table, headers...)
	suite.OutWriter.AssertCalled(suite.T(), "Write", fmt.Sprintln(expected[0]))
	suite.OutWriter.AssertCalled(suite.T(), "Write", fmt.Sprintln(expected[1]))
	suite.OutWriter.AssertCalled(suite.T(), "Write", fmt.Sprintln(expected[2]))
	suite.OutWriter.AssertCalled(suite.T(), "Write", fmt.Sprintln(expected[3]))
	suite.OutWriter.AssertNumberOfCalls(suite.T(), "Write", 4)
}

func (suite *PrinterSuite) TestTmplStencil() {
	id := "test id"
	data := map[string]string{"key": "value"}
	expected := "stencilled string"
	suite.Stenciller.On("UseTemplateStencil", id, data).Return(expected, nil)
	err := UseTemplateStencil(id, data)
	suite.NoError(err)
	suite.OutWriter.AssertCalled(suite.T(), "Write", fmt.Sprintln(expected))
}

func (suite *PrinterSuite) TestTmplStencilWithError() {
	id := "test id"
	data := map[string]string{"key": "value"}
	suite.Stenciller.On("UseTemplateStencil", id, data).Return("", errors.New("error"))
	err := UseTemplateStencil(id, data)
	suite.Error(err)
	suite.OutWriter.AssertNotCalled(suite.T(), "Write", mock.Anything)
}

func (suite *PrinterSuite) TestTableStencil() {
	id := "test id"
	rows := []map[string]string{{"key": "value"}}
	tableStencilResult := [][]string{{"row1a", "row1b"}, {"row2a", "row2b"}}
	tabulateResult := []string{"row1", "row2"}
	suite.Stenciller.On("UseTableStencil", id, rows).Return(tableStencilResult, nil)
	suite.Formatter.On("Tabulate", mock.Anything, mock.Anything).Return(tabulateResult)
	err := UseTableStencil(id, rows)
	suite.NoError(err)
	suite.OutWriter.AssertCalled(suite.T(), "Write", fmt.Sprintln(tabulateResult[0]))
	suite.OutWriter.AssertCalled(suite.T(), "Write", fmt.Sprintln(tabulateResult[1]))
}

func (suite *PrinterSuite) TestTableStencilWithError() {
	id := "test id"
	rows := []map[string]string{{"key": "value"}}
	suite.Stenciller.On("UseTableStencil", id, rows).Return([][]string{}, errors.New("error"))
	err := UseTableStencil(id, rows)
	suite.Error(err)
	suite.Formatter.AssertNotCalled(suite.T(), "Tabulate", mock.Anything)
	suite.OutWriter.AssertNotCalled(suite.T(), "Write", mock.Anything)
}

func TestPrinterSuite(t *testing.T) {
	suite.Run(t, new(PrinterSuite))
}
