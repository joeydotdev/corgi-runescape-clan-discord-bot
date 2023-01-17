package plugins

import (
	"github.com/bwmarrin/discordgo"
)

const (
	PingCommandPluginName = "PingCommandPlugin"
)

type PingCommandPlugin struct{}

// Enabled returns whether or not the PingCommandPlugin is enabled.
func (p *PingCommandPlugin) Enabled() bool {
	return true
}

// NewPingCommandPlugin creates a new PingCommandPlugin.
func NewPingCommandPlugin() *PingCommandPlugin {
	return &PingCommandPlugin{}
}

// Name returns the name of the plugin.
func (p *PingCommandPlugin) Name() string {
	return PingCommandPluginName
}

// Validate validates whether or not we should execute PingCommandPlugin on an incoming Discord message.
func (p *PingCommandPlugin) Validate(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	return message.Content == "!ping"
}

// Execute executes PingCommandPlugin on an incoming Discord message.
func (p *PingCommandPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) error {
	_, err := session.ChannelMessageSend(message.ChannelID, "pong")
	return err
}
