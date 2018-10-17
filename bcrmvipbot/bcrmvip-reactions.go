// Package main provides ...
package main

import (
	"github.com/bwmarrin/discordgo"
)

var (
	reactionCheck = "\U00002705"
)

func reactionAddEventHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	vipFlagEmojiID := bcrmVipConfig.Discord.VipFlagEmojiID
	emoji := r.Emoji
	emojiAPIName := emoji.APIName()
	emojiID := r.Emoji.ID
	opUserID := r.UserID
	msgID := r.MessageID
	channelID := r.ChannelID
	channelPoints, err := dbBcrm.VipChannelPointsByChannelID(channelID)
	if err != nil {
		if emojiID == vipFlagEmojiID {
			err = s.MessageReactionRemove(channelID, msgID, emojiAPIName, opUserID)
			if err != nil {
				log.Error(err)
			}
		}
		return
	}
	rolePointsList, err := dbBcrm.VipRolePointsList()
	if len(rolePointsList) == 0 {
		if emojiID == vipFlagEmojiID {
			err = s.MessageReactionRemove(channelID, msgID, emojiAPIName, opUserID)
			if err != nil {
				log.Error(err)
			}
		}
		return
	}

	channel, err := channel(s, channelID)
	if err != nil {
		return
	}
	if channel.Type == discordgo.ChannelTypeDM {
		return
	}
	emojiName := r.Emoji.Name
	if emojiName == reactionCheck && opUserID != discordSession.State.User.ID {
		err = s.MessageReactionRemove(channelID, msgID, reactionCheck, opUserID)
		if err != nil {
			log.Error(err)
		}
		return
	}
	if emojiID != bcrmVipConfig.Discord.VipFlagEmojiID {
		return
	}
	msg, err := message(s, channelID, msgID)
	if err != nil {
		return
	}
	msgAuthor := msg.Author
	if msgAuthor.Bot {
		return
	}
	for _, reactions := range msg.Reactions {
		if reactions.Emoji.Name == reactionCheck && reactions.Me {
			return
		}
	}
	guildID := channel.GuildID
	guild, err := guild(s, guildID)
	if err != nil {
		return
	}
	opMember, err := member(s, guild.ID, opUserID)
	if opMember.User.Bot {
		return
	}
	gc := gc(guildID)
	if gc == nil {
		return
	}
	if !gc.isBotManager(s, guild, opUserID) {
		return
	}
	userPoints, err := dbBcrm.VipUserPointsChange(msgAuthor.ID, channelPoints.Points)
	if err != nil {
		return
	}
	err = setVipUserRole(s, guildID, msgAuthor.ID, userPoints)
	if err != nil {
		return
	}
	err = s.MessageReactionAdd(channelID, msgID, reactionCheck)
	if err != nil {
		log.Error(err)
	}
}

func reactionRemoveEventHandler(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	channelID := r.ChannelID
	channelPoints, err := dbBcrm.VipChannelPointsByChannelID(channelID)
	if err != nil {
		return
	}
	rolePointsList, err := dbBcrm.VipRolePointsList()
	if len(rolePointsList) == 0 {
		return
	}

	if r.Emoji.ID != bcrmVipConfig.Discord.VipFlagEmojiID {
		return
	}
	channel, err := channel(s, channelID)
	if err != nil {
		return
	}
	if channel.Type == discordgo.ChannelTypeDM {
		return
	}
	msgID := r.MessageID
	msg, err := message(s, channelID, msgID)
	if err != nil {
		return
	}
	msgAuthor := msg.Author
	if msgAuthor.Bot {
		return
	}
	isHadVipBotReaction := false
	for _, reactions := range msg.Reactions {
		if reactions.Emoji.Name == reactionCheck && reactions.Me {
			isHadVipBotReaction = true
			break
		}
	}
	if !isHadVipBotReaction {
		return
	}
	guildID := channel.GuildID
	guild, err := guild(s, guildID)
	if err != nil {
		return
	}
	opUserID := r.UserID
	opMember, err := member(s, guild.ID, opUserID)
	if opMember.User.Bot {
		return
	}
	gc := gc(guildID)
	if gc == nil {
		return
	}
	if !gc.isBotManager(s, guild, opUserID) {
		return
	}
	userPoints, err := dbBcrm.VipUserPointsChange(msgAuthor.ID, -channelPoints.Points)
	if err != nil {
		return
	}
	err = setVipUserRole(s, guildID, msgAuthor.ID, userPoints)
	if err != nil {
		return
	}
	err = s.MessageReactionRemove(channelID, msgID, reactionCheck, s.State.User.ID)
	if err != nil {
		log.Error(err)
	}

}
