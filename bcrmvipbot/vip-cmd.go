// Package main provides ...
package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/duminghui/go-tipservice/db"
)

func registerSymbolVipCmds() {
	registerSymbolCmd("vipHelp", true, false, true, cmdVipHelpHandler)
	registerSymbolCmd("vip", true, false, true, cmdVipHandler)
	registerSymbolCmd("vipTop", true, false, true, cmdVipTopHandler)
	registerSymbolCmd("vipRoles", true, false, true, cmdVipRolesHandler)
	registerSymbolCmd("vipRolePoints", false, true, true, cmdVipRolePointsHandler)
	registerSymbolCmd("vipRoleScan", false, true, true, cmdVipRoleScanHandler)
	registerSymbolCmd("vipChannels", false, true, true, cmdVipChannelsHandler)
	registerSymbolCmd("vipChannelPoints", false, true, true, cmdVipChannelPointsHandler)
	registerSymbolCmd("vipPoints", false, true, true, cmdVipPointsHandler)
	registerSymbolCmd("vipEmoji", false, true, true, cmdVipEmojiHandler)
}

var cmdVipRoleScanHandler = (*guildSymbolPresenter).cmdVipRoleScanHandler

func (p *guildSymbolPresenter) cmdVipRoleScanHandler(parts *msgParts) {
	count, err := p.dbSymbol.VipUserPointsCount()
	if err != nil {
		parts.channelMessageSend("Bot Error")
		return
	}
	if count == 0 {
		parts.channelMessageSend("No User have VIP Points")
		return
	}
	start := 0
	size := 10
	msgID := ""
	msgTmpl := "%s VIP role scan will to process %d peoples, processed %d peoples"
	for {
		userPointsList, err := p.dbSymbol.VipUserPointsList(start, size)
		if err != nil {
			parts.channelMessageSend("Bot Error")
			return
		}
		userPointsListLen := len(userPointsList)
		if userPointsListLen == 0 {
			return
		}
		isErr := false
		for _, userPoints := range userPointsList {
			err := setVipUserRole(parts.s, parts.guild.ID, userPoints)
			if err != nil {
				isErr = true
				break
			}
		}
		if isErr {
			parts.channelMessageSend("Don't have config VIP role points")
			return
		}

		start += len(userPointsList)
		msgContent := fmt.Sprintf(msgTmpl, parts.m.Author.Mention(), count, start)
		if msgID == "" {
			msg, err := parts.channelMessageSend(msgContent)
			if err == nil {
				msgID = msg.ID
			}
		} else {
			_, err = parts.s.ChannelMessageEdit(parts.channel.ID, msgID, msgContent)
			if err != nil {
				log.Infof("[%s]VipRoleScan Error:%s", p.symbol, err)
			}
		}
	}
}

var cmdVipTopHandler = (*guildSymbolPresenter).cmdVipTopHandler

func (p *guildSymbolPresenter) cmdVipTopHandler(parts *msgParts) {
	var err error
	num := int64(0)
	if len(parts.contents) == 0 {
		num = 0
	} else {
		num, err = strconv.ParseInt(parts.contents[0], 10, 32)
		if err != nil {
			num = 0
		}
	}
	start := int(num - 1)
	if start < 0 {
		start = 0
	}
	size := 10
	userPointsList, err := p.dbSymbol.VipUserPointsList(start, size)
	embedFields := embedFields(len(userPointsList))
	end := start + 1
	desc := ""
	if len(userPointsList) == 0 {
		desc = "No VIP Peoples"
	} else {
		buf := new(bytes.Buffer)
		for i, userPoints := range userPointsList {
			member, err := member(parts.s, parts.guild.ID, userPoints.UserID)
			if err != nil {
				continue
			}
			roleName := userVipRoleName(parts.s, parts.guild.ID, userPoints)
			fieldTitle := fmt.Sprintf("**#%d %s#%s @%s Points:%d**\n", start+i+1, member.User.Username, member.User.Discriminator, roleName, userPoints.Points)
			// fieldContent := fmt.Sprintf("Points: %d | VIP Role: @%s", userPoints.Points, roleName)
			buf.WriteString(fieldTitle)
			embedFields = append(embedFields, mef(fieldTitle, " x ", false))
		}
		desc = buf.String()
		end = start + len(userPointsList)
	}
	title := fmt.Sprintf("%s VIP Leaderboard[ %d - %d ]", p.coinInfo.Name, start+1, end)
	embedInfo := &embedInfo{
		title: title,
		desc:  desc,
		color: 0x00FF00,
	}
	embedThumbnail := embedThumbnail(p.coinInfo.IconURL)
	embed := embed(embedInfo, nil, embedThumbnail, nil, nil, nil)
	_, err = parts.channelMessageSendEmbed(embed)
	if err != nil {
		log.Errorf("VipTop Error:%s", err)
	}
}

