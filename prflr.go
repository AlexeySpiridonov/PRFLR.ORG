package main

import (
	"./collector"
	"./web"
	"github.com/op/go-logging"
	"github.com/yvasiyarov/gorelic"
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

	//log to syslog
	slog, _ := logging.NewSyslogBackend("")
	logLeveledF := logging.NewBackendFormatter(slog, format)

	//setup logs
	logLeveled := logging.AddModuleLevel(logLeveledF)
	logLeveled.SetLevel(logging.INFO, "")
	logging.SetBackend(logLeveled)

	log.Notice("Logs ok")

}
