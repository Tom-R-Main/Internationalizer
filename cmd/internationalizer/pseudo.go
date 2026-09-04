package main

import (
	"fmt"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/pseudo"
	"github.com/spf13/cobra"
)

func newPseudoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pseudo",
		Short: "Generate deterministic accented or bidi pseudolocales",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			strategyValue, _ := cmd.Flags().GetString("strategy")
			locale, _ := cmd.Flags().GetString("locale")
			force, _ := cmd.Flags().GetBool("force")
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			results, err := pseudo.Generate(cfg, pseudo.GenerateOptions{
				Strategy: pseudo.Strategy(strategyValue),
				Locale:   locale,
				Force:    force,
				DryRun:   dryRun,
			})
			if err != nil {
				return err
			}
			for _, result := range results {
				action := "generated"
				if dryRun {
					action = "would generate"
				}
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s %s/%s: %d units -> %s\n", action, result.Bundle, result.Locale, result.Units, result.TargetPath); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringP("config", "c", "", "path to config file (default: .internationalizer.yml)")
	cmd.Flags().String("strategy", string(pseudo.Accented), "pseudo strategy (accented or bidi)")
	cmd.Flags().StringP("locale", "l", "", "output locale (defaults to en-XA or ar-XB)")
	cmd.Flags().Bool("force", false, "overwrite an existing artifact not owned by the pseudo generator")
	cmd.Flags().Bool("dry-run", false, "show outputs without writing files or manifest state")
	return cmd
}
