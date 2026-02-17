package pricing

type ModelPricing struct {
	Provider     string
	Model        string
	InputPer1M   float64
	OutputPer1M  float64
	CachedPer1M  float64
}

var Models = []ModelPricing{
	// OpenAI
	{"openai", "gpt-4o", 2.50, 10.00, 1.25},
	{"openai", "gpt-4o-mini", 0.15, 0.60, 0.075},
	{"openai", "gpt-4.1", 2.00, 8.00, 0.50},
	{"openai", "gpt-4.1-mini", 0.40, 1.60, 0.10},
	{"openai", "gpt-4.1-nano", 0.10, 0.40, 0.025},
	{"openai", "o3", 10.00, 40.00, 2.50},
	{"openai", "o3-mini", 1.10, 4.40, 0.55},
	{"openai", "o4-mini", 1.10, 4.40, 0.275},

	// Anthropic
	{"anthropic", "claude-opus-4-6", 15.00, 75.00, 7.50},
	{"anthropic", "claude-sonnet-4-5-20250929", 3.00, 15.00, 1.50},
	{"anthropic", "claude-haiku-4-5-20251001", 0.80, 4.00, 0.40},

	// Google
	{"google", "gemini-2.5-pro", 1.25, 10.00, 0.315},
	{"google", "gemini-2.5-flash", 0.15, 0.60, 0.0375},
	{"google", "gemini-2.0-flash", 0.10, 0.40, 0.025},

	// Groq
	{"groq", "llama-3.3-70b", 0.59, 0.79, 0},
	{"groq", "llama-3.1-8b", 0.05, 0.08, 0},
	{"groq", "llama-4-scout", 0.11, 0.34, 0},
	{"groq", "deepseek-r1-distill-70b", 0.75, 0.99, 0},

	// DeepSeek
	{"deepseek", "deepseek-r1", 0.55, 2.19, 0.14},
	{"deepseek", "deepseek-v3", 0.27, 1.10, 0.07},
}
