// Package main provides ...
package main

import "github.com/bwmarrin/discordgo"

func reactionAddEventHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	log.Infof("----------add-----------")
	log.Info("UserID:", r.UserID)
	log.Info("ChannelID:", r.ChannelID)
	log.Info("MessageID:", r.MessageID)
	log.Infof("###:%#v\n", r.Emoji)
	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Error(err)
	}
	log.Info(message.Author.Username)
}

func reactionRemoveEventHandler(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	log.Infof("reactionRemove:%#v", r)
}
