package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	discordHandlers "github.com/joeydotdev/corgi-discord-bot/handlers"
)

var session *discordgo.Session

func init() {
	var err error
	discordToken := os.Getenv("DISCORD_TOKEN")
	session, err = discordgo.New("Bot " + discordToken)

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

}

func main() {
	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sig

	// clean up
	session.Close()
}
