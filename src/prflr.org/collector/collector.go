package collector

import (
    "prflr.org/config"
    "prflr.org/structures"
    "prflr.org/db"
    "labix.org/v2/mgo"
	"log"
	"net"
	"strconv"
    "strings"
)

/* Starting UDP Server */
func Start() {
	// @TODO: add here UDP aggregator  in  different thread
	laddr, err := net.ResolveUDPAddr("udp", config.UDPPort)
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatal(err)
	}

	session := db.GetConnection()
    db      := session.DB(config.DBName)
    dbc     := db.C(config.DBTimers)

	// is Buffer enough?!?!
    var buffer [500]byte
    for {
        n, _, err := l.ReadFromUDP(buffer[0:])
        if err != nil {
            log.Panic(err)
        }
        go saveMessage(dbc, string(buffer[0:n]))
    }
}

/* UDP Handlers */
func saveMessage(dbc *mgo.Collection, msg string) {
	err := dbc.Insert(prepareMessage(msg))
	if err != nil {
		log.Panic(err)
	}
}

func prepareMessage(msg string) (timer structures.Timer) {
	fields := strings.Split(msg, "|")

	time, err := strconv.ParseFloat(fields[3], 32)
	if err != nil {
		log.Panic(err)
	}

	//return structures.Timer{fields[0][0:16], fields[1][0:16], fields[2][0:48], float32(time), fields[4][0:16]}
	//TODO add check for apikey and crop for fields lenght
	return structures.Timer{fields[0], fields[1], fields[2], float32(time), fields[4], fields[5]}
}

