// groupdecl-vet is a go vet tool enforcing a style for const and var
// declarations: each constant and variable is declared with its own keyword and
// on its own line, not inside a parenthesized group.
//
//	go install github.com/pin2t/govets/groupdecl-vet@latest
//	go vet -vettool=$(go env GOPATH)/bin/groupdecl-vet ./...
//
// A -vettool replaces vet's own analyzers rather than adding to them, so this
// runs as a second go vet beside the plain one, not instead of it.
//
// Every const or var declared inside a parenthesized group is reported, each at
// its own declaration: at package level, inside a function, inside a nested
// block, or inside a function literal. Type and import groups are not const or
// var declarations and are left alone. Generated files are skipped.
package main

import "go/ast"
import "go/token"
import "golang.org/x/tools/go/analysis"
import "golang.org/x/tools/go/analysis/unitchecker"

var analyzer = &analysis.Analyzer{
	Name: "groupdecl",
	Doc:  "check that const and var declarations are not grouped",
	Run:  run,
}

func main() {
	unitchecker.Main(analyzer)
}

// ast.Inspect reaches every block, so a group is found wherever it sits: at
// package level, inside a function, a nested block or a function literal.
func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		if ast.IsGenerated(f) {
			continue
		}
		ast.Inspect(f, func(n ast.Node) bool {
			var g, ok = n.(*ast.GenDecl)
			if !ok || !g.Lparen.IsValid() {
				return true
			}
			var noun, keyword string
			switch g.Tok {
			case token.CONST:
				noun, keyword = "constant", "const"
			case token.VAR:
				noun, keyword = "variable", "var"
			default:
				return true
			}
			for _, s := range g.Specs {
				pass.Reportf(s.Pos(), "grouped %s declaration: give each %s its own %s keyword and line", keyword, noun, keyword)
			}
			return true
		})
	}
	return nil, nil
}
