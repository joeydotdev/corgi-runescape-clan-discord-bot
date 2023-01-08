package worldtracker

import (
	"testing"
)

type TestCase struct {
	input    []string
	expected []string
}

func TestAdaptDiscordArgsIntoWorldTrackerOpts(t *testing.T) {
	t.Parallel()

	args := []string{"--threshold", "12", "--filter", "f2p", "--time", "12"}

	ret, err := AdaptDiscordArgsIntoWorldTrackerOpts(args)
	if err != nil {
		t.Error(err)
	}
	if ret.Filter != "f2p" {
		t.Errorf("Expected filter to be f2p, got %s", ret.Filter)
	}
	if ret.Threshold != 12 {
		t.Errorf("Expected threshold to be 12, got %d", ret.Threshold)
	}
	if ret.Time != 12 {
		t.Errorf("Expected time to be 12, got %d", ret.Time)
	}
}
