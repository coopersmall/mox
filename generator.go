package moxie

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

type generator interface {
	generateCode(objects ...object)
}

type generatorImpl struct {
	template *template.Template
}

func newGenerator() generator {
	return &generatorImpl{
		template: template.Must(template.New("mock").Parse(mockTemplate)),
	}
}

func (g *generatorImpl) generateCode(objects ...object) {
	for _, obj := range objects {
		var s strings.Builder
		if err := g.template.Execute(&s, obj); err != nil {
			fmt.Printf("Error generating mock code for %s: %v\n", obj.Name, err)
			continue
		}

		mockFileName := strings.ToLower(obj.Name) + "_mock.go"
		if err := os.WriteFile(mockFileName, []byte(s.String()), 0644); err != nil {
			fmt.Printf("Error writing mock file for %s: %v\n", obj.Name, err)
		} else {
			fmt.Printf("Generated mock file: %s\n", mockFileName)
		}
	}
}
