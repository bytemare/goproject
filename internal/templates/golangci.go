// Package templates holds the template and project building functions
package templates

const golangciIdentifier = "golangci"

// golangciConstructor returns the file content populated with the relevant values
func golangciConstructor(project *Project) (*file, error) { //nolint:unparam // project is not needed when no variables
	i, d, t := golangciValues()
	return newProjectFile(newFile(golangciIdentifier, i, d, t))
}

func golangciValues() (f, d, t string) { //nolint:funlen // length is due to a constant, no complexity here
	const filename = ".golangci.yml"

	const directory = "."

	const template = `linters:
  disable-all: true
  enable:
    - bodyclose
    - deadcode
    - depguard
    - dogsled
    - dupl
    - errcheck
    - funlen
    - gochecknoglobals
    - gochecknoinits
    - gocognit
    - goconst
    - gocritic
    - gocyclo
    - godox
    - gofmt
    - goimports
    - golint
    - gomnd
    - goprintffuncname
    - gosec
    - gosimple
    - govet
    - ineffassign
    - interfacer
    - lll
    - maligned
    - megacheck
    - misspell
    - nakedret
    - prealloc
    - rowserrcheck
    - scopelint
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - varcheck
    - whitespace
    - wsl
  presets:
    - bugs
    - unused
  fast: false


linters-settings:
  dogsled:
    max-blank-identifiers: 2
  dupl:
    threshold: 100
  errcheck:
    check-type-assertions: true
    check-blank: true
  funlen:
    lines: 100
    statements: 50
  gocognit:
    min-complexity: 15
  goconst:
    min-len: 1
    min-occurrences: 1
  gocritic:
    enabled-tags:
      - diagnostic
      - experimental
      - opinionated
      - performance
      - style

    settings: # settings passed to gocritic
      captLocal: # must be valid enabled check name
        paramsOnly: true
      rangeValCopy:
        sizeThreshold: 32
  gocyclo:
    min-complexity: 15
  godox:
    keywords: # default keywords are TODO, BUG, and FIXME, these can be overwritten by this setting
      - NOTE
      - OPTIMIZE
      - HACK
  gofmt:
    simplify: true
  goimports:
    local-prefixes: github.com/bytemare/goproject
  golint:
    min-confidence: 0
  gomnd:
    settings:
      mnd:
        checks: argument,case,condition,operation,return,assign
  govet:
    check-shadowing: true

    # settings per analyzer
    settings:
      printf: # analyzer name, run go tool vet help to see all analyzers
        funcs: # run go tool vet help printf to see available settings for printf analyzer
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Infof
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Warnf
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Errorf
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Fatalf

    # enable or disable analyzers by name
    enable:
      - atomicalign
    #enable-all: false
    # disable:
    #   - shadow
    disable-all: false
  #depguard:
  #  list-type: blacklist
  #  include-go-root: false
  #  packages:
  #    - github.com/sirupsen/logrus
  #  packages-with-error-message:
  #    # specify an error message to output when a blacklisted package is used
  #    - github.com/sirupsen/logrus: "logging is allowed only by logutils.Log"
  lll:
    line-length: 140
    # tab width ('\t') in spaces. Default to 1.
    tab-width: 1
  maligned:
    suggest-new: true
  misspell:
    locale: EN
    ignore-words:
      - someword
  nakedret:
    # make an issue if func has more lines of code than this setting and it has naked returns; default is 30
    max-func-lines: 30
  prealloc:
    # XXX: we don't recommend using this linter before doing performance profiling.
    # For most programs usage of prealloc will be a premature optimization.

    # Report preallocation suggestions only on simple loops that have no returns/breaks/continues/gotos in them.
    # True by default.
    simple: true
    range-loops: true # Report preallocation suggestions on range loops, true by default
    for-loops: true # Report preallocation suggestions on for loops, false by default
  rowserrcheck:
    packages:
      - github.com/jmoiron/sqlx
  unused:
    # treat code as a program (not a library) and report unused exported identifiers; default is false.
    # XXX: if you enable this setting, unused will report a lot of false-positives in text editors:
    # if it's called for subdir of a project it can't find funcs usages. All text editor integrations
    # with golangci-lint call it on a directory with the changed file.
    check-exported: false
  whitespace:
    multi-if: false   # Enforces newlines (or comments) after every multi-line if statement
    multi-func: false # Enforces newlines (or comments) after every multi-line function signature
  wsl:
    # If true append is only allowed to be cuddled if appending value is
    # matching variables, fields or types on line above. Default is true.
    strict-append: true
    # Allow calls and assignments to be cuddled as long as the lines have any
    # matching variables, fields or types. Default is true.
    allow-assign-and-call: true
    # Allow multiline assignments to be cuddled. Default is true.
    allow-multiline-assign: true
    # Allow declarations (var) to be cuddled.
    allow-cuddle-declarations: false
    # Allow trailing comments in ending of blocks
    allow-trailing-comment: false
    # Force newlines in end of case at this limit (0 = never).
    force-case-trailing-whitespace: 0
    # Force cuddling of err checks with err var assignment
    force-err-cuddling: false
    # Allow leading comments to be separated with empty liens
    allow-separated-leading-comment: false

issues:
  # List of regexps of issue texts to exclude, empty list by default.
  # But independently from this option we use default exclude patterns,
  # it can be disabled by exclude-use-default: false. To list all
  # excluded by default patterns execute golangci-lint run --help
  exclude:
    - abcdef

  # Excluding configuration per-path, per-linter, per-text and per-source
  exclude-rules:
    # Exclude some linters from running on tests files.
    - path: _test\.go
      linters:
        - gocyclo
        - errcheck
        - dupl
        - gosec

    # Exclude some staticcheck messages
    - linters:
        - staticcheck
      text: "SA9003:"

    # Exclude lll issues for long lines with go:generate
    - linters:
        - lll
      source: "^//go:generate "

  # Independently from option exclude we use default exclude patterns,
  # it can be disabled by this option. To list all
  # excluded by default patterns execute golangci-lint run --help.
  # Default value for this option is true.
  exclude-use-default: false

  # Maximum issues count per one linter. Set to 0 to disable. Default is 50.
  max-issues-per-linter: 0

  # Maximum count of issues with the same text. Set to 0 to disable. Default is 3.
  max-same-issues: 0

`

	return filename, directory, template
}
