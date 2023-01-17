package plugins

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const (
	MissingMembersPluginName = "MissingMembersPlugin"
)

type MissingMembersPlugin struct{}

// Enabled returns whether or not the MissingMembersPlugin is enabled.
func (m *MissingMembersPlugin) Enabled() bool {
	return true
}

// NewMissingMembersPlugin creates a new MissingMembersPlugin.
func NewMissingMembersPlugin() *MissingMembersPlugin {
	return &MissingMembersPlugin{}
}

// Name returns the name of the plugin.
func (m *MissingMembersPlugin) Name() string {
	return MissingMembersPluginName
}

// Validate validates whether or not we should execute MissingMembersPlugin on an incoming Discord message.
func (m *MissingMembersPlugin) Validate(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	fmt.Println(message.Content)
	return message.Content == "!missing"
}

// Execute executes MissingMembersPlugin on an incoming Discord message.
func (m *MissingMembersPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) error {
	_, err := session.ChannelMessageSend(message.ChannelID, "missing plugin")
	return err
}
