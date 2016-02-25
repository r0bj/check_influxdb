# check_influxdb
InfluxDB nagios check

Usage of ./check_influxdb:
  -H string
    	influxdb host (default "localhost")
  -P string
    	influxdb port (default "8086")
  -c int
    	critical threshold for number of returned records (default 100)
  -d string
    	influxdb database name (default "telegraf")
  -f string
    	influxdb measurement field (default "usage_system")
  -m string
    	influxdb measurement (default "cpu")
  -p string
    	influxdb password (default "admin")
  -t string
    	time range in influxdb query syntax: u microseconds, s seconds, m minutes, h hours, d days, w weeks (default "5m")
  -u string
    	influxdb username (default "admin")
  -w int
    	warning threshold for number of returned records (default 10000)
