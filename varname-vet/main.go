// varname-vet is a go vet tool enforcing a style for local variable names: a
// variable declared inside a function body is named with two camel-case words
// at most, activeAddr rather than activeUserAddr.
//
//	go install github.com/pin2t/govets/varname-vet@latest
//	go vet -vettool=$(go env GOPATH)/bin/varname-vet ./...
//
// A -vettool replaces vet's own analyzers rather than adding to them, so this
// runs as a second go vet beside the plain one, not instead of it.
//
// Every variable a statement declares is checked: var, :=, a range clause, a
// type switch's x := y.(type) and a select case, in a function declaration or
// a function literal, wherever the literal sits. Package-level variables are
// not, and neither are parameters, results, receivers, struct fields or
// constants. A variable is reported once, where it is declared, and generated
// files are skipped.
package main

import "go/ast"
import "go/types"
import "strings"
import "unicode"
import "golang.org/x/tools/go/analysis"
import "golang.org/x/tools/go/analysis/unitchecker"

var analyzer = &analysis.Analyzer{
	Name: "varname",
	Doc:  "check that variables inside function bodies are named with two camel-case words at most",
	Run:  run,
}

func main() {
	unitchecker.Main(analyzer)
}

// Every name a field list declares is a parameter, result, receiver, type
// parameter, struct field or method, so field lists are skipped whole. A type
// switch's symbol has no object of its own in Defs, so it is checked from the
// statement.
func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		if ast.IsGenerated(f) {
			continue
		}
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FieldList:
				return false
			case *ast.TypeSwitchStmt:
				if a, ok := x.Assign.(*ast.AssignStmt); ok {
					check(pass, a.Lhs[0].(*ast.Ident))
				}
			case *ast.Ident:
				if v, ok := pass.TypesInfo.Defs[x].(*types.Var); ok && v.Parent() != pass.Pkg.Scope() {
					check(pass, x)
				}
			}
			return true
		})
	}
	return nil, nil
}

func check(pass *analysis.Pass, id *ast.Ident) {
	var w = words(id.Name)
	if len(w) > 2 {
		pass.Reportf(id.Pos(), "variable name has %d words (%s): use two at most", len(w), strings.Join(w, ", "))
	}
}

// words splits a name into camel-case words, breaking at underscores as well.
// A new word starts at a capital that follows a lower-case letter or a digit,
// so digits stay with the word before them: utf8Reader is utf8 and Reader. A
// run of capitals is one word, an acronym, but its last capital starts the next
// word when a lower-case letter follows it: parsedJSONDocument is parsed, JSON
// and Document. A lone s closing an acronym is its plural, as in userIDs.
func words(name string) []string {
	var runes = []rune(name)
	var result []string
	var word []rune
	for i, c := range runes {
		var prev, next rune
		if i > 0 {
			prev = runes[i-1]
		}
		if i+1 < len(runes) {
			next = runes[i+1]
		}
		var plural = next == 's' && (i+2 == len(runes) || !unicode.IsLower(runes[i+2]))
		var starts = unicode.IsUpper(c) && (unicode.IsLower(prev) || unicode.IsDigit(prev) ||
			unicode.IsUpper(prev) && unicode.IsLower(next) && !plural)
		if (c == '_' || starts) && len(word) > 0 {
			result = append(result, string(word))
			word = nil
		}
		if c != '_' {
			word = append(word, c)
		}
	}
	if len(word) > 0 {
		result = append(result, string(word))
	}
	return result
}
