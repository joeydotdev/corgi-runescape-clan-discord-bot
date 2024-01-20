package plugins

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/joeydotdev/corgi-discord-bot/internal/discord"
	"github.com/joeydotdev/corgi-discord-bot/internal/memberlist"
)

const (
	MissingSignupsPluginName = "MissingSignupsPlugin"
	EventsCategoryChannelID  = "1071951449652740126"
	YesEmoji                 = "✅"
	NoEmoji                  = "❌"
)

type MissingSignupsPlugin struct{}

// Enabled returns whether or not the MissingSignupsPlugin is enabled.
func (p *MissingSignupsPlugin) Enabled() bool {
	return true
}

// NewMissingSignupsPlugin creates a new MissingSignupsPlugin.
func NewMissingSignupsPlugin() *MissingSignupsPlugin {
	return &MissingSignupsPlugin{}
}

// Name returns the name of the plugin.
func (m *MissingSignupsPlugin) Name() string {
	return MissingSignupsPluginName
}

// Validate validates whether or not we should execute MissingSignupsPlugin on an incoming Discord message.
func (m *MissingSignupsPlugin) Validate(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	return message.Content == "!missingsignups" && message.ChannelID == discord.AdminNotificationsChannelID
}

// Fetches all signup channels from the Terror server.
// We define a signup channel as follows:
// 1. The channel is a child of the Events category channel
// 2. The channel name contains the word "signup"
func getSignupChannels(session *discordgo.Session) ([]*discordgo.Channel, error) {
	channels, err := session.GuildChannels(discord.GuildID)
	if err != nil {
		return nil, err
	}

	signupChannels := make([]*discordgo.Channel, 0)
	for _, channel := range channels {
		if channel == nil {
			continue
		}
		if channel.ParentID == EventsCategoryChannelID && strings.Contains(channel.Name, "signup") {
			signupChannels = append(signupChannels, channel)
		}
	}

	return signupChannels, nil
}

// Fetches the signup message from a signup channel.
// A signup message is defined as a message that has both the ✅ and ❌ reactions.
// For a given signup channel, there should only ever be one signup message.
func getSignupMessage(session *discordgo.Session, channel *discordgo.Channel) (*discordgo.Message, error) {
	messages, err := session.ChannelMessages(channel.ID, 100, "", "", "")
	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		if message.Reactions == nil || len(message.Reactions) == 0 {
			continue
		}

		hasYesEmoji := false
		hasNoEmoji := false

		for _, reaction := range message.Reactions {
			if reaction.Emoji.Name == YesEmoji {
				hasYesEmoji = true
			}
			if reaction.Emoji.Name == NoEmoji {
				hasNoEmoji = true
			}
		}

		if hasYesEmoji && hasNoEmoji {
			// We found the signup message
			return message, nil
		}
	}

	return nil, errors.New("could not find signup message for channel " + channel.Name)
}

func getAllTerrorMembers(session *discordgo.Session) ([]*discordgo.Member, error) {
	members, err := session.GuildMembers(discord.GuildID, "", 1000)
	if err != nil {
		return nil, err
	}

	terrorMembers := make([]*discordgo.Member, 0)
	for _, member := range members {
		if member == nil {
			continue
		}
		rank, _ := memberlist.GetDiscordMemberClanRank(member)
		if rank == nil {
			continue
		}

		terrorMembers = append(terrorMembers, member)
	}

	return terrorMembers, nil
}

func getSignedUpMembers(session *discordgo.Session, signupMessage *discordgo.Message) []*discordgo.Member {
	memberChan := make(chan *discordgo.Member)

	var wg sync.WaitGroup
	fetchMember := func(user *discordgo.User) {
		defer wg.Done()
		if user == nil {
			return
		}

		member, err := session.GuildMember(discord.GuildID, user.ID)
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to get guild member (%s): %v", user.Username, err.Error()))
			return
		}

		memberChan <- member
	}

	processVotes := func(userVotes []*discordgo.User) {
		for _, user := range userVotes {
			wg.Add(1)
			go fetchMember(user)
		}
	}

	yesUserVotes, _ := session.MessageReactions(signupMessage.ChannelID, signupMessage.ID, YesEmoji, 100, "", "")
	noUserVotes, _ := session.MessageReactions(signupMessage.ChannelID, signupMessage.ID, NoEmoji, 100, "", "")

	processVotes(yesUserVotes)
	processVotes(noUserVotes)

	go func() {
		wg.Wait()
		close(memberChan)
	}()

	// Collect all members
	signedUpMembers := make([]*discordgo.Member, 0)
	for member := range memberChan {
		signedUpMembers = append(signedUpMembers, member)
	}

	return signedUpMembers
}

func buildChunkedMessageContent(missingMembers []*discordgo.Member) []string {
	messageContent := make([]string, 0)
	currentMessage := ""
	for _, member := range missingMembers {
		if len(currentMessage) > 1800 {
			messageContent = append(messageContent, currentMessage)
			currentMessage = ""
		}
		currentMessage += fmt.Sprintf("%s ", member.User.Mention())
	}
	messageContent = append(messageContent, currentMessage)
	return messageContent
}

func processSignupChannel(session *discordgo.Session, channel *discordgo.Channel) {
	signupMessage, err := getSignupMessage(session, channel)
	if err != nil {
		return
	}

	signedUpMembers := getSignedUpMembers(session, signupMessage)
	allTerrorMembers, err := getAllTerrorMembers(session)
	if err != nil {
		return
	}

	missingMembers := make([]*discordgo.Member, 0)
	for _, member := range allTerrorMembers {
		found := false
		for _, signedUpMember := range signedUpMembers {
			if member.User.ID == signedUpMember.User.ID {
				found = true
				break
			}
		}

		if !found {
			missingMembers = append(missingMembers, member)
		}
	}

	session.ChannelMessageSend(discord.AdminNotificationsChannelID, fmt.Sprintf("Missing signups for channel %s", channel.Name))
	messages := buildChunkedMessageContent(missingMembers)

	for _, msg := range messages {
		_, err := session.ChannelMessageSend(discord.AdminNotificationsChannelID, msg)
		if err != nil {
			fmt.Println("Failed to emit message: ", err)
		}
	}
}

// Execute executes MissingSignupsPlugin on an incoming Discord message.
func (m *MissingSignupsPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) error {
	signupChannels, err := getSignupChannels(session)
	if err != nil {
		return err
	}

	if len(signupChannels) > 0 {
		session.ChannelMessageSendReply(message.ChannelID, "Processing signup channels. This will take a second...", message.Reference())
	}

	for _, channel := range signupChannels {
		go processSignupChannel(session, channel)
	}

	return nil
}
