# Home

![build_release](https://github.com/galexrt/dellhw_exporter/workflows/build_release/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/galexrt/dellhw_exporter)](https://goreportcard.com/report/github.com/galexrt/dellhw_exporter)

Prometheus exporter for Dell Hardware components using OMSA.

The exporter was originally made by [PrFalken](https://github.com/PrFalken). Due to some issues in the code, I rewrote the whole exporter using the ["node_exporter"](https://github.com/prometheus/node_exporter) pattern and therefore moved it from being a fork out, to a standalone repository.

Omreport parsing functions were borrowed from the [Bosun project](https://github.com/bosun-monitor/bosun/blob/master/cmd/scollector/collectors/dell_hw.go), thank you very much for that, they are the most tedious part of the job.

This exporter wraps the "omreport" command from Dell OMSA. If you can't run omreport on your system, the exporter won't export any metrics.

## Compatibility

### Tested Dell OMSA Compatibility

The dellhw_exporter has been tested with the following OMSA versions:

* `7.4`
* `8.4`
* `9.1`

### Kernel Compatibility

**Please note that only kernel versions that are supported by DELL DSU / OMSA tools are working!**

**State 07.06.2019**: Dell OMSA `DSU_19.05.00` is not compatible with 5.x kernel it seems (e.g., Fedora uses that kernel).

Should you run into issues when using the Docker image, please follow the [Troubleshooting - No metrics being exported](#no-metrics-being-exported).
