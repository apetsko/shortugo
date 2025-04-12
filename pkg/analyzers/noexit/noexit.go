// Package noexit provides a static analysis tool that prohibits direct calls
// to os.Exit in the main function of a Go program. This ensures better control
// over program termination and encourages the use of structured error handling.
package noexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer defines the static analysis tool that checks for direct calls
// to os.Exit in the main function.
var Analyzer = &analysis.Analyzer{
	Name: "noexit",                                                 // The name of the analyzer.
	Doc:  "prohibits direct calls to os.Exit in the main function", // Description of the analyzer.
	Run:  run,                                                      // The function that performs the analysis.
}

// run is the main function of the analyzer. It inspects the AST of the package
// to find and report any direct calls to os.Exit in the main function.
func run(pass *analysis.Pass) (interface{}, error) {
	// Skip analysis if the package is not "main".
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	// Iterate over all files in the package.
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// Look for the main function declaration.
			fn, ok := n.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "main" || fn.Body == nil {
				return true
			}

			// Inspect the body of the main function for calls to os.Exit.
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				// Check if the call is to os.Exit.
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok || sel.Sel.Name != "Exit" {
					return true
				}

				// Verify that the selector is from the "os" package.
				if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "os" {
					// Report the usage of os.Exit.
					pass.Reportf(call.Pos(), "do not use os.Exit in main")
				}

				return true
			})
			return false
		})
	}
	return nil, nil
}
