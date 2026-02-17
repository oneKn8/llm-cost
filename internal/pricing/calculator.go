package pricing

import "strings"

func Calculate(provider, model string, inputTokens, outputTokens, cachedTokens int) float64 {
	p := findModel(provider, model)
	if p == nil {
		return 0
	}

	regularInput := inputTokens - cachedTokens
	if regularInput < 0 {
		regularInput = 0
	}

	cost := float64(regularInput) / 1_000_000 * p.InputPer1M
	cost += float64(outputTokens) / 1_000_000 * p.OutputPer1M
	cost += float64(cachedTokens) / 1_000_000 * p.CachedPer1M

	return cost
}

func findModel(provider, model string) *ModelPricing {
	provider = strings.ToLower(provider)
	model = strings.ToLower(model)

	for _, m := range Models {
		if strings.ToLower(m.Provider) == provider && strings.ToLower(m.Model) == model {
			return &m
		}
	}

	// Fuzzy match: check if model name contains the search term
	for _, m := range Models {
		if strings.ToLower(m.Provider) == provider && strings.Contains(strings.ToLower(m.Model), model) {
			return &m
		}
	}

	return nil
}
