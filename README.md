# dellhw_exporter
[![CircleCI branch](https://img.shields.io/circleci/project/github/RedSparr0w/node-csgo-parser/master.svg)]() [![Docker Repository on Quay](https://quay.io/repository/galexrt/dellhw_exporter/status "Docker Repository on Quay")](https://quay.io/repository/galexrt/dellhw_exporter) [![Go Report Card](https://goreportcard.com/badge/github.com/galexrt/dellhw_exporter)](https://goreportcard.com/report/github.com/galexrt/dellhw_exporter) [![license](https://img.shields.io/github/license/mashape/apistatus.svg)]()
Prometheus exporter for Dell Hardware components using OMSA.

The exporter was originally made by [PrFalken](https://github.com/PrFalken). Due to some issues in the code, I rewrote the whole exporter using the "node_exporter" pattern and therefore moved it from being a fork out, to a standalone repository.

Omreport parsing functions were borrowed from the [Bosun project](https://github.com/bosun-monitor/bosun/blob/master/cmd/scollector/collectors/dell_hw.go), thank you very much for that, they are the most tedious part of the job.

This exporter wraps the "omreport" command from Dell OMSA. If you can't run omreport on your system, the exporter won't export any metrics.

## Dell OMSA Compatibility
* `7.4`
* `8.4`

## Collectors
Which collectors are enabled is controlled by the `--colectors.enabled` flag.

### Enabled by default
All collectors are enabled by default right now.
Name     | Description
---------|-------------
chassis | Overall status of chassis components.
fans | Overall status of system fans.
memory | System RAM DIMM status.
processors | Overall status of CPUs.
ps_amps_sysboard_pwr | System board power usage.
ps | Overall status of power supplies.
storage_battery | Status of storage controller backup batteries.
storage_controller | Overall status of storage controllers.
storage_enclosure | Overall status of storage enclosures.
storage_vdisk | Overall status of virtual disks.
system | Overall status of system components.
temps | Overall temperatures and status of system temperature readings.
volts | Overall volts and status of power supply volt readings.

## Help Menu
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
