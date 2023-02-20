package discord

import (
	"discordbotgo/pkg/calendar"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"log"
	"time"
)

type DBot struct {
	token       string
	testChannel string
	mainChannel string
	codeInvite  string
	dg          *discordgo.Session
	m           *discordgo.MessageCreate
}

func NewDBot(token string) *DBot {
	return &DBot{token: token}
}

func (b *DBot) Start() error {
	session, err := discordgo.New("Bot " + b.token)
	if err != nil {
		return err
	}
	
	b.dg = session
	
	b.dg.AddHandler(eventCreate)
	b.dg.AddHandler(isALife)
	b.dg.AddHandler(message)
	
	err = b.dg.Open()
	if err != nil {
		return err
	}
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	
	return nil
}

func (b *DBot) HandlerGetNews() {
	for range time.Tick(time.Minute * 10) {
		sendNews(b.dg)
	}
}

func (b *DBot) HandlerGetHoliday() {
	newsChannel := viper.GetString("discord.newschannel")
	testChannel := viper.GetString("discord.testchannel")
	for range time.Tick(time.Hour) {
		if time.Now().Hour() == 8 {
			text, err := calendar.CalendarReq()
			if err != nil {
				b.dg.ChannelMessageSend(testChannel, err.Error())
			}
			b.dg.ChannelMessageSend(newsChannel, text)
		}
	}
}

func (b *DBot) Stop() error {
	return b.dg.Close()
}
