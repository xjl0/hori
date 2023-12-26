package main

import (
	"discordbotgo/pkg/discord"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	dBot := discord.NewDBot(os.Getenv("DGU_TOKEN"))
	if err := dBot.Start(); err != nil {
		log.Fatal(err)
	}
	defer dBot.Stop()

	go dBot.HandlerGetNews()
	go dBot.HandlerGetHoliday()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-c
}
