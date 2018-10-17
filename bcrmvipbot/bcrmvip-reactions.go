// Package main provides ...
package main

import (
	"github.com/bwmarrin/discordgo"
)

var (
	reactionCheck = "\U00002705"
)

func reactionAddEventHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.Emoji.ID != bcrmVipConfig.Discord.VipFlagEmojiID {
		return
	}
	log.Info("##########1")
	channelID := r.ChannelID
	channel, err := channel(s, channelID)
	if err != nil {
		return
	}
	log.Info("##########2")
	if channel.Type == discordgo.ChannelTypeDM {
		return
	}
	msgID := r.MessageID
	msg, err := message(s, channelID, msgID)
	log.Info("##########3")
	if err != nil {
		return
	}
	msgAuthor := msg.Author
	log.Info("##########4")
	if msgAuthor.Bot {
		return
	}
	log.Info("##########5")
	for _, reactions := range msg.Reactions {
		if reactions.Emoji.Name == reactionCheck && reactions.Me {
			return
		}
	}
	log.Info("##########6")
	guildID := channel.GuildID
	guild, err := guild(s, guildID)
	if err != nil {
		return
	}
	log.Info("##########8")
	opUserID := r.UserID
	opMember, err := member(s, guild.ID, opUserID)
	if opMember.User.Bot {
		return
	}
	gc := gc(guildID)
	log.Info("##########9")
	if gc == nil {
		return
	}
	log.Info("##########10")
	if !gc.isBotManager(s, guild, opUserID) {
		return
	}
	log.Info("##########11")
	channelPoints, err := dbBcrm.VipChannelPointsByChannelID(channelID)
	if err != nil {
		return
	}
	log.Info("##########12")
	userPoints, err := dbBcrm.VipUserPointsChange(msgAuthor.ID, channelPoints.Points)
	if err != nil {
		return
	}
	log.Info("##########13")
	err = setVipUserRole(s, guildID, msgAuthor.ID, userPoints)
	if err != nil {
		return
	}
	log.Info("##########14")
	err = s.MessageReactionAdd(channelID, msgID, reactionCheck)
	if err != nil {
		log.Error(err)
	}
}

func reactionRemoveEventHandler(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
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
	channelPoints, err := dbBcrm.VipChannelPointsByChannelID(channelID)
	if err != nil {
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
