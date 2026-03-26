package main

import (
	"fmt"

	"github.com/Tom-R-Main/Internationalizer/internal/detect"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newDetectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "detect",
		Short: "Auto-detect project type and suggest config",
		Long:  "Scan the current directory to detect the i18n framework and suggest a configuration.",
		RunE: func(cmd *cobra.Command, args []string) error {
			d := detect.Detect(".")

			if d.Type == detect.Unknown {
				fmt.Println("Could not detect project type.")
				fmt.Println("Create a .internationalizer.yml config file manually.")
				return nil
			}

			fmt.Printf("Detected: %s (confidence: %.0f%%)\n\n", d.Type, d.Confidence*100)
			fmt.Println("Suggested configuration:")

			suggested := map[string]interface{}{
				"source_locale": d.SourceLocale,
				"source_path":   d.SourcePath,
			}
			if len(d.TargetLocales) > 0 {
				suggested["target_locales"] = d.TargetLocales
			}
			suggested["llm"] = map[string]string{
				"provider":    "gemini",
				"model":       "gemini-3.1-pro-preview",
				"api_key_env": "GOOGLE_AI_STUDIO_API_KEY",
			}

			out, _ := yaml.Marshal(suggested)
			fmt.Println(string(out))
			return nil
		},
	}
}
