package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/review"
	"github.com/Tom-R-Main/Internationalizer/internal/state"
	"github.com/spf13/cobra"
)

func newReviewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "review",
		Short: "Inspect and approve exact translation artifacts",
	}
	cmd.AddCommand(newReviewListCmd(), newReviewApproveCmd())
	return cmd
}

func newReviewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tracked translations and review status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			manifest, err := state.Load(cfg.ManifestPath)
			if err != nil {
				return err
			}
			locale, _ := cmd.Flags().GetString("locale")
			bundle, _ := cmd.Flags().GetString("bundle")
			statusValue, _ := cmd.Flags().GetString("status")
			entries, err := review.List(manifest, review.Filter{Locale: locale, Bundle: bundle, Status: state.ReviewStatus(statusValue)})
			if err != nil {
				return err
			}
			asJSON, _ := cmd.Flags().GetBool("json")
			if asJSON {
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetIndent("", "  ")
				return encoder.Encode(entries)
			}
			for _, entry := range entries {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\t%s\t%s\n", entry.ReviewStatus, entry.Locale, entry.Bundle, entry.Key, entry.Origin); err != nil {
					return err
				}
			}
			return nil
		},
	}
	addReviewConfigFlag(cmd)
	cmd.Flags().StringP("locale", "l", "", "filter by target locale")
	cmd.Flags().String("bundle", "", "filter by bundle ID")
	cmd.Flags().String("status", "", "filter by review status (needs_review or approved)")
	cmd.Flags().Bool("json", false, "output entries as JSON")
	return cmd
}

func newReviewApproveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve",
		Short: "Approve current translations after validation",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			locale, _ := cmd.Flags().GetString("locale")
			bundle, _ := cmd.Flags().GetString("bundle")
			keys, _ := cmd.Flags().GetStringSlice("key")
			all, _ := cmd.Flags().GetBool("all")
			approved, err := review.Approve(cfg, review.Filter{Locale: locale, Bundle: bundle, Keys: keys, All: all}, time.Now())
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Approved %d translation(s) for %s.\n", len(approved), locale)
			return err
		},
	}
	addReviewConfigFlag(cmd)
	cmd.Flags().StringP("locale", "l", "", "target locale to approve")
	cmd.Flags().String("bundle", "", "bundle ID (required with --key)")
	cmd.Flags().StringSlice("key", nil, "individual translation key(s) to approve")
	cmd.Flags().Bool("all", false, "approve every matching current translation")
	_ = cmd.MarkFlagRequired("locale")
	cmd.MarkFlagsMutuallyExclusive("key", "all")
	return cmd
}

func addReviewConfigFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("config", "c", "", "path to config file (default: .internationalizer.yml)")
}
