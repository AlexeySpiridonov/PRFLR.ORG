package collector

import (
    "prflr.org/config"
    "prflr.org/timer"
    "log"
    "net"
    "strconv"
    "strings"
    "errors"
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
            //log.Print("! Collector.go::ReadFromUDP !")
            log.Panic(err)
        }
        go saveMessage(string(buffer[0:n]))
    }
}

/* UDP Handlers */
func saveMessage(msg string) {
    timer, err := parseStringToTimer(msg)
    if err != nil {
        log.Print(err)
        //log.Panic(err)
    } else {
        err = timer.Save()
        if err != nil {
            log.Print(err)
            //log.Panic(err)
        }
    }
}

func parseStringToTimer(msg string) (*timer.Timer, error) {
    fields := strings.Split(msg, "|")

    if (len(fields) < 6) {
        return nil, errors.New("Ivalid format")
    }

    time, err := strconv.ParseFloat(fields[3], 32)
    if err != nil {
        log.Print(err)
        //log.Panic(err)
        return nil, errors.New("Cannot parse string " + msg)
    }

    //TODO add check for apikey and crop for fields lenght
    return &timer.Timer{fields[0], fields[1], fields[2], float32(time), fields[4], fields[5]}, nil
}
