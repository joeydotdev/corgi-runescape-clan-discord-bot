package plugins

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	memberlistentity "github.com/joeydotdev/corgi-discord-bot/internal/memberlist"
)

const (
	MissingMembersPluginName = "MissingMembersPlugin"
)

var InvalidPlatformError error = errors.New("Invalid platform. Valid platforms are `discord` and `teamspeak`")

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
	return strings.HasPrefix(message.Content, "!missing")
}

func (m *MissingMembersPlugin) isValidPlatform(platform string) bool {
	return platform == "discord" || platform == "teamspeak"
}

func (m *MissingMembersPlugin) findMissingDiscordMembers(session *discordgo.Session, message *discordgo.MessageCreate) error {
	members := _memberlist.GetMembers()
	missingMembers := []memberlistentity.Member{}
	guildMemberIDsToMembersInVoice := make(map[string]*discordgo.Member)
	guild, err := session.Guild(message.GuildID)
	if err != nil {
		return err
	}

	for _, voiceState := range guild.VoiceStates {
		if voiceState.ChannelID != "" {
			guildMemberIDsToMembersInVoice[voiceState.UserID] = voiceState.Member
		}
	}

	for _, member := range members {
		if _, ok := guildMemberIDsToMembersInVoice[member.DiscordID]; !ok {
			missingMembers = append(missingMembers, member)
		}
	}

	if len(missingMembers) == 0 {
		_, err = session.ChannelMessageSend(message.ChannelID, "No missing members found.")
		return err
	}

	missingMembersString := "Missing members:\n"
	matchedDiscordMembers := 0
	for _, member := range missingMembers {
		var discordMemberInstance *discordgo.Member
		// Find the member in the guild.
		for _, guildMember := range guild.Members {
			if guildMember.User.ID == member.DiscordID {
				discordMemberInstance = guildMember
				matchedDiscordMembers += 1
				break
			}
		}

		if discordMemberInstance == nil {
			log.Println("Could not find member in guild: " + member.Name)
			continue
		}

		missingMembersString += fmt.Sprintf("%s (%s)\n", member.Name, discordMemberInstance.User.Username)
	}

	if matchedDiscordMembers == 0 {
		return errors.New("No members found in Discord guild. Please make sure the bot is in the guild has required permissions.")
	}

	_, err = session.ChannelMessageSend(message.ChannelID, missingMembersString)
	return err
}

func (m *MissingMembersPlugin) findMissingTeamspeakMembers(session *discordgo.Session, message *discordgo.MessageCreate) error {
	return errors.New("Not implemented")
}

// Execute executes MissingMembersPlugin on an incoming Discord message.
func (m *MissingMembersPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) error {
	segments := strings.Split(message.Content, " ")
	if len(segments) < 2 {
		return TooFewArgumentsError
	}

	platform := segments[1]
	if !m.isValidPlatform(platform) {
		return InvalidPlatformError
	}

	var err error
	switch platform {
	case "discord":
		err = m.findMissingDiscordMembers(session, message)
	case "teamspeak":
		err = m.findMissingTeamspeakMembers(session, message)
	}

	return err
}
