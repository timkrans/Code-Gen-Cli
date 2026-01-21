package factory

import (
    "code-gen-cli/internal/llm"
    "code-gen-cli/internal/llm/providers"
)

func NewClient(cfg llm.Config) llm.LLMClient {
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
