# varname-vet

A `go vet` tool that checks every variable declared inside a function body is
named with two camel-case words at most: `activeAddr` rather than
`activeUserAddr`.

## What it reports

Every local variable whose name has more than two words:

```go
func run(jobs []Job) {
    var address = 123                // fine: one word
    var activeAddr = true            // fine: two words
    var config = Config{}            // fine: one word
    var parsedJSONDocument = parse() // variable name has 3 words (parsed, JSON, Document): use two at most
    userConfigParam := true          // variable name has 3 words (user, Config, Param): use two at most
    for jobIndexValue := range jobs { // variable name has 3 words (job, Index, Value): use two at most
    }
}
```

Every way a statement declares a variable is checked: `var`, `:=`, `for` and
`range` headers, `if` and `switch` init statements, a type switch's
`x := y.(type)` and a `select` case. So is a variable in a function literal's
body, including a literal assigned to a package-level variable.

Each variable is reported once, where it is declared. A later `:=` that reuses
the name, like `n, err := g()` after an earlier `n`, isn't reported again.

## What it doesn't check

```go
var parsedJSONDocument = 1                  // package-level variable

const defaultListenAddress = ":8080"        // constant, here or in a function

type serverConfig struct {
    listenAddressPort int                   // struct field
}

func (currentServerConfig *serverConfig) load(configFilePath string) (loadedConfigValue serverConfig, loadErrorValue error) {
    var handler = func(requestBodyText string) {}   // function literal parameter
    // ...
}
```

Parameters, results and receivers belong to the function's signature rather
than its body. That holds for function literals too: a literal's parameters
aren't checked, but the variables in its body are. Labels, type names and
constants aren't variables, so they're not checked either.

Generated files (with a `// Code generated ... DO NOT EDIT.` header) are skipped.

## How words are counted

- A capital letter after a lower-case letter or a digit starts a new word:
  `activeAddr` is `active` and `Addr`, and `utf8Reader` is `utf8` and `Reader`.
- A run of capitals is one word, an acronym. Its last capital starts the next
  word when lower-case letters follow: `parsedJSONDocument` is `parsed`, `JSON`
  and `Document`, and `HTTPServer` is `HTTP` and `Server`.
- A lone `s` after an acronym is its plural: `userIDs` is `user` and `IDs`.
- Digits belong to the word before them: `sha256Sum` and `HTTP2Server` are two
  words each.
- Underscores separate words, so `user_config_param` is three words.

The report lists the words it counted, so you can see how a name was split.

## Install

```sh
go install github.com/pin2t/govets/varname-vet@latest
```

This needs Go 1.25 or newer. To pin the version in your `go.mod` instead (Go
1.24+ tool directives):

```sh
go get -tool github.com/pin2t/govets/varname-vet@latest
```

## Run

```sh
go vet -vettool="$(go env GOPATH)/bin/varname-vet" ./...
```

Or, with the tool directive:

```sh
go vet -vettool="$(go tool -n varname-vet)" ./...
```

Findings look like this, and `go vet` exits with status 1:

```
# example.com/consumer
# [example.com/consumer]
./main.go:10:6: variable name has 3 words (parsed, JSON, Document): use two at most
./main.go:11:6: variable name has 3 words (user, Config, Param): use two at most
```

In CI (GitHub Actions):

```yaml
- run: go vet ./...
- run: go install github.com/pin2t/govets/varname-vet@latest
- run: go vet -vettool="$(go env GOPATH)/bin/varname-vet" ./...
```

Things to know:

- **Keep running plain `go vet ./...` as well.** `-vettool` replaces vet's
  built-in analyzers instead of adding to them.
- **`-vettool` needs a path.** A bare `-vettool=varname-vet` fails even when
  the binary is on your `PATH`.
- **`-json` exits with status 0** even when there are findings.
- `_test.go` files are checked. `testdata` directories are skipped, as they are
  by `go vet` in general.
- Word splitting can't read intent. An acronym with a lower-case letter
  inside it, like `IPv4Addr`, splits as `I`, `Pv4`, `Addr`. Name it `ipv4Addr`
  instead.
- It doesn't rename anything for you, and there is no way to silence a single
  finding.

## Tests

`go test ./varname-vet/` runs the analyzer with `analysistest` over
`testdata/src/`, and checks word splitting on its own in `TestWords`:

- `clean` holds one- and two-word local names, acronyms and their plurals, and
  long names the tool doesn't check: package-level variables, constants, struct
  fields, parameters, results, receivers, type parameters and labels. It must
  produce no findings.
- `long` holds local variables of three or more words, each marked with a
  `// want` comment. They're declared every way a statement can declare one,
  including in a function literal's body at package level and inside a function.
- `generated` holds a generated file that must be skipped.
