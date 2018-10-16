// Package main provides ...
package main

import (
	"github.com/bwmarrin/discordgo"
)

var (
	reactionVipID = "501300710839943169"
	reactionCheck = "\U00002705"
)

func reactionAddEventHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.Emoji.ID != bcrmVipConfig.Discord.VipFlagEmojiID {
		return
	}
	channelID := r.ChannelID
	channel, err := channel(s, channelID)
	if err != nil {
		return
	}
	if channel.Type == discordgo.ChannelTypeDM {
		return
	}
	guildID := channel.GuildID
	guild, err := guild(s, guildID)
	if err != nil {
		return
	}
	userID := r.UserID
	member, err := member(s, guild.ID, userID)
	if member.User.Bot {
		return
	}
	gc := gc(guildID)
	if gc == nil {
		return
	}
	if !gc.isBotManager(s, guild, userID) {
		return
	}
	msgID := r.MessageID
	msg, err := message(s, channelID, msgID)
	if err != nil {
		return
	}
	err = s.MessageReactionAdd(channelID, msgID, reactionCheck)
	if err != nil {
		log.Info(err)
	}
	reactionss := msg.Reactions
	log.Info(len(reactionss))
	for _, v := range msg.Reactions {
		log.Infof("%#v", v.Emoji)
		log.Info(v.Count)
		log.Info(v.Me)
	}
	// userScore, err := dbBcrm.VipUserPointsChange(r.UserID, 100)
	// if err != nil {
	// 	log.Info("Errorrrr:", err)
	// }
	// log.Infof("User:%#v", userScore)
}

func reactionRemoveEventHandler(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	log.Infof("reactionRemove:%#v", r)
	// userID := r.UserID
	// msgID := r.MessageID
	// botUserID := s.State.User.ID
	// // don't process self's reaction
	// if userID == botUserID {
	// 	return
	// }
}
