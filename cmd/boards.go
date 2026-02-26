package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/the20100/trello-cli/internal/output"
)

var boardsCmd = &cobra.Command{
	Use:   "boards",
	Short: "Manage Trello boards",
}

// ---- boards list ----

var boardsListFilter string

var boardsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List boards for the authenticated member",
	Long: `List all Trello boards for the authenticated member.

Examples:
  trello boards list
  trello boards list --filter open
  trello boards list --filter all
  trello boards list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		boards, err := client.GetMyBoards(boardsListFilter)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(boards, output.IsPretty(cmd))
		}

		if len(boards) == 0 {
			fmt.Println("No boards found.")
			return nil
		}

		headers := []string{"ID", "NAME", "WORKSPACE", "LAST ACTIVITY", "CLOSED"}
		rows := make([][]string, len(boards))
		for i, b := range boards {
			rows[i] = []string{
				b.ID,
				output.Truncate(b.Name, 40),
				output.Truncate(b.IDOrganization, 24),
				output.FormatTime(b.DateLastActivity),
				output.FormatBool(b.Closed),
			}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

// ---- boards get ----

var boardsGetCmd = &cobra.Command{
	Use:   "get <board-id>",
	Short: "Get details of a specific board",
	Long: `Get full details of a Trello board by its ID or short link.

Examples:
  trello boards get abc123
  trello boards get abc123 --pretty`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		board, err := client.GetBoard(args[0], nil)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(board, output.IsPretty(cmd))
		}

		output.PrintKeyValue([][]string{
			{"ID", board.ID},
			{"Name", board.Name},
			{"Description", output.Truncate(board.Desc, 80)},
			{"Workspace", board.IDOrganization},
			{"URL", board.ShortURL},
			{"Last Activity", output.FormatTime(board.DateLastActivity)},
			{"Closed", output.FormatBool(board.Closed)},
			{"Permission", board.Prefs.PermissionLevel},
		})
		return nil
	},
}

// ---- boards create ----

var (
	boardsCreateDesc string
	boardsCreateOrg  string
	boardsCreatePriv string
)

var boardsCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new board",
	Long: `Create a new Trello board.

Examples:
  trello boards create "My Project"
  trello boards create "My Project" --desc "Project description"
  trello boards create "My Project" --workspace abc123
  trello boards create "My Project" --privacy private`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		board, err := client.CreateBoard(args[0], boardsCreateDesc, boardsCreateOrg, nil)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(board, output.IsPretty(cmd))
		}

		fmt.Printf("Board created: %s\n", board.Name)
		fmt.Printf("ID:  %s\n", board.ID)
		fmt.Printf("URL: %s\n", board.ShortURL)
		return nil
	},
}

// ---- boards update ----

var (
	boardsUpdateName   string
	boardsUpdateDesc   string
	boardsUpdateClosed bool
)

var boardsUpdateCmd = &cobra.Command{
	Use:   "update <board-id>",
	Short: "Update a board",
	Long: `Update a Trello board's name, description, or state.

Examples:
  trello boards update abc123 --name "New Name"
  trello boards update abc123 --desc "Updated description"
  trello boards update abc123 --closed`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"name", boardsUpdateName,
			"desc", boardsUpdateDesc,
		)
		if cmd.Flags().Changed("closed") {
			params.Set("closed", output.FormatBool(boardsUpdateClosed))
		}

		board, err := client.UpdateBoard(args[0], params)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(board, output.IsPretty(cmd))
		}

		fmt.Printf("Board updated: %s\n", board.Name)
		fmt.Printf("ID:  %s\n", board.ID)
		fmt.Printf("URL: %s\n", board.ShortURL)
		return nil
	},
}

// ---- boards delete ----

var boardsDeleteCmd = &cobra.Command{
	Use:   "delete <board-id>",
	Short: "Delete a board",
	Long: `Permanently delete a Trello board.

This action cannot be undone. The board and all its cards will be removed.

Examples:
  trello boards delete abc123`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteBoard(args[0]); err != nil {
			return err
		}
		fmt.Printf("Board %s deleted.\n", args[0])
		return nil
	},
}

// ---- boards members ----

var boardsMembersCmd = &cobra.Command{
	Use:   "members <board-id>",
	Short: "List members of a board",
	Long: `List all members of a Trello board.

Examples:
  trello boards members abc123
  trello boards members abc123 --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		members, err := client.GetBoardMembers(args[0])
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(members, output.IsPretty(cmd))
		}

		if len(members) == 0 {
			fmt.Println("No members found.")
			return nil
		}

		headers := []string{"ID", "NAME", "USERNAME"}
		rows := make([][]string, len(members))
		for i, m := range members {
			rows[i] = []string{m.ID, m.FullName, m.Username}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

// ---- boards labels ----

var boardsLabelsCmd = &cobra.Command{
	Use:   "labels <board-id>",
	Short: "List labels on a board",
	Long: `List all labels defined on a Trello board.

Examples:
  trello boards labels abc123
  trello boards labels abc123 --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		labels, err := client.GetBoardLabels(args[0])
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(labels, output.IsPretty(cmd))
		}

		if len(labels) == 0 {
			fmt.Println("No labels found.")
			return nil
		}

		headers := []string{"ID", "NAME", "COLOR"}
		rows := make([][]string, len(labels))
		for i, l := range labels {
			name := l.Name
			if name == "" {
				name = "-"
			}
			rows[i] = []string{l.ID, name, l.Color}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

func init() {
	// boards list flags
	boardsListCmd.Flags().StringVar(&boardsListFilter, "filter", "open", "Filter boards: open, closed, all, members, organization, public, starred")

	// boards create flags
	boardsCreateCmd.Flags().StringVar(&boardsCreateDesc, "desc", "", "Board description")
	boardsCreateCmd.Flags().StringVar(&boardsCreateOrg, "workspace", "", "Workspace/organization ID to create the board in")
	boardsCreateCmd.Flags().StringVar(&boardsCreatePriv, "privacy", "private", "Privacy level: private, public, org")

	// boards update flags
	boardsUpdateCmd.Flags().StringVar(&boardsUpdateName, "name", "", "New board name")
	boardsUpdateCmd.Flags().StringVar(&boardsUpdateDesc, "desc", "", "New board description")
	boardsUpdateCmd.Flags().BoolVar(&boardsUpdateClosed, "closed", false, "Archive the board")

	boardsCmd.AddCommand(
		boardsListCmd,
		boardsGetCmd,
		boardsCreateCmd,
		boardsUpdateCmd,
		boardsDeleteCmd,
		boardsMembersCmd,
		boardsLabelsCmd,
	)
	rootCmd.AddCommand(boardsCmd)
}
