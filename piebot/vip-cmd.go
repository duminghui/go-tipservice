// Package main provides ...
package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/duminghui/go-tipservice/db"
)

func registerSymbolVipCmds() {
	registerSymbolCmd("vip", true, false, true, cmdVipHandler)
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

func (p *guildSymbolPresenter) userPoints2embed(userPoints *db.VipUserPoints, userName string) *discordgo.MessageEmbed {
	showPoints := userPoints.Points
	embedAuthor := embedAuthor(p.coinInfo.Name, p.coinInfo.Website, "")
	embedThumbnail := embedThumbnail(p.coinInfo.IconURL)
	title := fmt.Sprintf("%s's VIP Points", userName)
	fields := embedFields(2)
	pointsField := fmt.Sprintf("%d", showPoints)
	roleField := fmt.Sprintf("@%s", userPoints.RoleName)
	fields = append(fields, mef("Points", pointsField, true))
	fields = append(fields, mef("VIP Role", roleField, true))
	embed := embed(&embedInfo{
		title: title,
		color: 0x00ff00,
	}, embedAuthor, embedThumbnail, nil, nil, fields)
	return embed
}
