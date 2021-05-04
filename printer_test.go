package printer

import (
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

type PrinterSuite struct {
	suite.Suite
	Printer
	Writer *MockWriter
}

func (suite *PrinterSuite) SetupTest() {
	suite.Writer = new(MockWriter)
	suite.Writer.On("Write", mock.Anything).Return(0, nil)
	suite.Printer = Printer{suite.Writer}
}

func (suite *PrinterSuite) TestMessage() {
	msg := "test message"
	expected := fmt.Sprintf("%v\n", msg)
	suite.Printer.Msg(msg)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestMessageWithArgument() {
	msg := "test message %v"
	arg := "arg"
	expected := fmt.Sprintf("%v\n", fmt.Sprintf(msg, arg))
	suite.Printer.Msg(msg, arg)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestError() {
	errMsg := "test error message"
	expected := fmt.Sprintf("Error: %v\n", errMsg)
	suite.Printer.Error(errMsg)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestErrorWithArgument() {
	errMsg := "test message %v"
	arg := "arg"
	expected := fmt.Sprintf("Error: %v\n", fmt.Sprintf(errMsg, arg))
	suite.Printer.Error(errMsg, arg)
	suite.Writer.AssertCalled(suite.T(), "Write", expected)
}

func (suite *PrinterSuite) TestTabulate() {
	table := [][]string{
		{"The", "first", "row"},
		{"This", "is", "another", "row"},
		{"The", "tertiary", "row"},
	}
	expected := []string{
		"The     first       row",
		"This    is          another    row",
		"The     tertiary    row",
	}
	actual, err := suite.Printer.Tabulate(table)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func TestPrinterSuite(t *testing.T) {
	suite.Run(t, new(PrinterSuite))
}
