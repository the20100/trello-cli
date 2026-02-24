package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vincentmaurin/trello-cli/internal/api"
	"github.com/vincentmaurin/trello-cli/internal/config"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage Trello authentication",
}

var authSetupCmd = &cobra.Command{
	Use:   "setup <api-key> <api-token>",
	Short: "Save Trello API credentials to the config file",
	Long: `Save your Trello API key and token to the local config file.

To get your credentials:
  1. Go to https://trello.com/power-ups/admin and create a Power-Up (or use an existing one)
  2. Generate an API key from the Power-Up settings
  3. Generate a token by visiting:
     https://trello.com/1/authorize?expiration=never&scope=read,write&response_type=token&key=YOUR_KEY

The credentials are stored at:
  macOS:   ~/Library/Application Support/trello/config.json
  Linux:   ~/.config/trello/config.json
  Windows: %AppData%\trello\config.json

You can also set env vars instead of using this command:
  export TRELLO_API_KEY=your-key
  export TRELLO_API_TOKEN=your-token`,
	Args:    cobra.ExactArgs(2),
	RunE:    runAuthSetup,
	Example: "  trello auth setup your_api_key your_api_token",
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current authentication status",
	RunE:  runAuthStatus,
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove saved credentials from the config file",
	RunE:  runAuthLogout,
}

func init() {
	authCmd.AddCommand(authSetupCmd, authStatusCmd, authLogoutCmd)
	rootCmd.AddCommand(authCmd)
}

func runAuthSetup(cmd *cobra.Command, args []string) error {
	apiKey := args[0]
	apiToken := args[1]

	if len(apiKey) < 8 {
		return fmt.Errorf("API key looks too short — check your key at https://trello.com/power-ups/admin")
	}
	if len(apiToken) < 8 {
		return fmt.Errorf("API token looks too short — re-generate it at https://trello.com/1/authorize?expiration=never&scope=read,write&response_type=token&key=%s", apiKey)
	}

	// Validate by fetching the authenticated member
	c := api.NewClient(apiKey, apiToken)
	member, err := c.GetMember("me", nil)
	if err != nil {
		return fmt.Errorf("credentials validation failed: %w", err)
	}

	cfg := &config.Config{
		APIKey:   apiKey,
		APIToken: apiToken,
		MemberID: member.ID,
		FullName: member.FullName,
		Username: member.Username,
	}
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Printf("Credentials saved to %s\n", config.Path())
	fmt.Printf("Authenticated as: %s (@%s)\n", member.FullName, member.Username)
	fmt.Printf("API key:          %s\n", maskOrEmpty(apiKey))
	fmt.Printf("API token:        %s\n", maskOrEmpty(apiToken))
	return nil
}

func runAuthStatus(cmd *cobra.Command, args []string) error {
	c, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	fmt.Printf("Config: %s\n", config.Path())
	fmt.Println()

	envKey := os.Getenv("TRELLO_API_KEY")
	envToken := os.Getenv("TRELLO_API_TOKEN")

	if envKey != "" && envToken != "" {
		fmt.Println("Credential source: env vars (take priority over config)")
		fmt.Printf("TRELLO_API_KEY:   %s\n", maskOrEmpty(envKey))
		fmt.Printf("TRELLO_API_TOKEN: %s\n", maskOrEmpty(envToken))
	} else if c.APIKey != "" && c.APIToken != "" {
		fmt.Println("Credential source: config file")
		fmt.Printf("API key:   %s\n", maskOrEmpty(c.APIKey))
		fmt.Printf("API token: %s\n", maskOrEmpty(c.APIToken))
		if c.FullName != "" {
			fmt.Printf("User:      %s (@%s)\n", c.FullName, c.Username)
		}
	} else {
		fmt.Println("Status: not authenticated")
		fmt.Println()
		fmt.Println("Run: trello auth setup <api-key> <api-token>")
		fmt.Println("Or set env vars:")
		fmt.Println("  export TRELLO_API_KEY=your-key")
		fmt.Println("  export TRELLO_API_TOKEN=your-token")
	}
	return nil
}

func runAuthLogout(cmd *cobra.Command, args []string) error {
	if err := config.Clear(); err != nil {
		return fmt.Errorf("removing config: %w", err)
	}
	fmt.Println("Credentials removed from config.")
	fmt.Println("Set TRELLO_API_KEY and TRELLO_API_TOKEN env vars if you still need access.")
	return nil
}
