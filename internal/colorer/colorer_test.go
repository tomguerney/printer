package colorer

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ColorerSuite struct {
	suite.Suite
	Colorer *Colorer
}

func (suite *ColorerSuite) SetupTest() {
	suite.Colorer = &Colorer{}
}

func (suite *ColorerSuite) TestCorrectColorNames() {
	c := Colorer{}
	const text = "text"
	var colorTests = []struct {
		color string
	}{
		{"black"},
		{"red"},
		{"green"},
		{"yellow"},
		{"blue"},
		{"magenta"},
		{"cyan"},
		{"white"},
	}
	for _, tt := range colorTests {
		suite.Run(tt.color, func() {
			_, ok := c.Color(text, tt.color)
			if !ok {
				suite.T().Fatalf("Colorize failed with color %v", tt.color)
			}
		})
	}
}

func (suite *ColorerSuite) TestIncorrectColorName() {
	c := Colorer{}
	_, ok := c.Color("not a color", "text")
	if ok {
		suite.T().Fatal("\"not a color\" should not return ok")
	}
}

func TestColorerSuite(t *testing.T) {
	suite.Run(t, new(ColorerSuite))
}
