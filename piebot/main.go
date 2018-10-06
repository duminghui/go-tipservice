// Package main provides ...
package main

import (
	"flag"
	"os"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	rpcclient "github.com/duminghui/go-rpcclient"
	"github.com/duminghui/go-tipservice/config"
	"github.com/duminghui/go-tipservice/db"
	"github.com/duminghui/go-util/ulog"
	"github.com/duminghui/go-util/umgo"
	daemon "github.com/sevlyar/go-daemon"
	"github.com/sirupsen/logrus"
)

var (
	cmdFlag    = flag.String("s", "", "send signal to the daemon\nstop - fast shutdown")
	configFile = flag.String("c", "config.json", "config file path")
)

func main() {
	flag.Parse()
	initConfigLog()
	daemon.AddCommand(daemon.StringFlag(cmdFlag, "stop"), syscall.SIGTERM, termHandler)

	cntxt := &daemon.Context{
		PidFileName: piebotConfig.PidFile,
		PidFilePerm: 0644,
		LogFileName: piebotConfig.Log.LogFile,
		LogFilePerm: 0640,
		WorkDir:     piebotConfig.WorkDir,
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

	initPresenter()

	go terminateHelper()

	discordSession, err = discordgo.New("Bot " + piebotConfig.Discord.Token)
	if err != nil {
		log.Fatalf("Createing Discrod Session Error: %s", err)
	}

	discordSession.AddHandler(messageCreate)

	err = discordSession.Open()
	if err != nil {
		log.Fatalf("Opening Discord connection error:%s", err)
	}

	log.Info("Discord Bot is now running...")

	go discordStopHelper()
	for _, p := range coinPresenters {
		p.rpc.Start()
	}
	reigsterBotCmdHandler()
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

type coinPresenter struct {
	db   *db.DB
	rpc  *rpcclient.Client
	coin *config.CoinInfo
}

var coinPresenters = make(map[symbolWrap]*coinPresenter)

func initPresenter() {
	mgoSession, err := umgo.NewSession(dbConfig)
	if err != nil {
		log.Fatalln("Init Mongodb Error:", err)
	}
	db.SetLog(log)
	db.SetSession(mgoSession)
	rpcclient.SetLog(log)
	for k, v := range coinInfos {
		sblWrap := symbolWrap(k)
		presenter := new(coinPresenter)
		presenter.db = db.New(v.Symbol, v.Database)
		presenter.rpc = rpcclient.New(v.RPC)
		presenter.coin = v
		coinPresenters[sblWrap] = presenter
	}
	initGuildConfig()
}

func initConfigLog() {
	err := readConfig(*configFile)
	if err != nil {
		logrus.Fatalf("Read config file error: %s", err)
	}
	logTmp, err := ulog.NewSingle(piebotConfig.Log)
	if err != nil {
		logrus.Fatalln("Init Log Error:", err)
	}
	log = logTmp
}
