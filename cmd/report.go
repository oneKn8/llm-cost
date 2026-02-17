package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/oneKn8/llm-cost/internal/display"
	"github.com/oneKn8/llm-cost/internal/storage"
	"github.com/spf13/cobra"
)

var (
	reportPeriod   string
	reportProvider string
	reportModel    string
	reportFormat   string
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Show cost reports",
	Example: `  llm-cost report --period today
  llm-cost report --period month --provider anthropic
  llm-cost report --period week --format json`,
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now().UTC()
		var since time.Time

		switch reportPeriod {
		case "today":
			since = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		case "week":
			since = now.AddDate(0, 0, -7)
		case "month":
			since = now.AddDate(0, -1, 0)
		default:
			since = time.Time{}
		}

		filters := storage.QueryFilters{
			Since:    since,
			Provider: reportProvider,
			Model:    reportModel,
		}

		entries, err := db.QueryUsage(filters)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error querying usage: %v\n", err)
			os.Exit(1)
		}

		switch reportFormat {
		case "json":
			data, _ := json.MarshalIndent(entries, "", "  ")
			fmt.Println(string(data))
		case "csv":
			w := csv.NewWriter(os.Stdout)
			w.Write([]string{"timestamp", "provider", "model", "input_tokens", "output_tokens", "cached_tokens", "cost"})
			for _, e := range entries {
				w.Write([]string{
					e.Timestamp.Format(time.RFC3339),
					e.Provider, e.Model,
					fmt.Sprintf("%d", e.InputTokens),
					fmt.Sprintf("%d", e.OutputTokens),
					fmt.Sprintf("%d", e.CachedTokens),
					fmt.Sprintf("%.6f", e.Cost),
				})
			}
			w.Flush()
		default:
			display.RenderReport(entries, reportPeriod)
		}
	},
}

func init() {
	reportCmd.Flags().StringVar(&reportPeriod, "period", "month", "Time period (today, week, month, all)")
	reportCmd.Flags().StringVar(&reportProvider, "provider", "", "Filter by provider")
	reportCmd.Flags().StringVar(&reportModel, "model", "", "Filter by model")
	reportCmd.Flags().StringVar(&reportFormat, "format", "table", "Output format (table, json, csv)")
	rootCmd.AddCommand(reportCmd)
}
