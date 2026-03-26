package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/validate"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate locale files against the source locale",
		Long:  "Check all target locales for missing keys, extra keys, and interpolation mismatches.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}

			reports, err := validate.Validate(cfg)
			if err != nil {
				return err
			}

			asJSON, _ := cmd.Flags().GetBool("json")
			quiet, _ := cmd.Flags().GetBool("quiet")

			if asJSON {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(reports)
			}

			if !quiet {
				fmt.Print(validate.FormatHuman(reports))
			}

			// Exit 1 if any locale has errors.
			for _, r := range reports {
				if len(r.Missing) > 0 || len(r.Mismatches) > 0 {
					os.Exit(1)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringP("config", "c", "", "path to config file (default: .internationalizer.yml)")
	cmd.Flags().Bool("json", false, "output report as JSON")
	cmd.Flags().BoolP("quiet", "q", false, "exit code only, no output")

	return cmd
}
