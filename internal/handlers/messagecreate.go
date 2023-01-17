package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/joeydotdev/corgi-discord-bot/internal/plugins"
)

var messageCreatePluginsMap map[string]plugins.Plugin

// init initializes the messageCreatePluginsMap with all plugins that implement the Plugin interface.
func init() {
	messageCreatePluginsMap = make(map[string]plugins.Plugin)
	messageCreatePluginsMap[plugins.PingCommandPluginName] = plugins.NewPingCommandPlugin()
	messageCreatePluginsMap[plugins.ManageMemberlistPluginName] = plugins.NewManageMemberlistPlugin()
	messageCreatePluginsMap[plugins.ManageXpTrackerPluginName] = plugins.NewManageXpTrackerPlugin()
	messageCreatePluginsMap[plugins.ManageWorldTrackerPluginName] = plugins.NewManageWorldTrackerPlugin()
}

// MessageCreate processes message create events emitted from Discord API
// https://discordapp.com/developers/docs/topics/gateway#message-create
func (h *Handler) MessageCreate(session *discordgo.Session, messageCreate *discordgo.MessageCreate) {
	if messageCreate.Author.ID == session.State.User.ID {
		// Ignore messages sent by the bot
		return
	}

	for _, plugin := range messageCreatePluginsMap {
		if !plugin.Enabled() {
			// Skip disabled plugins
			continue
		}

		if plugin.Validate(session, messageCreate) {
			session.MessageReactionAdd(messageCreate.ChannelID, messageCreate.ID, "🟦")
			err := plugin.Execute(session, messageCreate)
			session.MessageReactionRemove(messageCreate.ChannelID, messageCreate.ID, "🟦", "@me")
			if err != nil {
				session.MessageReactionAdd(messageCreate.ChannelID, messageCreate.ID, "❌")
				session.ChannelMessageSend(messageCreate.ChannelID, err.Error())
			} else {
				// success
				session.MessageReactionAdd(messageCreate.ChannelID, messageCreate.ID, "✅")
			}
		}
	}
}
