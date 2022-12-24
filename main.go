package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joeydotdev/corgi-discord-bot/config"
	discordHandlers "github.com/joeydotdev/corgi-discord-bot/handlers"
	"github.com/joeydotdev/corgi-discord-bot/storage"
)

var session *discordgo.Session

func init() {
	var err error
	config := config.Load()
	session, err = discordgo.New("Bot " + config.Token)

	if err != nil {
		panic("failed to initalize bot")
	}

	handlers := discordHandlers.New()
	session.AddHandler(handlers.MessageCreate)
	session.AddHandler(handlers.Ready)

	err = session.Open()
	if err != nil {
		panic(err)
	}

	err = storage.InitializeS3()
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sig

	// clean up
	session.Close()
}
