package influx

import (
	"github.com/influxdata/influxdb/client/v2"
	"github.com/op/go-logging"
	"prflr.org/timer"
	"time"
)

var log = logging.MustGetLogger("influx")

var iclient client.Client
var batch client.BatchPoints

var timers = make(chan *timer.Timer, 1000000)

var i = 0

func init() {
	log.Info("Stat InfluxDB Client")
	var err error
	iclient, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "prflr",
		Password: "prflr",
	})
	if err != nil {
		log.Panic(err)
	}

	bpCreate()

	go worker()

}

func bpCreate() {
	log.Debug("Create batch point")
	batch, _ = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "prflr",
		Precision: "us", //TODO ???
	})
}

func bpClose() {
	log.Debug("Close batch point")
	err := iclient.Write(batch)
	if err != nil {
		log.Error(err)
	}
}

func Save(t *timer.Timer) {
	timers <- t
}

func worker() {
	log.Info("Start worker")
	for {
		select {
		case t := <-timers:
			saveB(t)
		}
	}
}

func saveB(t *timer.Timer) {

	if i > 1000 {
		bpClose()
		bpCreate()
		i = 0
	}

	i++

	tags := map[string]string{
		"src":   t.Src,
		"timer": t.Timer,
		"thrd":  t.Thrd,
		"info":  t.Info,
	}

	fields := map[string]interface{}{
		"Exec": t.Time,
	}

	pt, err := client.NewPoint(t.Apikey, tags, fields, time.Now())
	if err != nil {
		log.Error(err)
	}
	batch.AddPoint(pt)

}
