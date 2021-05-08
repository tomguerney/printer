package domain

type Formatter interface {
	Msg(string, ...interface{}) string
	Error(string, ...interface{}) string
	Tabulate([][]string) []string
}

// Stenciller takes structs and formats their field values to predefined templates and colors
type Stenciller interface {
	AddStencil(id, template string, colors map[string]string)
	UseStencil(id string, s interface{}) (string, error)
}
