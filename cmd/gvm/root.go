package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "gvm",
	Short:   "Go Version Manager",
	Long:    `gvm is a Go Version Manager that helps you manage multiple Go versions.`,
	Version: version,
}

func Execute() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("gvm version %s (commit: %s, date: %s)\n", version, commit, date))
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
