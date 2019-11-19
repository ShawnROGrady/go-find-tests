# go-find-tests
Tool to determine which tests cover a block of code

**NOTE:** This tool is still in beta and there may be backwards incompatible changes prior to v1.0.0
## Overview
`go-find-tests` finds test functions which cover a position specified by a file path, line, and optionally column. 
Covering test are written to stdout and any encountered errors are written to stderr.

Sample usage:
```
$ go-find-tests ./cover/profile.go:155.12 
TestCovers
TestParseLine
```
## Options
### Behaviour

1. `-include-subs`: Find specific sub-tests which cover the specified block (default = false)
2. `-print-positions`: Print the positions of the found tests (default false)
    - **NOTE:** subtests will not have position information
3. `-run regexp`: Check only tests matching the regular expression (default = '.')
4. `-short`: Sets '-short' flag when testing for coverage (default = false)
    - see `go help testflag` for info
5. `-seq`: Run all tests sequentially. Greatly reduces performance but may be neccessary for integration tests (default = false)
6. `-h|-help`: Print a help message and exit (default = false)
### Formatting

1. `-json`: Print the output in json format instead of as a newline separated list (default = false)
2. `-line-fmt string`: With `-print-positions` - the fmt to use when writing the postions of found test (defualt = `%t:%f:%l:%c:%s`)
    - `%t`: test name
    - `%f`: file
    - `%l`: line
    - `%c`: column
    - `%o`: offset
    - `%s`: subtests

## Troubleshooting
Please try the following, if the problem persists feel free to open an issue or submit a pull request.
### I'm seeing an error
* Does the error start with "Error constructing tester"?
    - make sure provided filepath begins with "." in positional arg
* Does the error start with "Error determining covering tests"?
    - do the tests pass with `-count=1` and `-race` set?
        - these flags aren't set while determining coverage but they may indicate an underlying problem
    - do the tests have any dependency on the file structure (e.g. loading a file from a relative path)?
        - this may cause issues since for performance reasons this tool first compiles a test binary instead of running `go test` directly
        - if this ends up being the issue please consider opening an issue since the tool should be able to handle this
    - do the tests rely on a connection to some external process (e.g. db, http connections)
        - by default this tool runs each test in a separate go routine for performance reasons, which may cause conflicts when establishing these connections.
        - the `-seq` flag will result in tests being ran sequentially instead

### Unexpected results (no covering tests, specific test not returned)
* do coverage visualization tools (such as `go tool cover`) mark the specified position as covered?
    - often things like brackets an parentheses won't be marked
* do the tests pass with `-count=1` and `-race` set?
    - these flags aren't set while determining coverage but they may indicate an underlying problem
* was a column provided, along with the file and line, to the tool?
    - often just a line is sufficient, but a column should be provided for truly accurate results

## Project Status
This project is still in "beta" since I want to be able to quickly change the public API in order to enable additional tooling such as editor plugins.
