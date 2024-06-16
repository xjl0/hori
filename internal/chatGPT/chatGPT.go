package chatGPT

import (
	"github.com/sashabaranov/go-openai"
	"net/http"
	"net/url"
)

type state struct {
	message []openai.ChatCompletionMessage
}

type GPT struct {
	History     map[string]*state
	Client      *openai.Client
	sizeContext int
}

func NewGPTClient(token, proxy string, sizeContext int) *GPT {
	if proxy == "" {
		return &GPT{Client: openai.NewClient(token), History: make(map[string]*state), sizeContext: sizeContext}
	}
	config := openai.DefaultConfig(token)
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		panic(err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
	}

	return &GPT{Client: openai.NewClientWithConfig(config), History: make(map[string]*state), sizeContext: sizeContext}
}

func (g *GPT) AddHistory(channelID, message, role, userName string) {
	history, exists := g.History[channelID]
	if !exists {
		history = &state{message: make([]openai.ChatCompletionMessage, 0, g.sizeContext)}
		history.message = append(history.message, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Тебя зовут Hori. Тебе 18 лет. Ты милая японская девочка. Ты увлекаешься аниме. Ты ведёшь всегда себя мило. Ты иногда пишешь на конце милые эмодзи. Ты не должна вести себя как бот.",
			Name:    "Hori",
		})
		g.History[channelID] = history
	}

	if len(history.message) > 1 && len(history.message) >= g.sizeContext {
		history.message = append(history.message[:1], history.message[2:]...)
	}

	history.message = append(history.message, openai.ChatCompletionMessage{
		Role:    role,
		Content: message,
		Name:    userName,
	})
}

func (g *GPT) GetHistory(channelID string) []openai.ChatCompletionMessage {
	if _, exists := g.History[channelID]; !exists {
		return nil
	}
	return g.History[channelID].message
}
