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
	populationThreshold int
	timeWindow          int
}

type World struct {
	WorldNumber     int  `json:"world_number"`
	WorldPopulation int  `json:"world_population"`
	IsPVP           bool `json:"is_pvp"`
	Members         bool `json:"members"`
}

var worldsMap = map[int]World{}

func NewWorldTracker(populationThreshold int, timeWindow int) *WorldTracker {
	return &WorldTracker{
		populationThreshold: populationThreshold,
		timeWindow:          timeWindow,
	}
}

func (w *WorldTracker) PollAndCompare() {
	c := colly.NewCollector()

	c.OnHTML(".server-list__body", func(el *colly.HTMLElement) {
		currentWorldsMap := map[int]World{}
		fmt.Println("Polling RuneScape world population...")

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
			if populationDifference < w.populationThreshold {
				// The population difference is less than the threshold, so we don't care.
				return
			}

			if isIncrease {
				fmt.Printf("World %d has increased by %d players.\n", world.WorldNumber, int(populationDifference))
			} else {
				fmt.Printf("World %d has decreased by %d players.\n", world.WorldNumber, int(populationDifference))
			}

			worldsMap = currentWorldsMap
		})
	})

	c.Visit(WorldPopulationURL)
}
