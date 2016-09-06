package main

import (
	"PRFLR.ORG/collector"
	"PRFLR.ORG/web"
	"github.com/op/go-logging"
	"github.com/yvasiyarov/gorelic"
	"os"
)

var log = logging.MustGetLogger("main")

func main() {

	initGoRelic()
	initLogs()

	/* init HTTP Server and Handlers */
	web.Start()

	/* init UDP  Server and Handlers */
	collector.Start()
}

func initGoRelic() {
	agent := gorelic.NewAgent()
	agent.NewrelicLicense = "6d91ca13798027e532d8a67132d52ba34eba28bb"
	agent.NewrelicName = "PRFLR"
	agent.Run()
}

func initLogs() {
	format := logging.MustStringFormatter("PRFLR> %{module} %{shortfile} > %{level:.7s} > %{message}")
	//file to stdout
	log1 := logging.NewLogBackend(os.Stderr, "", 0)
	file, err := os.OpenFile("/var/log/prflr.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Open log file fail: /var/log/prflr.log")
	}
	log1F := logging.NewBackendFormatter(log1, format)

	//log to file
	log2 := logging.NewLogBackend(file, "", 0)
	log2F := logging.NewBackendFormatter(log2, format)

	//log to syslog
	log3, _ := logging.NewSyslogBackend("")
	log3LeveledF := logging.NewBackendFormatter(log3, format)

	//setup logs
	//if "prod" != "prod" {
	//	log3Leveled := logging.AddModuleLevel(log3LeveledF)
	//	log3Leveled.SetLevel(logging.INFO, "")
	//	logging.SetBackend(log3Leveled)
	//} else {
	logging.SetBackend(log1F, log2F, log3LeveledF)
	//}

	log.Info("Logs ok")

}
