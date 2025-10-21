# FAQ

Some frequently asked questions about running and using the `dellhw_exporter`.

## `dell_hw_chassis_temps` is reporting `0` and not a temperature?

The `dell_hw_chassis_temps` metric is reporting the [**status** (as a number)](metrics.md#what-do-the-metrics-mean) reported for each component.
The temperatures are reported by the other metrics, e.g., `dell_hw_chassis_temps_reading`.

## `lsmod: command not found` errors on start of exporter

Especially when running the`dellhw_exporter` using the container image, these error lines can appear.

The reason for that is that the dellhw_exporter image entrypoint script runs the Dell HW OMSA services start script. This start script attempts to check if some specific kernel modules are loaded/available.
As the container image doesn't contain the `lsmod` command, these "errors" can safely be ignored.

## How to disable specific collectors if not applicable to the system?

The `chassis_batteries` collector might not be available on all Dell systems, due to:

> Removal of CMOS battery sensor:
> As part of implementing the logging changes above, the CMOS battery sensor has been disabled to avoid inconsistencies in health reporting. As a result, the CMOS battery no longer appears in any management interface, including:
iDRAC web UI
>
> Source: <https://www.dell.com/support/kbdoc/en-uk/000227413/14g-intel-poweredge-coin-cell-battery-changes-in-august-2024-firmware>

To avoid error logs about collectors not being applicable to the system, the new flag `--collectors-check` can be used to specify a comma separated list of collectors to check for applicability.
Please note though that currently only the `chassis_batteries` collector is supported and is currently not checked by default.
To enable having the exporter check if `chassis_batteries` is available, you need to add the flag `--collectors-check=chassis_batteries` or env var `DELLHW_EXPORTER_COLLECTORS_CHECK=chassis_batteries`.

For more information regarding configuration of collectors, please see the [Configuration](configuration.md#collectors-configuration) documentation page.
