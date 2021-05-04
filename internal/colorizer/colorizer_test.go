package colorizer

import (
	"testing"
)

func TestCorrectColorNames(t *testing.T) {
	c := Colorizer{}
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
			_, err := c.Colorize(tt.color, text)
			if err != nil {
				t.Fatalf("error returned from Colorize with color %v", tt.color)
			}
		})
	}
}

func TestIncorrectColorName(t *testing.T) {
	c := Colorizer{}
	_, err := c.Colorize("not a color", "text")
	if err == nil {
		t.Fatal("\"not a color\" should return error")
	}
}
