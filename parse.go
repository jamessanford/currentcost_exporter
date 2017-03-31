package main

// Current Cost CC128 Real-Time power usage exported as Prometheus metrics.

// Reads serial port data as per http://www.currentcost.com/cc128/xml.htm

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"time"

	"github.com/prometheus/common/log"
	"github.com/tarm/serial"
	"golang.org/x/net/trace"
)

var (
	serialDev = flag.String("dev", "/dev/ttyUSB0", "serial port")
	baudRate  = flag.Int("baud", 57600, "serial port baud rate")
)

type wattUsage struct {
	// use strings so setGauge() can check for presence vs zero
	XMLName xml.Name `xml:"msg"`
	TempC   string   `xml:"tmpr"`
	TempF   string   `xml:"tmprF"`
	Watts1  string   `xml:"ch1>watts"`
	Watts2  string   `xml:"ch2>watts"`
}

// not exactly idiomatic go
func updateMetrics(tr trace.Trace, w *wattUsage) {
	e := func(err error) {
		if err != nil {
			parseErrors.Inc()
			tr.LazyPrintf("%v", err)
			tr.SetError()
		}
	}
	e(setGauge(tempC, w.TempC))
	e(setGauge(tempF, w.TempF))
	e(setGauge(watts.WithLabelValues("ch1"), w.Watts1))
	e(setGauge(watts.WithLabelValues("ch2"), w.Watts2))
}

// NOTE: be careful to not retain references to 'd' buffer
//       (consider changing this to take a string)
func parseWattLine(d []byte) {
	tr := trace.New("parseWattLine", "xml")
	defer tr.Finish()

	tr.LazyPrintf(fmt.Sprintf("input: %s", d)) // don't leak 'd'.

	w := new(wattUsage)
	err := xml.Unmarshal(d, &w)
	if err != nil {
		parseErrors.Inc()
		log.Errorf("xml: %v", err)
		tr.LazyPrintf("xml: %v", err)
		tr.SetError()
		return
	}

	tr.LazyPrintf("%+v", w)
	updateMetrics(tr, w)
}

func readWattLine(d []byte) {
	// ignore blank lines
	if len(d) == 0 {
		return
	}
	readCount.Inc()
	readTime.Set(float64(time.Now().Unix()))
	parseWattLine(d)
}

func readInput() {
	tr := trace.NewEventLog("readInput", *serialDev)
	defer tr.Finish()

	f, err := serial.OpenPort(
		&serial.Config{
			Name: *serialDev,
			Baud: *baudRate,
		})
	if err != nil {
		readErrors.Inc()
		tr.Errorf("open: %v", err)
		log.Error(err)
		return
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		readWattLine(s.Bytes())
	}

	err = s.Err()
	if err != nil {
		readErrors.Inc()
		tr.Errorf("scanner: %v", err)
		log.Errorf("scanner: %v", err)
	}
}
