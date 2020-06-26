package main

import (
	"fmt"

	"github.com/btmorr/gvm/internal/fetch"
)

// version determines selects the Go version specified by go.sum in the current
// working directory (looking up the most recent full semantic version, if not
// fully specified)
func main() {
	report := fetch.BuildReport(".")
	fmt.Println(report.Version)
}
