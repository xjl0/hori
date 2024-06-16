package main

import (
	"discordbotgo/internal/calendar"
	"discordbotgo/internal/chatGPT"
	"discordbotgo/internal/config"
	"discordbotgo/internal/discord"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()

	gptClient := chatGPT.NewGPTClient(cfg.OpenApiToken, cfg.Proxy, cfg.SizeContext)

	bot := discord.NewBot(cfg.DiscordToken, gptClient)

	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 0, 0, now.Location())

			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}

			duration := next.Sub(now)
			log.Printf("Next execution at: %s\n", next.Format("2006-01-02 15:04:05"))

			timer := time.NewTimer(duration)
			<-timer.C

			message, err := calendar.CalendarReq()
			if err != nil {
				log.Println(err)
			}
			bot.Send(cfg.ChannelDashboard, message)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c

	if err := bot.Stop(); err != nil {
		panic(err)
	}
}
