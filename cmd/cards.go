package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/vincentmaurin/trello-cli/internal/output"
)

var cardsCmd = &cobra.Command{
	Use:   "cards",
	Short: "Manage Trello cards",
}

// ---- cards list ----

var (
	cardsListBoardID string
	cardsListListID  string
	cardsListFilter  string
)

var cardsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List cards on a board or in a list",
	Long: `List Trello cards. Provide either --board or --list.

Examples:
  trello cards list --board <board-id>
  trello cards list --list <list-id>
  trello cards list --board <board-id> --filter all
  trello cards list --board <board-id> --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cardsListListID == "" && cardsListBoardID == "" {
			return fmt.Errorf("provide --board <board-id> or --list <list-id>")
		}

		if cardsListListID != "" {
			c, err := client.GetListCards(cardsListListID, cardsListFilter)
			if err != nil {
				return err
			}
			if output.IsJSON(cmd) {
				return output.PrintJSON(c, output.IsPretty(cmd))
			}
			printAPICardsTable(c)
			return nil
		}

		c, err := client.GetBoardCards(cardsListBoardID, cardsListFilter)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(c, output.IsPretty(cmd))
		}
		printAPICardsTable(c)
		return nil
	},
}

// ---- cards get ----

var cardsGetCmd = &cobra.Command{
	Use:   "get <card-id>",
	Short: "Get details of a specific card",
	Long: `Get full details of a Trello card by its ID or short link.

Examples:
  trello cards get abc123
  trello cards get abc123 --pretty`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		card, err := client.GetCard(args[0], nil)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(card, output.IsPretty(cmd))
		}

		labelNames := make([]string, len(card.Labels))
		for i, l := range card.Labels {
			if l.Name != "" {
				labelNames[i] = l.Name
			} else {
				labelNames[i] = l.Color
			}
		}

		checklistSummary := "-"
		if card.Badges.CheckItems > 0 {
			checklistSummary = fmt.Sprintf("%d/%d", card.Badges.CheckItemsChecked, card.Badges.CheckItems)
		}

		output.PrintKeyValue([][]string{
			{"ID", card.ID},
			{"#", fmt.Sprintf("%d", card.IDShort)},
			{"Name", card.Name},
			{"Description", output.Truncate(card.Desc, 80)},
			{"List", card.IDList},
			{"Board", card.IDBoard},
			{"URL", card.ShortURL},
			{"Due", output.FormatDate(card.Due)},
			{"Due complete", output.FormatBool(card.DueComplete)},
			{"Labels", output.FormatLabels(labelNames)},
			{"Checklists", checklistSummary},
			{"Attachments", fmt.Sprintf("%d", card.Badges.Attachments)},
			{"Comments", fmt.Sprintf("%d", card.Badges.Comments)},
			{"Last Activity", output.FormatTime(card.DateLastActivity)},
			{"Closed", output.FormatBool(card.Closed)},
		})
		return nil
	},
}

// ---- cards create ----

var (
	cardsCreateListID string
	cardsCreateDesc   string
	cardsCreateDue    string
	cardsCreateLabels string
	cardsCreatePos    string
)

var cardsCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new card",
	Long: `Create a new Trello card in a list.

Examples:
  trello cards create "Fix the bug" --list <list-id>
  trello cards create "Deploy v2" --list <list-id> --desc "Deploy new version"
  trello cards create "Review PR" --list <list-id> --due 2024-12-31
  trello cards create "Task" --list <list-id> --pos top`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cardsCreateListID == "" {
			return fmt.Errorf("--list is required")
		}

		extra := url.Values{}
		if cardsCreateDue != "" {
			extra.Set("due", cardsCreateDue)
		}
		if cardsCreatePos != "" {
			extra.Set("pos", cardsCreatePos)
		}

		card, err := client.CreateCard(cardsCreateListID, args[0], cardsCreateDesc, extra)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(card, output.IsPretty(cmd))
		}

		fmt.Printf("Card created: %s\n", card.Name)
		fmt.Printf("ID:  %s\n", card.ID)
		fmt.Printf("#%d  %s\n", card.IDShort, card.ShortURL)
		return nil
	},
}

// ---- cards update ----

var (
	cardsUpdateName    string
	cardsUpdateDesc    string
	cardsUpdateDue     string
	cardsUpdateClosed  bool
	cardsUpdateDueComplete bool
)

var cardsUpdateCmd = &cobra.Command{
	Use:   "update <card-id>",
	Short: "Update a card",
	Long: `Update a Trello card's name, description, due date, or state.

