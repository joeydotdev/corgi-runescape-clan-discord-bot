package plugins

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	MassPMCommandPluginName = "MassPMCommandPlugin"
	MassPMChannelID         = "1070035899221545171"
)

var (
	ROLE_IDS []string = []string{
		"692876249118539817",  // Leader
		"1070076373781205143", // PK Leader
		// "692879184024043540",  // High Council
		// "692879285417017375",  // Council
		// "817499802148274226",  // Leadership
		// "692879600380018699",  // Officer
		// "692879942777569312",  // Legend
		// "1024119526801023006", // Old School
		// "692880299855446106",  // Veteran
		// "699354924185682031",  // Advanced
		// "692880390440091659",  // Member
		// "773216677423874048",  // Applicant
	}
	EXCLUDED_USER_IDS []string = []string{
		"223169696055296011", // joey
	}
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
	// parse everything after !masspm
	segments := strings.Split(message.Content, " ")
	if len(segments) < 2 {
		return TooFewArgumentsError
	}

	// get the message to send
	messageToSend := strings.Join(segments[1:], " ")
	if len(messageToSend) == 0 {
		return TooFewArgumentsError
	}

	members, err := session.GuildMembers(message.GuildID, "", 1000)
	if err != nil {
		return err
	}

	messageDispatcher := func(member *discordgo.Member) {
		for _, excludedUserID := range EXCLUDED_USER_IDS {
			if member.User.ID == excludedUserID {
				return
			}
		}

		hasAttemptedToMessageMember := false
		for _, role := range member.Roles {
			if hasAttemptedToMessageMember {
				break
			}
			for _, roleID := range ROLE_IDS {
				if hasAttemptedToMessageMember {
					break
				}
				if role == roleID {
					channel, err := session.UserChannelCreate(member.User.ID)
					if channel == nil || err != nil {
						fmt.Println("Failed to create channel with member: ", member.User.ID)
						hasAttemptedToMessageMember = true
						continue
					}

					_, err = session.ChannelMessageSend(channel.ID, messageToSend)
					if err != nil {
						fmt.Println("Failed to send message to member: ", member.User.ID)
						fmt.Println("error: ", err)
					}

					hasAttemptedToMessageMember = true
				}
			}
		}
	}

	for _, member := range members {
		go messageDispatcher(member)
	}

	_, err = session.ChannelMessageSend(message.ChannelID, "Mass PM sent!")
	return err
}
