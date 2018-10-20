// Package main provides ...
package main

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/duminghui/go-tipservice/db"
)

var pieAutoScanWG sync.WaitGroup

var pieAutoScanStopChan = make(chan struct{})

func pieAutoScanStart() {
	ticker := time.NewTicker(60 * time.Second)
	pieAutoScanWG.Add(1)
	go func(ticker *time.Ticker) {
		defer func() {
			ticker.Stop()
			pieAutoScanWG.Done()
		}()
		for {
			select {
			case <-ticker.C:
				pieAutoList, err := dbGuild.PieAutoProcessLists(10)
				if err != nil {
					log.Error("Pie Auto Scan error:", err)
					continue
				}
				if len(pieAutoList) == 0 {
					// log.Info("Pie Auto List Empyt...")
					continue
				}
				pieAutoListSend(pieAutoList)
			case <-pieAutoScanStopChan:
				log.Infof("Pie Auto Scan ticker stop")
				return
			}
		}
	}(ticker)
	log.Info("Pie Auto Scan start")
}

func pieAutoListSend(list []*db.PieAuto) {
	for _, pieAuto := range list {
		pieSend(pieAuto)
	}
	if len(list) > 0 {
		pieAutoList, err := dbGuild.PieAutoProcessLists(10)
		if err != nil {
			log.Error("Pie Auto Scan error:", err)
			return
		}
		if len(pieAutoList) == 0 {
			return
		}
		pieAutoListSend(pieAutoList)
	}
}

type pieAutoReceiverGenerator struct {
	pieAutoTaskID string
	guildID       string
	channelID     string
	pieerID       string
	roleID        string
	isOnlineUser  bool
}

func (r *pieAutoReceiverGenerator) Receivers() ([]*discordgo.User, error) {
	receivers := make([]*discordgo.User, 0)
	s := discordSession
	guildID := r.guildID
	guild, err := guild(s, guildID)
	if err != nil {
		return nil, err
	}
	if r.roleID != "" {
		_, err := role(s, guildID, r.roleID)
		if err != nil {
			return receivers, nil
		}
	}
	guildName := guild.Name
	gc := guildConfigs.gc(guildID)
	excludeRoles := strings.Join(gc.ExcludeRoles, "|")
	log.Infof("PieAutoTask #%s will send to [%s#%s] %d members", r.pieAutoTaskID, guildName, guildID, len(guild.Members))
	for _, member := range guild.Members {
		userID := member.User.ID
		switch {
		case member.User.Bot:
			fallthrough
		case userID == r.pieerID:
			continue
		}
		userName := member.User.Username
		userPermission, err := userChannelPermissions(s, userID, r.channelID)
		if err != nil {
			log.Errorf("Pie Auto taks get premission Error:%s[%s()][%s(%s)]", err, userName, guildName, guildID)
			continue
		}
		if (userPermission & discordgo.PermissionReadMessages) != discordgo.PermissionReadMessages {
			continue
		}
		isInExcludeRoles := false
		isInRoles := false
		for _, role := range member.Roles {
			if strings.Contains(excludeRoles, role) {
				isInExcludeRoles = true
				break
			}
			if r.roleID == role {
				isInRoles = true
			}
		}

		if r.roleID == "" {
			isInRoles = true
		}

		if isInExcludeRoles {
			continue
		}
		if !isInRoles {
			continue
		}
		if !r.isOnlineUser {
			receivers = append(receivers, member.User)
			continue
		}
		isOnline := false
		presence, err := presence(s, guild.ID, userID)
		if err != nil {
			log.Errorf("PieAuto get Presence Error:%s[%s(%s)][%s(%s)]", err, member.User.Username, userID, guild.Name, guild.ID)
			continue
		}
		if presence.Status == discordgo.StatusOnline || presence.Status == discordgo.StatusIdle {
			isOnline = true
		}
		if isOnline {
			receivers = append(receivers, member.User)
		}
	}
	return receivers, nil
}

