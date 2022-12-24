package plugins

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	memberlistentity "github.com/joeydotdev/corgi-discord-bot/memberlist"
)

const (
	ManageMemberlistPluginName = "ManageMemberlistCommand"
)

var TooFewArgumentsError error = errors.New("Too few arguments")
var InvalidOperationError error = errors.New("Invalid operation. Valid operations are: add, remove, update")
var NoDiscordUsernameAndDiscriminatorError error = errors.New("No Discord username and discriminator provided.")

type ManageMemberlistPlugin struct{}

var memberlist *memberlistentity.Memberlist

func init() {
	memberlist = &memberlistentity.Memberlist{}
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

func isValidOperation(operation string) bool {
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

func handleAddMember(discordHandle, runescapeName string) error {
	member := memberlistentity.Member{
		DiscordHandle: discordHandle,
		RuneScapeName: runescapeName,
	}

	memberlist.AddMember(member)
	return nil
}

func handleUpdateMember(discordHandle, runescapeName string) error {
	updatedMember := memberlistentity.Member{
		DiscordHandle: discordHandle,
		RuneScapeName: runescapeName,
	}

	memberlist.UpdateMemberByDiscordHandle(discordHandle, updatedMember)
	return nil
}

// Execute executes ManageMemberlistPlugin on an incoming Discord message.
func (m *ManageMemberlistPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) error {
	segments := strings.Split(message.Content, " ")
	if len(segments) < 2 {
		session.ChannelMessageSendReply(message.ChannelID, "Please provide an operation.", message.Reference())
		return TooFewArgumentsError
	}

	operation := segments[1]
	if !isValidOperation(operation) {
		session.ChannelMessageSendReply(message.ChannelID, "Invalid operation.", message.Reference())
		return InvalidOperationError
	}

	discordHandle, runescapeName, err := getDiscordAndRuneScapeName(segments[2:])
	if err != nil {
		session.ChannelMessageSendReply(message.ChannelID, "Please provide a Discord username and a RuneScape name.", message.Reference())
		return err
	}

	switch operation {
	case "add":
		err = handleAddMember(discordHandle, runescapeName)
	case "remove":
		// TODO
	case "update":
		// TODO
	}

	return nil
}
