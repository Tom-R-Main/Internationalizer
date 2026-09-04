package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type commandContract struct {
	Argv         []string       `json:"argv"`
	Description  string         `json:"description"`
	Arguments    []flagContract `json:"arguments"`
	SideEffects  []string       `json:"side_effects"`
	Network      string         `json:"network"`
	InputSchema  map[string]any `json:"input_schema,omitempty"`
	OutputSchema map[string]any `json:"output_schema,omitempty"`
	OutputFormat string         `json:"output_format"`
	Next         [][]string     `json:"next,omitempty"`
}
type flagContract struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Default     string `json:"default"`
	Description string `json:"description"`
}

func newCommandsCmd(root *cobra.Command) *cobra.Command {
	var asJSON bool
	var selected string
	var limit int
	cmd := &cobra.Command{Use: "commands", Short: "Describe installed workflow entry points, arguments, and effects", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		if limit < 0 {
			return codedError("invalid_arguments", fmt.Errorf("limit must be nonnegative"))
		}
		contracts := []commandContract{}
		var visit func(*cobra.Command)
		visit = func(c *cobra.Command) {
			if c.Hidden {
				return
			}
			if c.RunE != nil || c.Run != nil {
				path := c.CommandPath()
				item := commandContract{Argv: strings.Fields(path), Description: c.Short, Arguments: []flagContract{}, SideEffects: []string{}, Network: "none", OutputFormat: "human-readable"}
				c.Flags().VisitAll(func(f *pflag.Flag) {
					if !f.Hidden {
						item.Arguments = append(item.Arguments, flagContract{Name: "--" + f.Name, Type: f.Value.Type(), Default: f.DefValue, Description: f.Usage})
					}
				})
				switch strings.TrimPrefix(path, root.Name()+" ") {
				case "detect":
					item.Next = [][]string{{root.Name(), "config", "check", "--json"}, {root.Name(), "config", "plan", "--help"}}
				case "config check":
					item.Next = [][]string{{root.Name(), "config", "plan", "--help"}, {root.Name(), "translate", "--dry-run", "--json"}}
				case "config plan":
					item.SideEffects = []string{"writes a new plan file only when --out is supplied"}
					item.Next = [][]string{{root.Name(), "config", "apply", "--help"}}
				case "config apply":
					item.SideEffects = []string{"writes only the selected plan's configuration file after integrity and drift checks"}
					item.Next = [][]string{{root.Name(), "config", "check", "--json"}, {root.Name(), "translate", "--dry-run", "--json"}}
				case "translate":
					item.Network = "provider calls except --dry-run and --adopt-existing"
					item.SideEffects = []string{"may write catalogs, translation memory, and state; --dry-run writes nothing"}
				case "validate":
					item.Next = [][]string{{root.Name(), "review", "--help"}}
				case "commands":
				default:
					item.SideEffects = []string{"mode-dependent; inspect command help before execution"}
					item.Network = "mode-dependent; inspect command help"
				}
				if c.Flags().Lookup("json") != nil {
					switch strings.TrimPrefix(path, root.Name()+" ") {
					case "detect", "commands", "config check", "config plan", "config apply", "translate", "validate":
						item.OutputFormat = "internationalizer.cli.v1 (--json)"
						item.InputSchema = workflowInputSchema(c)
						item.OutputSchema = workflowOutputSchema(path)
					default:
						item.OutputFormat = "command-specific JSON (--json)"
					}
				}
				contracts = append(contracts, item)
			}
			for _, child := range c.Commands() {
				visit(child)
			}
		}
		visit(root)
		total := len(contracts)
		if selected != "" {
			filtered := contracts[:0]
			for _, item := range contracts {
				if strings.Join(item.Argv[1:], " ") == selected {
					filtered = append(filtered, item)
				}
			}
			contracts = filtered
			if len(contracts) == 0 {
				return codedError("invalid_arguments", fmt.Errorf("unknown command selection %q", selected))
			}
		}
		matched := len(contracts)
		if limit > 0 && len(contracts) > limit {
			contracts = contracts[:limit]
		}
		if asJSON {
			return emitJSON(cmd, "ok", map[string]any{"cli_version": version, "commands": contracts, "total": total, "matched": matched, "truncated": len(contracts) < matched, "exit_codes": map[string]string{"0": "command completed; inspect status and diagnostics for unresolved decisions", "1": "command failed or validation blocked"}, "states": []string{"configuration_checked", "planned", "applied", "generated", "structurally_valid", "human_approved"}, "authorization": "config apply is an explicit mutation request; --no-input disables prompts only"}, nil)
		}
		for _, c := range contracts {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s — %s\n", strings.Join(c.Argv, " "), c.Description); err != nil {
				return err
			}
		}
		return nil
	}}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Emit versioned machine-readable command contracts")
	cmd.Flags().StringVar(&selected, "command", "", "Return one command contract, for example 'config plan'")
	cmd.Flags().IntVar(&limit, "limit", 50, "Maximum command contracts (0 for all)")
	return cmd
}
