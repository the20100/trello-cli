package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vincentmaurin/trello-cli/internal/output"
)

var membersCmd = &cobra.Command{
	Use:   "members",
	Short: "View Trello member profiles",
}

// ---- members me ----

var membersMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Show the authenticated member's profile",
	Long: `Show the authenticated member's profile.

Examples:
  trello members me
  trello members me --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		member, err := client.GetMember("me", nil)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(member, output.IsPretty(cmd))
		}

		output.PrintKeyValue([][]string{
			{"ID", member.ID},
			{"Full Name", member.FullName},
			{"Username", member.Username},
			{"Email", member.Email},
			{"Bio", output.Truncate(member.Bio, 80)},
			{"URL", member.URL},
			{"Boards", fmt.Sprintf("%d", len(member.IDBoards))},
		})
		return nil
	},
}

// ---- members get ----

var membersGetCmd = &cobra.Command{
	Use:   "get <id-or-username>",
	Short: "Get a member's profile",
	Long: `Get the profile of any Trello member by their ID or username.

Examples:
  trello members get johndoe
  trello members get 5e7d1a2b3c4d5e6f7a8b9c0d
  trello members get johndoe --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		member, err := client.GetMember(args[0], nil)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(member, output.IsPretty(cmd))
		}

		output.PrintKeyValue([][]string{
			{"ID", member.ID},
			{"Full Name", member.FullName},
			{"Username", member.Username},
			{"Bio", output.Truncate(member.Bio, 80)},
			{"URL", member.URL},
		})
		return nil
	},
}

// ---- members boards ----

var (
	membersBoardsMember string
	membersBoardsFilter string
)

var membersBoardsCmd = &cobra.Command{
	Use:   "boards [id-or-username]",
	Short: "List boards for a member (default: self)",
	Long: `List boards for a Trello member. Defaults to the authenticated member.

Examples:
  trello members boards
  trello members boards johndoe
  trello members boards --filter all
  trello members boards --json`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := "me"
		if len(args) > 0 {
			target = args[0]
		}

		boards, err := client.GetMemberBoards(target, membersBoardsFilter)
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

		headers := []string{"ID", "NAME", "LAST ACTIVITY", "CLOSED"}
		rows := make([][]string, len(boards))
		for i, b := range boards {
			rows[i] = []string{
				b.ID,
				output.Truncate(b.Name, 44),
				output.FormatTime(b.DateLastActivity),
				output.FormatBool(b.Closed),
			}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

// ---- members cards ----

var (
	memberCardsTarget string
	memberCardsFilter string
)

var membersCardsCmd = &cobra.Command{
	Use:   "cards [id-or-username]",
	Short: "List cards assigned to a member (default: self)",
	Long: `List all cards assigned to a Trello member. Defaults to the authenticated member.

Examples:
  trello members cards
  trello members cards johndoe
  trello members cards --filter all
  trello members cards --json`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := "me"
		if len(args) > 0 {
			target = args[0]
		}

		cards, err := client.GetMemberCards(target, memberCardsFilter)
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

		headers := []string{"ID", "#", "NAME", "BOARD", "DUE"}
		rows := make([][]string, len(cards))
		for i, c := range cards {
			rows[i] = []string{
				c.ID,
				fmt.Sprintf("%d", c.IDShort),
				output.Truncate(c.Name, 44),
				output.Truncate(c.IDBoard, 24),
				output.FormatDate(c.Due),
			}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

// ---- members workspaces ----

var membersWorkspacesCmd = &cobra.Command{
	Use:   "workspaces [id-or-username]",
	Short: "List workspaces for a member (default: self)",
	Long: `List all Trello workspaces (organizations) for a member.

Examples:
  trello members workspaces
  trello members workspaces johndoe
  trello members workspaces --json`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := "me"
		if len(args) > 0 {
			target = args[0]
		}

		orgs, err := client.GetMemberOrganizations(target)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(orgs, output.IsPretty(cmd))
		}

		if len(orgs) == 0 {
			fmt.Println("No workspaces found.")
			return nil
		}

		headers := []string{"ID", "NAME", "DISPLAY NAME", "BOARDS"}
		rows := make([][]string, len(orgs))
		for i, o := range orgs {
			rows[i] = []string{
				o.ID,
				output.Truncate(o.Name, 24),
				output.Truncate(o.DisplayName, 30),
				fmt.Sprintf("%d", len(o.IDBoards)),
			}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

func init() {
	// members boards flags
	membersBoardsCmd.Flags().StringVar(&membersBoardsFilter, "filter", "open", "Filter: open, closed, all, members, organization, public, starred")

	// members cards flags
	membersCardsCmd.Flags().StringVar(&memberCardsFilter, "filter", "open", "Filter: open, closed, all, visible")

	membersCmd.AddCommand(
		membersMeCmd,
		membersGetCmd,
		membersBoardsCmd,
		membersCardsCmd,
		membersWorkspacesCmd,
	)
	rootCmd.AddCommand(membersCmd)
}
