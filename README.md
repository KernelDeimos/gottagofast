# GottaGoFast

## Introduction
Do you ever need to perform some complicated string manipulation or parsing
function, that isn't in the standard library and really shouldn't be, and
don't want to port code from another language or write it from scratch?

This may be the repo for you!

GottaGoFast includes such utilities that were written for my interpreted code
generator and other tools. This includes things like performing multiple
string replacements "atomically" (i.e. without them affecting each other)
and parsing (lists with (nested) brackets).

GottaGoFast is named as such because it's a great package to get your
ideas into code as quickly as possible. While there is currently no
benchmarking, contributions of benchmarking tests and optimizations are
welcome and encouraged.

## Package Naming
Packages are named in the format `typesubject`, where `type` is one of the
following and `subject` describes what the contents of the package are used for:
- `sc`   - Shortcut functions. These provide no real functionality.
- `util` - General-purpose utility functions, like atomic replace for strings.
- `data` - Definitions or data structures. These will usually exist for use
           by tool* packages.
- `tool` - Specific implementations that could be useful across many projects.

## Contribution Ideas

### Code Generation
Some of the packages here are implementations of data structures. These use an
empty interface{} and you will need to perform type assertions. In some cases
it may be preferable to use code generation instead.

A CLI tool that generates each of these data structures for a specified type
would be a convenient solution, and could be used with `go:generate`.
