package main

import (
    "prflr.org/collector"
    "prflr.org/web"
    "github.com/yvasiyarov/gorelic"
)

func main() {
	
    initGoRelic()

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
