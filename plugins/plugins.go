package plugins

import "github.com/bwmarrin/discordgo"

// Plugin is an interface that all plugins must implement.
type Plugin interface {
	Name() string
	Validate(session *discordgo.Session, message *discordgo.MessageCreate) bool
	Execute(session *discordgo.Session, message *discordgo.MessageCreate) error
	Enabled() bool
}
