# vardecl-vet

A `go vet` tool that checks every variable is declared with the `var` keyword:
`var a = f()` rather than `a := f()`.

## What it reports

Every short variable declaration (`:=`) that could be written with `var`:

```go
a := f()        // short variable declaration: declare it with the var keyword
b, err := g()   // short variable declaration: declare it with the var keyword
flush := func() {
    // ...
}               // reported at flush
```

Write these instead:

```go
var a = f()
var b, err = g()
var flush = func() {
    // ...
}
```

`var c int`, `var c int = 1` and grouped `var ( ... )` blocks are all fine.

## What it allows

Go's grammar doesn't allow `var` in these places, so `:=` is allowed there:

```go
for i := 0; i < n; i++ {}          // for init statement
for i, v := range s {}             // range clause
if v, ok := m[k]; ok {}            // if init statement, including else if
switch n := len(s); n {}           // switch init statement
switch t := x.(type) {}            // type switch
select {
case v, ok := <-ch:                // select case receiving into new variables
}
```

Only the `:=` in the statement's header is allowed. A `:=` in the body of the
`for`, `if`, `switch` or `select` is reported like any other.

**Redeclaring a variable from the same scope is allowed.** Here `err` already
exists, so `var n, err = g()` would not compile (`err redeclared in this block`):

```go
var err error
// ...
n, err := g()     // allowed: err is reused, only n is new
```

**Shadowing in a nested block is reported.** A `:=` inside a nested block
declares a new variable even if the name exists outside, so `var` works there:

```go
if n == 0 {
    err := check() // short variable declaration: declare it with the var keyword
}
```

Generated files (with a `// Code generated ... DO NOT EDIT.` header) are skipped.

## Install

```sh
go install github.com/pin2t/govets/vardecl-vet@latest
```

This needs Go 1.25 or newer. To pin the version in your `go.mod` instead (Go
1.24+ tool directives):

```sh
go get -tool github.com/pin2t/govets/vardecl-vet@latest
```

## Run

```sh
go vet -vettool="$(go env GOPATH)/bin/vardecl-vet" ./...
```

Or, with the tool directive:

```sh
go vet -vettool="$(go tool -n vardecl-vet)" ./...
```

Findings look like this, and `go vet` exits with status 1:

```
# example.com/consumer
# [example.com/consumer]
./main.go:6:2: short variable declaration: declare it with the var keyword
./main.go:10:3: short variable declaration: declare it with the var keyword
```

In CI (GitHub Actions):

```yaml
- run: go vet ./...
- run: go install github.com/pin2t/govets/vardecl-vet@latest
- run: go vet -vettool="$(go env GOPATH)/bin/vardecl-vet" ./...
```

Things to know:

- **Keep running plain `go vet ./...` as well.** `-vettool` replaces vet's
  built-in analyzers instead of adding to them.
- **`-vettool` needs a path.** A bare `-vettool=vardecl-vet` fails even when
  the binary is on your `PATH`.
- **`-json` exits with status 0** even when there are findings.
- `_test.go` files are checked. `testdata` directories are skipped, as they are
  by `go vet` in general.
- It doesn't fix anything for you, and there is no way to silence a single
  finding.

## Tests

`go test ./vardecl-vet/` runs the analyzer with `analysistest` over
`testdata/src/`:

- `clean` holds `var` declarations and every allowed `:=` form, and must produce
  no findings.
- `short` holds `:=` declarations that must each be reported, each marked with a
  `// want` comment. These include `:=` in the bodies of `if`, `for`, `switch`,
  `select` and a function literal, and a nested `:=` that shadows an outer name.
- `generated` holds a generated file that must be skipped.
