package xptracker

import (
	"log"
	"time"

	"github.com/google/uuid"
	memberlistentity "github.com/joeydotdev/corgi-discord-bot/memberlist"
	hiscores "github.com/joeydotdev/osrs-hiscores"
)

// XpTable is a map of skills to their xp.
type XpTable = map[string]int64

type Participant struct {
	// Name is the name of the participant.
	Name string `json:"name"`
	// RuneScapeName is the runescape name of the participant.
	RuneScapeName string `json:"runescape_name"`
	// InitialXpTable is the initial xp state of the participant.
	InitialXpTable XpTable `json:"xp_table"`
	// XpGainedTable is the xp gained by the participant after event has been concluded.
	XpGainedTable XpTable `json:"xp_gained_table"`
}

type XpTrackerEvent struct {
	// Uuid is the uuid of the event.
	Uuid string `json:"uuid"`
	// Name is the name of the event.
	Name string `json:"name"`
	// IsActive is whether or not the given event is active.
	IsActive bool `json:"is_active"`
	// InitialMembersXpList is a list of members and their initial combat xps.
	Participants []Participant `json:"participants"`
	// StartDate is the start date of the event.
	StartDate string `json:"start_date"`
	// EndDate is the end date of the event.
	EndDate string `json:"end_date"`
}

// NewXpTrackerEvent creates a new xp tracker event.
func NewXpTrackerEvent(name string, members []memberlistentity.Member) *XpTrackerEvent {
	hiscores := hiscores.NewHiscores()
	participants := []Participant{}

	for _, v := range members {
		attackXp, err := hiscores.GetPlayerSkillXp(v.RuneScapeName, "attack")
		strengthXp, err := hiscores.GetPlayerSkillXp(v.RuneScapeName, "strength")
		defenceXp, err := hiscores.GetPlayerSkillXp(v.RuneScapeName, "defence")
		rangedXp, err := hiscores.GetPlayerSkillXp(v.RuneScapeName, "ranged")
		magicXp, err := hiscores.GetPlayerSkillXp(v.RuneScapeName, "magic")
		hitpointsXp, err := hiscores.GetPlayerSkillXp(v.RuneScapeName, "hitpoints")
		if err != nil {
			log.Printf(err.Error())
			continue
		}

		participants = append(participants, Participant{
			Name:          v.Name,
			RuneScapeName: v.RuneScapeName,
			InitialXpTable: XpTable{
				"attack":    attackXp,
				"strength":  strengthXp,
				"defence":   defenceXp,
				"ranged":    rangedXp,
				"magic":     magicXp,
				"hitpoints": hitpointsXp,
			},
		})
	}

	return &XpTrackerEvent{
		Uuid:         uuid.New().String(),
		Name:         name,
		IsActive:     true,
		Participants: participants,
		StartDate:    time.Now().Format(time.RFC3339),
		EndDate:      "",
	}
}

// GetParticipantCount returns the number of participants in the event.
func (x *XpTrackerEvent) GetParticipantCount() int {
	return len(x.Participants)
}

// GetEventDuration returns the duration of the event.
func (x *XpTrackerEvent) GetEventDuration() string {
	start, err := time.Parse(time.RFC3339, x.StartDate)
	if err != nil {
		log.Printf(err.Error())
		return ""
	}

	end, err := time.Parse(time.RFC3339, x.EndDate)
	if err != nil {
		log.Printf(err.Error())
		return ""
	}

	return end.Sub(start).String()
}

// GetParticipantXpGain returns the xp gained by a participant.
func (x *XpTrackerEvent) GetParticipantXpGain(participantName string) (XpTable, error) {
	var selectedParticipant Participant
	for _, v := range x.Participants {
		if v.Name == participantName {
			selectedParticipant = v
			break
		}
	}

	if selectedParticipant.Name == "" {
		return nil, nil
	}

	hiscores := hiscores.NewHiscores()
	attackXp, err := hiscores.GetPlayerSkillXp(selectedParticipant.RuneScapeName, "attack")
	strengthXp, err := hiscores.GetPlayerSkillXp(selectedParticipant.RuneScapeName, "strength")
	defenceXp, err := hiscores.GetPlayerSkillXp(selectedParticipant.RuneScapeName, "defence")
	rangedXp, err := hiscores.GetPlayerSkillXp(selectedParticipant.RuneScapeName, "ranged")
	magicXp, err := hiscores.GetPlayerSkillXp(selectedParticipant.RuneScapeName, "magic")
	hitpointsXp, err := hiscores.GetPlayerSkillXp(selectedParticipant.RuneScapeName, "hitpoints")
	if err != nil {
		return nil, err
	}

	return XpTable{
		"attack":    attackXp - selectedParticipant.InitialXpTable["attack"],
		"strength":  strengthXp - selectedParticipant.InitialXpTable["strength"],
		"defence":   defenceXp - selectedParticipant.InitialXpTable["defence"],
		"ranged":    rangedXp - selectedParticipant.InitialXpTable["ranged"],
		"magic":     magicXp - selectedParticipant.InitialXpTable["magic"],
		"hitpoints": hitpointsXp - selectedParticipant.InitialXpTable["hitpoints"],
	}, nil
}

// EndEvent ends the event.
func (x *XpTrackerEvent) EndEvent() {
	x.IsActive = false
	x.EndDate = time.Now().Format(time.RFC3339)

	for i, v := range x.Participants {
		xpGained, err := x.GetParticipantXpGain(v.Name)
		if err != nil {
			log.Printf(err.Error())
			continue
		}

		x.Participants[i].XpGainedTable = xpGained
	}
}
