# FAQ

Some frequently asked questions about running and using the `dellhw_exporter`.

## `dell_hw_chassis_temps` is reporting `0` and not a temperature?

The `dell_hw_chassis_temps` metric is reporting the [**status** (as a number)](metrics.md#what-do-the-metrics-mean) reported for each component.
The temperatures are reported by the other metrics, e.g., `dell_hw_chassis_temps_reading`.

## `lsmod: command not found` errors on start of exporter

Especially when running the`dellhw_exporter` using the container image, these error lines can appear.

The reason for that is that the dellhw_exporter image entrypoint script runs the Dell HW OMSA services start script. This start script attempts to check if some specific kernel modules are loaded/available.
As the container image doesn't contain the `lsmod` command, these "errors" can safely be ignored.