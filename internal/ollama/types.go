package ollama

type (
	PromptReq struct {
		Model  string `json:"model"`
		System string `json:"system"`
		Prompt string `json:"prompt"`
		Stream bool   `json:"stream"`
	}

	PromptResp struct {
		Response string `json:"response"`
	}
)
