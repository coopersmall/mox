package mox

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

const (
	version = "0.1.0"
	help    = `
Moxie is a tool for generating mock implementations of Go interfaces.
    Usage: moxie [file]...
    Flags:
        -h, --help      Show this help message
        -v, --version   Show the version of moxie`
)

// Run runs the moxie generator.
func Run() {
	args := args()
	generator := newGenerator()

	files := getFiles(args)
	for _, file := range files {
		generator.generateCode(getSpecs(file))
	}
}

func args() []string {
	args := os.Args[1:]

	for _, arg := range args {
		switch arg {
		case "-h", "--help":
			fmt.Println(help)
			os.Exit(0)
		case "-v", "--version":
			fmt.Println(fmt.Sprintf("moxie version %s", version))
			os.Exit(0)
		default:
			if !strings.HasPrefix(arg, ".go") {
				fmt.Printf("Error: only .go files are supported, got %s\n", arg)
				os.Exit(1)
			}
		}
	}

	return args
}

func getFiles(args []string) []string {
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

func getSpecs(filename string) []spec {
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
	specs := make([]spec, 0)

	// Inspect type declarations and identify interfaces and structs
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		// Inspect type specs and identify interfaces
		for _, s := range genDecl.Specs {
			typeSpec, ok := s.(*ast.TypeSpec)
			if !ok {
				continue
			}

			if _, isInterface := typeSpec.Type.(*ast.InterfaceType); !isInterface {
				continue
			}

			spec := spec{
				Package: packageName,
				Imports: imports,
				Name:    typeSpec.Name.Name,
				Methods: specMethods(typeSpec.Type.(*ast.InterfaceType).Methods.List),
			}

			specs = append(specs, spec)
		}

	}
	return specs
}

func specMethods(methods []*ast.Field) []method {
	methodsList := make([]method, 0)

	for _, m := range methods {
		methodsList = append(methodsList, method{
			Name:    m.Names[0].Name,
			Params:  specMethodParams(m.Type.(*ast.FuncType).Params.List),
			Returns: specMethodReturns(m.Type.(*ast.FuncType).Results.List),
		})
	}

	return methodsList
}

func specMethodParams(params []*ast.Field) []param {
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

func specMethodReturns(results []*ast.Field) []ret {
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