func pieSend(pieAuto *db.PieAuto) {
	symbol := symbol(pieAuto.Symbol)
	guildID := pieAuto.GuildID
	channelID := pieAuto.ChannelID
	roleID := pieAuto.RoleID
	isOnlineUser := pieAuto.IsOnlineUser
	pieerID := pieAuto.UserID
	pieerName := pieAuto.UserName
	sendAmount := pieAuto.Amount
	log.Infof("[%s]PieAuto Task #%s Send %s(%s):#%s(%s) times:#%d", symbol, pieAuto.ID.Hex(), pieAuto.GuildName, guildID, pieAuto.ChannelName, channelID, pieAuto.RunnedTimes)
	generator := &pieAutoReceiverGenerator{
		guildID:      guildID,
		channelID:    channelID,
		pieerID:      pieerID,
		roleID:       roleID,
		isOnlineUser: isOnlineUser,
	}
	pie := &pie{
		symbol:            symbol,
		userID:            pieerID,
		userName:          pieerName,
		amount:            sendAmount,
		receiverGenerator: generator,
	}
	report, err := pie.pie()
	pieAutoID := pieAuto.ID.Hex()
	if err != nil {
		log.Errorf("[%s]PieAuto Error:%s", symbol, err)
		switch err {
		case errPieNoSymbol,
			errPieAmountMin,
			errPieUserNotExists:
			dbGuild.PieAutoRemove(pieAuto.UserID, pieAutoID)
			return
		case errPieNotEnoughAmount,
			errPieNotEnoughEachAmount:
			pieer := report.pieer
			err := dbGuild.PieAutoRemove(pieAuto.UserID, pieAutoID)
			if err != nil {
				log.Errorf("[%s]PieAutoRemoveError:[%s]", symbol, err)
			}
			msg := msgFromTmpl("pieAutoTaskNoAmountRemoveInfo", tmplValueMap{
				"AutoPieID":         pieAutoID,
				"Amount":            pieer.Amount,
				"UnconfirmedAmount": pieer.UnconfirmedAmount,
				"Symbol":            symbol,
			})
			pieAutoEmbed := pieAuto2Embed(discordSession, pieAuto)
			directMessageComplx(discordSession, pieAuto.UserID, msg, pieAutoEmbed)
			return
		case errPieNoReceiver:
		}
	}
	dbGuild.PieAutoUpTimes(pieAutoID)
	if err != nil && err == errPieNoReceiver {
		return
	}
	eachMsgReceiverNum := piebotConfig.Discord.EachPieMsgReceiversLimit
	receiverCount := report.receiverCount
	receivers := report.receivers
	receiversMap := make(map[int][]string)
	for i, receiver := range receivers {
		//msg index
		index := int(math.Floor(float64(i) / float64(eachMsgReceiverNum)))
		receiversMap[index] = append(receiversMap[index], receiver.Mention())
	}

	tmplValue := tmplValueMap{
		"UserMention":   fmt.Sprintf("<@%s>", pieerID),
		"Amount":        sendAmount,
		"Symbol":        symbol,
		"ReceiverCount": receiverCount,
	}
	sendCountMsg := msgFromTmpl("pieSendCountHint", tmplValue)
	discordSession.ChannelMessageSend(channelID, sendCountMsg)

	coinInfo := coinInfos[string(symbol)]
	roleName := ""
	if roleID != "" {
		role, err := role(discordSession, guildID, roleID)
		if err == nil {
			roleName = role.Name
		}
	}
	for _, receivers := range receiversMap {
		msg := msgFromTmpl("pieSuccess", tmplValueMap{
			"CoinName":      fmt.Sprintf("%s Auto", coinInfo.Name),
			"AmountEach":    report.eachAmount,
			"Symbol":        symbol,
			"Receivers":     receivers,
			"ReceiverCount": receiverCount,
			// "ShowAllPeople": receiverCount > eachMsgReceiverNum,
			"ShowAllPeople": true,
			"RoleName":      roleName,
		})
		discordSession.ChannelMessageSend(channelID, msg)
	}
	log.Infof("[%s]%s(%s) Pie Auto %s to %d peoples in [%s(%s) #%s(%s)]", symbol, pieerName, pieerID, sendAmount, receiverCount, pieAuto.GuildName, guildID, pieAuto.ChannelName, pieAuto.ChannelID)
}

func pieAutoScanStop() {
	pieAutoScanStopChan <- struct{}{}
}

func pieAutoScanWait() {
	pieAutoScanWG.Wait()
}
