package mox

const mockTemplate = `package {{ .Package }}

import (
    . "github.com/coopersmall/penthouse/mock"
    {{- range .Imports }}{{ . }}{{- end }}
)

type {{ .Name }}Mock struct {
    Mock
}

func New{{ .Name }}Mock() *{{ .Name }}Mock {
    return &{{ .Name }}Mock{
        Mock: NewMock(),
    }
}
{{ range .Methods }}
func (m *{{ $.Name }}Mock) {{ .Name }}({{ $numParams := len .Params}}{{ range $idx, $param := .Params }}{{ if and (ne $idx $numParams) (ne $idx 0) }}, {{end}}{{ .Name }} {{ .Type }}{{ end }}) {{ if ne 0 (len .Returns) }}({{$numRets := len .Returns}}{{ range $idx, $ret := .Returns }}{{if and (ne $idx $numRets) (ne $idx 0)}}, {{end}}{{ .Type }}{{ end }}){{ end }} {
    rets := m.Mock.CallMethod("{{.Name}}", {{ range $idx, $params := .Params }}{{if and (ne $idx $numParams) (ne $idx 0)}}, {{end}}{{ $params.Name }}{{ end }})
    return {{ $numReturns := len .Returns }}{{ range $idx, $ret := .Returns }}{{ if and (ne $idx $numReturns) (ne $idx 0 ) }}, {{end}}{{if eq "error" $ret.Type}}Error(rets[{{ $idx }}]){{ else }}rets[{{ $idx }}].({{ $ret.Type }}){{ end }}{{ end }}
}
{{ end }}
`
