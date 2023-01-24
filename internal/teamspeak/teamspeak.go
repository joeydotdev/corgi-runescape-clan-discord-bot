package teamspeak

import (
	"log"
	"os"

	"github.com/multiplay/go-ts3"
)

var TS3Client *ts3.Client

func GetClientsInEventChannels() []*ts3.OnlineClient {
	EVENT_CHANNEL_IDS := []int{3, 4, 17, 5, 6}

	onlineUsers, _ := TS3Client.Server.ClientList()
	attendees := make([]*ts3.OnlineClient, 0, len(onlineUsers))

	clients, err := TS3Client.Server.ClientList()
	if err != nil {
		panic(err)
	}

	for _, c := range clients {
		for _, id := range EVENT_CHANNEL_IDS {
			if c.ChannelID == id {
				attendees = append(attendees, c)
			}
		}
	}

	return attendees
}

func init() {
	serverAddress := os.Getenv("TS3_SERVER_QUERY_ADDRESS")
	serverQueryUsername := os.Getenv("TS3_SERVER_QUERY_USERNAME")
	serverQueryPassword := os.Getenv("TS3_SERVER_QUERY_PASSWORD")

	if serverAddress == "" || serverQueryUsername == "" || serverQueryPassword == "" {
		panic("TS3_SERVER_QUERY_ADDRESS, TS3_SERVER_QUERY_USERNAME, and TS3_SERVER_QUERY_PASSWORD must be set")
	}

	c, err := ts3.NewClient(serverAddress)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	if err := c.Login(serverQueryUsername, serverQueryPassword); err != nil {
		panic(err)
	}

	if err := c.Use(1); err != nil {
		panic(err)
	}

	log.Println("Connected to teamspeak server")

	TS3Client = c
}