Examples:
  trello cards update abc123 --name "New title"
  trello cards update abc123 --desc "Updated description"
  trello cards update abc123 --due 2024-12-31
  trello cards update abc123 --due-complete
  trello cards update abc123 --closed`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"name", cardsUpdateName,
			"desc", cardsUpdateDesc,
			"due", cardsUpdateDue,
		)
		if cmd.Flags().Changed("closed") {
			params.Set("closed", output.FormatBool(cardsUpdateClosed))
		}
		if cmd.Flags().Changed("due-complete") {
			params.Set("dueComplete", output.FormatBool(cardsUpdateDueComplete))
		}

		card, err := client.UpdateCard(args[0], params)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(card, output.IsPretty(cmd))
		}

		fmt.Printf("Card updated: %s\n", card.Name)
		fmt.Printf("ID:  %s\n", card.ID)
		return nil
	},
}

// ---- cards delete ----

var cardsDeleteCmd = &cobra.Command{
	Use:   "delete <card-id>",
	Short: "Delete a card",
	Long: `Permanently delete a Trello card.

This action cannot be undone.

Examples:
  trello cards delete abc123`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteCard(args[0]); err != nil {
			return err
		}
		fmt.Printf("Card %s deleted.\n", args[0])
		return nil
	},
}

// ---- cards move ----

var (
	cardsMoveListID  string
	cardsMoveBoard   string
)

var cardsMoveCmd = &cobra.Command{
	Use:   "move <card-id>",
	Short: "Move a card to a different list",
	Long: `Move a Trello card to a different list (and optionally a different board).

Examples:
  trello cards move abc123 --list <list-id>
  trello cards move abc123 --list <list-id> --board <board-id>`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cardsMoveListID == "" {
			return fmt.Errorf("--list is required")
		}

		card, err := client.MoveCard(args[0], cardsMoveListID, cardsMoveBoard)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(card, output.IsPretty(cmd))
		}

		fmt.Printf("Card moved: %s\n", card.Name)
		fmt.Printf("New list: %s\n", card.IDList)
		return nil
	},
}

// ---- cards archive ----

var cardsArchiveCmd = &cobra.Command{
	Use:   "archive <card-id>",
	Short: "Archive a card",
	Long: `Archive (close) a Trello card.

Examples:
  trello cards archive abc123`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams("closed", "true")
		card, err := client.UpdateCard(args[0], params)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(card, output.IsPretty(cmd))
		}

		fmt.Printf("Card archived: %s\n", card.Name)
		return nil
	},
}

// ---- cards comment ----

var cardsCommentCmd = &cobra.Command{
	Use:   "comment <card-id> <text>",
	Short: "Add a comment to a card",
	Long: `Add a comment to a Trello card.

Examples:
  trello cards comment abc123 "This is a comment"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		action, err := client.AddComment(args[0], args[1])
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(action, output.IsPretty(cmd))
		}

		fmt.Printf("Comment added to card %s.\n", args[0])
		fmt.Printf("Action ID: %s\n", action.ID)
		return nil
	},
}

// ---- cards checklists ----

var cardsChecklistsCmd = &cobra.Command{
	Use:   "checklists <card-id>",
	Short: "List checklists on a card",
	Long: `List all checklists on a Trello card, with their check items.

Examples:
  trello cards checklists abc123
  trello cards checklists abc123 --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		checklists, err := client.GetCardChecklists(args[0])
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(checklists, output.IsPretty(cmd))
		}

		if len(checklists) == 0 {
			fmt.Println("No checklists found.")
			return nil
		}

		for _, cl := range checklists {
			fmt.Printf("\n%s (ID: %s)\n", cl.Name, cl.ID)
			if len(cl.CheckItems) == 0 {
				fmt.Println("  (empty)")
				continue
			}
			for _, item := range cl.CheckItems {
				mark := "[ ]"
				if item.State == "complete" {
					mark = "[x]"
				}
				fmt.Printf("  %s %s  (ID: %s)\n", mark, item.Name, item.ID)
			}
		}
		return nil
	},
}

// ---- cards attachments ----

var cardsAttachmentsCmd = &cobra.Command{
	Use:   "attachments <card-id>",
	Short: "List attachments on a card",
	Long: `List all attachments on a Trello card.

Examples:
  trello cards attachments abc123
  trello cards attachments abc123 --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		attachments, err := client.GetCardAttachments(args[0])
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(attachments, output.IsPretty(cmd))
		}

		if len(attachments) == 0 {
			fmt.Println("No attachments found.")
			return nil
		}

		headers := []string{"ID", "NAME", "URL", "DATE"}
		rows := make([][]string, len(attachments))
		for i, a := range attachments {
			rows[i] = []string{
				a.ID,
				output.Truncate(a.Name, 30),
				output.Truncate(a.URL, 50),
				output.FormatTime(a.Date),
			}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

// ---- cards label ----

var (
	cardsLabelAdd    string
	cardsLabelRemove string
)

var cardsLabelCmd = &cobra.Command{
	Use:   "label <card-id>",
	Short: "Add or remove labels on a card",
	Long: `Add or remove labels on a Trello card.

Examples:
  trello cards label abc123 --add <label-id>
  trello cards label abc123 --remove <label-id>`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cardsLabelAdd == "" && cardsLabelRemove == "" {
			return fmt.Errorf("provide --add <label-id> or --remove <label-id>")
		}

		if cardsLabelAdd != "" {
			if err := client.AddLabelToCard(args[0], cardsLabelAdd); err != nil {
				return err
			}
			fmt.Printf("Label %s added to card %s.\n", cardsLabelAdd, args[0])
		}

		if cardsLabelRemove != "" {
			if err := client.RemoveLabelFromCard(args[0], cardsLabelRemove); err != nil {
				return err
			}
			fmt.Printf("Label %s removed from card %s.\n", cardsLabelRemove, args[0])
		}

		return nil
	},
}

