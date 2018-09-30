package main

import (
	"bufio"
	"flag"
	"os"
	"path/filepath"
	"runtime"
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

var log = logrus.New()

type txProcess struct {
	symbol                   string
	scanTxDir                string
	depositMinConfirmactions int64
	db                       *db.DB
	rpc                      *rpcclient.Client
	wg                       *sync.WaitGroup
	fileScanStop             chan struct{}
	txProcessStop            chan struct{}
}

func (txPrs *txProcess) start() {
	txPrs.rpc.Start()
	txPrs.fileScanStart()
	txPrs.txProcessStart()
}

func (txPrs *txProcess) stop() {
	txPrs.fileScanStop <- struct{}{}
	txPrs.txProcessStop <- struct{}{}
}

func (txPrs *txProcess) fileScanStart() {
	ticker := time.NewTicker(60 * time.Second)
	txPrs.wg.Add(1)
	go func(ticker *time.Ticker) {
		defer func() {
			ticker.Stop()
			txPrs.wg.Done()
		}()
		for {
			select {
			case <-ticker.C:
				txPrs.fileScan()
			case <-txPrs.fileScanStop:
				log.Infof("[%s]File scan ticker stop ", txPrs.symbol)
				return
			}
		}
	}(ticker)
	log.Infof("[%s]File Scan Start", txPrs.symbol)
}

func (txPrs *txProcess) processFile(filepath string) {
	newFilepath := filepath + ".process"
	err := os.Rename(filepath, newFilepath)
	if err != nil {
		log.Errorf("[%s]Move File Error:%s[%s]", txPrs.symbol, err, filepath)
		return
	}
	file, err := os.Open(newFilepath)
	if err != nil {
		log.Errorf("[%s]Load File Error:%s[%s]", txPrs.symbol, err, filepath)
		return
	}
	isSuccess := false
	defer func() {
		if isSuccess {
			err := os.Remove(newFilepath)
			if err != nil {
				log.Errorf("[%s]Remove Process File Error:%s[%s]", txPrs.symbol, err, newFilepath)
			}
		} else {
			err := os.Rename(newFilepath, filepath)
			if err != nil {
				log.Errorf("[%s]Move File Error:%s[%s]", txPrs.symbol, err, newFilepath)
				return
			}
		}
	}()
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		symbolTxID := strings.Split(scanner.Text(), ",")
		err := txPrs.db.SaveTxProcess(symbolTxID[0], symbolTxID[1])
		if err != nil {
			log.Errorf("[%s]SaveTxProcess Error:%s[%s]", txPrs.symbol, err, symbolTxID[1])
		} else {
			log.Infof("[%s]SaveTxProcess success[%s]", txPrs.symbol, symbolTxID[1])
		}
	}
	isSuccess = true
}

func (txPrs *txProcess) fileScan() {
	filepath.Walk(txPrs.scanTxDir,
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
			log.Infof("[%s]Will process file:%s", txPrs.symbol, path)
			txPrs.processFile(path)
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

func (txPrs *txProcess) txProcessStart() {
	ticker := time.NewTicker(30 * time.Second)
	txPrs.wg.Add(1)
	go func(ticker *time.Ticker) {
		defer func() {
			ticker.Stop()
			txPrs.wg.Done()
		}()
		for {
			select {
			case <-ticker.C:
				txPrs.txProcessInfos()
			case <-txPrs.txProcessStop:
				log.Infof("[%s]Tx process ticker stop ", txPrs.symbol)
				return
			}
		}
	}(ticker)
	log.Infof("[%s]Tx Process Start", txPrs.symbol)
}

func (txPrs *txProcess) txProcessInfos() {
	_, err := txPrs.rpc.GetConnectionCount()
	if err != nil {
		log.Errorf("[%s]TxProcessInfos Error:%s", txPrs.symbol, err)
		return
	}
	txs, err := txPrs.db.TxProcessInfos()
	if err != nil {
		log.Errorf("[%s]txProcessList Error:%s", txPrs.symbol, err)
		return
	}
	for _, tx := range txs {
		err := txPrs.db.TxProcessUpdate(nil, tx.TxID, 0, 10)
		if err != nil {
			log.Errorf("[%s]txProcessList Error:%s[%s]", txPrs.symbol, err, tx.TxID)
			continue
		}
		txPrs.txProcessInfo(tx.TxID)
	}
	if len(txs) > 0 {
		txPrs.txProcessInfos()
	}
}

func (txPrs *txProcess) txProcessInfo(txID string) {
	txInfo, err := txPrs.rpc.GetTransaction(txID, nil)
	symbol := txPrs.symbol
	if err != nil {
		log.Errorf("[%s]GetTransaction Error:[%s][%s]", symbol, err, txID)
		txPrs.db.TxProcessUpdate(nil, txID, 0, 30)
		return
	}
	if txInfo.Amount <= 0.0 {
		// log.Infof("[%s]Ingore prcess TX Amount is < 0.0 [%s]", symbol, txID)
		txPrs.db.TxProcessUpdate(nil, txID, 1, 0)
		return
	}
	isConfirmed := txInfo.Confirmations >= txPrs.depositMinConfirmactions
	for _, txDetail := range txInfo.Details {
		amount := txDetail.Amount
		address := txDetail.Address
		err = txPrs.db.Deposit(txDetail.Address, txID, amount, isConfirmed)
		if err != nil {
			log.Errorf("[%s]Deposit Error:[%s][%s][%s][%f]", symbol, err, txID, address, amount)
			continue
		}
		log.Infof("[%s]Process Deposit Success:[%s][%s][%f][confirmed:%d]", symbol, txID, address, amount, txInfo.Confirmations)
	}

}

var txPrses []*txProcess

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
	db.SetLog(log)
	db.SetSession(mgoSession)
	rpcclient.SetLog(log)
	txPrses = make([]*txProcess, 0, len(config.Infos))
	for k, v := range config.Infos {
		txPrs := new(txProcess)
		txPrs.symbol = k
		txPrs.db = db.New(v.Symbol, v.Database)
		txPrs.rpc = rpcclient.New(v.RPC)
		txPrs.depositMinConfirmactions = v.MinConfirmations4Deposit
		txPrs.scanTxDir = v.ScanTxDir
		txPrs.wg = new(sync.WaitGroup)
		txPrs.fileScanStop = make(chan struct{})
		txPrs.txProcessStop = make(chan struct{})
		txPrses = append(txPrses, txPrs)
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

	for _, txprs := range txPrses {
		txprs.start()
	}

	go terminateHelper()

	err = daemon.ServeSignals()
	if err != nil {
		log.Info("daemon terminate Error:", err)
	}
	log.Println("daemon terminated")
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
	for _, v := range txPrses {
		v.stop()
	}
	for _, v := range txPrses {
		v.wg.Wait()
	}
	stop <- struct{}{}
	<-done
	return daemon.ErrStop
}
