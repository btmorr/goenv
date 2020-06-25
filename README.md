# gvm

[![Build Status][build-badge]][build]
[![Coverage Status][coverage]][coverage-badge]
[![License][license-badge]][license]

A version manager for Go

The goal of this project is to make it possible to have different project in a system using different versions of Go and automatically (as much as possible) use the version specified by each project's "go.mod" file. Inspiration is drawn heavily from [rbenv]

To build, run tests, and run the app:

```
make
make run
```

Or to do so manually:

```
go build -o gvm
./gvm
```

[rbenv]: https://github.com/rbenv/rbenv#how-rbenv-hooks-into-your-shell

[build]: https://travis-ci.com/btmorr/gvm
[build-badge]: https://travis-ci.com/btmorr/gvm.svg?branch=edge
[coverage]: https://coveralls.io/repos/github/btmorr/gvm/badge.svg?branch=edge
[coverage-badge]: https://coveralls.io/github/btmorr/gvm?branch=edge
[license]: https://github.com/btmorr/gvm/LICENSE
[license-badge]: https://img.shields.io/github/license/btmorr/gvm.svg
