# govets

Useful vet tools for Go. Each one is a `go/analysis` analyzer built as a
standalone `go vet -vettool` binary, so it needs nothing but the Go toolchain
to run in any project.

| Tool | What it enforces |
| --- | --- |
| [`imports-vet`](imports-vet/main.go) | Each import has its own `import` keyword and its own line, with no blank line between the first import and the last |
| [`funcbody-vet`](funcbody-vet/main.go) | No blank lines and no comments inside function bodies |
| [`vardecl-vet`](vardecl-vet/README.md) | Every variable is declared with the `var` keyword, not `:=`, wherever Go allows it |
| [`varname-vet`](varname-vet/README.md) | Variables inside function bodies are named with two camel-case words at most |

All of them skip generated files (those with a `// Code generated ... DO NOT EDIT.` header).

## imports-vet

Reports a grouped `import (...)` block, two imports on the same line, and blank
lines between imports. A run of blank lines is reported once, at its first line.

```go
import (            // grouped imports: give each import its own import keyword and line
    "fmt"
                    // blank line in the import block
    "os"
)
```

Write this instead:

```go
import "fmt"
// A comment between imports is fine: it is not a blank line.
import "os"
```

## funcbody-vet

Reports every comment and every run of blank lines between a function's `{` and
`}`. If a comment explains something, it belongs in the doc comment above the
function.

```go
func run() {
    // load the config      // comment inside a function: move it to the doc comment
    var c = load()
                            // blank line inside a function
    start(c)
}
```

What counts as inside a body:

- A function literal is part of the body around it, so a comment in a closure is
  reported once. A literal outside any function, such as a package-level `var`,
  is a body of its own.
- A blank line inside a multi-line raw string belongs to the string and is not
  reported.
- A doc comment, a comment after the closing `}`, and a blank line in a
  parameter list that spans several lines are all outside the body.

## vardecl-vet

Reports every `:=` that could be written with `var`. It allows `:=` where Go
doesn't allow `var`: `for`, `range`, `if` and `switch` headers, type switches
and `select` cases. It also allows a `:=` that reuses a variable already
declared in the same scope.

```go
a := f()              // short variable declaration: declare it with the var keyword
var b = f()           // fine
for i := range s {}   // fine: var isn't allowed here
```

See [vardecl-vet/README.md](vardecl-vet/README.md) for the full rules.

## varname-vet

Reports every variable declared inside a function body whose name has more than
two camel-case words. Package-level variables, parameters, results, struct
fields and constants aren't checked.

```go
var parsedJSONDocument = 1     // not checked: a package-level variable

func run() {
    var address = 123          // fine: one word
    var activeAddr = true      // fine: two words
    var userConfigParam = true // variable name has 3 words (user, Config, Param): use two at most
}
```

An acronym counts as one word, and so does its plural: `userIDs` and
`HTTPServer` are two words each. Underscores separate words too.

See [varname-vet/README.md](varname-vet/README.md) for the full rules.

## Using the tools in your project

### 1. Install

```sh
go install github.com/pin2t/govets/imports-vet@latest
go install github.com/pin2t/govets/funcbody-vet@latest
go install github.com/pin2t/govets/vardecl-vet@latest
go install github.com/pin2t/govets/varname-vet@latest
```

This needs Go 1.25 or newer.

To pin the version in your own `go.mod` instead (Go 1.24+ tool directives):

```sh
go get -tool github.com/pin2t/govets/imports-vet@latest
go get -tool github.com/pin2t/govets/funcbody-vet@latest
go get -tool github.com/pin2t/govets/vardecl-vet@latest
go get -tool github.com/pin2t/govets/varname-vet@latest
```

This adds `golang.org/x/tools` to your module graph as an indirect dependency.

### 2. Run

