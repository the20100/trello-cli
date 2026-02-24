package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vincentmaurin/trello-cli/internal/output"
)

var (
	searchTypes []string
	searchLimit int
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search across Trello cards, boards, and members",
	Long: `Search Trello for cards, boards, and members matching a query.

By default searches all types. Use --type to narrow the search.

Examples:
  trello search "deploy"
  trello search "bug fix" --type cards
  trello search "John" --type members
  trello search "project" --type boards,cards
  trello search "deploy" --limit 5
  trello search "deploy" --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		results, err := client.Search(args[0], searchTypes, searchLimit)
		if err != nil {
			return err
		}

		if output.IsJSON(cmd) {
			return output.PrintJSON(results, output.IsPretty(cmd))
		}

		totalCards := len(results.Cards)
		totalBoards := len(results.Boards)
		totalMembers := len(results.Members)

		if totalCards+totalBoards+totalMembers == 0 {
			fmt.Println("No results found.")
			return nil
		}

		// Cards
		if totalCards > 0 {
			fmt.Printf("\nCards (%d)\n", totalCards)
			headers := []string{"ID", "#", "NAME", "DUE", "LABELS"}
			rows := make([][]string, totalCards)
			for i, c := range results.Cards {
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
					output.Truncate(c.Name, 44),
					output.FormatDate(c.Due),
					output.FormatLabels(labelNames),
				}
			}
			output.PrintTable(headers, rows)
		}

		// Boards
		if totalBoards > 0 {
			fmt.Printf("\nBoards (%d)\n", totalBoards)
			headers := []string{"ID", "NAME", "URL", "CLOSED"}
			rows := make([][]string, totalBoards)
			for i, b := range results.Boards {
				rows[i] = []string{
					b.ID,
					output.Truncate(b.Name, 44),
					b.ShortURL,
					output.FormatBool(b.Closed),
				}
			}
			output.PrintTable(headers, rows)
		}

		// Members
		if totalMembers > 0 {
			fmt.Printf("\nMembers (%d)\n", totalMembers)
			headers := []string{"ID", "NAME", "USERNAME"}
			rows := make([][]string, totalMembers)
			for i, m := range results.Members {
				rows[i] = []string{m.ID, m.FullName, m.Username}
			}
			output.PrintTable(headers, rows)
		}

		return nil
	},
}

func init() {
	searchCmd.Flags().StringArrayVar(&searchTypes, "type", nil, "Limit to: cards, boards, members (can be repeated or comma-separated)")
	searchCmd.Flags().IntVar(&searchLimit, "limit", 10, "Max results per type (1-1000)")
	rootCmd.AddCommand(searchCmd)
}
