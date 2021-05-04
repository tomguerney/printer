package colorizer

import (
	"fmt"

	"github.com/fatih/color"
)

// Colorizer encapsulates colorizing text
type Colorizer struct{}

// Colorize transforms the text into one of the available colors
func (c *Colorizer) Colorize(colorName string, text string) (string, error) {
	switch colorName {
	case "black":
		return c.Black(text), nil
	case "red":
		return c.Red(text), nil
	case "green":
		return c.Green(text), nil
	case "yellow":
		return c.Yellow(text), nil
	case "blue":
		return c.Blue(text), nil
	case "magenta":
		return c.Magenta(text), nil
	case "cyan":
		return c.Cyan(text), nil
	case "white":
		return c.White(text), nil
	}
	return "", fmt.Errorf("color \"%v\" not available", text)
}

// Black returns black text
func (c *Colorizer) Black(text string) string {
	return color.BlackString(text)
}

// Red returns red text
func (c *Colorizer) Red(text string) string {
	return color.RedString(text)
}

// Green returns green text
func (c *Colorizer) Green(text string) string {
	return color.GreenString(text)
}

// Yellow returns yellow text
func (c *Colorizer) Yellow(text string) string {
	return color.YellowString(text)
}

// Blue returns blue text
func (c *Colorizer) Blue(text string) string {
	return color.BlueString(text)
}

// Magenta returns magenta text
func (c *Colorizer) Magenta(text string) string {
	return color.MagentaString(text)
}

// Cyan returns cyan text
func (c *Colorizer) Cyan(text string) string {
	return color.CyanString(text)
}

// White returns white text
func (c *Colorizer) White(text string) string {
	return color.WhiteString(text)
}
