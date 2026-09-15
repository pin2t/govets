// vardecl-vet is a go vet tool enforcing a style for variable declarations:
// every variable is declared with the var keyword, var a = f() rather than
// a := f().
//
//	go install github.com/pin2t/govets/vardecl-vet@latest
//	go vet -vettool=$(go env GOPATH)/bin/vardecl-vet ./...
//
// A -vettool replaces vet's own analyzers rather than adding to them, so this
// runs as a second go vet beside the plain one, not instead of it.
//
// A short variable declaration is left alone where Go's grammar has no room
// for var: the init statement of an if, for or switch, a for's range clause, a
// type switch's x := y.(type), and a select case receiving into new variables.
// So is one redeclaring a variable from the same scope, like a, err := g()
// after an earlier err, because var would not compile there. A := inside a
// nested block that shadows an outer variable declares a new one, and is
// reported. Generated files are skipped.
package main

import "go/ast"
import "go/token"
import "golang.org/x/tools/go/analysis"
import "golang.org/x/tools/go/analysis/unitchecker"

var analyzer = &analysis.Analyzer{
	Name: "vardecl",
	Doc:  "check that variables are declared with the var keyword rather than :=",
	Run:  run,
}

func main() {
	unitchecker.Main(analyzer)
}

// ast.Inspect visits a statement before its children, so the init, type switch
// and select case statements are marked exempt before they are reached.
func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		if ast.IsGenerated(f) {
			continue
		}
		var exempt = map[ast.Stmt]bool{}
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.IfStmt:
				exempt[x.Init] = true
			case *ast.ForStmt:
				exempt[x.Init] = true
			case *ast.SwitchStmt:
				exempt[x.Init] = true
			case *ast.TypeSwitchStmt:
				exempt[x.Init] = true
				exempt[x.Assign] = true
			case *ast.CommClause:
				exempt[x.Comm] = true
			case *ast.AssignStmt:
				if x.Tok != token.DEFINE || exempt[x] {
					return true
				}
				for _, l := range x.Lhs {
					if id, ok := l.(*ast.Ident); ok && pass.TypesInfo.Uses[id] != nil {
						return true
					}
				}
				pass.Reportf(x.Pos(), "short variable declaration: declare it with the var keyword")
			}
			return true
		})
	}
	return nil, nil
}
