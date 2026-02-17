package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/oneKn8/llm-cost/internal/pricing"
	"github.com/oneKn8/llm-cost/internal/storage"
	"github.com/spf13/cobra"
)

var (
	trackProvider     string
	trackModel        string
	trackInputTokens  int
	trackOutputTokens int
	trackCachedTokens int
)

var trackCmd = &cobra.Command{
	Use:   "track",
	Short: "Record a usage entry",
	Example: `  llm-cost track --provider anthropic --model claude-sonnet-4-5-20250929 --input-tokens 1500 --output-tokens 500
  llm-cost track -p openai -m gpt-4o -i 10000 -o 2000`,
	Run: func(cmd *cobra.Command, args []string) {
		cost := pricing.Calculate(trackProvider, trackModel, trackInputTokens, trackOutputTokens, trackCachedTokens)

		entry := storage.UsageEntry{
			Timestamp:    time.Now().UTC(),
			Provider:     trackProvider,
			Model:        trackModel,
			InputTokens:  trackInputTokens,
			OutputTokens: trackOutputTokens,
			CachedTokens: trackCachedTokens,
			Cost:         cost,
		}

		if err := db.RecordUsage(entry); err != nil {
			fmt.Fprintf(os.Stderr, "Error recording usage: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Recorded: %s/%s | %d in + %d out = $%.6f\n",
			trackProvider, trackModel, trackInputTokens, trackOutputTokens, cost)
	},
}

func init() {
	trackCmd.Flags().StringVarP(&trackProvider, "provider", "p", "", "Provider (openai, anthropic, groq, google, deepseek)")
	trackCmd.Flags().StringVarP(&trackModel, "model", "m", "", "Model name")
	trackCmd.Flags().IntVarP(&trackInputTokens, "input-tokens", "i", 0, "Input token count")
	trackCmd.Flags().IntVarP(&trackOutputTokens, "output-tokens", "o", 0, "Output token count")
	trackCmd.Flags().IntVarP(&trackCachedTokens, "cached-tokens", "c", 0, "Cached input token count")
	_ = trackCmd.MarkFlagRequired("provider")
	_ = trackCmd.MarkFlagRequired("model")
	_ = trackCmd.MarkFlagRequired("input-tokens")
	_ = trackCmd.MarkFlagRequired("output-tokens")
	rootCmd.AddCommand(trackCmd)
}
