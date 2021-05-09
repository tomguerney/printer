package colorer

import (
	"testing"
)

func TestCorrectColorNames(t *testing.T) {
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
		t.Run(tt.color, func(t *testing.T) {
			_, ok := c.Color(text, tt.color)
			if !ok {
				t.Fatalf("Colorize failed with color %v", tt.color)
			}
		})
	}
}

func TestIncorrectColorName(t *testing.T) {
	c := Colorer{}
	_, ok := c.Color("not a color", "text")
	if ok {
		t.Fatal("\"not a color\" should not return ok")
	}
}
