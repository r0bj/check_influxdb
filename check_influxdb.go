package main

import (
	"fmt"
	"flag"
	"encoding/json"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/olorin/nagiosplugin"
)

var (
	defaultHost = "localhost"
	defaultPort = "8086"
	defaultUsername = "admin"
	defaultPassword = "admin"
	defaultDB = "telegraf"
	defaultMeasurement = "cpu"
	defaultField = "usage_system"
	defaultTimeRange = "5m"
	defaultWarningThreshold = 10000
	defaultCriticalThreshold = 100
)

func queryDB(c client.Client, db string, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: db,
	}
	if response, err := c.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}


func main() {
	host := flag.String("H", defaultHost, "influxdb host")
	port := flag.String("P", defaultPort, "influxdb port")
	username := flag.String("u", defaultUsername, "influxdb username")
	password := flag.String("p", defaultPassword, "influxdb password")
	db := flag.String("d", defaultDB, "influxdb database name")
	measurement := flag.String("m", defaultMeasurement, "influxdb measurement")
	field := flag.String("f", defaultField, "influxdb measurement field")
	timeRange := flag.String("t", defaultTimeRange, "time range in influxdb query syntax: u microseconds, s seconds, m minutes, h hours, d days, w weeks")
	warningThreshold := flag.Int("w", defaultWarningThreshold, "warning threshold for number of returned records")
	criticalThreshold := flag.Int("c", defaultCriticalThreshold, "critical threshold for number of returned records")
	flag.Parse()

	check := nagiosplugin.NewCheck()
	defer check.Finish()

	if *warningThreshold < *criticalThreshold {
		check.AddResult(nagiosplugin.UNKNOWN, "warning threshold lower than critical")
		return
	}

	url := fmt.Sprintf("http://%s:%s", *host, *port)
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: url,
		Username: *username,
		Password: *password,
	})
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, "error connecting to influxdb")
		return
	}

	q := fmt.Sprintf("SELECT count(%s) FROM %s WHERE time > now() - %s", *field, *measurement, *timeRange)
	res, err := queryDB(c, *db, q)
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, "query to influxdb failed")
		return
	}
	count := res[0].Series[0].Values[0][1]

	var i int
	err = json.Unmarshal([]byte(count.(json.Number)), &i)
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, "unmarshal error")
		return
	}

	if i > *warningThreshold {
		check.AddResult(nagiosplugin.OK, fmt.Sprintf("fresh data present: %d records", i))
	} else if i < *warningThreshold && i > *criticalThreshold {
		check.AddResult(nagiosplugin.WARNING, fmt.Sprintf("number of records below warning threshold: %d", i))
	} else if i < *criticalThreshold {
		check.AddResult(nagiosplugin.CRITICAL, fmt.Sprintf("number of records below critical threshold: %d", i))
	}
}
