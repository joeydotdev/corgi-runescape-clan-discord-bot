package memberlist

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

type Rank struct {
	Name   string
	RoleID string
}

var RANKS []Rank = []Rank{
	{
		Name:   "Leader",
		RoleID: "692876249118539817",
	},
	{
		Name:   "High Council",
		RoleID: "692879184024043540",
	},
	{
		Name:   "Council",
		RoleID: "692879285417017375",
	},
	{
		Name:   "Leadership",
		RoleID: "817499802148274226",
	},
	{
		Name:   "Officer",
		RoleID: "692879600380018699",
	},
	{
		Name:   "Legend",
		RoleID: "692879942777569312",
	},
	{
		Name:   "Old School",
		RoleID: "1024119526801023006",
	},
	{
		Name:   "Veteran",
		RoleID: "692880299855446106",
	},
	{
		Name:   "Advanced",
		RoleID: "699354924185682031",
	},
	{
		Name:   "Member",
		RoleID: "692880390440091659",
	},
	{
		Name:   "Applicant",
		RoleID: "773216677423874048",
	},
}

var ErrMemberNotInClan error = errors.New("discord member is not in the clan")

// GetDiscordMemberRank returns the clan rank of a Discord member. If the member is not in the clan, nil is returned.
func GetDiscordMemberClanRank(member *discordgo.Member) (*Rank, error) {
	if member == nil {
		return nil, errors.New("nil member")
	}

	for _, v := range RANKS {
		for _, roleID := range member.Roles {
			if roleID == v.RoleID {
				return &v, nil
			}
		}
	}

	return nil, ErrMemberNotInClan
}
