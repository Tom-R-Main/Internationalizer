package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	if err := execute(newRootCmd(), os.Args[1:]); err != nil {
		var reported reportedError
		if !errors.Is(err, errValidationFailed) && !errors.As(err, &reported) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
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
		newPseudoCmd(),
		newConfigCmd(),
	)
	rootCmd.AddCommand(newCommandsCmd(rootCmd))
	return rootCmd
}
