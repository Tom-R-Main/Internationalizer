package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:     "internationalizer",
		Short:   "AI-native i18n CLI tool",
		Long:    "Translate, validate, and manage internationalization files using LLMs.",
		Version: version,
	}

	rootCmd.AddCommand(
		newTranslateCmd(),
		newDetectCmd(),
		newGlossaryCmd(),
		newTmCmd(),
		newValidateCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
