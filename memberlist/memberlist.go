package memberlist

import (
	"github.com/joeydotdev/corgi-discord-bot/storage"
)

// Member is a member of the clan.
type Member struct {
	// Name is the name of the member.
	Name string `json:"name"`
	// Rank is the rank of the member.
	Rank string `json:"rank"`
	// DiscordID is the Discord ID of the member.
	DiscordID string `json:"discord_id"`
	// RuneScape Name is the RuneScape name of the member.
	RuneScapeName string `json:"runescape_name"`
}

type Memberlist struct {
	// Members is a list of members.
	Members []Member `json:"members"`
}

// NewMemberlist creates a new memberlist.
func NewMemberlist() *Memberlist {
	m := &Memberlist{}
	m.hydrate()
	return m
}

// hydrate hydrates the memberlist from the data store.
func (m *Memberlist) hydrate() error {
	err := storage.DownloadJSON("memberlist.json", m)
	return err
}

// sync syncs the memberlist to the data store.
func (m *Memberlist) sync() error {
	err := storage.UploadJSON("memberlist.json", m)
	return err
}

// AddMember adds a member to the memberlist.
func (m *Memberlist) AddMember(member Member) {
	m.Members = append(m.Members, member)
	m.sync()
}

// AddMembers adds multiple members to the memberlist.
func (m *Memberlist) AddMembers(members []Member) {
	m.Members = append(m.Members, members...)
	m.sync()
}

// RemoveMember removes a member from the memberlist.
func (m *Memberlist) RemoveMember(member Member) {
	for i, v := range m.Members {
		if v == member {
			m.Members = append(m.Members[:i], m.Members[i+1:]...)
			m.sync()
			return
		}
	}
}

// RemoveMembers removes multiple members from the memberlist.
func (m *Memberlist) RemoveMembers(members []Member) {
	for _, member := range members {
		m.RemoveMember(member)
	}
}

// GetMemberByName gets a member from the memberlist by their name.
func (m *Memberlist) GetMemberByName(name string) *Member {
	for _, v := range m.Members {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

// GetMemberByDiscordID gets a member from the memberlist by their Discord ID.
func (m *Memberlist) GetMemberByDiscordID(discordID string) *Member {
	for _, v := range m.Members {
		if v.DiscordID == discordID {
			return &v
		}
	}
	return nil
}

// GetMemberByRuneScapeName gets a member from the memberlist by their RuneScape name.
func (m *Memberlist) GetMemberByRuneScapeName(runescapeName string) *Member {
	for _, v := range m.Members {
		if v.RuneScapeName == runescapeName {
			return &v
		}
	}
	return nil
}
