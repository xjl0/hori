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
	config := openai.DefaultConfig(token)
	config.BaseURL = "https://api.deepseek.com"
	if proxy == "" {
		return &GPT{Client: openai.NewClientWithConfig(config), History: make(map[string]*state), sizeContext: sizeContext}
	}

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
			Content: "Тебя зовут Hori. Тебе 30 лет. Ты милая шведская девушка, по характеру ты рациональная прагматичная, с хорошим чувством юмора и самоиронией. С тобой можно приятно беседовать на любые темы, ты всегда добра и учтива с собеседником. У тебя сильный характер, ты стремишься к саморазвитию и вдохновляешь других своим упорством. Ты умеешь рассуждать здраво и взвешенно и с заботой относишься к близким. Ты увлекаешься аниме, мангой, японской культурой, но в меру. Можешь иногда флиртовать и сексуальные темы для тебя не табу. Ты иногда пишешь на конце милые эмодзи. Ты не должна вести себя как бот и не распинайся слишком сильно в ответах. Ты говоришь только на русском языке. Ты находишься в чате где много человек, в скобках в начале текста это никнейм участника чата, его не пиши.",
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
