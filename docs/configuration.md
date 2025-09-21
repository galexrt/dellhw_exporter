The dellhw_exporter can be configured using flags or environment variables.
In case of the container image there are certain specific environment variables, to help running inside a containerized environment.

## Flags

```console
$ dellhw_exporter --help
Usage of dellhw_exporter:
      --cache-duration int              Cache duration in seconds (default 20)
      --cache-enabled                   Enable metrics caching to reduce load
      --collectors-additional strings   Comma separated list of collectors to enable additionally to the collectors-enabled list
      --collectors-cmd-timeout int      Command execution timeout for omreport (default 15)
      --collectors-enabled strings      Comma separated list of active collectors (default [chassis,chassis_batteries,fans,firmwares,memory,nics,processors,ps,ps_amps_sysboard_pwr,storage_battery,storage_controller,storage_enclosure,storage_pdisk,storage_vdisk,system,temps,version,volts])
      --collectors-omreport string      Path to the omreport executable (based on the OS (linux or windows) default paths are used if unset) (default "/opt/dell/srvadmin/bin/omreport")
      --collectors-print                If true, print available collectors and exit.
      --log-level string                Set log level (default "INFO")
      --monitored-nics strings          Comma separated list of nics to monitor (default, empty list, is to monitor all)
      --version                         Show version information
      --web-config-file string          [EXPERIMENTAL] Path to configuration file that can enable TLS or authentication.
      --web-listen-address string       The address to listen on for HTTP requests (default ":9137")
      --web-telemetry-path string       Path the metrics will be exposed under (default "/metrics")
```

The `--web-config-file` instructs the exporter to load a separate YAML config file that provides the following abilities:

- HTTPS
- TLS cert authentication
- HTTP2
- Basic Authentication
- TLS versions and cipher suites
- Headers like `Strict-Transport-Security`, `X-XSS-Protection`, `X-Frame-Options`, etc.

The exact format of the file and all its options can be found [here](https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md).

## Environment Variables

For the description of the env vars, see the above equivalent flags (and their defaults).

```console
DELLHW_EXPORTER_CACHE_DURATION
DELLHW_EXPORTER_CACHE_ENABLED
DELLHW_EXPORTER_COLLECTORS_ADDITIONAL
DELLHW_EXPORTER_COLLECTORS_CMD_TIMEOUT
DELLHW_EXPORTER_COLLECTORS_ENABLED
DELLHW_EXPORTER_COLLECTORS_OMREPORT
DELLHW_EXPORTER_LOG_LEVEL
DELLHW_EXPORTER_MONITORED_NICS
DELLHW_EXPORTER_WEB_LISTEN_ADDRESS
DELLHW_EXPORTER_WEB_TELEMETRY_PATH
DELLHW_EXPORTER_WEB_CONFIG_FILE
```

### Container Image specific Environment Variables

| Env                            | Default | Description                                                                             |
| ------------------------------ | ------- | --------------------------------------------------------------------------------------- |
| `START_DELL_SRVADMIN_SERVICES` | `true`  | Set to false if you don't want the srvadmin services to be started inside the container |
