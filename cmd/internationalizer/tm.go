package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/tm"
	"github.com/spf13/cobra"
)

func newTmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tm",
		Short: "Manage translation memory",
		Long:  "View stats, export, or clear the translation memory cache.",
	}

	stats := &cobra.Command{
		Use:   "stats",
		Short: "Show translation memory statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := tmPath(cmd)
			if err != nil {
				return err
			}
			memory, err := tm.Load(path)
			if err != nil {
				return err
			}
			s := memory.Stats()
			fmt.Printf("Translation Memory: %s\n", path)
			fmt.Printf("Total records: %d\n", s.TotalRecords)
			fmt.Printf("File size: %d bytes\n", s.FileSize)
			if len(s.ByLocale) > 0 {
				fmt.Println("\nBy locale:")
				for locale, count := range s.ByLocale {
					fmt.Printf("  %s: %d\n", locale, count)
				}
			}
			return nil
		},
	}
	stats.Flags().StringP("config", "c", "", "path to config file")

	clear := &cobra.Command{
		Use:   "clear",
		Short: "Clear the translation memory",
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("force")
			path, err := tmPath(cmd)
			if err != nil {
				return err
			}
			memory, err := tm.Load(path)
			if err != nil {
				return err
			}

			if !force {
				s := memory.Stats()
				fmt.Printf("This will delete %d records from %s.\n", s.TotalRecords, path)
				fmt.Print("Continue? [y/N] ")
				reader := bufio.NewReader(os.Stdin)
				answer, _ := reader.ReadString('\n')
				if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(answer)), "y") {
					fmt.Println("Aborted.")
					return nil
				}
			}

			if err := memory.Clear(); err != nil {
				return err
			}
			fmt.Println("Translation memory cleared.")
			return nil
		},
	}
	clear.Flags().StringP("config", "c", "", "path to config file")
	clear.Flags().Bool("force", false, "skip confirmation prompt")

	export := &cobra.Command{
		Use:   "export",
		Short: "Export translation memory as JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := tmPath(cmd)
			if err != nil {
				return err
			}
			memory, err := tm.Load(path)
			if err != nil {
				return err
			}
			return memory.Export(os.Stdout)
		},
	}
	export.Flags().StringP("config", "c", "", "path to config file")

	cmd.AddCommand(stats, clear, export)
	return cmd
}

func tmPath(cmd *cobra.Command) (string, error) {
	cfgPath, _ := cmd.Flags().GetString("config")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return ".internationalizer/tm.jsonl", nil
	}
	return cfg.TMPath, nil
}
