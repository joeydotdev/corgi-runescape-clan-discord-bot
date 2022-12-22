package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joeydotdev/corgi-discord-bot/worldtracker"
)

var Token string

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
	fmt.Printf("Token: %s\n", Token)
}

func main() {
	sig := make(chan os.Signal, 1)
	timeWindowInSeconds := 12
	ticker := time.NewTicker(time.Duration(timeWindowInSeconds) * time.Second)
	worldTracker := worldtracker.NewWorldTracker(5, timeWindowInSeconds)
	worldTracker.PollAndCompare()
	for {
		select {
		case <-ticker.C:
			worldTracker.PollAndCompare()
		case <-sig:
			return
		}
	}
}
