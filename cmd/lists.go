package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vincentmaurin/trello-cli/internal/output"
)

var listsCmd = &cobra.Command{
	Use:   "lists",
	Short: "Manage Trello lists (columns)",
}

// ---- lists list ----

var (
	listsListBoardID string
	listsListFilter  string
)

var listsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List lists on a board",
	Long: `List all lists (columns) on a Trello board.

Examples:
  trello lists list --board <board-id>
  trello lists list --board <board-id> --filter all
  trello lists list --board <board-id> --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if listsListBoardID == "" {
			return fmt.Errorf("--board is required")
		}

		lists, err := client.GetBoardLists(listsListBoardID, listsListFilter)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(lists, output.IsPretty(cmd))
		}

		if len(lists) == 0 {
			fmt.Println("No lists found.")
			return nil
		}

		headers := []string{"ID", "NAME", "CLOSED"}
		rows := make([][]string, len(lists))
		for i, l := range lists {
			rows[i] = []string{
				l.ID,
				output.Truncate(l.Name, 50),
				output.FormatBool(l.Closed),
			}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

// ---- lists get ----

var listsGetCmd = &cobra.Command{
	Use:   "get <list-id>",
	Short: "Get details of a specific list",
	Long: `Get full details of a Trello list by its ID.

Examples:
  trello lists get abc123
  trello lists get abc123 --pretty`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		list, err := client.GetList(args[0])
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(list, output.IsPretty(cmd))
		}

		output.PrintKeyValue([][]string{
			{"ID", list.ID},
			{"Name", list.Name},
			{"Board", list.IDBoard},
			{"Closed", output.FormatBool(list.Closed)},
		})
		return nil
	},
}

// ---- lists create ----

var (
	listsCreateBoardID string
	listsCreatePos     string
)

var listsCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new list on a board",
	Long: `Create a new Trello list (column) on a board.

Examples:
  trello lists create "To Do" --board <board-id>
  trello lists create "To Do" --board <board-id> --pos top
  trello lists create "Done" --board <board-id> --pos bottom`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if listsCreateBoardID == "" {
			return fmt.Errorf("--board is required")
		}

		list, err := client.CreateList(args[0], listsCreateBoardID, listsCreatePos)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(list, output.IsPretty(cmd))
		}

		fmt.Printf("List created: %s\n", list.Name)
		fmt.Printf("ID:    %s\n", list.ID)
		fmt.Printf("Board: %s\n", list.IDBoard)
		return nil
	},
}

// ---- lists rename ----

var listsRenameCmd = &cobra.Command{
	Use:   "rename <list-id> <new-name>",
	Short: "Rename a list",
	Long: `Rename a Trello list.

Examples:
  trello lists rename abc123 "In Progress"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams("name", args[1])
		list, err := client.UpdateList(args[0], params)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(list, output.IsPretty(cmd))
		}

		fmt.Printf("List renamed to: %s\n", list.Name)
		return nil
	},
}

// ---- lists archive ----

var listsArchiveCmd = &cobra.Command{
	Use:   "archive <list-id>",
	Short: "Archive a list",
	Long: `Archive (close) a Trello list. Cards are preserved.

Examples:
  trello lists archive abc123`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		list, err := client.ArchiveList(args[0], true)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(list, output.IsPretty(cmd))
		}

		fmt.Printf("List archived: %s\n", list.Name)
		return nil
	},
}

// ---- lists unarchive ----

var listsUnarchiveCmd = &cobra.Command{
	Use:   "unarchive <list-id>",
	Short: "Unarchive a list",
	Long: `Unarchive (reopen) a Trello list.

Examples:
  trello lists unarchive abc123`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		list, err := client.ArchiveList(args[0], false)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(list, output.IsPretty(cmd))
		}

		fmt.Printf("List unarchived: %s\n", list.Name)
		return nil
	},
}

// ---- lists cards ----

var (
	listsCardsFilter string
)

var listsCardsCmd = &cobra.Command{
	Use:   "cards <list-id>",
	Short: "List cards in a list",
	Long: `List all cards in a Trello list.

Examples:
  trello lists cards abc123
  trello lists cards abc123 --filter all
  trello lists cards abc123 --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cards, err := client.GetListCards(args[0], listsCardsFilter)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(cards, output.IsPretty(cmd))
		}

		if len(cards) == 0 {
			fmt.Println("No cards found.")
			return nil
		}

		headers := []string{"ID", "#", "NAME", "DUE", "LABELS"}
		rows := make([][]string, len(cards))
		for i, c := range cards {
			labelNames := make([]string, len(c.Labels))
			for j, l := range c.Labels {
				if l.Name != "" {
					labelNames[j] = l.Name
				} else {
					labelNames[j] = l.Color
				}
			}
			rows[i] = []string{
				c.ID,
				fmt.Sprintf("%d", c.IDShort),
				output.Truncate(c.Name, 50),
				output.FormatDate(c.Due),
				output.FormatLabels(labelNames),
			}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

func init() {
	// lists list flags
	listsListCmd.Flags().StringVar(&listsListBoardID, "board", "", "Board ID (required)")
	listsListCmd.Flags().StringVar(&listsListFilter, "filter", "open", "Filter: open, closed, all")

	// lists create flags
	listsCreateCmd.Flags().StringVar(&listsCreateBoardID, "board", "", "Board ID (required)")
	listsCreateCmd.Flags().StringVar(&listsCreatePos, "pos", "", "Position: top, bottom, or a positive float")

	// lists cards flags
	listsCardsCmd.Flags().StringVar(&listsCardsFilter, "filter", "open", "Filter: open, closed, all")

	listsCmd.AddCommand(
		listsListCmd,
		listsGetCmd,
		listsCreateCmd,
		listsRenameCmd,
		listsArchiveCmd,
		listsUnarchiveCmd,
		listsCardsCmd,
	)
	rootCmd.AddCommand(listsCmd)
}
