package moxie

type object struct {
	Package string
	Imports []string
	Name    string
	Methods []method
}

type method struct {
	Name    string
	Params  []param
	Returns []ret
}

type param struct {
	Name string
	Type string
}

type ret struct {
	Type string
}
