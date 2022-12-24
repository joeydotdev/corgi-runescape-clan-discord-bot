package plugins

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	ManageMemberlistPluginName = "ManageMemberlistCommand"
)

type ManageMemberlistPlugin struct{}

// Enabled returns whether or not the ManageMemberlistPlugin is enabled.
func (m *ManageMemberlistPlugin) Enabled() bool {
	return true
}

// NewManageMemberlistPlugin creates a new ManageMemberlistCo.
func NewManageMemberlistPlugin() *ManageMemberlistPlugin {
	return &ManageMemberlistPlugin{}
}

// Name returns the name of the plugin.
func (m *ManageMemberlistPlugin) Name() string {
	return ManageMemberlistPluginName
}

// Validate validates whether or not we should execute ManageMemberlistPlugin on an incoming Discord message.
func (m *ManageMemberlistPlugin) Validate(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	return strings.HasPrefix(message.Content, "!memberlist")
}

// Execute executes ManageMemberlistPlugin on an incoming Discord message.
func (m *ManageMemberlistPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) error {
	_, err := session.ChannelMessageSend(message.ChannelID, "memberlist")
	return err
}
