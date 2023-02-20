package main

import (
	"discordbotgo/pkg/discord"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}
	if err := initConfig(); err != nil {
		log.Fatalf("Error init config file: %s", err.Error())
	}
	
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

func initConfig() error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}
