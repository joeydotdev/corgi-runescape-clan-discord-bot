package plugins

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	MassPMCommandPluginName = "MassPMCommandPlugin"
	MassPMChannelID         = "1070035899221545171"
)

type MassPMCommandPlugin struct{}

// Enabled returns whether or not the MassPMCommandPlugin is enabled.
func (p *MassPMCommandPlugin) Enabled() bool {
	return true
}

// NewMassPMCommandPlugin creates a new MassPMCommandPlugin.
func NewMassPMCommandPlugin() *MassPMCommandPlugin {
	return &MassPMCommandPlugin{}
}

// Name returns the name of the plugin.
func (p *MassPMCommandPlugin) Name() string {
	return MassPMCommandPluginName
}

// Validate validates whether or not we should execute MassPMCommandPlugin on an incoming Discord message.
func (p *MassPMCommandPlugin) Validate(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	return strings.HasPrefix(message.Content, "!masspm") && message.ChannelID == MassPMChannelID
}

// Execute executes MassPMCommandPlugin on an incoming Discord message.
func (p *MassPMCommandPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) error {
	_, err := session.ChannelMessageSend(message.ChannelID, "pong")
	return err
}
