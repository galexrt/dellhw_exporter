# dellhw_exporter

[![CircleCI branch](https://img.shields.io/circleci/project/github/galexrt/dellhw_exporter/master.svg)]() [![Docker Repository on Quay](https://quay.io/repository/galexrt/dellhw_exporter/status "Docker Repository on Quay")](https://quay.io/repository/galexrt/dellhw_exporter) [![Go Report Card](https://goreportcard.com/badge/github.com/galexrt/dellhw_exporter)](https://goreportcard.com/report/github.com/galexrt/dellhw_exporter)

Prometheus exporter for Dell Hardware components using OMSA.

The exporter was originally made by [PrFalken](https://github.com/PrFalken). Due to some issues in the code, I rewrote the whole exporter using the ["node_exporter"](https://github.com/prometheus/node_exporter) pattern and therefore moved it from being a fork out, to a standalone repository.

Omreport parsing functions were borrowed from the [Bosun project](https://github.com/bosun-monitor/bosun/blob/master/cmd/scollector/collectors/dell_hw.go), thank you very much for that, they are the most tedious part of the job.

This exporter wraps the "omreport" command from Dell OMSA. If you can't run omreport on your system, the exporter won't export any metrics.

## Tested Dell OMSA Compatibility

The dellhw_exporter has been tested with the following OMSA versions:

* `7.4`
* `8.4`
* `9.1`

### Kernel Compatibility

**Please note that only kernel versions that are supported by DELL DSU / OMSA tools are working!**

**State 07.06.2019**: Dell OMSA `DSU_19.05.00` is not compatible with 5.x kernel it seems (e.g., Fedora uses that kernel).

Should you run into issues when using the Docker image, please follow the [Troubleshooting - No metrics being exported](#no-metrics-being-exported).

## Collectors

Which collectors are enabled is controlled by the `--colectors.enabled` flag.

### Enabled by default

All collectors are enabled by default. You can disable collectors by specifying the whole list of collectors through the `--collectors.enabled` flag.

| Name                 | Description                                                                      |
| -------------------- | -------------------------------------------------------------------------------- |
| chassis              | Overall status of chassis components.                                            |
| chassis_batteries    | Overall status of chassis CMOS batteries.                                        |
| fans                 | Overall status of system fans.                                                   |
| memory               | System RAM DIMM status.                                                          |
| nics                 | NICs connection status.                                                          |
| processors           | Overall status of CPUs.                                                          |
| ps                   | Overall status of power supplies.                                                |
| ps_amps_sysboard_pwr | System board power usage.                                                        |
| storage_battery      | Status of storage controller backup batteries.                                   |
| storage_controller   | Overall status of storage controllers.                                           |
| storage_enclosure    | Overall status of storage enclosures.                                            |
| storage_pdisk        | Overall status of physical disks.                                                |
| storage_vdisk        | Overall status of virtual disks.                                                 |
| system               | Overall status of system components.                                             |
| temps                | Overall temperatures (**in Celsius**) and status of system temperature readings. |
| volts                | Overall volts and status of power supply volt readings.                          |

### What do the metrics mean?

Most metrics returned besides temperature, volts, fans RPM count and others, are state indicators which can have the following of the four states:

* `0` - `OK`, the component should be fine.
* `1` - `Critical`, the component is not okay / has potentially failed / `Unknown` status.
* `2` - `Non-Critical`, the component is not okay, but not critical.

## Configuration

### Flags

```cosnole
$ dellhw_exporter --help
unknown flag: --debug
Usage of dellhw_exporter:
      --collectors-cmd-timeout int   Command execution timeout for omreport (default 15)
      --collectors-enabled string    Comma separated list of active collectors (default "chassis,chassis_batteries,fans,memory,nics,processors,ps,ps_amps_sysboard_pwr,storage_battery,storage_controller,storage_enclosure,storage_pdisk,storage_vdisk,system,temps,volts")
      --collectors-omreport string   Path to the omReport executable (default "/opt/dell/srvadmin/bin/omreport")
      --collectors-print             If true, print available collectors and exit.
      --log-level string             Set log level (default "INFO")
      --version                      Show version information
      --web-listen-address string    The address to listen on for HTTP requests (default ":9137")
      --web-telemetry-path string    Path the metrics will be exposed under (default "/metrics")
unknown flag: --debug
```

### Environment variables

For the description of the env vars, see the above equivalent flags.

```console
DELLHW_EXPORTER_COLLECTORS_CMD_TIMEOUT
DELLHW_EXPORTER_COLLECTORS_ENABLED
DELLHW_EXPORTER_COLLECTORS_OMREPORT
DELLHW_EXPORTER_COLLECTORS_PRINT
DELLHW_EXPORTER_LOG_LEVEL
DELLHW_EXPORTER_HELP
DELLHW_EXPORTER_VERSION
DELLHW_EXPORTER_WEB_LISTEN_ADDRESS
DELLHW_EXPORTER_WEB_TELEMETRY_PATH
```

#### Docker specific environment variables

```console
START_DELL_SRVADMIN_SERVICES # Defaults to `true`, toggle if the srvadmin services are started inside the container
```

## Running in Docker

The container image is available from [Docker Hub](https://hub.docker.com/) and [Quay.io](https://quay.io/):

### Pull the Docker image

#### Docker Hub

```console
docker pull galexrt/dellhw_exporter
```

#### Quay.io

```console
docker pull quay.io/galexrt/dellhw_exporter
```

### Run the Docker image

> **NOTE** The `--privileged` flag is required as the OMSA needs to access the host's devices and other components.

```console
docker run -d --name dellhw_exporter --privileged -p 9137:9137 galexrt/dellhw_exporter
# or for quay.io
docker run -d --name dellhw_exporter --privileged -p 9137:9137 quay.io/galexrt/dellhw_exporter
```

## Running without Docker

To run without Docker either download a [release binary](https://github.com/galexrt/dellhw_exporter/releases) or build it (using `make build`):

```console
./dellhw_exporter
./dellhw_exporter --help
./dellhw_exporter YOUR_FLAGS
```

**The DELL OMSA services must already be running for the exporter to be able to collect metrics!**

E.g., run `/opt/dell/srvadmin/sbin/srvadmin-services.sh start` and / or `systemctl start SERVICE_NAME` (to enable autostart use `systemctl enable SERVICE_NAME`; where `SERVICE_NAME` [are the DELL OMSA service(s) you installed](http://linux.dell.com/repo/hardware/omsa.html)).

## Prometheus

The exporter runs on TCP port `9137`.

Example static Prometheus Job config:

```yaml
[...]
  - job_name: 'dellhw_exporter'
    # Override the global default and scrape targets from this job every 60 seconds.
    scrape_interval: 60s
    static_configs:
      - targets:
        - 'YOUR_SERVER_HERE:9137'
[...]
```

## Monitoring

Checkout the files in the [`contrib/monitoring/`](contrib/monitoring/) directory.

## Troubleshooting

### No metrics being exported

If you are not running the Docker container, it is probably that your OMSA / srvadmin services are not running. Start them using the following commands:

```console
/opt/dell/srvadmin/sbin/srvadmin-services.sh status
/opt/dell/srvadmin/sbin/srvadmin-services.sh start
echo "return code: $?"
```
Please note that the return code should be `0`, if not please investigate the logs of srvadmin services.

When running inside the container this most of the time means
Be sure to enter the container and run the following commands to verify if the kernel modules have been loaded:

```console
/usr/libexec/instsvcdrv-helper status
lsmod | grep -iE 'dell|dsu'
```

Should the `lsmod` not contain any module named after `dell_` and / or `dsu_`, be sure to add the following read-only mounts depending on your OS, for the kernel modules directory (`/lib/modules`) and / or the kernel source / headers directory (depends hardly on the OS your are using) to the `dellhw_exporter` Docker container using `-v HOST_PATH:CONTAINER_PATH:ro` flag.

## Development

### Dependencies

`dep` is used for vendoring the dependencies.
