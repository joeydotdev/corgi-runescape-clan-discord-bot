package worldtracker

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

const (
	// WorldPopulationURL is the URL for RuneScape's world population page.
	WorldPopulationURL = "https://oldschool.runescape.com/slu"
)

type WorldTracker struct {
	PopulationThreshold int
	TimeWindow          int
}

type WorldTrackerConfiguration struct {
	PopulationThreshold int
	TimeWindow          int
}

type WorldTrackerSpikeEvent struct {
	WorldNumber      int
	PlayerSpikeCount int
	Members          bool
	IsPVP            bool
}

type World struct {
	WorldNumber     int  `json:"world_number"`
	WorldPopulation int  `json:"world_population"`
	IsPVP           bool `json:"is_pvp"`
	Members         bool `json:"members"`
}

var worldsMap = map[int]World{}

// NewWorldTracker creates a new WorldTracker.
func NewWorldTracker(config *WorldTrackerConfiguration) *WorldTracker {
	return &WorldTracker{
		PopulationThreshold: config.PopulationThreshold,
		TimeWindow:          config.TimeWindow,
	}
}

// PollAndCompare polls the RuneScape world population page and compares the current world data to the previous world data.
func (w *WorldTracker) PollAndCompare() []WorldTrackerSpikeEvent {
	c := colly.NewCollector()
	events := []WorldTrackerSpikeEvent{}

	c.OnHTML(".server-list__body", func(el *colly.HTMLElement) {
		currentWorldsMap := map[int]World{}

		el.ForEach("tr.server-list__row", func(_ int, el *colly.HTMLElement) {
			world := World{}
			world.IsPVP = strings.Contains(el.Attr("class"), "pvp")
			el.ForEach("td", func(_ int, el *colly.HTMLElement) {
				switch el.Index {
				case 0:
					var sanitizedWorldString string
					sanitizedWorldString = strings.TrimSpace(el.Text)
					sanitizedWorldString = strings.TrimPrefix(sanitizedWorldString, "OldSchool")
					sanitizedWorldString = strings.TrimPrefix(sanitizedWorldString, "Old School")
					sanitizedWorldString = strings.TrimSpace(sanitizedWorldString)
					worldNumber, err := strconv.Atoi(sanitizedWorldString)
					if err != nil {
						fmt.Println("Error parsing world number: " + sanitizedWorldString)
						break
					}
					world.WorldNumber = worldNumber
				case 1:
					var sanitizedPopulationString string
					sanitizedPopulationString = strings.TrimSpace(el.Text)
					sanitizedPopulationString = strings.TrimSuffix(sanitizedPopulationString, " players")
					worldPopulation, err := strconv.Atoi(sanitizedPopulationString)
					if err != nil {
						fmt.Printf("Error parsing population (world %d): %s\n", world.WorldNumber, sanitizedPopulationString)
						break
					}
					world.WorldPopulation = worldPopulation
				case 3:
					world.Members = strings.Contains("Members", el.Text)
				}
			})

			currentWorldsMap[world.WorldNumber] = world

			// Compare the current world data to the previous world data.
			previousWorld, ok := worldsMap[world.WorldNumber]
			if !ok {
				// The world wasn't in the previous map, so we don't have any data to compare against.
				return
			}

			populationDifference := int(math.Abs(float64(world.WorldPopulation) - float64(previousWorld.WorldPopulation)))
			isIncrease := world.WorldPopulation > previousWorld.WorldPopulation
			if populationDifference < w.PopulationThreshold {
				// The population difference is less than the threshold, so we don't care.
				return
			}

			// worldLabel := "F2P"
			// if world.Members {
			// 	worldLabel = "P2P"
			// }

			// if world.IsPVP {
			// 	worldLabel = fmt.Sprintf("%s · PVP", worldLabel)
			// }

			// if isIncrease {
			// 	fmt.Printf("World %d (%s) has increased by %d players.\n", world.WorldNumber, worldLabel, int(populationDifference))
			// } else {
			// 	fmt.Printf("World %d (%s) has decreased by %d players.\n", world.WorldNumber, worldLabel, int(populationDifference))
			// }

			spikeCount := populationDifference
			if !isIncrease {
				spikeCount = -spikeCount
			}

			events = append(events, WorldTrackerSpikeEvent{
				WorldNumber:      world.WorldNumber,
				PlayerSpikeCount: spikeCount,
				Members:          world.Members,
				IsPVP:            world.IsPVP,
			})
		})

		worldsMap = currentWorldsMap
	})

	c.Visit(WorldPopulationURL)
	return events
}
