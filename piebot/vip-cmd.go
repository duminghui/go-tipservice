// Package main provides ...
package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/duminghui/go-tipservice/db"
)

func registerSymbolVipCmds() {
	registerSymbolCmd("vip", true, false, true, cmdVipHandler)
	registerSymbolCmd("vipRoles", true, false, true, cmdVipRolesHandler)
	registerSymbolCmd("vipRolePoints", false, true, true, cmdVipRolePointsHandler)
	registerSymbolCmd("vipChannels", false, true, true, cmdVipChannelsHandler)
	registerSymbolCmd("vipChannelPoints", false, true, true, cmdVipChannelPointsHandler)
	registerSymbolCmd("vipPoints", false, true, true, cmdVipPointsHandler)
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
	rolePointList, err := p.dbSymbol.VipRolePointsList()
	embed := new(discordgo.MessageEmbed)
	embed.Title = "Role's Points"
	embed.Color = 0x00FF00
	coinInfo := p.coinInfo
	embed.Author = embedAuthor(coinInfo.Name, coinInfo.Website, "")
	embed.Thumbnail = embedThumbnail(coinInfo.IconURL)
	if err != nil {
		embed.Description = "Bot Error"
		return embed
	}
	rolePointCount := len(rolePointList)
	if rolePointCount == 0 {
		embed.Description = "No Role's Points"
		return embed
	}
	fields := embedFields(len(rolePointList))
	for _, rolePoint := range rolePointList {
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
	channelPointsList, err := p.dbSymbol.VipChannelPointsList()
	embed := new(discordgo.MessageEmbed)
	embed.Title = "Channel's Points"
	embed.Color = 0x00FF00
	coinInfo := p.coinInfo
	embed.Author = embedAuthor(coinInfo.Name, coinInfo.Website, "")
	embed.Thumbnail = embedThumbnail(coinInfo.IconURL)
	if err != nil {
		embed.Description = "Bot Error"
		return embed
	}
	rolePointCount := len(channelPointsList)
	if rolePointCount == 0 {
		embed.Description = "No Channel's Points"
		return embed
	}
	fields := embedFields(len(channelPointsList))
	for _, channelPoints := range channelPointsList {
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
	p.setVipUserRole(parts.s, parts.guild.ID, userID, userPoints)
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
	fields = append(fields, mef("Points", pointsField, true))
	fields = append(fields, mef("VIP Role", userPoints.RoleName, true))
	embed := embed(&embedInfo{
		title: title,
		color: 0x00ff00,
	}, embedAuthor, embedThumbnail, nil, nil, fields)
	return embed
}
