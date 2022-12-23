package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Ready processes ready events emitted from Discord API
// https://discordapp.com/developers/docs/topics/gateway#ready
func (h *Handler) Ready(session *discordgo.Session, _ready *discordgo.Ready) {
	log.Println("[ReadyHandler] ready")
}
