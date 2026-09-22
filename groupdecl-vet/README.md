# groupdecl-vet

A `go vet` tool that checks const and var declarations are not grouped. Each
constant and variable is declared with its own keyword and on its own line.

## What it reports

Every `const` or `var` declared inside a parenthesized group, reported at each
declaration in the group:

```go
const (
    a = 1 // grouped const declaration: give each constant its own const keyword and line
    b = 2 // grouped const declaration: give each constant its own const keyword and line
)

var (
    c = 3 // grouped var declaration: give each variable its own var keyword and line
    d = 4 // grouped var declaration: give each variable its own var keyword and line
)
```

Write these instead:

```go
const a = 1
const b = 2

var c = 3
var d = 4
```

A group with a single declaration is still a group and is reported. `type (...)` and
`import (...)` blocks are not const or var declarations and are left alone.

## Where it checks

The group can be anywhere Go allows a declaration:

```go
const a = 1 // fine: package level, own keyword

var closure = func() {
    const (
        b = 2 // grouped const declaration: ...
    )
}

func run() {
    var c = 3 // fine
    if c > 0 {
        var (
            d = 4 // grouped var declaration: ...
        )
        _ = d
    }
}
```

A group is reported whether it sits at package level, inside a function, inside
a nested block, or inside a function literal.

Generated files (with a `// Code generated ... DO NOT EDIT.` header) are skipped.

## Install

```sh
go install github.com/pin2t/govets/groupdecl-vet@latest
```

This needs Go 1.25 or newer. To pin the version in your `go.mod` instead (Go
1.24+ tool directives):

```sh
go get -tool github.com/pin2t/govets/groupdecl-vet@latest
```

## Run

```sh
go vet -vettool="$(go env GOPATH)/bin/groupdecl-vet" ./...
```

Or, with the tool directive:

```sh
go vet -vettool="$(go tool -n groupdecl-vet)" ./...
```

Findings look like this, and `go vet` exits with status 1:

```
# example.com/consumer
# [example.com/consumer]
./main.go:3:2: grouped const declaration: give each constant its own const keyword and line
./main.go:7:2: grouped var declaration: give each variable its own var keyword and line
```

In CI (GitHub Actions):

```yaml
- run: go vet ./...
- run: go install github.com/pin2t/govets/groupdecl-vet@latest
- run: go vet -vettool="$(go env GOPATH)/bin/groupdecl-vet" ./...
```

Things to know:

- **Keep running plain `go vet ./...` as well.** `-vettool` replaces vet's
  built-in analyzers instead of adding to them.
- **`-vettool` needs a path.** A bare `-vettool=groupdecl-vet` fails even when
  the binary is on your `PATH`.
- **`-json` exits with status 0** even when there are findings.
- `_test.go` files are checked. `testdata` directories are skipped, as they are
  by `go vet` in general.
- It doesn't fix anything for you, and there is no way to silence a single
  finding.

## Tests

`go test ./groupdecl-vet/` runs the analyzer with `analysistest` over
`testdata/src/`:

- `clean` holds ungrouped `const` and `var` declarations at package level, in
  functions, nested blocks and function literals, plus a `type` group, and must
  produce no findings.
- `grouped` holds grouped `const` and `var` declarations in all those places,
  each marked with a `// want` comment.
- `generated` holds a generated file that must be skipped.
