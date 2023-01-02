package plugins

import (
	"github.com/bwmarrin/discordgo"
)

const (
	ManageWorldTrackerPluginName = "ManageWorldTrackerPlugin"
)

type ManageWorldTrackerPlugin struct{}

// Enabled returns whether or not the ManageWorldTrackerPlugin is enabled.
func (p *ManageWorldTrackerPlugin) Enabled() bool {
	return true
}

// NewManageWorldTrackerPlugin creates a new ManageWorldTrackerPlugin.
func NewManageWorldTrackerPlugin() *ManageWorldTrackerPlugin {
	return &ManageWorldTrackerPlugin{}
}

// Name returns the name of the plugin.
func (p *ManageWorldTrackerPlugin) Name() string {
	return PingCommandPluginName
}

// Validate validates whether or not we should execute ManageWorldTrackerPlugin on an incoming Discord message.
func (p *ManageWorldTrackerPlugin) Validate(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	return message.Content == "!worldtracker"
}

// Execute executes ManageWorldTrackerPlugin on an incoming Discord message.
func (p *ManageWorldTrackerPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) error {
	_, err := session.ChannelMessageSend(message.ChannelID, "worldtracker command")
	return err
}
