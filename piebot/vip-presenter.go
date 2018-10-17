// Package main provides ...
package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/duminghui/go-tipservice/db"
)

func (p *guildSymbolPresenter) setVipUserRole(s *discordgo.Session, guildID, userID string, userPoints *db.VipUserPoints) {
	rolePointsList, _ := p.dbSymbol.VipRolePointsList()
	if len(rolePointsList) == 0 {
		return
	}
	roleIndex := 0
	var roleUse *discordgo.Role
	roleExist := make([]string, 0)
	for _, rolePoints := range rolePointsList {
		role, err := role(s, guildID, rolePoints.RoleID)
		if err != nil {
			log.Errorf("setVipUserRole role error:%s", err)
			continue
		}
		if userPoints.Points >= rolePoints.Points {
			roleUse = role
			roleIndex++
		}
		roleExist = append(roleExist, role.ID)
	}
	roleName := "Not VIP"
	if roleIndex > 0 {
		roleName = roleUse.Name
	}
	userPoints.RoleName = roleName
	p.dbSymbol.VipUserPointsRoleName(userID, roleName)
	member, err := member(s, guildID, userID)
	if err != nil {
		return
	}
	// log.Infof("setVipUserRole %#v", roleExist)
	// log.Infof("setVipUserRole add %#v", roleExist[:roleIndex])
	for _, role := range roleExist[:roleIndex] {
		isUserHadRole := false
		for _, userRole := range member.Roles {
			if role == userRole {
				isUserHadRole = true
				break
			}
		}
		if !isUserHadRole {
			err := s.GuildMemberRoleAdd(guildID, userID, role)
			if err != nil {
				log.Errorf("setVipUserRole Add Error:%s,gid:%s,uid:%s,role:%s#%s", err, guildID, userID, roleName, role)
			}
		}
	}
	// log.Infof("setVipUserRole del %#v", roleExist[roleIndex:])
	for _, role := range roleExist[roleIndex:] {
		isUserHadRole := false
		for _, userRole := range member.Roles {
			if role == userRole {
				isUserHadRole = true
				break
			}
		}
		if isUserHadRole {
			s.GuildMemberRoleRemove(guildID, userID, role)
			if err != nil {
				log.Errorf("setVipUserRole Remove Error:%s,gid:%s:uid:%s,role:%s#%s", err, guildID, userID, roleName, role)
			}
		}
	}
}