// ---- cards member ----

var (
	cardsMemberAdd    string
	cardsMemberRemove string
)

var cardsMemberCmd = &cobra.Command{
	Use:   "member <card-id>",
	Short: "Add or remove members from a card",
	Long: `Assign or unassign members from a Trello card.

Examples:
  trello cards member abc123 --add <member-id>
  trello cards member abc123 --remove <member-id>`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cardsMemberAdd == "" && cardsMemberRemove == "" {
			return fmt.Errorf("provide --add <member-id> or --remove <member-id>")
		}

		if cardsMemberAdd != "" {
			if err := client.AddMemberToCard(args[0], cardsMemberAdd); err != nil {
				return err
			}
			fmt.Printf("Member %s added to card %s.\n", cardsMemberAdd, args[0])
		}

		if cardsMemberRemove != "" {
			if err := client.RemoveMemberFromCard(args[0], cardsMemberRemove); err != nil {
				return err
			}
			fmt.Printf("Member %s removed from card %s.\n", cardsMemberRemove, args[0])
		}

		return nil
	},
}

func init() {
	// cards list flags
	cardsListCmd.Flags().StringVar(&cardsListBoardID, "board", "", "Board ID")
	cardsListCmd.Flags().StringVar(&cardsListListID, "list", "", "List ID")
	cardsListCmd.Flags().StringVar(&cardsListFilter, "filter", "open", "Filter: open, closed, all, visible")

	// cards create flags
	cardsCreateCmd.Flags().StringVar(&cardsCreateListID, "list", "", "List ID (required)")
	cardsCreateCmd.Flags().StringVar(&cardsCreateDesc, "desc", "", "Card description")
	cardsCreateCmd.Flags().StringVar(&cardsCreateDue, "due", "", "Due date (ISO-8601, e.g. 2024-12-31)")
	cardsCreateCmd.Flags().StringVar(&cardsCreatePos, "pos", "", "Position: top, bottom, or a positive float")
	cardsCreateCmd.Flags().StringVar(&cardsCreateLabels, "labels", "", "Comma-separated label IDs to add")

	// cards update flags
	cardsUpdateCmd.Flags().StringVar(&cardsUpdateName, "name", "", "New card name")
	cardsUpdateCmd.Flags().StringVar(&cardsUpdateDesc, "desc", "", "New card description")
	cardsUpdateCmd.Flags().StringVar(&cardsUpdateDue, "due", "", "Due date (ISO-8601)")
	cardsUpdateCmd.Flags().BoolVar(&cardsUpdateClosed, "closed", false, "Archive the card")
	cardsUpdateCmd.Flags().BoolVar(&cardsUpdateDueComplete, "due-complete", false, "Mark due date as complete")

	// cards move flags
	cardsMoveCmd.Flags().StringVar(&cardsMoveListID, "list", "", "Target list ID (required)")
	cardsMoveCmd.Flags().StringVar(&cardsMoveBoard, "board", "", "Target board ID (optional, for cross-board moves)")

	// cards label flags
	cardsLabelCmd.Flags().StringVar(&cardsLabelAdd, "add", "", "Label ID to add")
	cardsLabelCmd.Flags().StringVar(&cardsLabelRemove, "remove", "", "Label ID to remove")

	// cards member flags
	cardsMemberCmd.Flags().StringVar(&cardsMemberAdd, "add", "", "Member ID to add")
	cardsMemberCmd.Flags().StringVar(&cardsMemberRemove, "remove", "", "Member ID to remove")

	cardsCmd.AddCommand(
		cardsListCmd,
		cardsGetCmd,
		cardsCreateCmd,
		cardsUpdateCmd,
		cardsDeleteCmd,
		cardsMoveCmd,
		cardsArchiveCmd,
		cardsCommentCmd,
		cardsChecklistsCmd,
		cardsAttachmentsCmd,
		cardsLabelCmd,
		cardsMemberCmd,
	)
	rootCmd.AddCommand(cardsCmd)
}
