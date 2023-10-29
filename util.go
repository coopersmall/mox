package moxie

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// Run runs the moxie generator.
func Run() {
	generator := newGenerator()

	files := getFiles(args()...)
	for _, file := range files {
		generator.generateCode(getObjects(file)...)
	}
}

func args() []string {
	return os.Args[1:]
}

func getFiles(args ...string) []string {
	if len(args) >= 1 {
		return args
	}

	files := make([]string, 0)

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fs, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range fs {
		if strings.HasSuffix(file.Name(), ".go") {
			files = append(files, file.Name())
		}
	}

	return files
}

func getObjects(filename string) []object {
	// Parse the source file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing source file: %v\n", err)
		return nil
	}

	// get the package name
	packageName := node.Name.Name

	// get imports
	imports := make([]string, 0)
	for _, imp := range node.Imports {
		imports = append(imports, imp.Path.Value)
	}

	// Create a map to store interfaces and structs
	objects := make([]object, 0)

	// Inspect type declarations and identify interfaces and structs
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		// Inspect type specs and identify interfaces
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			if _, isInterface := typeSpec.Type.(*ast.InterfaceType); !isInterface {
				continue
			}

			obj := object{
				Package: packageName,
				Imports: imports,
				Name:    typeSpec.Name.Name,
				Methods: getMethods(typeSpec.Type.(*ast.InterfaceType).Methods.List),
			}

			objects = append(objects, obj)
		}

	}
	return objects
}

func getMethods(methods []*ast.Field) []method {
	methodsList := make([]method, 0)

	for _, m := range methods {
		methodsList = append(methodsList, method{
			Name:    m.Names[0].Name,
			Params:  getParams(m.Type.(*ast.FuncType).Params.List),
			Returns: getReturns(m.Type.(*ast.FuncType).Results.List),
		})
	}

	return methodsList
}

func getParams(params []*ast.Field) []param {
	paramsList := make([]param, 0)

	for _, p := range params {
		for _, paramName := range p.Names {
			paramsList = append(paramsList, param{
				Name: paramName.Name,
				Type: p.Type.(*ast.Ident).Name,
			})
		}
	}

	return paramsList
}

func getReturns(results []*ast.Field) []ret {
	returns := make([]ret, 0)

	for _, r := range results {
		if r == nil {
			continue
		}

		returns = append(returns, ret{
			Type: r.Type.(*ast.Ident).Name,
		})
	}

	return returns
}
