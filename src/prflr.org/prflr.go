package main

import (
    "prflr.org/collector"
    "prflr.org/web"
//    "prflr.org/db"
//    "log"
)

func main() {
	/*session, _ := db.GetConnection()

	session, _ = db.GetConnection()

	session, _ = db.GetConnection()

	log.Print(session)

	session.Close()

	log.Print(session)*/

    /* init HTTP Server and Handlers */
    web.Start()

    /* init UDP  Server and Handlers */
    collector.Start()
}
