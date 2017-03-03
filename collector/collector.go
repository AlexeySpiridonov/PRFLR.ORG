package collector

import (
	"../config"
	"../db"
	"../timer"
	"errors"
	"github.com/op/go-logging"
	"net"
	"strconv"
	"strings"
	"time"
)

var log = logging.MustGetLogger("collector")

var timers = make(chan string, 1000000)

/* Starting UDP Server */
func Start() {

	go worker()

	// @TODO: add here UDP aggregator  in  different thread
	laddr, err := net.ResolveUDPAddr("udp", config.UDPPort)
	if err != nil {
		log.Fatal(err.Error())
	}

	l, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatal(err.Error())
	}

	db.Init()

	// is Buffer enough?!?!
	var buffer [500]byte
	for {
		n, _, err := l.ReadFromUDP(buffer[0:])
		if err != nil {
			log.Error(err.Error())
			continue
		}

		timers <- string(buffer[0:n])

	}
}

func worker() {
	for {
		select {
			case t := <-timers:
				saveMessage(t)
		}
	}
}

/* UDP Handlers */
func saveMessage(msg string) {
	nowMS := time.Now().UnixNano() / int64(time.Millisecond)

	timer, err := parseStringToTimer(msg)

	// Couldn't Parse? Wrong Format? Just skip it!!!
	if err != nil {
		log.Warning(err.Error())
		return
	}

	timer.Timestamp = (nowMS - int64(timer.Time)) / 1000 // format to Seconds

	err = timer.Save()
	if err != nil {
		log.Error(err.Error())
	}
}

func parseStringToTimer(msg string) (*timer.Timer, error) {
	fields := strings.Split(msg, "|")

	if len(fields) < 6 {
		return nil, errors.New("Invalid message format: " + msg)
	}

	// Validate Thread, 32 chars
	if len(fields[0]) > 32 {
		fields[0] = fields[0][:32]
	}
	if len(fields[0]) == 0 {
		return nil, errors.New("Invalid format: Thread field is not specified: " + msg)
	}
	// Validate Source
	if len(fields[1]) == 0 {
		return nil, errors.New("Invalid format: Source field is not specified: " + msg)
	}
	// Validate Timer, 48 chars
	if len(fields[2]) > 48 {
		fields[2] = fields[2][:48]
	}
	if len(fields[2]) == 0 {
		return nil, errors.New("Invalid format: Timer field is not specified: " + msg)
	}
	// Validate Info, 32 chars
	if len(fields[4]) > 32 {
		fields[4] = fields[4][:32]
	}

	// Validate Api Key, 32 chars
	if len(fields[5]) > 32 {
		fields[5] = fields[5][:32]
	}
	if len(fields[5]) == 0 {
		return nil, errors.New("Invalid format: Api Key field is not specified: " + msg)
	}

	// cut protocol from API Key if any
	apiKey := ""
	apiKeyParts := strings.Split(fields[5], "://")
	if len(apiKey) > 1 {
		apiKey = apiKeyParts[1]
	} else {
		apiKey = apiKeyParts[0]
	}

	// Validate Duration
	time, err := strconv.ParseFloat(fields[3], 32)
	if err != nil {
		return nil, errors.New("Cannot parse string, Duration is Invalid: " + msg)
	}

	//TODO add check for apikey and crop for fields lenght
	return &timer.Timer{Thrd: fields[0], Src: fields[1], Timer: fields[2], Time: float32(time), Info: fields[4], Apikey: apiKey}, nil
}
