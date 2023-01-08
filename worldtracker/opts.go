package worldtracker

import (
	"github.com/jessevdk/go-flags"
)

type WorldTrackerOpts struct {
	Threshold int    `short:"t" long:"threshold" description:"The threshold for the number of players to trigger a world tracker spike event" default:"12" required:"true"`
	Filter    string `short:"f" long:"filter" description:"The filter for the world tracker spike event" default:"f2p"`
	Time      int    `short:"i" long:"time" description:"The time in seconds to wait before checking for spikes again" default:"12" required:"true"`
}

func AdaptDiscordArgsIntoWorldTrackerOpts(segments []string) (*WorldTrackerOpts, error) {
	opts := &WorldTrackerOpts{}
	_, err := flags.ParseArgs(&opts, segments)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
