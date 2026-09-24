package ai

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ThinkingConfig struct {
	Type string `json:"type"`
}

type DeepSeekRequest struct {
	Model    string         `json:"model"`
	Messages []Message      `json:"messages"`
	Thinking ThinkingConfig `json:"thinking"`
	Stream   bool           `json:"stream"`
}

type ResponseChoice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type DeepSeekResponse struct {
	Choices []ResponseChoice `json:"choices"`
}

type DeepSeekErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}
