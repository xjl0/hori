package discord

import (
	"context"
	"discordbotgo/internal/chatGPT"
	"errors"
	"fmt"
	"github.com/alexsergivan/transliterator"
	"github.com/bwmarrin/discordgo"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
	"regexp"
	"strings"
	"time"
)

type Bot struct {
	dBot *discordgo.Session
}

func NewBot(token string, clientGPT *chatGPT.GPT) *Bot {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	trans := transliterator.NewTransliterator(nil)

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		MessageHandler(s, m, clientGPT, trans)
	})

	if err := session.Open(); err != nil {
		panic(err)
	}

	return &Bot{dBot: session}
}

func (b *Bot) Stop() error {
	return b.dBot.Close()
}

func (b *Bot) Send(channelID, message string) {
	b.dBot.ChannelMessageSend(channelID, message)
}

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate, gpt *chatGPT.GPT, trans *transliterator.Transliterator) {
	if m.Author.Bot {
		return
	}

	reg, err := regexp.Compile("[^a-zA-Z]+")
	if err != nil {
		return
	}

	firstName := reg.ReplaceAllString(trans.Transliterate(m.Author.Username, "en"), "")
	if m.Author.GlobalName != "" {
		firstName = reg.ReplaceAllString(trans.Transliterate(m.Author.GlobalName, "en"), "")
	}

	if trimTag(m.Content) != "" {
		if m.Message.ReferencedMessage != nil && m.Message.ReferencedMessage.Author.Bot {
			gpt.AddHistory(m.ChannelID, fmt.Sprintf("(%s) %s", "Hori", trimTag(m.Message.ReferencedMessage.Content)), openai.ChatMessageRoleAssistant, "Hori")
		}
		gpt.AddHistory(m.ChannelID, fmt.Sprintf("(%s) %s", firstName, trimTag(m.Content)), openai.ChatMessageRoleUser, firstName)
	}

	if !isMessageForHori(s, m) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	var messages []openai.ChatCompletionMessage

	messages = gpt.GetHistory(m.ChannelID)

	stream, err := gpt.Client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model:    "deepseek-chat",
			Messages: messages,
			Stream:   true,
			TopP:     1,
			N:        1,
		},
	)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Прости, я сегодня уже не могу ответить, занята, давай завтра \n ```%s\n```", err.Error()), &discordgo.MessageReference{
			MessageID: m.Message.ID,
			ChannelID: m.ChannelID,
			GuildID:   m.GuildID,
		})
		log.Printf("ChatCompletion error: %v\n", err)
		return
	}
	if stream == nil {
		return
	}
	defer stream.Close()

	throttleTime := time.Now()
	content := ""
	msgID := ""

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Printf("\nStream error: %v\n", err)
			return
		}

		if response.Choices[0].Delta.Content == "" {
			continue
		}

		content += response.Choices[0].Delta.Content

		if msgID == "" {
			message, _ := s.ChannelMessageSendReply(m.ChannelID, content, &discordgo.MessageReference{
				MessageID: m.Message.ID,
				ChannelID: m.ChannelID,
				GuildID:   m.GuildID,
			})
			msgID = message.ID
		} else {
			if throttleTime.Add(time.Second).Before(time.Now()) {
				s.ChannelMessageEdit(m.ChannelID, msgID, content)
				throttleTime = time.Now()
			}
		}
	}

	s.ChannelMessageEdit(m.ChannelID, msgID, content)

	gpt.AddHistory(m.ChannelID, content, openai.ChatMessageRoleAssistant, "Hori")
	return
}

func isMessageForHori(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if m.Message == nil || m.Message.Mentions == nil || len(m.Message.Mentions) == 0 || m.Author.ID == s.State.User.ID {
		return false
	}
	for _, mention := range m.Message.Mentions {
		if mention.Username == "Hori" {
			break
		}
		return false
	}

	return true
}

func trimTag(message string) string {
	return strings.TrimSpace(regexp.MustCompile("<[^>]*>").ReplaceAllString(message, ""))
}
