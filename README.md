# go-find-tests
Tool to determine which tests cover a block of code

**NOTE:** This tool is still in beta and there may be backwards incompatible changes prior to v1.0.0
## Overview
`go-find-tests` finds test functions which cover a position specified by a file path, line, and optionally column. 
Covering test are written to stdout and any encountered errors are written to stderr.

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

## Project Status
This project is still in "beta" since I want to be able to quickly change the public API in order to enable additional tooling such as editor plugins.
