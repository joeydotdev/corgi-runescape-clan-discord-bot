package plugins

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	teamspeakentity "github.com/joeydotdev/corgi-discord-bot/internal/teamspeak"
)

const (
	AttendanceCommandPluginName = "AttendanceCommandPlugin"
)

type AttendanceCommandPlugin struct{}

var ts3client *teamspeakentity.TeamSpeakClient

// Enabled returns whether or not the AttendanceCommandPlugin is enabled.
func (a *AttendanceCommandPlugin) Enabled() bool {
	return false
}

// NewAttendanceCommandPlugin creates a new AttendanceCommandPlugin.
func NewAttendanceCommandPlugin() *AttendanceCommandPlugin {
	var err error
	ts3client, err = teamspeakentity.NewTeamSpeakClient()
	if err != nil {
		fmt.Println("Failed to create new TeamSpeak client: ", err)
		return nil
	}

	return &AttendanceCommandPlugin{}
}

// Name returns the name of the plugin.
func (a *AttendanceCommandPlugin) Name() string {
	return AttendanceCommandPluginName
}

// Validate validates whether or not we should execute AttendanceCommandPlugin on an incoming Discord message.
func (a *AttendanceCommandPlugin) Validate(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	return strings.HasPrefix(message.Content, "!attendance")
}

// Execute executes AttendanceCommandPlugin on an incoming Discord message.
func (a *AttendanceCommandPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) error {
	segments := strings.Split(message.Content, " ")
	if len(segments) < 2 {
		return TooFewArgumentsError
	}

	attendanceSnapshotName := strings.Join(segments[1:], " ")
	if len(attendanceSnapshotName) == 0 {
		return TooFewArgumentsError
	}

	messageString := fmt.Sprintf("Attendance for **%s**:\n", attendanceSnapshotName)
	clients, err := ts3client.GetClientsInEventChannels()
	if err != nil {
		return err
	}

	for _, client := range clients {
		messageString += fmt.Sprintf("%s\n", client.Nickname)
	}

	_, err = session.ChannelMessageSend(message.ChannelID, messageString)
	return err
}
