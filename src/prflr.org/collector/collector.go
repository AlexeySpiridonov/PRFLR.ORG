package collector

import (
    "prflr.org/config"
    "prflr.org/timer"
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

    // is Buffer enough?!?!
    var buffer [500]byte
    for {
        n, _, err := l.ReadFromUDP(buffer[0:])
        if err != nil {
            log.Print("! Collector.go::ReadFromUDP !")
            log.Panic(err)
        }
        go saveMessage(string(buffer[0:n]))
    }
}

/* UDP Handlers */
func saveMessage(msg string) {
    timer := parseStringToTimer(msg)
    err   := timer.Save()
    if err != nil {
        log.Print("! Collector.go::saveMessage !")
        log.Panic(err)
    }
}

func parseStringToTimer(msg string) timer.Timer {
    fields := strings.Split(msg, "|")

    time, err := strconv.ParseFloat(fields[3], 32)
    if err != nil {
        log.Print("! Collector.go::parseStringToTimer !")
        log.Panic(err)
    }

    //TODO add check for apikey and crop for fields lenght
    return timer.Timer{fields[0], fields[1], fields[2], float32(time), fields[4], fields[5]}
}
