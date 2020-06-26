package main

import (
	"fmt"
	"os"

	"github.com/btmorr/gvm/internal/fetch"
)

// version determines selects the Go version specified by go.sum in the current
// working directory (looking up the most recent full semantic version, if not
// fully specified)
func main() {
	defer func() {
		if a := recover(); a != nil {
			os.Exit(1)
		}
	}()
	report := fetch.BuildReport(".")
	fmt.Println(report.Version)
}
