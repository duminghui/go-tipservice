// Package main provides ...
package main

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/duminghui/go-tipservice/db"
)

func userVipRoleName(s *discordgo.Session, guildID string, userPoints *db.VipUserPoints) string {
	roleName := "Not VIP"
	if len(vipRolePointses) == 0 {
		return roleName
	}

	for _, rolePoints := range vipRolePointses {
		role, err := role(s, guildID, rolePoints.RoleID)
		if err != nil {
			log.Errorf("userVipRoleName error:%s", err)
			continue
		}
		if userPoints.Points < rolePoints.Points {
			break
		}
		roleName = role.Name
	}
	return roleName
}

func setVipUserRole(s *discordgo.Session, guildID, userID string, userPoints *db.VipUserPoints) error {
	if len(vipRolePointses) == 0 {
		return errors.New("No VipRolePointsList")
	}
	roleIndex := 0
	var roleUse *discordgo.Role
	roleExist := make([]string, 0)
	for _, rolePoints := range vipRolePointses {
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
	dbBcrm.VipUserPointsRoleName(userID, roleName)
	member, err := member(s, guildID, userID)
	if err != nil {
		return err
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
	return nil
}

var vipEmoji *db.VipEmoji

func readVipEmojiFromDB() {
	emoji, err := dbBcrm.VipEmoji()
	if err == nil {
		vipEmoji = emoji
	}
}

var vipChannelPointsMap map[string]*db.VipChannelPoints
var vipChannelPointses []*db.VipChannelPoints

func readVipChannelPointsFromDB() {
	vipChannelPointsMap = make(map[string]*db.VipChannelPoints)
	vipChannelPointses = make([]*db.VipChannelPoints, 0)
	channelPointses, err := dbBcrm.VipChannelPointsList()
	if err != nil {
		return
	}
	vipChannelPointses = channelPointses
	for _, channelPoints := range channelPointses {
		vipChannelPointsMap[channelPoints.ChannelID] = channelPoints
	}
}

var vipRolePointses []*db.VipRolePoints

func readVipRolePointsFromDB() {
	vipRolePointses = make([]*db.VipRolePoints, 0)
	rolePointses, err := dbBcrm.VipRolePointsList()
	if err != nil {
		return
	}
	vipRolePointses = rolePointses
}
