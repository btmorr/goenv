package main

import (
	"fmt"

	"github.com/btmorr/gvm/internal/fetch"
)

// dry-run determines what steps will be taken when "gvm install" is invoked
// and prints a report
func main() {
	report := fetch.BuildReport(".")
	fmt.Printf(report.Archive)
}
