package llm

import( 
	"net/http"

)

type LLMClient interface {
    Generate(prompt string) (*http.Response, error)
}
