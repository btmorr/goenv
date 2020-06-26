package main

import (
	"fmt"

	"github.com/btmorr/gvm/internal/fetch"
)

// fetch determines the current OS, architecture, and selects the Go version
// specified by go.sum in the current working directory (looking up the most
// recent full semantic version, if not fully specified), then fetches the
// corresponding archive from the Golang download server to a temporary file,
// and prints the path to that file.
func main() {
	report := fetch.BuildReport(".")
	tmpArchive := fetch.FetchArchive(report.Archive)
	fmt.Printf(tmpArchive)
}
