// Package main is the treemux CLI entrypoint.
package main

import (
	"fmt"
	"os"

	"github.com/ian-howell/treemux/internal/cli"
)

// main runs the treemux CLI and exits on error.
func main() {
	if err := cli.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
