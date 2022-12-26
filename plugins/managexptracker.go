package plugins

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joeydotdev/corgi-discord-bot/xptracker"
)

const (
	ManageXpTrackerPluginName = "XpTrackerPlugin"
)

type ManageXpTrackerPlugin struct{}

// activeXpTrackerEvent is the currently active tracker event.
var activeXpTrackerEvent *xptracker.XpTrackerEvent

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
	return strings.HasPrefix(message.Content, "!xptracker")
}

func (m *ManageXpTrackerPlugin) isValidOperation(operation string) bool {
	return operation == "start" || operation == "stop" || operation == "status"
}

func (m *ManageXpTrackerPlugin) start(args []string, session *discordgo.Session, message *discordgo.MessageCreate) error {
	if activeXpTrackerEvent != nil {
		return ActiveOngoingEventError
	}

	if len(args) < 1 {
		return TooFewArgumentsError
	}

	name := strings.Join(args, " ")
	members := getMemberlist().GetMembers()
	activeXpTrackerEvent = xptracker.NewXpTrackerEvent(name, members)
	_, err := session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Successfully started event. Use `!xptracker status %s` to track the event.", activeXpTrackerEvent.Uuid))
	return err
}

func (m *ManageXpTrackerPlugin) stop(session *discordgo.Session, message *discordgo.MessageCreate) error {
	if activeXpTrackerEvent == nil {
		return NoEventError
	}

	activeXpTrackerEvent.EndEvent()
	_, err := session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Successfully ended event. Use `!xptracker status %s` to see the results.", activeXpTrackerEvent.Uuid))
	return err
}

func (m *ManageXpTrackerPlugin) status(args []string, session *discordgo.Session, message *discordgo.MessageCreate) error {
	var targetEvent *xptracker.XpTrackerEvent
	var err error

	uuid := args[0]
	if len(uuid) == 0 && activeXpTrackerEvent == nil {
		return NoEventError
	}

	if len(uuid) == 0 {
		targetEvent = activeXpTrackerEvent
	} else {
		targetEvent, err = xptracker.GetXpTrackerEventByUUID(uuid)
		if err != nil {
			return err
		}
	}

	if targetEvent == nil {
		return NoEventError
	}

	_, err = session.ChannelMessageSend(message.ChannelID, fmt.Sprintf(`
Event Name: %s
Event UUID: %s
Event Started: %s
Event Ended: %s
Event Participants: %d
		`, targetEvent.Name, targetEvent.Uuid, targetEvent.StartDate, targetEvent.EndDate, len(targetEvent.Participants)))

	return err
}

// Execute executes ManageXpTrackerPlugin on an incoming Discord message.
func (m *ManageXpTrackerPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) error {
	segments := strings.Split(message.Content, " ")
	if len(segments) < 2 {
		return TooFewArgumentsError
	}

	operation := segments[1]
	if !m.isValidOperation(operation) {
		return InvalidOperationError
	}
	args := segments[2:]

	var err error
	switch operation {
	case "start":
		err = m.start(args, session, message)
	case "stop":
		err = m.stop(session, message)
	case "status":
		err = m.status(args, session, message)
	}

	return err
}
