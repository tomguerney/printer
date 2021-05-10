package printer

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

type PrinterSuite struct {
	suite.Suite
	Writer *MockWriter
}

func (suite *PrinterSuite) SetupTest() {
	suite.Writer = new(MockWriter)
	suite.Writer.On("Write", mock.Anything).Return(0, nil)
	setter.Writer = suite.Writer
}

func (suite *PrinterSuite) TestMessage() {
	Msg("Hello")
	expected := "Hello\n"
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestSMessage() {
	msg := "Hello"
	expected := "Hello\n"
	actual := SMsg(msg)
	suite.Equal(expected, actual)
}

func (suite *PrinterSuite) TestMessageWithArgument() {
	msg := "Hello %v"
	args := "World"
	expected := "Hello World\n"
	Msg(msg, args)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

// func (suite *PrinterSuite) TestSMessageWithArgument() {
// 	msg := "test message"
// 	args := "test args"
// 	expected := "formatted string"
// 	suite.Formatter.On("Msg", mock.Anything, mock.Anything).Return(expected)
// 	actual := suite.Setter.SMsg(msg, args)
// 	suite.Equal(expected, actual)
// }

// func (suite *PrinterSuite) TestError() {
// 	errMsg := "test error message"
// 	expected := "formatted error string"
// 	suite.Formatter.On("Error", errMsg, mock.Anything).Return(expected)
// 	suite.Setter.Error(errMsg)
// 	suite.Writer.AssertCalled(suite.T(), "Write", expected)
// }

// func (suite *PrinterSuite) TestSError() {
// 	errMsg := "test error message"
// 	expected := "formatted error string"
// 	suite.Formatter.On("Error", errMsg, mock.Anything).Return(expected)
// 	actual := suite.Setter.SError(errMsg)
// 	suite.Equal(expected, actual)
// }

// func (suite *PrinterSuite) TestErrorWithArgument() {
// 	errMsg := "test message %v"
// 	args := "test arg"
// 	expected := "formatted error string"
// 	suite.Formatter.On("Error", errMsg, mock.Anything).Return(expected)
// 	suite.Setter.Error(errMsg, args)
// 	suite.Writer.AssertCalled(suite.T(), "Write", expected)
// }

// func (suite *PrinterSuite) TestSErrorWithArgument() {
// 	errMsg := "test message %v"
// 	args := "test arg"
// 	expected := "formatted error string"
// 	suite.Formatter.On("Error", errMsg, mock.Anything).Return(expected)
// 	actual := suite.Setter.SError(errMsg, args)
// 	suite.Equal(expected, actual)
// }

// func (suite *PrinterSuite) TestTabulate() {
// 	table := [][]string{
// 		{"The", "first", "row"},
// 		{"This", "is", "another", "row"},
// 		{"The", "tertiary", "row"},
// 	}
// 	expected := []string{
// 		"row1",
// 		"row2",
// 		"row3",
// 	}
// 	suite.Formatter.On("Tabulate", table).Return(expected)
// 	suite.Setter.Tabulate(table)
// 	suite.Writer.AssertCalled(suite.T(), "Write", expected[0])
// 	suite.Writer.AssertCalled(suite.T(), "Write", expected[1])
// 	suite.Writer.AssertCalled(suite.T(), "Write", expected[2])
// 	suite.Writer.AssertNumberOfCalls(suite.T(), "Write", 3)
// }
// func (suite *PrinterSuite) TestSTabulate() {
// 	table := [][]string{
// 		{"The", "first", "row"},
// 		{"This", "is", "another", "row"},
// 		{"The", "tertiary", "row"},
// 	}
// 	expected := []string{
// 		"row1",
// 		"row2",
// 		"row3",
// 	}
// 	suite.Formatter.On("Tabulate", table).Return(expected)
// 	actual := suite.Setter.STabulate(table)
// 	suite.Equal(expected, actual)
// }
func TestPrinterSuite(t *testing.T) {
	suite.Run(t, new(PrinterSuite))
}
