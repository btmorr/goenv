# gvm

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
