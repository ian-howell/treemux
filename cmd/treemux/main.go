// Package main is the treemux CLI entrypoint.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ian-howell/treemux/internal/cli"
)

// main runs the treemux CLI and exits on error.
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// run executes the treemux CLI and returns any errors encountered.
func run() error {
	var configFilePatch string
	flag.StringVar(&configFilePatch, "config-file", "", "Path to a treemux configuration file.")
	flag.Parse()

	var config cli.Config
	if configFilePatch != "" {
		var err error
		config, err = cli.LoadConfig(configFilePatch)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
	}

	if err := cli.Run(config); err != nil {
		return fmt.Errorf("running treemux: %w", err)
	}

	return nil
}
