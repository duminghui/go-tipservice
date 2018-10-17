// Package main provides ...
package main

import (
	"flag"
	"os"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	rpcclient "github.com/duminghui/go-rpcclient"
	"github.com/duminghui/go-tipservice/db"
	"github.com/duminghui/go-util/ulog"
	"github.com/duminghui/go-util/umgo"
	daemon "github.com/sevlyar/go-daemon"
	"github.com/sirupsen/logrus"
)

var (
	cmdFlag    = flag.String("s", "", "send signal to the daemon\nstop - fast shutdown")
	configFile = flag.String("c", "bcrmvipbot.json", "config file path")
)

func main() {
	flag.Parse()
	initConfigLog()
	daemon.AddCommand(daemon.StringFlag(cmdFlag, "stop"), syscall.SIGTERM, termHandler)

	cntxt := &daemon.Context{
		PidFileName: bcrmVipConfig.PidFile,
		PidFilePerm: 0644,
		LogFileName: bcrmVipConfig.Log.LogFile,
		LogFilePerm: 0640,
		WorkDir:     bcrmVipConfig.WorkDir,
		Umask:       027,
		// Args:        []string{"[piebot]"},
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			logrus.Fatalln("Unable send signal to the daemon:", err)
		}
		daemon.SendCommands(d)
		return
	}

	d, err := cntxt.Reborn()
	if err != nil {
		logrus.Fatalln("Reborn Error:", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()
	log.Info("-----------------------")
	log.Info("daemon started")

	initRunEnv()

	go terminateHelper()

	discordSession, err = discordgo.New("Bot " + bcrmVipConfig.Discord.Token)
	if err != nil {
		log.Fatalf("Createing Discrod Session Error: %s", err)
	}

	// discordSession.State.MaxMessageCount = 200
	discordSession.AddHandler(reactionAddEventHandler)
	discordSession.AddHandler(reactionRemoveEventHandler)

	err = discordSession.Open()
	if err != nil {
		log.Fatalf("Opening Discord connection error:%s", err)
	}

	log.Info("Discord Bot is now running...")

	go discordStopHelper()
	err = daemon.ServeSignals()
	if err != nil {
		log.Info("daemon terminate Error:", err)
	}
	log.Println("daemon terminated")
}

var (
	discordSession *discordgo.Session
	dgStop         = make(chan struct{})
	dgStopDone     = make(chan struct{})
)

func discordStopHelper() {
	for {
		<-dgStop
		err := discordSession.Close()
		if err != nil {
			log.Errorf("Discord Stop Error:%s", err)
		} else {
			log.Info("Discord Stop Success")
		}
		dgStopDone <- struct{}{}
		return
	}
}

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func terminateHelper() {
	func() {
		for {
			time.Sleep(time.Second)
			select {
			case <-stop:
				return
			default:
			}
		}
	}()
	done <- struct{}{}
}

func termHandler(sig os.Signal) error {
	// log.Info("terminating...")
	log.Info("terminating...")
	dgStop <- struct{}{}
	<-dgStopDone
	stop <- struct{}{}
	<-done
	return daemon.ErrStop
}

var log = logrus.New()

func initRunEnv() {
	mgoSession, err := umgo.NewSession(dbConfig)
	if err != nil {
		log.Fatalln("Init Mongodb Error:", err)
	}
	db.SetLog(log)
	db.SetSession(mgoSession)
	rpcclient.SetLog(log)
}

func initConfigLog() {
	err := readConfig(*configFile)
	if err != nil {
		logrus.Fatalf("Read config file error: %s", err)
	}
	logTmp, err := ulog.NewSingle(bcrmVipConfig.Log)
	if err != nil {
		logrus.Fatalln("Init Log Error:", err)
	}
	log = logTmp
}
