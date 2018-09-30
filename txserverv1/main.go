package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/duminghui/go-tipservice/config"
	"github.com/duminghui/go-tipservice/db"
	"github.com/duminghui/go-tipservice/dbrpcmanager"
	"github.com/duminghui/go-util/ulog"
	"github.com/duminghui/go-util/umgo"
	"github.com/sevlyar/go-daemon"
	"github.com/sirupsen/logrus"
)

type TxInfo struct {
	Symbol string `json:"symbol"`
	TxID   string `json:"txid"`
}

var log = logrus.New()

var dbrpcmanagers = make(map[string]*dbrpcmanager.DBRpcManager)
var coinConfigs = make(map[string]*config.CoinInfo)

func init() {
	config, err := config.New("./config.json")
	if err != nil {
		logrus.Fatalf("Read config file error: %s", err)
	}
	logTmp, err := ulog.New(config.Log)
	if err != nil {
		logrus.Fatalln("Init Log Error:", err)
	}
	log = logTmp
	mgoSession, err := umgo.NewSession(config.Mongodb)
	if err != nil {
		log.Fatalln("Init Mongodb Error:", err)
	}
	for k, v := range config.Infos {
		dbrpcmanagers[k] = &dbrpcmanager.DBRpcManager{}
		dbrpcmanagers[k].Init(log, mgoSession, v)
		coinConfigs[k] = v
	}
}

func panicHelper() {
	if p := recover(); p != nil {
		var buf [4096]byte
		n := runtime.Stack(buf[:], false)
		log.Fatalf("%s", buf[:n])
	}
}

var (
	cmdFlag = flag.String("s", "", `send signal to the daemon
			stop - fast shutdown`)
)

// func main() {
// 	http.HandleFunc("/", handler)
// 	log.Fatal(http.ListenAndServe("localhost:8085", nil))
// }

func main() {

	flag.Parse()
	daemon.AddCommand(daemon.StringFlag(cmdFlag, "stop"), syscall.SIGTERM, termHandler)

	cntxt := &daemon.Context{
		PidFileName: "pid",
		PidFilePerm: 0644,
		LogFileName: "log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[txserver]"},
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalln("Unable send signal to the daemon:", err)
		}
		daemon.SendCommands(d)
		return
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln("Reborn Error:", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()
	log.Info("-----------------------")
	log.Info("daemon started")

	for _, m := range dbrpcmanagers {
		m.RPC.Start()
	}

	go terminateHelper()
	go httpServerProcess()
	go httpServerStopHelper()
	go txQueueHandler()
	txProcessHandler()

	err = daemon.ServeSignals()
	if err != nil {
		log.Info("daemon terminate Error:", err)
	}
	log.Println("daemon terminated")
}

var server *http.Server

func httpServerProcess() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	server = &http.Server{
		Addr:         "127.0.0.1:8085",
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}

	log.Info("Starting Http Server")
	err := server.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			log.Info("Http Server closed")
		} else {
			log.Fatal("Http Server closed unexpected: ", err)
		}
	}
	stopHTTPServerDone <- struct{}{}
}

var (
	stopHTTPServer     = make(chan struct{})
	stopHTTPServerDone = make(chan struct{})
)

func httpServerStopHelper() {
	defer panicHelper()
	<-stopHTTPServer
	log.Info("Http server start shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	// err := server.Shutdown(nil)
	if err != nil {
		log.Info("Http server shutdown error:", err)
	}
	log.Info("Http server shutdowning...")
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

func handler(w http.ResponseWriter, r *http.Request) {
	defer panicHelper()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Server Read Body Error", err)
		fmt.Fprintf(w, "Server Read Body Error: %s\n", err)
		return
	}
	var txInfo TxInfo
	err = json.Unmarshal(bytes, &txInfo)
	if err != nil {
		log.Errorf("Unmarshal Error: %s (err:%s)", bytes, err)
		fmt.Fprintf(w, "Unmarshal Error: %s (err:%s)", bytes, err)
		return
	}
	if _, ok := dbrpcmanagers[txInfo.Symbol]; !ok {
		log.Errorf("Dont had '%s''s config", txInfo.Symbol)
		fmt.Fprintf(w, "Dont had '%s''s config\n", txInfo.Symbol)
		return
	}
	db.SaveTxProcess(txInfo.Symbol, txInfo.TxID)
	// txQueue <- &txInfo
	fmt.Fprintln(w, "accept success")
}

var tickerStop = make(chan struct{})

func txProcessHandler() {
	ticker := time.NewTicker(30 * time.Second)
	go func(ticker *time.Ticker) {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				txs, err := db.TxProcessInfos()
				if err == nil {
					for _, tx := range txs {
						err := db.TxProcessUpdate(nil, tx.Symbol, tx.TxID, 0, 10)
						if err != nil {
							continue
						}
						txQueue <- &TxInfo{
							Symbol: tx.Symbol,
							TxID:   tx.TxID,
						}
					}
				}
			case <-tickerStop:
				log.Info("Ticker Stop")
				return
			}
		}
	}(ticker)
}

var (
	txQueue      = make(chan *TxInfo, 100)
	txProcessSum = make(chan struct{}, 10)
)

func txQueueHandler() {
	defer panicHelper()
	for tx := range txQueue {
		txProcessSum <- struct{}{}
		go func(tx *TxInfo) {
			defer panicHelper()
			txProcess(tx)
			<-txProcessSum
		}(tx)
	}
}

func txProcess(tx *TxInfo) {
	dbrpcmanager := dbrpcmanagers[tx.Symbol]
	txInfo, err := dbrpcmanager.RPC.GetTransaction(tx.TxID, nil)
	if err != nil {
		log.Errorf("[%s]GetTransaction Error:[%s][%s]", tx.Symbol, err, tx.TxID)
		db.TxProcessUpdate(nil, tx.Symbol, tx.TxID, 0, 30)
		return
	}
	if txInfo.Amount <= 0.0 {
		log.Infof("[%s]Ingore prcess TX Amount is < 0.0 [%s]", tx.Symbol, tx.TxID)
		db.TxProcessUpdate(nil, tx.Symbol, tx.TxID, 1, 0)
		return
	}
	coinConfig := coinConfigs[tx.Symbol]
	isConfirmed := txInfo.Confirmations >= coinConfig.MinConfirmations4Deposit
	for _, txDetail := range txInfo.Details {
		amount := txDetail.Amount
		address := txDetail.Address
		err = dbrpcmanager.DB.Deposit(txDetail.Address, tx.TxID, amount, isConfirmed)
		if err != nil {
			log.Errorf("[%s]Deposit Error:[%s][%s][%s][%f]", tx.Symbol, err, tx.TxID, address, amount)
			continue
		}
		log.Infof("[%s]Process Deposit Success:[%s][%s][%f][confirmed:%d]", tx.Symbol, tx.TxID, address, amount, txInfo.Confirmations)
	}
}

func termHandler(sig os.Signal) error {
	// log.Info("terminating...")
	log.Info("terminating...")
	tickerStop <- struct{}{}
	stopHTTPServer <- struct{}{}
	<-stopHTTPServerDone
	stop <- struct{}{}
	<-done
	return daemon.ErrStop
}
