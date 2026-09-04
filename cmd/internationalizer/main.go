package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:           "internationalizer",
		Short:         "AI-native i18n CLI tool",
		Long:          "Translate, validate, and manage internationalization files using LLMs.",
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	rootCmd.AddCommand(
		newTranslateCmd(),
		newDetectCmd(),
		newGlossaryCmd(),
		newTmCmd(),
		newValidateCmd(),
		newReviewCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		if !errors.Is(err, errValidationFailed) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
