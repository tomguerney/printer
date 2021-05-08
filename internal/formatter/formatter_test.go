package formatter

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type FormatterSuite struct {
	suite.Suite
	Formatter *Formatter
}

func (suite *FormatterSuite) SetupTest() {
	suite.Formatter = &Formatter{}
}

func (suite *FormatterSuite) TestMessage() {
	msg := "test message"
	expected := "test message\n"
	actual := suite.Formatter.Msg(msg)
	suite.
		Equal(expected, actual)
}

func (suite *FormatterSuite) TestMessageWithArgs() {
	msg := "%v test message %v"
	arg1 := "arg1"
	arg2 := "arg2"
	expected := "arg1 test message arg2\n"
	actual := suite.Formatter.Msg(msg, arg1, arg2)
	suite.Equal(expected, actual)
}

func (suite *FormatterSuite) TestError() {
	msg := "test message"
	expected := "Error: test message\n"
	actual := suite.Formatter.Error(msg)
	suite.Equal(expected, actual)
}

func (suite *FormatterSuite) TestErrorWithArgs() {
	msg := "%v test message %v"
	arg1 := "arg1"
	arg2 := "arg2"
	expected := "Error: arg1 test message arg2\n"
	actual := suite.Formatter.Error(msg, arg1, arg2)
	suite.Equal(expected, actual)
}

func (suite *FormatterSuite) TestTabulate() {
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
	actual := suite.Formatter.Tabulate(table)
	suite.Equal(expected, actual)
}

func TestFormatterSuite(t *testing.T) {
	suite.Run(t, new(FormatterSuite))
}
