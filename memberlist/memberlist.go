package memberlist

import (
	"errors"

	hiscores "github.com/joeydotdev/osrs-hiscores"
)

type RuneScapeAccounts struct {
	LPC  string `json:"lpc"`
	XLPC string `json:"xlpc"`
}

// Member is a member of the clan.
type Member struct {
	// Uuid is the UUID of the member.
	Uuid string `json:"uuid"`
	// Name is the name of the member.
	Name string `json:"name"`
	// Rank is the rank of the member.
	Rank string `json:"rank"`
	// DiscordHandle is the Dicscord username and discriminator of the member.
	DiscordID string `json:"discord_id"`
	// TeamSpeakID is the TeamSpeak ID of the member.
	TeamSpeakID string `json:"teamspeak_id"`
	// RuneScapeAccounts is a list of the member's RuneScape accounts.
	Accounts RuneScapeAccounts `json:"runescape_accounts"`
}

type Memberlist struct {
	// Members is a list of members.
	Members []Member `json:"members"`
}

var DuplicateInMemberlistError error = errors.New("Member already exists in memberlist. Try updating instead.")

// NewMemberlist creates a new memberlist.
func NewMemberlist() *Memberlist {
	m := &Memberlist{
		Members: []Member{},
	}
	m.hydrate()

	return m
}

// hydrate hydrates the memberlist from the data store.
func (m *Memberlist) hydrate() error {
	resp, err := GetMemberlistSheet()
	if err != nil {
		return err
	}

	for _, v := range resp.Values {
		if len(v) < 5 {
			continue
		}

		member := Member{
			Uuid:        v[0].(string),
			Name:        v[1].(string),
			DiscordID:   v[2].(string),
			TeamSpeakID: v[3].(string),
			Accounts: RuneScapeAccounts{
				XLPC: v[4].(string),
				LPC:  v[5].(string),
			},
			Rank: v[6].(string),
		}

		m.Members = append(m.Members, member)
	}

	return nil
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
func (m *Memberlist) GetMemberByDiscordID(discordId string) *Member {
	for _, v := range m.Members {
		if v.DiscordID == discordId {
			return &v
		}
	}
	return nil
}

// GetMemberByRuneScapeName gets a member from the memberlist by their RuneScape name.
func (m *Memberlist) GetMemberByRuneScapeName(runescapeName string) *Member {
	for _, v := range m.Members {
		if v.Accounts.LPC == runescapeName || v.Accounts.XLPC == runescapeName {
			return &v
		}
	}
	return nil
}

// GetMembersWithInvalidXLPCRSNs gets a list of members with invalid XLPC RSNs.
func (m *Memberlist) GetMembersWithInvalidXLPCRSNs() []Member {
	hiscores := hiscores.NewHiscores()
	var members []Member

	for _, v := range m.Members {
		if len(v.Accounts.XLPC) == 0 {
			members = append(members, v)
			continue
		}

		overallLevel, err := hiscores.GetPlayerSkillLevel(v.Accounts.XLPC, "overall")
		if err != nil || overallLevel < 0 {
			members = append(members, v)
		}
	}

	return members
}

// GetMembersWithInvalidLPCRSNs gets a list of members with invalid LPC RSNs.
func (m *Memberlist) GetMembersWithInvalidLPCRSNs() []Member {
	hiscores := hiscores.NewHiscores()
	var members []Member

	for _, v := range m.Members {
		if len(v.Accounts.LPC) == 0 {
			members = append(members, v)
			continue
		}

		overallLevel, err := hiscores.GetPlayerSkillLevel(v.Accounts.LPC, "overall")
		if err != nil || overallLevel < 0 {
			members = append(members, v)
		}
	}

	return members
}

func (m *Memberlist) GetMembers() []Member {
	return m.Members
}
