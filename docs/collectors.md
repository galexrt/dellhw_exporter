Which collectors are enabled is controlled by the `--collectors-enabled` flag.

## Enabled by default

All collectors are enabled by default. You can disable collectors by specifying the whole list of collectors through the `--collectors-enabled` flag.

| Name                   | Description                                                                      |
| ---------------------- | -------------------------------------------------------------------------------- |
| `chassis`              | Overall status of chassis components.                                            |
| `chassis_batteries`    | Overall status of chassis CMOS batteries.                                        |
| `fans`                 | Overall status of system fans.                                                   |
| `firmwares`            | Information about some firmware versions (DRAC, BIOS)                            |
| `memory`               | System RAM DIMM status.                                                          |
| `nics`                 | NICs connection status.                                                          |
| `processors`           | Overall status of CPUs.                                                          |
| `ps`                   | Overall status of power supplies.                                                |
| `ps_amps_sysboard_pwr` | System board power usage.                                                        |
| `storage_battery`      | Status of storage controller backup batteries.                                   |
| `storage_controller`   | Overall status of storage controllers.                                           |
| `storage_enclosure`    | Overall status of storage enclosures.                                            |
| `storage_pdisk`        | Overall status of physical disks + failure prediction (if available).            |
| `storage_vdisk`        | Overall status of virtual disks.                                                 |
| `system`               | Overall status of system components.                                             |
| `temps`                | Overall temperatures (**in Celsius**) and status of system temperature readings. |
| `volts`                | Overall volts and status of power supply volt readings.                          |
