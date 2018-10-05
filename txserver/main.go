package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	rpcclient "github.com/duminghui/go-rpcclient"
	"github.com/duminghui/go-tipservice/config"
	"github.com/duminghui/go-tipservice/db"
	"github.com/duminghui/go-util/ulog"
	"github.com/duminghui/go-util/umgo"
	"github.com/sevlyar/go-daemon"
	"github.com/sirupsen/logrus"
)

type txInfo struct {
	Symbol string `json:"symbol"`
	TxID   string `json:"txid"`
}

var log = logrus.New()

type processPresenter struct {
	symbol                   string
	scanTxDir                string
	depositMinConfirmactions int64
	db                       *db.DB
	rpc                      *rpcclient.Client
	wg                       *sync.WaitGroup
	fileScanStop             chan struct{}
	txProcessStop            chan struct{}
}

func (p *processPresenter) start() {
	p.rpc.Start()
	p.fileScanStart()
	p.txProcessDBStart()
}

func (p *processPresenter) stop() {
	p.fileScanStop <- struct{}{}
	p.txProcessStop <- struct{}{}
}

func (p *processPresenter) fileScanStart() {
	ticker := time.NewTicker(60 * time.Second)
	p.wg.Add(1)
	go func(ticker *time.Ticker) {
		defer func() {
			ticker.Stop()
			p.wg.Done()
		}()
		for {
			select {
			case <-ticker.C:
				p.fileScan()
			case <-p.fileScanStop:
				log.Infof("[%s]File scan ticker stop ", p.symbol)
				return
			}
		}
	}(ticker)
	log.Infof("[%s]File Scan Start", p.symbol)
}

func (p *processPresenter) processFile(filepath string) {
	newFilepath := filepath + ".process"
	err := os.Rename(filepath, newFilepath)
	if err != nil {
		log.Errorf("[%s]Move File Error:%s[%s]", p.symbol, err, filepath)
		return
	}
	file, err := os.Open(newFilepath)
	if err != nil {
		log.Errorf("[%s]Load File Error:%s[%s]", p.symbol, err, filepath)
		return
	}
	isSuccess := false
	defer func() {
		if isSuccess {
			err := os.Remove(newFilepath)
			if err != nil {
				log.Errorf("[%s]Remove Process File Error:%s[%s]", p.symbol, err, newFilepath)
			}
		} else {
			err := os.Rename(newFilepath, filepath)
			if err != nil {
				log.Errorf("[%s]Move File Error:%s[%s]", p.symbol, err, newFilepath)
				return
			}
		}
	}()
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		symbolTxID := strings.Split(scanner.Text(), ",")
		if len(symbolTxID) != 2 {
			continue
		}
		symbol := symbolTxID[0]
		txID := symbolTxID[1]
		txQueue <- &txInfo{
			Symbol: symbol,
			TxID:   txID,
		}
	}
	isSuccess = true
}

func (p *processPresenter) fileScan() {
	filepath.Walk(p.scanTxDir,
		func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !isCanProcessFile(info.Name()) {
				return nil
			}
			log.Infof("[%s]Will process file:%s", p.symbol, path)
			p.processFile(path)
			log.Infof("[%s]Process file done:%s", p.symbol, path)
			return nil
		})
}

const timeLayout = "200601021504"

var locationLocal, _ = time.LoadLocation("Local")

func isCanProcessFile(filename string) bool {
	filenameParts := strings.Split(filename, ".")
	if len(filenameParts) != 2 {
		// log.Errorf("FileName Len Error:%s", filename)
		return false
	}
	fileTime, err := time.ParseInLocation(timeLayout, filenameParts[0], locationLocal)
	if err != nil {
		log.Errorf("FileName Error::%s,[%s]", err, filename)
		return false
	}
	minuteRange, err := strconv.ParseInt(filenameParts[1], 10, 32)
	if err != nil {
		log.Errorf("FileName Error::%s,[%s]", err, filename)
		return false
	}
	processTime := fileTime.Add(time.Duration(minuteRange) * time.Minute)
	return time.Now().After(processTime)
}

