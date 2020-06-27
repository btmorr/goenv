# goenv

[![Build Status][build-badge]][build]
[![License][license-badge]][license]

A version manager for Go

The goal of this project is to make it possible to have different project in a system using different versions of Go and automatically (as much as possible) use the version specified by each project's "go.mod" file. Inspiration is drawn heavily from [rbenv].

When a user invokes `go` in a directory with no "go.mod" file, goenv will use the system installation of `go` (primarily when running something like `go mod init` or quickly doing `go run myfile.go` to test something out). When a user invokes `go` in a directory that does have a "go.mod" file, goenv will check the go version specified in "go.mod", install the corresponding version in a sandboxed environment if it is missing, and then use that version. When installing a new version for the first time, it will print a message like "goenv installing shim for go 1.14.4". Otherwise, it does not write anything to stdout/stderr--no need to create additional files or use commands to ensure that goenv is invoked and uses the correct version, this happens automatically.

## Build, install, run

To compile and run tests:

```
make
```

Or to do so manually:

```
go test ./...
```

To package and install goenv:

```
make build
./package/install-goenv.sh
```

And then follow the directions to add the goenv executables to your PATH. After that, whenever you invoke `go` in a directory with a "go.mod" file, goenv will automatically ensure the correct version is used.

## How it works

On install, goenv creates a hidden directory under the user's home directory (~/.goenv), and copies binaries and other operating files under that directory. It also attempts to detect the current shell and modify the user's shell profile configuration to add itself to the PATH (or the user can do this manually). Once goenv's bin directory is in the PATH, it will handle choosing which version of Go to use (note that for this to work, goenv's bin directory must come earlier in the PATH than any other directory containing a `go` executable).

When the user invokes `go`, this is first handled by ~/.goenv/bin/go, which checks for a "go.mod" file in the current directory. If it does not find one, the command is routed to the system installation of `go` (the path to which was saved when goenv was first installed). If it does find one, it checks the Go version specified in the file, checks if it is a fully-qualified semantic version and if not then selects the most recent version (e.g.: "1.12" turns into "1.12.17"), then checks if it already has a shim for that version available.

A shim is a copy of a particular version of Go for the current OS and architecture, installed under "~/.goenv/shims/<version>". If one does not already exist (such as if a project specifies "1.14", and the current most recent version is "1.14.4" and "~/.goenv/shims/1.14.4/" doesn't exist), then goenv downloads the appropriate archive from the official [Golang download page] and unzips it to "~/.goenv/shims/<version>"

~/.goenv/bin/go invokes the go binary in that shim, passing all other supplied arguments along unmodified, so any command that would work in the specified version is guaranteed to work as expected.

## Prior art

This project has taken a lot of inspiration from [rbenv], though there are some differences when working with Go as opposed to Ruby (most notably, Go projects already have a file specifying the Go version, so there is no need for an extra env/version file, and the Go toolchain is focused around a single executable--`go`--instead of several). This implementation is also starting from scratch, so not all techniques that are used in rbenv are employed here (such as rehashing).

Other full-featured analogous tools already exist for Go (notably [syndbg/goenv], which branched off of [pyenv]--itself forked from rbenv--and modified it to work for Go). The aim of this project is to develop something hyper-minimal, and to learn about how to manage these toolchains. There is definitely a lot of useful learning in that project (notable example: [this issue about security vulnerabilities] relating to PATH configuration). In relation to that issue, it turns out the Golang pkg installer for MacOS adds `go` to the PATH by adding "/etc/path.d/go" containing a path string, which then gets picked up by "/usr/libexec/path_helper" and added to the middle of the PATH (itself called by /etc/profile which is invoked on login). This means that it is not possible to safely install goenv on MacOS if Go was installed using the pkg installer. Recommended fix is to remove that installation and install Go by following the official instructions for [installing a Go tarball]. Alternately, leave the installation alone but delete "/etc/paths.d/go".

tl;dr: managing PATH on MacOS is particularly weird--make sure ~/.goenv/bin is earlier than any other path containing `go`, but try to put goenv and all others as late as possible in PATH to prevent security risks.

There are a couple of explanations of path_helper, notably [this](http://www.softec.lu/site/DevelopersCorner/MasteringThePathHelper) and [this](http://hea-www.harvard.edu/~fine/OSX/path_helper.html)

## Contributions

Contributions are welcome! Feel free to open an issue to report a bug, request a feature, or ask a question. Also feel free to open a PR to provide a fix or feature, or update documentation.

[rbenv]: https://github.com/rbenv/rbenv#how-rbenv-hooks-into-your-shell
[pyenv]: https://github.com/pyenv/pyenv
[syndbg/goenv]: https://github.com/syndbg/goenv
[this issue about security vulnerabilities]: https://github.com/syndbg/goenv/issues/99
[Golang download page]: https://golang.org/dl/
[installing a Go tarball]: https://golang.org/doc/install#tarball

[build]: https://travis-ci.com/btmorr/goenv
[build-badge]: https://travis-ci.com/btmorr/goenv.svg?branch=edge
[license]: https://github.com/btmorr/goenv/LICENSE
[license-badge]: https://img.shields.io/github/license/btmorr/goenv.svg
