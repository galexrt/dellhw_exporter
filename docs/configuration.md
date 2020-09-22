The dellhw_exporter can be configured using flags or environment variables.
In case of docker there are certain specific environment variables, to help running inside a containerized environment.

## Flags

```console
$ dellhw_exporter --help
Usage of dellhw_exporter:
      --collectors-cmd-timeout int   Command execution timeout for omreport (default 15)
      --collectors-enabled string    Comma separated list of active collectors (default "chassis,chassis_batteries,fans,firmwares,memory,nics,processors,ps,ps_amps_sysboard_pwr,storage_battery,storage_controller,storage_enclosure,storage_pdisk,storage_vdisk,system,temps,volts")
      --collectors-omreport string   Path to the omReport executable (default "/opt/dell/srvadmin/bin/omreport")
      --collectors-print             If true, print available collectors and exit.
      --log-level string             Set log level (default "INFO")
      --version                      Show version information
      --web-listen-address string    The address to listen on for HTTP requests (default ":9137")
      --web-telemetry-path string    Path the metrics will be exposed under (default "/metrics")
      --cache-enabled bool           Enable caching (default false)
      --cache-duration int           Duration in seconds for the cache lifetime (default 20)
```

## Environment Variables

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
DELLHW_EXPORTER_CACHE_ENABLED
DELLHW_EXPORTER_CACHE_DURATION
```

### Docker specific Environment Variables

```console
START_DELL_SRVADMIN_SERVICES # Defaults to `true`, toggle if the srvadmin services are started inside the container
```
