package display

import (
	"fmt"
	"strings"

	"github.com/oneKn8/llm-cost/internal/pricing"
	"github.com/oneKn8/llm-cost/internal/storage"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#58a6ff"))
	greenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#39d353"))
	dimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#8b949e"))
	warnStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#f0883e"))
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#f85149"))
	boldStyle   = lipgloss.NewStyle().Bold(true)
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#c9d1d9")).Underline(true)
)

func RenderReport(entries []storage.UsageEntry, period string) {
	fmt.Println()
	fmt.Println(titleStyle.Render(fmt.Sprintf("  Cost Report (%s)", period)))
	fmt.Println()

	if len(entries) == 0 {
		fmt.Println(dimStyle.Render("  No usage data found."))
		fmt.Println()
		return
	}

	var totalCost float64
	var totalInput, totalOutput int
	byProvider := map[string]float64{}
	byModel := map[string]float64{}

	for _, e := range entries {
		totalCost += e.Cost
		totalInput += e.InputTokens
		totalOutput += e.OutputTokens
		byProvider[e.Provider] += e.Cost
		byModel[e.Provider+"/"+e.Model] += e.Cost
	}

	fmt.Printf("  Total cost:      %s\n", greenStyle.Render(fmt.Sprintf("$%.4f", totalCost)))
	fmt.Printf("  Total tokens:    %s in / %s out\n",
		boldStyle.Render(fmt.Sprintf("%d", totalInput)),
		boldStyle.Render(fmt.Sprintf("%d", totalOutput)))
	fmt.Printf("  Entries:         %d\n", len(entries))
	fmt.Println()

	fmt.Println(headerStyle.Render("  By Provider"))
	for p, cost := range byProvider {
		pct := cost / totalCost * 100
		fmt.Printf("    %-15s %s  (%s)\n", boldStyle.Render(p), greenStyle.Render(fmt.Sprintf("$%.4f", cost)), dimStyle.Render(fmt.Sprintf("%.1f%%", pct)))
	}
	fmt.Println()

	fmt.Println(headerStyle.Render("  By Model"))
	for m, cost := range byModel {
		fmt.Printf("    %-40s %s\n", dimStyle.Render(m), greenStyle.Render(fmt.Sprintf("$%.4f", cost)))
	}
	fmt.Println()
}

func RenderBudget(spent, limit float64) {
	fmt.Println()
	fmt.Println(titleStyle.Render("  Budget Status"))
	fmt.Println()

	pct := spent / limit * 100
	barLen := 30
	filled := int(pct / 100 * float64(barLen))
	if filled > barLen {
		filled = barLen
	}

	barStyle := greenStyle
	if pct >= 80 {
		barStyle = warnStyle
	}
	if pct >= 100 {
		barStyle = errorStyle
	}

	bar := barStyle.Render(strings.Repeat("█", filled)) + dimStyle.Render(strings.Repeat("░", barLen-filled))

	fmt.Printf("  Spent:   %s / $%.2f\n", greenStyle.Render(fmt.Sprintf("$%.4f", spent)), limit)
	fmt.Printf("  Used:    [%s] %.1f%%\n", bar, pct)

	if pct >= 100 {
		fmt.Println(errorStyle.Render("  OVER BUDGET"))
	} else if pct >= 80 {
		fmt.Println(warnStyle.Render("  WARNING: Approaching budget limit"))
	}
	fmt.Println()
}

func RenderModelList() {
	fmt.Println()
	fmt.Println(titleStyle.Render("  Supported Models"))
	fmt.Println()

	currentProvider := ""
	for _, m := range pricing.Models {
		if m.Provider != currentProvider {
			currentProvider = m.Provider
			fmt.Println(headerStyle.Render(fmt.Sprintf("  %s", strings.ToUpper(currentProvider))))
		}
		cached := ""
		if m.CachedPer1M > 0 {
			cached = fmt.Sprintf(" (cached: $%.3f)", m.CachedPer1M)
		}
		fmt.Printf("    %-35s %s in / %s out%s\n",
			boldStyle.Render(m.Model),
			greenStyle.Render(fmt.Sprintf("$%.3f", m.InputPer1M)),
			greenStyle.Render(fmt.Sprintf("$%.3f", m.OutputPer1M)),
			dimStyle.Render(cached))
	}
	fmt.Println()
}
