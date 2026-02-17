package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/oneKn8/llm-cost/internal/display"
	"github.com/oneKn8/llm-cost/internal/storage"
	"github.com/spf13/cobra"
)

var budgetLimit float64

var budgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Manage monthly budget",
}

var budgetSetCmd = &cobra.Command{
	Use:     "set",
	Short:   "Set monthly budget limit",
	Example: "  llm-cost budget set --limit 50.00",
	Run: func(cmd *cobra.Command, args []string) {
		if err := db.SetBudget(budgetLimit); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting budget: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Monthly budget set to $%.2f\n", budgetLimit)
	},
}

var budgetStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show budget status",
	Run: func(cmd *cobra.Command, args []string) {
		limit, err := db.GetBudget()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting budget: %v\n", err)
			os.Exit(1)
		}
		if limit == 0 {
			fmt.Println("No budget set. Use: llm-cost budget set --limit <amount>")
			return
		}

		now := time.Now().UTC()
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		entries, err := db.QueryUsage(storage.QueryFilters{Since: monthStart})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error querying usage: %v\n", err)
			os.Exit(1)
		}

		var spent float64
		for _, e := range entries {
			spent += e.Cost
		}

		display.RenderBudget(spent, limit)
	},
}

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List supported models and pricing",
	Run: func(cmd *cobra.Command, args []string) {
		display.RenderModelList()
	},
}

func init() {
	budgetSetCmd.Flags().Float64Var(&budgetLimit, "limit", 0, "Monthly budget limit in USD")
	_ = budgetSetCmd.MarkFlagRequired("limit")
	budgetCmd.AddCommand(budgetSetCmd)
	budgetCmd.AddCommand(budgetStatusCmd)
	rootCmd.AddCommand(budgetCmd)
	rootCmd.AddCommand(modelsCmd)
}
