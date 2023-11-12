package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

const pollInterval = 2
const reportInterval = 10

var valuesGauge = map[string]float64{}
var pollCount uint64

func main() {

	pollTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			getMetrics()
		case <-reportTicker.C:
			postQueries()
		}
	}
}

func getMetrics() {
	var rtm runtime.MemStats

	pollCount += 1
	runtime.ReadMemStats(&rtm)

	valuesGauge["Alloc"] = float64(rtm.Alloc)
	valuesGauge["BuckHashSys"] = float64(rtm.BuckHashSys)
	valuesGauge["Frees"] = float64(rtm.Frees)
	valuesGauge["GCCPUFraction"] = float64(rtm.GCCPUFraction)
	valuesGauge["HeapAlloc"] = float64(rtm.HeapAlloc)
	valuesGauge["HeapIdle"] = float64(rtm.HeapIdle)
	valuesGauge["HeapInuse"] = float64(rtm.HeapInuse)
	valuesGauge["HeapObjects"] = float64(rtm.HeapObjects)
	valuesGauge["HeapReleased"] = float64(rtm.HeapReleased)
	valuesGauge["HeapSys"] = float64(rtm.HeapSys)
	valuesGauge["LastGC"] = float64(rtm.LastGC)
	valuesGauge["Lookups"] = float64(rtm.Lookups)
	valuesGauge["MCacheInuse"] = float64(rtm.MCacheInuse)
	valuesGauge["MCacheSys"] = float64(rtm.MCacheSys)
	valuesGauge["MSpanInuse"] = float64(rtm.MSpanInuse)
	valuesGauge["MSpanSys"] = float64(rtm.MSpanSys)
	valuesGauge["Mallocs"] = float64(rtm.Mallocs)
	valuesGauge["NextGC"] = float64(rtm.NextGC)
	valuesGauge["NumForcedGC"] = float64(rtm.NumForcedGC)
	valuesGauge["NumGC"] = float64(rtm.NumGC)
	valuesGauge["OtherSys"] = float64(rtm.OtherSys)
	valuesGauge["PauseTotalNs"] = float64(rtm.PauseTotalNs)
	valuesGauge["StackInuse"] = float64(rtm.StackInuse)
	valuesGauge["StackSys"] = float64(rtm.StackSys)
	valuesGauge["Sys"] = float64(rtm.Sys)
	valuesGauge["TotalAlloc"] = float64(rtm.TotalAlloc)
}

func postQueries() {
	for k, v := range valuesGauge {
		post("gauge", k, strconv.FormatFloat(v, 'f', -1, 64))
	}
	fmt.Println(pollCount)
	post("counter", "PollCount", strconv.FormatUint(pollCount, 10))
	post("gauge", "RandomValue", strconv.FormatFloat(rand.Float64(), 'f', -1, 64))
	pollCount = 0
}

func post(mType string, mName string, mValue string) {
	resp, err := http.Post(fmt.Sprintf("http://127.0.0.1:8080/update/%s/%s/%s", mType, mName, mValue), "text/plain", bytes.NewReader([]byte{}))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
