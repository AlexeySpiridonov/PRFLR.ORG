package influx

import (
	"../timer"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/op/go-logging"
	"time"
)

var log = logging.MustGetLogger("influx")

var iclient client.Client
var batch client.BatchPoints

var timers = make(chan *timer.Timer, 1000000)

var i = 0

func init() {
  var err error
	iclient, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "username",
		Password: "password",
	})
	if err != nil {
		log.Panic(err.Error())
	}

	bpCreate()

  go worker()

}

func bpCreate() {
	batch, _ = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "prflr",
		Precision: "us", //TODO ???
	})
}

func bpClose() {
	err := iclient.Write(batch)
	if err != nil {
		log.Error(err.Error())
	}
}

func Save(t *timer.Timer) {
	timers <- t
}

func worker() {
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
		"time": t.Time,
	}

  pt, err := client.NewPoint(t.Apikey, tags, fields, time.Now())
  if err != nil {
  		log.Error(err.Error())
  }
  batch.AddPoint(pt)

}
