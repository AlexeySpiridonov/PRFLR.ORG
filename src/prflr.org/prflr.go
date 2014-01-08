package main

import (
	//"prflr.org/db"
	"prflr.org/collector"
	"prflr.org/web"
)

func main() {
    /* init HTTP Server and Handlers */
	web.Start()

	/* init MongoDB Connect */
	// do we really need to init here?..
	//db.Init()

    /* init UDP Server and Handlers */
	collector.Start()
}
