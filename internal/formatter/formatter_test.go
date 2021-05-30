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
	tabwriterOptions := &TabwriterOptions{
		Minwidth: 0,
		Tabwidth: 8,
		Padding:  4,
		Padchar:  ' ',
		Divchar:  '-',
	}
	suite.Formatter = &Formatter{tabwriterOptions}
}

func (suite *FormatterSuite) TestText() {
	msg := "test message"
	expected := "test message\n"
	actual := suite.Formatter.Text(msg)
	suite.Equal(expected, actual)
}

func (suite *FormatterSuite) TestTextWithArgs() {
	msg := "%v test message %v"
	arg1 := "arg1"
	arg2 := "arg2"
	expected := "arg1 test message arg2\n"
	actual := suite.Formatter.Text(msg, arg1, arg2)
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

func (suite *FormatterSuite) TestTabulateWithHeaders() {
	table := [][]string{
		{"The", "first", "row"},
		{"This", "is", "another", "row"},
		{"The", "tertiary", "row"},
	}
	headers := []string{"header1", "h2", "this is header3", "4"}
	expected := []string{
		"header1    h2          this is header3    4",
		"-------    --------    ---------------    ---",
		"The        first       row",
		"This       is          another            row",
		"The        tertiary    row",
	}
	actual := suite.Formatter.Tabulate(table, headers...)
	suite.Equal(expected, actual)
}

func (suite *FormatterSuite) TestTabulateWithMorePadding() {
	suite.Formatter.TWOptions.Padding = 10
	table := [][]string{
		{"The", "first", "row"},
		{"This", "is", "another", "row"},
		{"The", "tertiary", "row"},
	}
	expected := []string{
		"The           first             row",
		"This          is                another          row",
		"The           tertiary          row",
	}
	actual := suite.Formatter.Tabulate(table)
	suite.Equal(expected, actual)
}

func TestFormatterSuite(t *testing.T) {
	suite.Run(t, new(FormatterSuite))
}
