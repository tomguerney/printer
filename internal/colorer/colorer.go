package colorer

import (
	"github.com/fatih/color"
)

// Colorer colors text for terminal output using the "github.com/fatih/color"
// package
type Colorer struct{}

// New returns a pointer to a new Colorer struct
func New() *Colorer {
	return &Colorer{}
}

// Color transforms a string into one of the available colors. If the color is
// not available the string will not be coloured and ok will be false. 
func (c *Colorer) Color(text string, colorName string) (colored string, ok bool) {
	switch colorName {
	case "black":
		return c.Black(text), true
	case "red":
		return c.Red(text), true
	case "green":
		return c.Green(text), true
	case "yellow":
		return c.Yellow(text), true
	case "blue":
		return c.Blue(text), true
	case "magenta":
		return c.Magenta(text), true
	case "cyan":
		return c.Cyan(text), true
	case "white":
		return c.White(text), true
	default:
		return text, false
	}
}

// Black returns black text
func (c *Colorer) Black(text string) string {
	return color.BlackString(text)
}

// Red returns red text
func (c *Colorer) Red(text string) string {
	return color.RedString(text)
}

// Green returns green text
func (c *Colorer) Green(text string) string {
	return color.GreenString(text)
}

// Yellow returns yellow text
func (c *Colorer) Yellow(text string) string {
	return color.YellowString(text)
}

// Blue returns blue text
func (c *Colorer) Blue(text string) string {
	return color.BlueString(text)
}

// Magenta returns magenta text
func (c *Colorer) Magenta(text string) string {
	return color.MagentaString(text)
}

// Cyan returns cyan text
func (c *Colorer) Cyan(text string) string {
	return color.CyanString(text)
}

// White returns white text
func (c *Colorer) White(text string) string {
	return color.WhiteString(text)
}