var cmdVipEmojiHandler = (*guildSymbolPresenter).cmdVipEmojiHandler

func (p *guildSymbolPresenter) cmdVipEmojiHandler(parts *msgParts) {
	userMention := parts.m.Author.Mention()
	cmdPrefix := parts.prefix
	tmplValue := &tmplValueMap{
		"IsShowUsageHint": true,
		"CmdName":         "vipEmoji",
		"UserMention":     userMention,
		"Prefix":          cmdPrefix,
	}
	usageMsg := msgFromTmpl("vipEmojiUsage", tmplValue)
	if len(parts.contents) == 0 {
		parts.channelMessageSend(usageMsg)
		return
	}
	id, name := emojiFromContent(parts.contents[0])
	if id == "" || name == "" {
		parts.channelMessageSend(usageMsg)
		return
	}
	_, err := p.dbSymbol.VipEmojiChange(id, name)
	if err != nil {
		parts.channelMessageSend("Save VIP Emoji failed")
		return
	}
	readVipEmojiFromDB()
	emoji := &discordgo.Emoji{
		ID:   vipEmoji.ID,
		Name: vipEmoji.Name,
	}
	successMsg := fmt.Sprintf("Now VIP Emoji is <:%s>", emoji.APIName())
	parts.channelMessageSend(successMsg)
}

var cmdVipHelpHandler = (*guildSymbolPresenter).cmdVipHelpHandler

func (p *guildSymbolPresenter) cmdVipHelpHandler(parts *msgParts) {
	cmdPrefix := parts.prefix
	sbl := p.symbol
	isManager := parts.isManager
	tmplValue := &tmplValueMap{
		"IsShowUsageHint": false,
		"CmdName":         "help",
		"UserMention":     parts.m.Author.Mention(),
		"Prefix":          cmdPrefix,
		"Symbol":          sbl,
		"IsManager":       isManager,
	}
	helpInfo := msgFromTmpl("vipHelpUsage", tmplValue)
	parts.channelMessageSend(helpInfo)
}

var cmdVipHandler = (*guildSymbolPresenter).cmdVipHandler

func (p *guildSymbolPresenter) cmdVipHandler(parts *msgParts) {
	userID := parts.m.Author.ID
	userName := parts.m.Author.Username
	userMention := parts.m.Author.Mention()
	userPoints, _ := p.dbSymbol.VipUserPointsByUserID(nil, userID)
	msg := fmt.Sprintf("%s Your VIP Points:", userMention)
	embed := p.userPoints2embed(userPoints, userName)
	parts.channelMessageSendComplex(msg, embed)
}

var cmdVipRolePointsHandler = (*guildSymbolPresenter).cmdVipRolePointsHandler

