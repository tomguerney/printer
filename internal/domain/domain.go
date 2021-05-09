package domain

// Formatter formats text for output
type Formatter interface {
	Msg(string, ...interface{}) string
	Error(string, ...interface{}) string
	Tabulate([][]string) []string
}

// Stenciller takes maps and creates strings from their values defined
// by predefined templates and colors
type Stenciller interface {
	AddStencil(id, template string, colors map[string]string) error
	UseStencil(id string, data map[string]string) (string, error)
}

// Colorer colors text
type Colorer interface {
	Color(text, color string) (string, bool)
}
