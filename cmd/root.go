package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/the20100/trello-cli/internal/api"
	"github.com/the20100/trello-cli/internal/config"
)

var (
	// Persistent flags
	jsonFlag   bool
	prettyFlag bool

	// Global API client, set in PersistentPreRunE
	client *api.Client

	// Global config, set in PersistentPreRunE
	cfg *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "trello",
	Short: "Trello CLI — manage boards, lists, and cards via the Trello API",
	Long: `trello is a CLI tool for the Trello API.

It outputs JSON when piped (for agent use) and human-readable tables in a terminal.

Authentication requires a Trello API key and token.
Get yours at: https://trello.com/power-ups/admin

Token resolution order:
  1. TRELLO_API_KEY + TRELLO_API_TOKEN env vars
  2. Config file  (~/.config/trello/config.json  via: trello auth setup)

Examples:
  trello auth setup
  trello boards list
  trello boards get <id>
  trello lists list --board <board-id>
  trello cards list --board <board-id>
  trello cards create --list <list-id> --name "My task"
  trello cards get <id>
  trello members me
  trello search "my query"`,
	SilenceUsage: true,
}

// Execute is the entrypoint called from main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Force JSON output")
	rootCmd.PersistentFlags().BoolVar(&prettyFlag, "pretty", false, "Force pretty-printed JSON output (implies --json)")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if isAuthCommand(cmd) || cmd.Name() == "info" {
			return nil
		}

		apiKey, apiToken, err := resolveCredentials()
		if err != nil {
			return err
		}

		client = api.NewClient(apiKey, apiToken)
		return nil
	}

	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show tool info: config path, auth status, and environment",
	Run: func(cmd *cobra.Command, args []string) {
		printInfo()
	},
}

func printInfo() {
	fmt.Println("trello — Trello CLI")
	fmt.Println()

	exe, _ := os.Executable()
	fmt.Printf("  binary:  %s\n", exe)
	fmt.Printf("  os/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println()

	fmt.Println("  config paths by OS:")
	fmt.Println("    macOS:    ~/Library/Application Support/trello/config.json")
	fmt.Println("    Linux:    ~/.config/trello/config.json")
	fmt.Println("    Windows:  %AppData%\\trello\\config.json")
	fmt.Printf("  config:   %s\n", config.Path())
	fmt.Println()

	keySource := "(not set)"
	if k := os.Getenv("TRELLO_API_KEY"); k != "" {
		keySource = "TRELLO_API_KEY env var"
	} else if c, err := config.Load(); err == nil && c.APIKey != "" {
		keySource = "config file"
	}
	fmt.Printf("  key source:   %s\n", keySource)
	fmt.Println()
	fmt.Println("  env vars:")
	fmt.Printf("    TRELLO_API_KEY   = %s\n", maskOrEmpty(os.Getenv("TRELLO_API_KEY")))
	fmt.Printf("    TRELLO_API_TOKEN = %s\n", maskOrEmpty(os.Getenv("TRELLO_API_TOKEN")))
	fmt.Println()
	fmt.Println("  credential resolution order:")
	fmt.Println("    1. TRELLO_API_KEY + TRELLO_API_TOKEN env vars")
	fmt.Println("    2. config file  (trello auth setup)")
}

func maskOrEmpty(v string) string {
	if v == "" {
		return "(not set)"
	}
	if len(v) <= 8 {
		return "***"
	}
	return v[:4] + "..." + v[len(v)-4:]
}

// resolveCredentials returns the best available API key and token.
func resolveCredentials() (string, string, error) {
	// 1. Env vars
	envKey := os.Getenv("TRELLO_API_KEY")
	envToken := os.Getenv("TRELLO_API_TOKEN")
	if envKey != "" && envToken != "" {
		return envKey, envToken, nil
	}

	// 2. Config file
	var err error
	cfg, err = config.Load()
	if err != nil {
		return "", "", fmt.Errorf("failed to load config: %w", err)
	}
	if cfg.APIKey != "" && cfg.APIToken != "" {
		return cfg.APIKey, cfg.APIToken, nil
	}

	return "", "", fmt.Errorf("not authenticated — run: trello auth setup\nor set TRELLO_API_KEY and TRELLO_API_TOKEN env vars")
}

// isAuthCommand returns true if cmd is a child of the "auth" command.
func isAuthCommand(cmd *cobra.Command) bool {
	if cmd.Name() == "auth" {
		return true
	}
	p := cmd.Parent()
	for p != nil {
		if p.Name() == "auth" {
			return true
		}
		p = p.Parent()
	}
	return false
}