func (p *guildSymbolPresenter) cmdVipRolePointsHandler(parts *msgParts) {
	contents := parts.contents
	contentCount := len(contents)
	userMention := parts.m.Author.Mention()
	cmdPrefix := parts.prefix
	tmplValue := &tmplValueMap{
		"IsShowUsageHint": true,
		"CmdName":         "vipRolePoints",
		"UserMention":     userMention,
		"Prefix":          string(cmdPrefix),
	}
	usageMsg := msgFromTmpl("vipRolePointsUsage", tmplValue)
	if contentCount == 0 {
		parts.channelMessageSend(usageMsg)
		return
	}
	if contentCount != 2 {
		parts.channelMessageSend(usageMsg)
		return
	}
	if len(parts.m.MentionRoles) != 1 {
		parts.channelMessageSend(usageMsg)
		return
	}
	roleID := parts.m.MentionRoles[0]
	points, err := strconv.ParseInt(contents[1], 10, 64)
	if err != nil {
		parts.channelMessageSend(usageMsg)
		return
	}
	if points < 0 {
		parts.channelMessageSend(usageMsg)
		return
	}
	err = p.dbSymbol.VipRolePointsSet(roleID, points)
	if err != nil {
		return
	}
	msgEmbed := p.vipRolePointsListEmbed(parts.s, parts.guild.ID)
	msg := fmt.Sprintf("%s now role's upgrade points is:", userMention)
	parts.channelMessageSendComplex(msg, msgEmbed)
}

var cmdVipRolesHandler = (*guildSymbolPresenter).cmdVipRolesHandler

func (p *guildSymbolPresenter) cmdVipRolesHandler(parts *msgParts) {
	msgEmbed := p.vipRolePointsListEmbed(parts.s, parts.guild.ID)
	parts.channelMessageSendEmbed(msgEmbed)
}

func (p *guildSymbolPresenter) vipRolePointsListEmbed(s *discordgo.Session, guildID string) *discordgo.MessageEmbed {
	readVipRolePointsFromDB()
	embed := new(discordgo.MessageEmbed)
	embed.Title = "Role's Points"
	embed.Color = 0x00FF00
	coinInfo := p.coinInfo
	embed.Author = embedAuthor(coinInfo.Name, coinInfo.Website, "")
	embed.Thumbnail = embedThumbnail(coinInfo.IconURL)
	rolePointCount := len(vipRolePointses)
	if rolePointCount == 0 {
		embed.Description = "No Role's Points"
		return embed
	}
	fields := embedFields(len(vipRolePointses))
	for _, rolePoint := range vipRolePointses {
		role, _ := role(s, guildID, rolePoint.RoleID)
		if role == nil {
			continue
		}
		fields = append(fields, mef(fmt.Sprintf("@%s", role.Name), fmt.Sprintf("%d", rolePoint.Points), true))
	}
	embed.Fields = fields
	return embed
}

var cmdVipChannelsHandler = (*guildSymbolPresenter).cmdVipChannelsHandler

func (p *guildSymbolPresenter) cmdVipChannelsHandler(parts *msgParts) {
	msgEmbed := p.vipChannelPointsListEmbed(parts.s)
	parts.channelMessageSendEmbed(msgEmbed)
}

func (p *guildSymbolPresenter) vipChannelPointsListEmbed(s *discordgo.Session) *discordgo.MessageEmbed {
	readVipChannelPointsFromDB()
	embed := new(discordgo.MessageEmbed)
	embed.Title = "Channel's Points"
	embed.Color = 0x00FF00
	coinInfo := p.coinInfo
	embed.Author = embedAuthor(coinInfo.Name, coinInfo.Website, "")
	embed.Thumbnail = embedThumbnail(coinInfo.IconURL)
	channelPointsCount := len(vipChannelPointses)
	if channelPointsCount == 0 {
		embed.Description = "No Channel's Points"
		return embed
	}
	fields := embedFields(len(vipChannelPointses))
	for _, channelPoints := range vipChannelPointses {
		channel, err := channel(s, channelPoints.ChannelID)
		if err != nil {
			continue
		}
		fields = append(fields, mef(fmt.Sprintf("#%s", channel.Name), fmt.Sprintf("%d", channelPoints.Points), true))
	}
	embed.Fields = fields
	return embed
}

var cmdVipChannelPointsHandler = (*guildSymbolPresenter).cmdVipChannelPointsHandler

