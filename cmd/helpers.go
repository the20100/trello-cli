package cmd

import (
	"fmt"
	"net/url"

	"github.com/vincentmaurin/trello-cli/internal/api"
	"github.com/vincentmaurin/trello-cli/internal/output"
)

// buildParams creates a url.Values from alternating key/value strings,
// skipping pairs where the value is empty.
func buildParams(pairs ...string) url.Values {
	p := url.Values{}
	for i := 0; i+1 < len(pairs); i += 2 {
		if pairs[i+1] != "" {
			p.Set(pairs[i], pairs[i+1])
		}
	}
	return p
}

// printAPICardsTable renders a slice of api.Card as a table.
func printAPICardsTable(cards []api.Card) {
	if len(cards) == 0 {
		fmt.Println("No cards found.")
		return
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
			output.Truncate(c.Name, 44),
			output.FormatDate(c.Due),
			output.FormatLabels(labelNames),
		}
	}
	output.PrintTable(headers, rows)
}
