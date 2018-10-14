// Package main provides ...
package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	//⏹  U+23F9
	reactionStop = "\U000023F9"
	//🔄 U+1F504
	reactionRefresh = "\U0001F504"
)

func reactionAddEventHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	userID := r.UserID
	msgID := r.MessageID
	botUserID := s.State.User.ID
	// don't process self's reaction
	if userID == botUserID {
		return
	}

	emojiName := r.Emoji.Name
	switch emojiName {
	case reactionStop,
		reactionRefresh:
	default:
		return
	}
	msg, err := message(s, r.ChannelID, msgID)
	if err != nil {
		log.Error(err)
		return
	}
	// don't process message that not bot's
	if msg.Author.ID != botUserID {
		return
	}
	switch emojiName {
	case reactionRefresh:
		pieAutoInfoRefresh(s, r.ChannelID, msgID, userID)
	case reactionStop:
		pieAutoInfoRemove(s, r.ChannelID, msgID, userID)
	}
}

func pieAutoInfoRemove(s *discordgo.Session, channelID, botMsgID, userID string) {
	pieAutoMsg, err := dbGuild.PieAutoMsg(botMsgID, userID)
	if err != nil {
		return
	}
	pieAutoID := pieAutoMsg.PieAutoID
	err = dbGuild.PieAutoRemove(userID, pieAutoID)
	if err == nil {
		dbGuild.PieAutoMsgRemove(botMsgID, pieAutoID)
		embed := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("PieAuto-Task#%s", pieAutoID),
			Description: msgFromTmpl("pieAutoIsRemoved", nil),
			Color:       0xFF0000,
		}
		channelMessageEditComplx(s, channelID, botMsgID, "", embed)
	}
	// Error:HTTP 403 Forbidden, {"code": 50003, "message": "Cannot execute action on a DM channel"}
	err = s.MessageReactionRemove(channelID, botMsgID, reactionStop, userID)
	if err != nil {
		log.Errorf("pieAutoInfoRemove:MessageReactionRemove Error:%s[cID:%s,mID:%s,uID:%s]", err, channelID, botMsgID, userID)
	}
}

func pieAutoInfoRefresh(s *discordgo.Session, channelID, botMsgID, userID string) {
	pieAutoMsg, err := dbGuild.PieAutoMsg(botMsgID, userID)
	if err != nil {
		return
	}
	pieAutoID := pieAutoMsg.PieAutoID
	pieAuto, err := dbGuild.PieAutoByID(nil, pieAutoID)
	if err != nil {
		return
	}
	if pieAuto == nil {
		embed := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("PieAuto-Task#%s", pieAutoID),
			Description: msgFromTmpl("pieAutoIsRemoved", nil),
			Color:       0xFF0000,
		}
		dbGuild.PieAutoMsgRemove(botMsgID, userID)
		_, err = channelMessageEditComplx(s, channelID, botMsgID, "", embed)
		if err != nil {
			return
		}
	} else {
		embed := pieAuto2Embed(s, pieAuto)
		_, err = channelMessageEditComplx(s, channelID, botMsgID, "", embed)
		if err != nil {
			return
		}
	}
	// Error:HTTP 403 Forbidden, {"code": 50003, "message": "Cannot execute action on a DM channel"}
	err = s.MessageReactionRemove(channelID, botMsgID, reactionRefresh, userID)
	if err != nil {
		log.Errorf("pieAutoInfoRefresh:MessageReactionRemove Error:%s[cID:%s,mID:%s,uID:%s]", err, channelID, botMsgID, userID)
	}
}

func reactionRemoveEventHandler(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	log.Infof("reactionRemove:%#v", r)
}
