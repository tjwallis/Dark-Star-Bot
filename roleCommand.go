package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func roleCommand(s *discordgo.Session, m *discordgo.MessageCreate) { // Add role to someone
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Print("Error getting current channel: ", err)
	}
	if len(splitMsgLowered) > 1 && channel.GuildID != "" { // If it just isnt `$role`

		for i := 0; i < len(splitMsgLowered)-1; i++ {
			assignRole(s, m, splitMsgLowered[i+1])
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please type a role name after `$role` !")
	}
}

func assignRole(s *discordgo.Session, m *discordgo.MessageCreate, givenRole string) {
	cfg := loadConfig("roles/config.json")
	channel, err := s.Channel(m.ChannelID)
	_ = channel
	if err != nil {
		fmt.Print("Error getting current channel: ", err)
	}

	tserver := loadServer(cfg, channel.GuildID)

	roleUsed := false

	// Handles the ingame roles
	for i := 0; i < len(tserver.Roles); i++ {
		for j := 0; j < len(tserver.Roles[i].Calls); j++ {
			switch givenRole {

			case tserver.Roles[i].Calls[j]:
				roleUsed = true
				giveRole(s, m, tserver.Roles[i].Role, tserver.Roles[i].Locked, tserver.Roles[i].GroupID)
			}
		}
	}

	if !roleUsed {
		s.ChannelMessageSend(m.ChannelID, "You cannot add that role! Use `$Roles` to see all available roles")
	}
}

func giveRole(s *discordgo.Session, m *discordgo.MessageCreate, roleNeeded string, locked bool, groupID string) { // Assigns the role based off of a role needed

	currentGuild := getGuild(s, m)
	currentMember := getMember(s, m)

	if currentGuild != nil && currentMember != nil && groupID != "Vote" { // groupID != "Vote" is temporary

		tempRoleID := findRoleID(roleNeeded, currentGuild)

		hasRole := memberHasRole(currentMember, tempRoleID)

		if !hasRole { // Give that guy a role
			err := s.GuildMemberRoleAdd(currentGuild.ID, m.Author.ID, tempRoleID) // Give the role
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Unable to add role! Message <@!121105861539135490> and tell him there's a problem.")
				log.Fatal(err)
			}
			if err == nil { // Didnt want this popping up if there was an error
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s role added!", roleNeeded))
			}
		} else { // Bye bye role
			if !locked {
				err := s.GuildMemberRoleRemove(currentGuild.ID, m.Author.ID, tempRoleID) // Remove the role
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Unable to remove role! Message <@!121105861539135490> and tell him there's a problem.")
					log.Fatal(err)
				}
				if err == nil { // Didnt want this popping up if there was an error
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s role removed!", roleNeeded))
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s role removed!", roleNeeded))
			}
		}

	}

}
