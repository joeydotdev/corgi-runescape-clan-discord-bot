package plugins

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	memberlistentity "github.com/joeydotdev/corgi-discord-bot/memberlist"
	hiscores "github.com/joeydotdev/osrs-hiscores"
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
	memberlist = memberlistentity.NewMemberlist()
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

func handleAddMember(segments []string) error {
	discordHandle, runescapeName, err := getDiscordAndRuneScapeName(segments)
	if err != nil {
		return err
	}

	member := memberlistentity.Member{
		DiscordHandle: discordHandle,
		RuneScapeName: runescapeName,
	}

	memberlist.AddMember(member)
	return nil
}

func handleUpdateMember(segments []string) error {
	discordHandle, runescapeName, err := getDiscordAndRuneScapeName(segments)
	if err != nil {
		return err
	}
	updatedMember := memberlistentity.Member{
		DiscordHandle: discordHandle,
		RuneScapeName: runescapeName,
	}

	err = memberlist.UpdateMemberByDiscordHandle(discordHandle, updatedMember)
	return err
}

func handleRemoveMember(segments []string) error {
	discordHandle, _, err := getDiscordAndRuneScapeName(segments)
	if err != nil {
		return err
	}
	err = memberlist.RemoveMemberByDiscordHandle(discordHandle)
	return err
}

func filterInvalidRSNs(members []memberlistentity.Member) []memberlistentity.Member {
	hiscores := hiscores.NewHiscores()
	invalidMembers := []memberlistentity.Member{}

	for _, member := range members {
		overallLevel, err := hiscores.GetPlayerSkillLevel(member.RuneScapeName, "overall")
		isValidRSN := err == nil && overallLevel > 0
		if !isValidRSN {
			invalidMembers = append(invalidMembers, member)
		}
	}

	return invalidMembers
}

// Execute executes ManageMemberlistPlugin on an incoming Discord message.
func (m *ManageMemberlistPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) error {
	segments := strings.Split(message.Content, " ")
	if len(segments) < 2 {
		// list members
		members := memberlist.GetMembers()
		memberString := ""
		for _, member := range members {
			memberString += member.DiscordHandle + " - " + member.RuneScapeName + "\n"
		}

		session.ChannelMessageSendReply(message.ChannelID, memberString, message.Reference())
		return nil
	}

	operation := segments[1]
	if !isValidOperation(operation) {
		return InvalidOperationError
	}

	session.MessageReactionAdd(message.ChannelID, message.ID, "🟦")

	args := segments[2:]
	var err error
	switch operation {
	case "add":
		err = handleAddMember(args)
	case "remove":
		err = handleRemoveMember(args)
	case "update":
		err = handleUpdateMember(args)
	}

	session.MessageReactionRemove(message.ChannelID, message.ID, "🟦", "@me")

	if err != nil {
		session.MessageReactionAdd(message.ChannelID, message.ID, "❌")
		return err
	}

	session.MessageReactionAdd(message.ChannelID, message.ID, "✅")

	return nil
}
