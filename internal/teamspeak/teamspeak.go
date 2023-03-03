package teamspeak

import (
	"errors"
	"log"
	"os"

	"github.com/multiplay/go-ts3"
)

type TeamSpeakClient struct {
	c *ts3.Client
}

func (t *TeamSpeakClient) GetClientsInEventChannels() ([]*ts3.OnlineClient, error) {
	EVENT_CHANNEL_IDS := []int{3, 4, 17, 5, 6}

	attendees := make([]*ts3.OnlineClient, 0)
	clients, err := t.c.Server.ClientList()
	if err != nil {
		return nil, err
	}

	for _, c := range clients {
		for _, id := range EVENT_CHANNEL_IDS {
			if c.ChannelID == id {
				attendees = append(attendees, c)
			}
		}
	}

	return attendees, nil
}

func NewTeamSpeakClient() (*TeamSpeakClient, error) {
	serverAddress := os.Getenv("TS3_SERVER_QUERY_ADDRESS")
	serverQueryUsername := os.Getenv("TS3_SERVER_QUERY_USERNAME")
	serverQueryPassword := os.Getenv("TS3_SERVER_QUERY_PASSWORD")

	if serverAddress == "" || serverQueryUsername == "" || serverQueryPassword == "" {
		return nil, errors.New("TS3_SERVER_QUERY_ADDRESS, TS3_SERVER_QUERY_USERNAME, and TS3_SERVER_QUERY_PASSWORD must be set")
	}

	c, err := ts3.NewClient(serverAddress)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	if err := c.Login(serverQueryUsername, serverQueryPassword); err != nil {
		return nil, err
	}

	if err := c.Use(1); err != nil {
		return nil, err
	}

	log.Println("Connected to teamspeak server")
	return &TeamSpeakClient{
		c: c,
	}, nil
}