func (p *guildSymbolPresenter) cmdVipChannelPointsHandler(parts *msgParts) {
	contents := parts.contents
	contentCount := len(contents)
	userMention := parts.m.Author.Mention()
	cmdPrefix := parts.prefix
	tmplValue := &tmplValueMap{
		"IsShowUsageHint": true,
		"CmdName":         "vipChannelPoints",
		"UserMention":     userMention,
		"Prefix":          string(cmdPrefix),
	}
	usageMsg := msgFromTmpl("vipChannelPointsUsage", tmplValue)
	if contentCount == 0 {
		parts.channelMessageSend(usageMsg)
		return
	}
	if contentCount != 2 {
		parts.channelMessageSend(usageMsg)
		return
	}
	channelIDs := channelIDsFromContent(strings.Join(contents, "|"))
	if len(channelIDs) != 1 {
		parts.channelMessageSend(usageMsg)
		return
	}
	channelID := channelIDs[0]
	_, err := channel(parts.s, channelID)
	if err != nil {
		parts.channelMessageSend(usageMsg)
		return
	}
	points, err := strconv.ParseInt(contents[1], 10, 64)
	if err != nil {
		parts.channelMessageSend(usageMsg)
		return
	}
	if points < 0 {
		parts.channelMessageSend(usageMsg)
		return
	}
	err = p.dbSymbol.VipChannelPointsSet(channelID, points)
	if err != nil {
		return
	}
	msgEmbed := p.vipChannelPointsListEmbed(parts.s)
	msg := fmt.Sprintf("%s now channel's give points is:", userMention)
	parts.channelMessageSendComplex(msg, msgEmbed)
}

var cmdVipPointsHandler = (*guildSymbolPresenter).cmdVipPointsHandler

func (p *guildSymbolPresenter) cmdVipPointsHandler(parts *msgParts) {
	contents := parts.contents
	contentCount := len(contents)
	userMention := parts.m.Author.Mention()
	cmdPrefix := parts.prefix
	tmplValue := &tmplValueMap{
		"IsShowUsageHint": true,
		"CmdName":         "vipPoints",
		"UserMention":     userMention,
		"Prefix":          string(cmdPrefix),
	}
	usageMsg := msgFromTmpl("vipPointsUsage", tmplValue)
	if contentCount == 0 {
		parts.channelMessageSend(usageMsg)
		return
	}
	if contentCount != 2 {
		parts.channelMessageSend(usageMsg)
		return
	}
	if len(parts.m.Mentions) != 1 {
		parts.channelMessageSend(usageMsg)
		return
	}
	user := parts.m.Mentions[0]
	if user.Bot {
		msg := fmt.Sprintf("%s can't give VIP points to a bot", userMention)
		parts.channelMessageSend(msg)
		return
	}
	userID := user.ID
	userName := user.Username

	points, err := strconv.ParseInt(contents[1], 10, 64)
	if err != nil {
		parts.channelMessageSend(usageMsg)
		return
	}
	userPoints, err := p.dbSymbol.VipUserPointsChange(userID, points)
	if err != nil {
		log.Info(err)
		return
	}
	setVipUserRole(parts.s, parts.guild.ID, userPoints)
	embed := p.userPoints2embed(userPoints, userName)
	parts.channelMessageSendComplex("", embed)
}

func (p *guildSymbolPresenter) userPoints2embed(userPoints *db.VipUserPoints, userName string) *discordgo.MessageEmbed {
	showPoints := userPoints.Points
	embedAuthor := embedAuthor(p.coinInfo.Name, p.coinInfo.Website, "")
	embedThumbnail := embedThumbnail(p.coinInfo.IconURL)
	title := fmt.Sprintf("%s's VIP Points", userName)
	fields := embedFields(2)
	pointsField := fmt.Sprintf("%d", showPoints)
	userRoleName := userVipRoleName(discordSession, p.guildID, userPoints)
	roleField := fmt.Sprintf("@%s", userRoleName)
	fields = append(fields, mef("Points", pointsField, true))
	fields = append(fields, mef("VIP Role", roleField, true))
	embed := embed(&embedInfo{
		title: title,
		color: 0x00ff00,
	}, embedAuthor, embedThumbnail, nil, nil, fields)
	return embed
}
