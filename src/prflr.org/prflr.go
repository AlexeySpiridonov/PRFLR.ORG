package main

import (
	"prflr.org/collector"
	"prflr.org/web"
)

func main() {
    /* init HTTP Server and Handlers */
	web.Start()

    /* init UDP  Server and Handlers */
	collector.Start()
}
