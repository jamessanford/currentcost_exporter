**Current Cost CC128 real-time power usage exporter for [Prometheus](https://github.com/prometheus/prometheus)**

Defaults to `-dev /dev/ttyUSB0 -baud 57600`

#### Install

    go get -u github.com/jamessanford/currentcost_exporter

#### Metrics exported

```
temperature_f 79.5
watts{channel="ch1"} 65
watts{channel="ch2"} 320

currentcost_read_count 30
currentcost_read_errors 0
currentcost_read_time 1.4612893e+09
currentcost_parse_errors 0
```

#### XML input from serial port

```
<msg><src>CC128-v0.15</src><dsb>01772</dsb><time>00:37:10</time><tmprF>79.5</tmprF><sensor>0</sensor><id>00865</id><type>1</type><ch1><watts>00065</watts></ch1><ch2><watts>00320</watts></ch2></msg>
```

XML format definition from http://www.currentcost.com/cc128/xml.htm

#### x/net/trace

Just for fun, the exporter registers [golang.org/x/net/trace](https://godoc.org/golang.org/x/net/trace) handlers at `/debug/events` and `/debug/requests`
