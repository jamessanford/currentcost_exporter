package main

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// collected data
var (
	tempC = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "temperature_c",
			Help: "temperature in degrees celsius",
		})

	tempF = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "temperature_f",
			Help: "temperature in degrees fahrenheit",
		})
	watts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "watts",
			Help: "current watts drawn labeled by channel",
		},
		[]string{"channel"})
)

// instrumentation
var (
	parseErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "currentcost_parse_errors",
			Help: "number of errors parsing XML data",
		})
	readErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "currentcost_read_errors",
			Help: "number of errors reading from serial port",
		})
	readCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "currentcost_read_count",
			Help: "number of watt usage data points read over serial",
		})
	readTime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "currentcost_read_time",
			Help: "last watt usage data point was at this many seconds since unix epoch",
		})
)

func init() {
	prometheus.MustRegister(tempC)
	prometheus.MustRegister(tempF)
	prometheus.MustRegister(watts)

	prometheus.MustRegister(parseErrors)
	prometheus.MustRegister(readErrors)
	prometheus.MustRegister(readCount)
	prometheus.MustRegister(readTime)
}

func setGauge(g prometheus.Gauge, s string) error {
	// ignore empty values
	if len(s) == 0 {
		return nil
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	g.Set(f)
	return nil
}
