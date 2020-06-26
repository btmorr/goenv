# gvm

[![Build Status][build-badge]][build]
[![License][license-badge]][license]

A version manager for Go

The goal of this project is to make it possible to have different project in a system using different versions of Go and automatically (as much as possible) use the version specified by each project's "go.mod" file. Inspiration is drawn heavily from [rbenv].

gmv works almost completely silently after installation. When a user invokes `go` in a directory with no "go.mod" file, gmv will use the system installation of `go` (primarily when running something like `go mod init` or quickly doing `go run myfile.go` to test something out). When a user invokes `go` in a directory that does have a "go.mod" file, gmv will check if a shim already exists for the go version specified in "go.mod", install a shim if it is missing, and then use the shim. When installing a shim, it will print a message like "gvm installing shim for go 1.14.4". Otherwise, it does not write anything to stdout/stderr--no need to create additional files or use commands to ensure that gvm is invoked and uses the correct version, this happens automatically and silently.

To compile and run tests:

```
make
```

Or to do so manually:

```
go test ./...
```

To package and install gvm:

```
make build
./package/install-gvm.sh
```

And then follow the directions to add the gvm executables to your PATH. After that, whenever you invoke `go` in a directory with a "go.mod" file, gvm will automatically install a shim (if not already installed) and use it. If `go` is invoked in a directory with no "go.mod" file, it will fall back to the system installation.

[rbenv]: https://github.com/rbenv/rbenv#how-rbenv-hooks-into-your-shell

[build]: https://travis-ci.com/btmorr/gvm
[build-badge]: https://travis-ci.com/btmorr/gvm.svg?branch=edge
[license]: https://github.com/btmorr/gvm/LICENSE
[license-badge]: https://img.shields.io/github/license/btmorr/gvm.svg
