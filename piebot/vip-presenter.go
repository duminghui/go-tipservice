// Package main provides ...
package main

import (
	"github.com/bwmarrin/discordgo"
)

func (p *guildSymbolPresenter) setVipUserRole(s *discordgo.Session, guildID, userID string, userPoints int64) {
	rolePointsList, _ := p.dbSymbol.VipRolePointsList()
	if len(rolePointsList) == 0 {
		return
	}
	// roleIndex := -1
	// for i, rolePoints := range rolePointsList {
	// 	if userPoints < rolePoints.Points {
	// 		break
	// 	}
	// 	roleIndex = i
	// }
}
