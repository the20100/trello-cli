package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vincentmaurin/trello-cli/internal/output"
)

var checklistsCmd = &cobra.Command{
	Use:   "checklists",
	Short: "Manage card checklists",
}

// ---- checklists create ----

var (
	checklistsCreateCardID string
)

var checklistsCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a checklist on a card",
	Long: `Create a new checklist on a Trello card.

Examples:
  trello checklists create "Acceptance Criteria" --card <card-id>`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if checklistsCreateCardID == "" {
			return fmt.Errorf("--card is required")
		}

		cl, err := client.CreateChecklist(checklistsCreateCardID, args[0])
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(cl, output.IsPretty(cmd))
		}

		fmt.Printf("Checklist created: %s\n", cl.Name)
		fmt.Printf("ID:   %s\n", cl.ID)
		fmt.Printf("Card: %s\n", cl.IDCard)
		return nil
	},
}

// ---- checklists delete ----

var checklistsDeleteCmd = &cobra.Command{
	Use:   "delete <checklist-id>",
	Short: "Delete a checklist",
	Long: `Delete a Trello checklist and all its items.

Examples:
  trello checklists delete abc123`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteChecklist(args[0]); err != nil {
			return err
		}
		fmt.Printf("Checklist %s deleted.\n", args[0])
		return nil
	},
}

// ---- checklists add-item ----

var (
	checklistsAddItemChecklist string
)

var checklistsAddItemCmd = &cobra.Command{
	Use:   "add-item <name>",
	Short: "Add an item to a checklist",
	Long: `Add a new item to a Trello checklist.

Examples:
  trello checklists add-item "Write tests" --checklist <checklist-id>`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if checklistsAddItemChecklist == "" {
			return fmt.Errorf("--checklist is required")
		}

		item, err := client.CreateCheckItem(checklistsAddItemChecklist, args[0])
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}

		fmt.Printf("Item added: %s\n", item.Name)
		fmt.Printf("ID: %s\n", item.ID)
		return nil
	},
}

// ---- checklists check ----

var (
	checklistsCheckCard      string
	checklistsCheckChecklist string
)

var checklistsCheckCmd = &cobra.Command{
	Use:   "check <check-item-id>",
	Short: "Mark a checklist item as complete",
	Long: `Mark a Trello checklist item as complete.

Examples:
  trello checklists check <item-id> --card <card-id> --checklist <checklist-id>`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if checklistsCheckCard == "" {
			return fmt.Errorf("--card is required")
		}
		if checklistsCheckChecklist == "" {
			return fmt.Errorf("--checklist is required")
		}

		item, err := client.UpdateCheckItem(checklistsCheckCard, checklistsCheckChecklist, args[0], "complete")
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}

		fmt.Printf("Item checked: %s\n", item.Name)
		return nil
	},
}

// ---- checklists uncheck ----

var (
	checklistsUncheckCard      string
	checklistsUncheckChecklist string
)

var checklistsUncheckCmd = &cobra.Command{
	Use:   "uncheck <check-item-id>",
	Short: "Mark a checklist item as incomplete",
	Long: `Mark a Trello checklist item as incomplete.

Examples:
  trello checklists uncheck <item-id> --card <card-id> --checklist <checklist-id>`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if checklistsUncheckCard == "" {
			return fmt.Errorf("--card is required")
		}
		if checklistsUncheckChecklist == "" {
			return fmt.Errorf("--checklist is required")
		}

		item, err := client.UpdateCheckItem(checklistsUncheckCard, checklistsUncheckChecklist, args[0], "incomplete")
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}

		fmt.Printf("Item unchecked: %s\n", item.Name)
		return nil
	},
}

func init() {
	// create flags
	checklistsCreateCmd.Flags().StringVar(&checklistsCreateCardID, "card", "", "Card ID (required)")

	// add-item flags
	checklistsAddItemCmd.Flags().StringVar(&checklistsAddItemChecklist, "checklist", "", "Checklist ID (required)")

	// check flags
	checklistsCheckCmd.Flags().StringVar(&checklistsCheckCard, "card", "", "Card ID (required)")
	checklistsCheckCmd.Flags().StringVar(&checklistsCheckChecklist, "checklist", "", "Checklist ID (required)")

	// uncheck flags
	checklistsUncheckCmd.Flags().StringVar(&checklistsUncheckCard, "card", "", "Card ID (required)")
	checklistsUncheckCmd.Flags().StringVar(&checklistsUncheckChecklist, "checklist", "", "Checklist ID (required)")

	checklistsCmd.AddCommand(
		checklistsCreateCmd,
		checklistsDeleteCmd,
		checklistsAddItemCmd,
		checklistsCheckCmd,
		checklistsUncheckCmd,
	)
	rootCmd.AddCommand(checklistsCmd)
}
