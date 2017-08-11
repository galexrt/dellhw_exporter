# dellhw_exporter

Prometheus exporter for Dell Hardware components
*Supports Dell OMSA 7.4*

This exporter wraps the "omreport" command from Dell OMSA. If you can't run omreport on your system, the exporter won't export any metrics.

Omreport parsing functions were borrowed from the [Bosun project](https://github.com/bosun-monitor/bosun/blob/master/cmd/scollector/collectors/dell_hw.go), thank you very much for that, they are the most tedious part of the job.

```
./dellhw_exporter [FLAGS]
  -collectors.enabled string
    	Comma separated list of active collectors (default "chassis,fans,memory,processors,ps,ps_amps_sysboard_pwr,storage_battery,storage_enclosure,storage_controller,storage_vdisk,system,temps,volts")
  -collectors.omr-report string
    	Path to the omReport executable (default "/opt/dell/srvadmin/bin/omreport")
  -collectors.print
    	If true, print available collectors and exit.
  -debug
    	Enable debug output
  -help
    	Show help menu
  -log.format value
    	Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true" (default "logger:stderr")
  -log.level value
    	Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal] (default "info")
  -version
    	Show version information
  -web.listen-address string
    	The address to listen on for HTTP requests (default ":9137")
  -web.telemetry-path string
    	Path the metrics will be exposed under (default "/metrics")
```
