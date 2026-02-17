package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/oneKn8/llm-cost/internal/storage"
	"github.com/spf13/cobra"
)

var dbPath string
var db *storage.DB

var rootCmd = &cobra.Command{
	Use:   "llm-cost",
	Short: "Multi-provider LLM API cost tracker",
	Long:  "Track token usage and costs across OpenAI, Anthropic, Groq, Google, and DeepSeek from the terminal.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "help" || cmd.Name() == "completion" {
			return nil
		}
		var err error
		db, err = storage.Open(dbPath)
		if err != nil {
			return fmt.Errorf("open database: %w", err)
		}
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if db != nil {
			db.Close()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	home, _ := os.UserHomeDir()
	defaultDB := filepath.Join(home, ".llm-cost.db")
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", defaultDB, "SQLite database path")
}
