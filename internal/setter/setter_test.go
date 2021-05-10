package setter

import (
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

func (m *MockFormatter) Tabulate(table [][]string) []string {
	args := m.Called(table)
	return args.Get(0).([]string)
}

// MockTemplater is a mock templater for testing
type MockStenciller struct {
	mock.Mock
}

func (m *MockStenciller) AddStencil(id, template string, colors map[string]string) error {
	args := m.Called(id, template, colors)
	return args.Error(0)
}

func (m *MockStenciller) UseStencil(id string, data map[string]string) (string, error) {
	args := m.Called(id, data)
	return args.String(0), args.Error(1)
}

type SetterSuite struct {
	suite.Suite
	Setter     *Setter
	Writer     *MockWriter
	Formatter  *MockFormatter
	Stenciller *MockStenciller
}

func (suite *SetterSuite) SetupTest() {
	suite.Writer = new(MockWriter)
	suite.Writer.On("Write", mock.Anything).Return(0, nil)
	suite.Formatter = new(MockFormatter)
	suite.Stenciller = new(MockStenciller)
	suite.Setter = &Setter{suite.Writer, suite.Formatter, suite.Stenciller}
}

func (suite *SetterSuite) TestMessage() {
	msg := "test message"
	expected := "formatted string"
	suite.Formatter.On("Msg", msg, mock.Anything).Return(expected)
	suite.Setter.Msg(msg)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *SetterSuite) TestSMessage() {
	msg := "test message"
	expected := "formatted string"
	suite.Formatter.On("Msg", msg, mock.Anything).Return(expected)
	actual := suite.Setter.SMsg(msg)
	suite.Equal(expected, actual)
}

func (suite *SetterSuite) TestMessageWithArgument() {
	msg := "test message"
	args := "test args"
	expected := "formatted string"
	suite.Formatter.On("Msg", mock.Anything, mock.Anything).Return(expected)
	suite.Setter.Msg(msg, args)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
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

func (suite *SetterSuite) TestSError() {
	errMsg := "test error message"
	expected := "formatted error string"
	suite.Formatter.On("Error", errMsg, mock.Anything).Return(expected)
	actual := suite.Setter.SError(errMsg)
	suite.Equal(expected, actual)
}

func (suite *SetterSuite) TestErrorWithArgument() {
	errMsg := "test message %v"
	args := "test arg"
	expected := "formatted error string"
	suite.Formatter.On("Error", errMsg, mock.Anything).Return(expected)
	suite.Setter.Error(errMsg, args)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
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
	suite.Formatter.On("Tabulate", table).Return(expected)
	suite.Setter.Tabulate(table)
	suite.Writer.AssertCalled(suite.T(), "Write", expected[0])
	suite.Writer.AssertCalled(suite.T(), "Write", expected[1])
	suite.Writer.AssertCalled(suite.T(), "Write", expected[2])
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
	suite.Formatter.On("Tabulate", table).Return(expected)
	actual := suite.Setter.STabulate(table)
	suite.Equal(expected, actual)
}
func TestSetterSuite(t *testing.T) {
	suite.Run(t, new(SetterSuite))
}