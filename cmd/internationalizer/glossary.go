package main

import (
	"fmt"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/glossary"
	"github.com/spf13/cobra"
)

func newGlossaryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "glossary",
		Short: "Manage per-language glossary terms",
		Long:  "List, add, or remove glossary terms that are enforced during translation.",
	}

	list := &cobra.Command{
		Use:   "list",
		Short: "List glossary terms for a locale",
		RunE: func(cmd *cobra.Command, args []string) error {
			locale, _ := cmd.Flags().GetString("locale")
			dir, err := glossaryDir(cmd)
			if err != nil {
				return err
			}

			terms, err := glossary.Load(dir, locale)
			if err != nil {
				return err
			}
			if len(terms) == 0 {
				fmt.Printf("No glossary terms for %s.\n", locale)
				return nil
			}

			fmt.Printf("Glossary for %s (%d terms):\n\n", locale, len(terms))
			fmt.Println(glossary.FormatForPrompt(terms))
			return nil
		},
	}
	list.Flags().String("locale", "", "target locale (required)")
	list.MarkFlagRequired("locale")
	list.Flags().StringP("config", "c", "", "path to config file")

	add := &cobra.Command{
		Use:   "add",
		Short: "Add a glossary term",
		RunE: func(cmd *cobra.Command, args []string) error {
			locale, _ := cmd.Flags().GetString("locale")
			source, _ := cmd.Flags().GetString("source")
			target, _ := cmd.Flags().GetString("target")
			dir, err := glossaryDir(cmd)
			if err != nil {
				return err
			}

			if err := glossary.Add(dir, locale, source, target); err != nil {
				return err
			}
			fmt.Printf("Added: %s -> %s (%s)\n", source, target, locale)
			return nil
		},
	}
	add.Flags().String("locale", "", "target locale (required)")
	add.MarkFlagRequired("locale")
	add.Flags().String("source", "", "source term (required)")
	add.MarkFlagRequired("source")
	add.Flags().String("target", "", "translation (required)")
	add.MarkFlagRequired("target")
	add.Flags().StringP("config", "c", "", "path to config file")

	remove := &cobra.Command{
		Use:   "remove",
		Short: "Remove a glossary term",
		RunE: func(cmd *cobra.Command, args []string) error {
			locale, _ := cmd.Flags().GetString("locale")
			source, _ := cmd.Flags().GetString("source")
			dir, err := glossaryDir(cmd)
			if err != nil {
				return err
			}

			if err := glossary.Remove(dir, locale, source); err != nil {
				return err
			}
			fmt.Printf("Removed: %s (%s)\n", source, locale)
			return nil
		},
	}
	remove.Flags().String("locale", "", "target locale (required)")
	remove.MarkFlagRequired("locale")
	remove.Flags().String("source", "", "source term to remove (required)")
	remove.MarkFlagRequired("source")
	remove.Flags().StringP("config", "c", "", "path to config file")

	cmd.AddCommand(list, add, remove)
	return cmd
}

func glossaryDir(cmd *cobra.Command) (string, error) {
	cfgPath, _ := cmd.Flags().GetString("config")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		// Fall back to default directory if no config.
		return "glossary", nil
	}
	return cfg.GlossaryDir, nil
}
