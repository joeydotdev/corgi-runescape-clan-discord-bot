package plugins

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joeydotdev/corgi-discord-bot/worldtracker"
)

const (
	ManageWorldTrackerPluginName = "ManageWorldTrackerPlugin"
	WorldTrackerChannelID        = "1042819494441124034"
)

type ManageWorldTrackerPlugin struct{}

const MINIMUM_TIME_WINDOW = 10
const MINIMUM_POPULATION_THRESHOLD = 8
const POSITIVE_COLOR = 0x86efac
const NEGATIVE_COLOR = 0x9f1239

var WorldTrackerAlreadyRunningError error = errors.New("World tracker is already running. Stop the current instance before starting a new one.")
var WorldTrackerMinimumTimeWindowError error = errors.New(fmt.Sprintf("Time window must be greater than %d seconds.", MINIMUM_TIME_WINDOW))
var WorldTrackerMinimumPopulationThresholdError error = errors.New(fmt.Sprintf("Population threshold must be greater than %d.", MINIMUM_POPULATION_THRESHOLD))

var activeWorldTrackerInstance *worldtracker.WorldTracker
var activeWorldTrackerKillSwitch chan bool

// Enabled returns whether or not the ManageWorldTrackerPlugin is enabled.
func (m *ManageWorldTrackerPlugin) Enabled() bool {
	return true
}

// isValidOperation returns whether or not the operation is valid.
func (m *ManageWorldTrackerPlugin) isValidOperation(operation string) bool {
	return operation == "start" || operation == "stop" || operation == "help"
}

// NewManageWorldTrackerPlugin creates a new ManageWorldTrackerPlugin.
func NewManageWorldTrackerPlugin() *ManageWorldTrackerPlugin {
	return &ManageWorldTrackerPlugin{}
}

// Name returns the name of the plugin.
func (m *ManageWorldTrackerPlugin) Name() string {
	return ManageWorldTrackerPluginName
}

// Validate validates whether or not we should execute ManageWorldTrackerPlugin on an incoming Discord message.
func (m *ManageWorldTrackerPlugin) Validate(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	return strings.HasPrefix(message.Content, "!worldtracker") && message.ChannelID == WorldTrackerChannelID
}

// sendTrackerEventMessages sends messages to Discord for each world tracker event.
func (m *ManageWorldTrackerPlugin) sendTrackerEventMessages(session *discordgo.Session, message *discordgo.MessageCreate, events []worldtracker.WorldTrackerSpikeEvent) {
	for _, event := range events {
		isIncrease := event.PlayerSpikeCount > 0
		worldLabel := "F2P"
		if event.Members {
			worldLabel = "P2P"
		}

		var m string
		var color int
		if isIncrease {
			m = fmt.Sprintf("World %d (%s) has increased by %d players.\n", event.WorldNumber, worldLabel, event.PlayerSpikeCount)
			color = POSITIVE_COLOR
		} else {
			m = fmt.Sprintf("World %d (%s) has decreased by %d players.\n", event.WorldNumber, worldLabel, event.PlayerSpikeCount)
			color = NEGATIVE_COLOR
		}

		session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
			Description: m,
			Color:       color,
		})
	}
}

// startTrackerJob starts a job that polls the world tracker and sends messages to Discord when a world's population changes.
func (m *ManageWorldTrackerPlugin) startTrackerJob(session *discordgo.Session, message *discordgo.MessageCreate) chan bool {
	stop := make(chan bool)
	go func() {
		for {
			events := activeWorldTrackerInstance.PollAndCompare()
			m.sendTrackerEventMessages(session, message, events)
			select {
			case <-time.After(time.Duration(activeWorldTrackerInstance.TimeWindow) * time.Second):
			case <-stop:
				return
			}
		}
	}()

	return stop
}

func (m *ManageWorldTrackerPlugin) start(opts *worldtracker.WorldTrackerOpts, session *discordgo.Session, message *discordgo.MessageCreate) error {
	if activeWorldTrackerInstance != nil {
		return WorldTrackerAlreadyRunningError
	}

	if opts.Time < MINIMUM_TIME_WINDOW {
		return WorldTrackerMinimumTimeWindowError
	}

	if opts.Threshold < MINIMUM_POPULATION_THRESHOLD {
		return WorldTrackerMinimumPopulationThresholdError
	}

	activeWorldTrackerInstance = worldtracker.NewWorldTracker(&worldtracker.WorldTrackerConfiguration{
		PopulationThreshold: opts.Threshold,
		TimeWindow:          opts.Time,
	})

	activeWorldTrackerKillSwitch = m.startTrackerJob(session, message)

	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("World tracker has started with population threshold of %d players and time window of %d seconds.", opts.Threshold, opts.Time))
	return nil
}

func (m *ManageWorldTrackerPlugin) stop(session *discordgo.Session, message *discordgo.MessageCreate) error {
	if activeWorldTrackerInstance == nil || activeWorldTrackerKillSwitch == nil {
		return errors.New("World tracker is not running. Start the world tracker before stopping it.")
	}
	activeWorldTrackerInstance = nil
	activeWorldTrackerKillSwitch <- true
	session.ChannelMessageSend(message.ChannelID, "World tracker has stopped.")
	return nil
}

func (m *ManageWorldTrackerPlugin) help(session *discordgo.Session, message *discordgo.MessageCreate) error {
	_, err := session.ChannelMessageSend(message.ChannelID, "Usage: `!worldtracker start <population threshold> <time window in seconds>`")
	return err
}

// Execute executes ManageWorldTrackerPlugin on an incoming Discord message.
func (m *ManageWorldTrackerPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) error {
	segments := strings.Split(message.Content, " ")
	if len(segments) < 2 {
		return TooFewArgumentsError
	}
	operation := segments[1]
	if !m.isValidOperation(operation) {
		return errors.New(fmt.Sprintf("Invalid operation: %s - valid operations are start, stop, help", operation))
	}
	args := segments[2:]
	opts, err := worldtracker.AdaptDiscordArgsIntoWorldTrackerOpts(args)
	if err != nil {
		return err
	}
	switch operation {
	case "start":
		err = m.start(opts, session, message)
	case "stop":
		err = m.stop(session, message)
	case "help":
		err = m.help(session, message)
	}

	return err
}