func (p *processPresenter) txProcessDBStart() {
	ticker := time.NewTicker(30 * time.Second)
	p.wg.Add(1)
	go func(ticker *time.Ticker) {
		defer func() {
			ticker.Stop()
			p.wg.Done()
		}()
		for {
			select {
			case <-ticker.C:
				txs, err := p.db.TxProcessInfos()
				if err != nil {
					log.Errorf("[%s]txProcessDB Error:%s", p.symbol, err)
					continue
				}
				if len(txs) == 0 {
					log.Infof("[%s]TxProcesss Empty...", p.symbol)
					continue
				}
				_, err = p.rpc.GetConnectionCount()
				if err != nil {
					log.Errorf("[%s]TxProcessInfos Error:%s", p.symbol, err)
				} else {
					p.txProcessDBInfos()
				}
			case <-p.txProcessStop:
				log.Infof("[%s]Tx process DB stop ", p.symbol)
				return
			}
		}
	}(ticker)
	log.Infof("[%s]Tx Process DB Start", p.symbol)
}

func (p *processPresenter) txProcessDBInfos() {
	txs, err := p.db.TxProcessInfos()
	if err != nil {
		log.Errorf("[%s]txProcessList Error:%s", p.symbol, err)
		return
	}
	for _, tx := range txs {
		err = p.db.UpsertTxProcess(nil, tx.Symbol, tx.TxID, db.TxProcessStatusWait, 30)
		txQueue <- &txInfo{
			Symbol: p.symbol,
			TxID:   tx.TxID,
		}
	}
	if len(txs) > 0 {
		p.txProcessDBInfos()
	}
}

var (
	txQueue      = make(chan *txInfo, 100)
	txProcessSum = make(chan struct{}, 1)
)

func txQueueHandler() {
	log.Info("Tx Queue Start")
	for tx := range txQueue {
		txProcessSum <- struct{}{}
		go func(tx *txInfo) {
			txProcessInfo(tx)
			<-txProcessSum
		}(tx)
	}
}

func txProcessInfo(tx *txInfo) {
	symbol := tx.Symbol
	p, _ := presenters[symbol]
	txID := tx.TxID
	isProcessDone, err := p.db.IsTxProcessDone(nil, txID)
	if err != nil {
		log.Errorf("[%s]txProcessInfo ITPD Error:%s[%s]", symbol, err, txID)
		return
	}
	if !isProcessDone {
		// this must status wait,txid may be not saved in db
		err = p.db.UpsertTxProcess(nil, symbol, txID, db.TxProcessStatusWait, 30)
		if err != nil {
			log.Errorf("[%s]txProcessInfo SaveTxProcess Error:%s[%s]", symbol, err, txID)
			return
		}
		// log.Infof("[%s]txProcessInfo SaveTxProcess success[%s]", symbol, txID)
	} else {
		log.Infof("[%s]txProcessInfo TxProcess is process Done [%s]", symbol, txID)
		return
	}
	txInfo, err := p.rpc.GetTransaction(txID, nil)
	if err != nil {
		log.Errorf("[%s]GetTransaction Error:[%s][%s]", symbol, err, txID)
		return
	}
	if txInfo.Amount <= 0.0 {
		// log.Infof("[%s]Ingore prcess TX Amount is < 0.0 [%s]", symbol, txID)
		p.db.UpsertTxProcess(nil, symbol, txID, db.TxProcessStatusDone, 0)
		return
	}
	isConfirmed := txInfo.Confirmations >= p.depositMinConfirmactions
	for _, txDetail := range txInfo.Details {
		amount := txDetail.Amount
		address := txDetail.Address
		err = p.db.Deposit(txDetail.Address, txID, amount, isConfirmed)
		if err != nil {
			log.Errorf("[%s]Deposit Error:[%s][%s][%s][%f]", symbol, err, txID, address, amount)
			continue
		}
		log.Infof("[%s]Process Deposit Success:[%s][%s][%f][confirmed:%d]", symbol, txID, address, amount, txInfo.Confirmations)
	}

}

