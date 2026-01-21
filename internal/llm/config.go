package llm

import "os"

type Config struct {
    Provider string
    Model    string

    OllamaBaseURL string

    HFBaseURL string
    HFAPIKey  string

    OpenAIBaseURL string
    OpenAIAPIKey  string

    GoogleAPIKey string
    GoogleModel  string

    AnthropicAPIKey  string
    AnthropicBaseURL string
}

func LoadConfig() Config {
    return Config{
        Provider: os.Getenv("LLM_PROVIDER"),
        Model:    os.Getenv("LLM_MODEL"),

        OllamaBaseURL: os.Getenv("OLLAMA_BASE_URL"),

        HFBaseURL: os.Getenv("HF_BASE_URL"),
        HFAPIKey:  os.Getenv("HF_API_KEY"),

        OpenAIBaseURL: os.Getenv("OPENAI_BASE_URL"),
        OpenAIAPIKey:  os.Getenv("OPENAI_API_KEY"),

        GoogleAPIKey: os.Getenv("GOOGLE_API_KEY"),
        GoogleModel:  os.Getenv("GOOGLE_MODEL"),

        AnthropicAPIKey:  os.Getenv("ANTHROPIC_API_KEY"),
        AnthropicBaseURL: os.Getenv("ANTHROPIC_BASE_URL"),
    }
}
