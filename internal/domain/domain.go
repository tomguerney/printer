package domain

type Printer interface {
	
}

type Templater interface {
	addTemplate(id, template string, colors map[string]string)
	useTemplate(id string, s interface{})
}