var presenters map[string]*processPresenter
var allConfig *config.Config

func initPresenter() {
	mgoSession, err := umgo.NewSession(allConfig.Mongodb)
	if err != nil {
		logrus.Fatalln("Init Mongodb Error:", err)
	}
	db.SetLog(log)
	db.SetSession(mgoSession)
	rpcclient.SetLog(log)
	presenters = make(map[string]*processPresenter)
	for k, v := range allConfig.Infos {
		p := new(processPresenter)
		p.symbol = k
		p.db = db.New(v.Symbol, v.Database)
		p.rpc = rpcclient.New(v.RPC)
		p.depositMinConfirmactions = v.MinConfirmations4Deposit
		p.scanTxDir = v.ScanTxDir
		p.wg = new(sync.WaitGroup)
		p.fileScanStop = make(chan struct{})
		p.txProcessStop = make(chan struct{})
		presenters[v.Symbol] = p
	}
}

func initConfig() {
	config, err := config.New(*configFile)
	if err != nil {
		logrus.Fatalf("Read config file error: %s", err)
	}
	allConfig = config
	logTmp, err := ulog.NewSingle(config.Env.Log)
	if err != nil {
		logrus.Fatalln("Init Log Error:", err)
	}
	log = logTmp
}

var (
	cmdFlag    = flag.String("s", "", "send signal to the daemon\nstop - fast shutdown")
	configFile = flag.String("c", "config.json", "config file path")
)

func main() {
	flag.Parse()
	initConfig()
	daemon.AddCommand(daemon.StringFlag(cmdFlag, "stop"), syscall.SIGTERM, termHandler)

	cntxt := &daemon.Context{
		PidFileName: allConfig.Env.PidFile,
		PidFilePerm: 0644,
		LogFileName: allConfig.Env.Log.LogFile,
		LogFilePerm: 0640,
		WorkDir:     allConfig.Env.WorkDir,
		Umask:       027,
		// Args:        []string{"[txserver]"},
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
	for _, txprs := range presenters {
		txprs.start()
	}

	go terminateHelper()
	go httpServerStart()
	go httpServerStopHelper()
	go txQueueHandler()

	err = daemon.ServeSignals()
	if err != nil {
		log.Info("daemon terminate Error:", err)
	}
	log.Println("daemon terminated")
}

var httpServer *http.Server

func httpServerStart() {
	mux := http.NewServeMux()
	mux.HandleFunc("/wallet", httpHandler)
	httpServer = &http.Server{
		Addr:         allConfig.TxServer.ListenerAddr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}

	log.Infof("Starting Http Server:%s", allConfig.TxServer.ListenerAddr)
	err := httpServer.ListenAndServe()
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
	<-stopHTTPServer
	log.Info("Http server start shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := httpServer.Shutdown(ctx)
	// err := server.Shutdown(nil)
	if err != nil {
		log.Info("Http server shutdown error:", err)
	}
	log.Info("Http server shutdowning...")
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Server Read Body Error", err)
		fmt.Fprintf(w, "Server Read Body Error: %s\n", err)
		return
	}
	var info txInfo
	err = json.Unmarshal(bytes, &info)
	if err != nil {
		log.Errorf("Unmarshal Error: %s (err:%s)", bytes, err)
		fmt.Fprintf(w, "Unmarshal Error: %s (err:%s)", bytes, err)
		return
	}
	symbol := info.Symbol
	if _, ok := presenters[symbol]; !ok {
		log.Errorf("Dont had '%s''s config", symbol)
		fmt.Fprintf(w, "Dont had '%s''s config\n", symbol)
		return
	}
	txQueue <- &info
	fmt.Fprintln(w, "accept success")
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
	stopHTTPServer <- struct{}{}
	<-stopHTTPServerDone
	for _, v := range presenters {
		v.stop()
	}
	for _, v := range presenters {
		v.wg.Wait()
	}
	stop <- struct{}{}
	<-done
	return daemon.ErrStop
}
