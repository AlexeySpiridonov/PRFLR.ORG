package collector

import (
    "prflr.org/config"
    "prflr.org/timer"
    "prflr.org/db"
    "prflr.org/PRFLRLogger"
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
        PRFLRLogger.Fatal(err)
    }

    l, err := net.ListenUDP("udp", laddr)
    if err != nil {
        PRFLRLogger.Fatal(err)
    }

    db.Init()

    // is Buffer enough?!?!
    var buffer [500]byte
    for {
        n, _, err := l.ReadFromUDP(buffer[0:])
        if err != nil {
            PRFLRLogger.Error(err)
            continue
        }

        go saveMessage(string(buffer[0:n]))
    }
}

/* UDP Handlers */
func saveMessage(msg string) {
    timer, err := parseStringToTimer(msg)

    // Couldn't Parse? Wrong Format? Just skip it!!!
    if err != nil {
        PRFLRLogger.Error(err)
        return
    }

    //PRFLRLogger.Debug("Saving timer: " + msg)

    err = timer.Save()
    if err != nil {
        PRFLRLogger.Error(err)
    }
}

func parseStringToTimer(msg string) (*timer.Timer, error) {
    fields := strings.Split(msg, "|")

    if (len(fields) < 6) {
        return nil, errors.New("Invalid format")
    }

    time, err := strconv.ParseFloat(fields[3], 32)
    if err != nil {
        PRFLRLogger.Error(err)
        return nil, errors.New("Cannot parse string " + msg)
    }

    //TODO add check for apikey and crop for fields lenght
    return &timer.Timer{fields[0], fields[1], fields[2], float32(time), fields[4], fields[5]}, nil
}
