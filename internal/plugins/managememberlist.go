package plugins

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	memberlistentity "github.com/joeydotdev/corgi-discord-bot/internal/memberlist"
)

const (
	ManageMemberlistPluginName = "ManageMemberlistCommand"
)

type ManageMemberlistPlugin struct{}

var _memberlist *memberlistentity.Memberlist

func init() {
	_memberlist = memberlistentity.NewMemberlist()
}

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

func (m *ManageMemberlistPlugin) isValidOperation(operation string) bool {
	return operation == "add" || operation == "remove" || operation == "update"
}

func getDiscordAndRuneScapeName(segments []string) (string, string, error) {
	if len(segments) < 2 {
		return "", "", TooFewArgumentsError
	}

	discriminatorIndex := -1
	for i, segment := range segments {
		if strings.Contains(segment, "#") {
			discriminatorIndex = i
			break
		}
	}

	if discriminatorIndex == -1 {
		return "", "", NoDiscordUsernameAndDiscriminatorError
	}

	discordName := strings.Join(segments[:discriminatorIndex+1], " ")
	runescapeName := strings.Join(segments[discriminatorIndex+1:], " ")

	return discordName, runescapeName, nil
}

func (m *ManageMemberlistPlugin) add(segments []string) error {
	return nil
}

func (m *ManageMemberlistPlugin) remove(segments []string) error {
	return nil
}

// Execute executes ManageMemberlistPlugin on an incoming Discord message.
func (m *ManageMemberlistPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) error {
	segments := strings.Split(message.Content, " ")
	if len(segments) < 2 {
		// list members
		members := _memberlist.GetMembers()
		memberString := ""
		for _, member := range members {
			memberString += member.Name + " - " + member.Accounts.LPC + "\n"
		}

		session.ChannelMessageSendReply(message.ChannelID, memberString, message.Reference())
		return nil
	}

	operation := segments[1]
	if !m.isValidOperation(operation) {
		return InvalidOperationError
	}

	args := segments[2:]
	var err error
	switch operation {
	case "add":
		err = m.add(args)
	case "remove":
		err = m.remove(args)
	case "update":
	}

	return err
}

func getMemberlist() *memberlistentity.Memberlist {
	return _memberlist
}
