package plugins

import (
	"github.com/bwmarrin/discordgo"
)

const (
	ManageXpTrackerPluginName = "XpTrackerPlugin"
)

type ManageXpTrackerPlugin struct{}

// Enabled returns whether or not the ManageXpTrackerPlugin is enabled.
func (m *ManageXpTrackerPlugin) Enabled() bool {
	return true
}

// NewManageXpTrackerPlugin creates a new ManageXpTrackerPlugin.
func NewManageXpTrackerPlugin() *ManageXpTrackerPlugin {
	return &ManageXpTrackerPlugin{}
}

// Name returns the name of the plugin.
func (m *ManageXpTrackerPlugin) Name() string {
	return PingCommandPluginName
}

// Validate validates whether or not we should execute ManageXpTrackerPlugin on an incoming Discord message.
func (m *ManageXpTrackerPlugin) Validate(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	return message.Content == "!xptracker"
}

// Execute executes ManageXpTrackerPlugin on an incoming Discord message.
func (m *ManageXpTrackerPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) error {
	_, err := session.ChannelMessageSend(message.ChannelID, "xp tracker")
	return err
}
