package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/joeydotdev/corgi-discord-bot/internal/discord"
	"github.com/joeydotdev/corgi-discord-bot/internal/memberlist"
)

func isWebClientOrigin(status discordgo.ClientStatus) bool {
	return len(status.Web) > 0 && status.Web != discordgo.StatusOffline
}

func getDiscordMemberFromDiscordUser(session *discordgo.Session, user *discordgo.User) (*discordgo.Member, error) {
	guild, err := session.Guild(discord.GuildID)
	if err != nil {
		return nil, err
	}

	member, err := session.GuildMember(guild.ID, user.ID)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (h *Handler) PresenceUpdate(session *discordgo.Session, presenceUpdate *discordgo.PresenceUpdate) {
	if presenceUpdate == nil || presenceUpdate.User == nil {
		return
	}

	if presenceUpdate.User.Username == "" || presenceUpdate.User.ID == session.State.User.ID {
		// Ignore presence updates that are from users without usernames or from the bot itself
		return
	}

	member, err := getDiscordMemberFromDiscordUser(session, presenceUpdate.User)
	if err != nil {
		fmt.Println("Failed to get member from user: ", err)
		return
	}

	if !isWebClientOrigin(presenceUpdate.Presence.ClientStatus) {
		// Ignore presence updates that are not from the web client
		return
	}

	rank, err := memberlist.GetDiscordMemberClanRank(member)
	if err != nil && err != memberlist.ErrMemberNotInClan {
		fmt.Println("Failed to get member clan rank: ", err)
		return
	}

	if rank == nil {
		// Not in Terror
		return
	}

	msg := fmt.Sprintf("%s has connected to Discord through a web browser", presenceUpdate.User.Mention())
	_, err = session.ChannelMessageSend(discord.AdminNotificationsChannelID, msg)

	if err != nil {
		fmt.Println("Failed to send message: ", err)
	}
}