```sh
go vet -vettool="$(go env GOPATH)/bin/imports-vet" ./...
go vet -vettool="$(go env GOPATH)/bin/funcbody-vet" ./...
go vet -vettool="$(go env GOPATH)/bin/vardecl-vet" ./...
go vet -vettool="$(go env GOPATH)/bin/varname-vet" ./...
```

With the tool directive, let `go tool -n` find the binary for you:

```sh
go vet -vettool="$(go tool -n imports-vet)" ./...
go vet -vettool="$(go tool -n funcbody-vet)" ./...
go vet -vettool="$(go tool -n vardecl-vet)" ./...
go vet -vettool="$(go tool -n varname-vet)" ./...
```

Findings look like this, and `go vet` exits with status 1:

```
# example.com/consumer
# [example.com/consumer]
./main.go:3:1: grouped imports: give each import its own import keyword and line
./main.go:5:1: blank line in the import block
```

Things to know:

- **Keep running plain `go vet ./...` as well.** `-vettool` *replaces* vet's
  built-in analyzers (`printf`, `copylocks` and the rest) instead of adding to
  them. That's why each tool runs as its own `go vet` command.
- **`-vettool` needs a path.** A bare `-vettool=imports-vet` fails even if the
  binary is on your `PATH`, so use `$(go env GOPATH)/bin/...`,
  `$(which imports-vet)` or `$(go tool -n imports-vet)`. If you set `GOBIN`, the
  binaries are installed there instead.
- **`-json` exits with status 0** even when there are findings, so check the
  output instead of the exit code.
- `_test.go` files are checked, while `testdata` directories are skipped, as they
  are by `go vet` in general.
- There is no way to silence a single finding. Only generated files are skipped.
- Tools that add imports for you, such as goimports or your editor's
  organize-imports, may create grouped blocks. `imports-vet` will flag those
  until you split them.

### 3. Add to CI

GitHub Actions:

```yaml
- uses: actions/setup-go@v6
  with:
    go-version: stable
- run: go vet ./...
- run: go install github.com/pin2t/govets/imports-vet@latest github.com/pin2t/govets/funcbody-vet@latest github.com/pin2t/govets/vardecl-vet@latest github.com/pin2t/govets/varname-vet@latest
- run: go vet -vettool="$(go env GOPATH)/bin/imports-vet" ./...
- run: go vet -vettool="$(go env GOPATH)/bin/funcbody-vet" ./...
- run: go vet -vettool="$(go env GOPATH)/bin/vardecl-vet" ./...
- run: go vet -vettool="$(go env GOPATH)/bin/varname-vet" ./...
```

Replace `@latest` with a commit hash or tag so a new rule doesn't break your
build without warning.

A Makefile target works the same way:

```make
GOBIN := $(shell go env GOPATH)/bin

vet:
	go vet ./...
	go vet -vettool=$(GOBIN)/imports-vet ./...
	go vet -vettool=$(GOBIN)/funcbody-vet ./...
	go vet -vettool=$(GOBIN)/vardecl-vet ./...
	go vet -vettool=$(GOBIN)/varname-vet ./...
```

## Development

```sh
go test -race ./...
```

Each tool's test runs the analyzer with `analysistest` over the packages in its
`testdata/src/`. Expected findings are marked with `// want "..."` comments.
funcbody-vet's fixtures can't put that comment inside a function, because it
would be a comment inside a function itself. They put it above the function
with a line offset instead: `// want +3 "..."` expects the finding three lines
below the comment.

CI ([`.github/workflows/ci.yml`](.github/workflows/ci.yml)) builds, vets and
tests the module on Go 1.25 and the latest stable Go. It also runs all the tools
over this repository, installed the same way as above.

**Don't reformat the files under `testdata/`.** They are deliberately badly
formatted: same-line imports, runs of blank lines, `// want` comments on exact
lines. gofmt, or an IDE's reformat or optimize-imports on commit, rewrites them
and breaks the tests. `go vet ./...` never reads `testdata`, so the fixtures
can't fail the vet steps. They only fail `go test`.
