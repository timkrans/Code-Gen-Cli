package llm

import (
    "code-gen-cli/internal/llm/providers"
    "code-gen-cli/internal/llm"
)

func NewClient(cfg Config) LLMClient {
    switch cfg.Provider {

    case "ollama":
        return providers.NewOllama(cfg)

    case "huggingface":
        return providers.NewHuggingFace(cfg)

    case "openai":
        return providers.NewOpenAI(cfg)

    case "google":
        return providers.NewGoogle(cfg)

    case "anthropic":
        return providers.NewAnthropic(cfg)

    default:
        panic("unsupported LLM provider: " + cfg.Provider)
    }
}
