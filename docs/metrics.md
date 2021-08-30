## What do the metrics mean?

Most metrics returned besides temperature, volts, fans RPM count and others, are state indicators which can have the following of the four states:

* `0` - `OK`, the component should be fine.
* `1` - `Critical`, the component is not okay / has potentially failed / `Unknown` status.
* `2` - `Non-Critical`, the component is not okay, but not critical.

## Example Metrics Output

```plain

# HELP dell_hw_bios Version info of firmwares/bios.
# TYPE dell_hw_bios gauge
dell_hw_bios{manufacturer="dell inc.",release_date="06/26/2020",version="2.8.1"} 0
# HELP dell_hw_chassis_current_reading System board power usage.
# TYPE dell_hw_chassis_current_reading gauge
dell_hw_chassis_current_reading{pwrsupply="PS1"} 0.4
dell_hw_chassis_current_reading{pwrsupply="PS2"} 0.4
# HELP dell_hw_chassis_fan_reading Overall status of system fans.
# TYPE dell_hw_chassis_fan_reading gauge
dell_hw_chassis_fan_reading{fan="System_Board_Fan1A"} 6840
dell_hw_chassis_fan_reading{fan="System_Board_Fan1B"} 6480
dell_hw_chassis_fan_reading{fan="System_Board_Fan2A"} 6840
dell_hw_chassis_fan_reading{fan="System_Board_Fan2B"} 6480
dell_hw_chassis_fan_reading{fan="System_Board_Fan3A"} 6840
dell_hw_chassis_fan_reading{fan="System_Board_Fan3B"} 6360
dell_hw_chassis_fan_reading{fan="System_Board_Fan4A"} 7200
dell_hw_chassis_fan_reading{fan="System_Board_Fan4B"} 6960
dell_hw_chassis_fan_reading{fan="System_Board_Fan5A"} 7080
dell_hw_chassis_fan_reading{fan="System_Board_Fan5B"} 6720
dell_hw_chassis_fan_reading{fan="System_Board_Fan6A"} 6960
dell_hw_chassis_fan_reading{fan="System_Board_Fan6B"} 6480
dell_hw_chassis_fan_reading{fan="System_Board_Fan7A"} 6840
dell_hw_chassis_fan_reading{fan="System_Board_Fan7B"} 6480
dell_hw_chassis_fan_reading{fan="System_Board_Fan8A"} 6840
dell_hw_chassis_fan_reading{fan="System_Board_Fan8B"} 6600
# HELP dell_hw_chassis_fan_status Overall status of system fans.
# TYPE dell_hw_chassis_fan_status gauge
dell_hw_chassis_fan_status{fan="System_Board_Fan1A"} 0
dell_hw_chassis_fan_status{fan="System_Board_Fan1B"} 0
dell_hw_chassis_fan_status{fan="System_Board_Fan2A"} 0
dell_hw_chassis_fan_status{fan="System_Board_Fan2B"} 0
dell_hw_chassis_fan_status{fan="System_Board_Fan3A"} 0
dell_hw_chassis_fan_status{fan="System_Board_Fan3B"} 0
dell_hw_chassis_fan_status{fan="System_Board_Fan4A"} 0
dell_hw_chassis_fan_status{fan="System_Board_Fan4B"} 0
dell_hw_chassis_fan_status{fan="System_Board_Fan5A"} 0
dell_hw_chassis_fan_status{fan="System_Board_Fan5B"} 0
dell_hw_chassis_fan_status{fan="System_Board_Fan6A"} 0
dell_hw_chassis_fan_status{fan="System_Board_Fan6B"} 0
dell_hw_chassis_fan_status{fan="System_Board_Fan7A"} 0
dell_hw_chassis_fan_status{fan="System_Board_Fan7B"} 0
dell_hw_chassis_fan_status{fan="System_Board_Fan8A"} 0
dell_hw_chassis_fan_status{fan="System_Board_Fan8B"} 0
# HELP dell_hw_chassis_memory_status System RAM DIMM status.
# TYPE dell_hw_chassis_memory_status gauge
dell_hw_chassis_memory_status{memory="A1"} 0
dell_hw_chassis_memory_status{memory="A2"} 0
dell_hw_chassis_memory_status{memory="A3"} 0
dell_hw_chassis_memory_status{memory="A4"} 0
dell_hw_chassis_memory_status{memory="A5"} 0
dell_hw_chassis_memory_status{memory="A6"} 0
dell_hw_chassis_memory_status{memory="B1"} 0
dell_hw_chassis_memory_status{memory="B2"} 0
dell_hw_chassis_memory_status{memory="B3"} 0
dell_hw_chassis_memory_status{memory="B4"} 0
dell_hw_chassis_memory_status{memory="B5"} 0
dell_hw_chassis_memory_status{memory="B6"} 0
# HELP dell_hw_chassis_power_fail_level System board power usage.
# TYPE dell_hw_chassis_power_fail_level gauge
dell_hw_chassis_power_fail_level 1300
# HELP dell_hw_chassis_power_reading System board power usage.
# TYPE dell_hw_chassis_power_reading gauge
dell_hw_chassis_power_reading 156
# HELP dell_hw_chassis_power_warn_level System board power usage.
# TYPE dell_hw_chassis_power_warn_level gauge
dell_hw_chassis_power_warn_level 1170
# HELP dell_hw_chassis_processor_status Overall status of CPUs.
# TYPE dell_hw_chassis_processor_status gauge
dell_hw_chassis_processor_status{processor="CPU1"} 0
dell_hw_chassis_processor_status{processor="CPU2"} 0
# HELP dell_hw_chassis_status Overall status of chassis components.
# TYPE dell_hw_chassis_status gauge
dell_hw_chassis_status{component="Batteries"} 0
dell_hw_chassis_status{component="Fans"} 0
dell_hw_chassis_status{component="Hardware_Log"} 0
dell_hw_chassis_status{component="Intrusion"} 0
dell_hw_chassis_status{component="Memory"} 0
dell_hw_chassis_status{component="Power_Management"} 0
dell_hw_chassis_status{component="Power_Supplies"} 0
dell_hw_chassis_status{component="Processors"} 0
dell_hw_chassis_status{component="Temperatures"} 0
dell_hw_chassis_status{component="Voltages"} 0
# HELP dell_hw_chassis_temps Overall temperatures and status of system temperature readings.
# TYPE dell_hw_chassis_temps gauge
dell_hw_chassis_temps{component="CPU1_Temp"} 0
dell_hw_chassis_temps{component="CPU2_Temp"} 0
dell_hw_chassis_temps{component="System_Board_Exhaust_Temp"} 0
dell_hw_chassis_temps{component="System_Board_Inlet_Temp"} 0
# HELP dell_hw_chassis_temps_max_failure Overall temperatures and status of system temperature readings.
# TYPE dell_hw_chassis_temps_max_failure gauge
dell_hw_chassis_temps_max_failure{component="CPU1_Temp"} 97
dell_hw_chassis_temps_max_failure{component="CPU2_Temp"} 97
dell_hw_chassis_temps_max_failure{component="System_Board_Exhaust_Temp"} 80
dell_hw_chassis_temps_max_failure{component="System_Board_Inlet_Temp"} 47
# HELP dell_hw_chassis_temps_max_warning Overall temperatures and status of system temperature readings.
# TYPE dell_hw_chassis_temps_max_warning gauge
dell_hw_chassis_temps_max_warning{component="System_Board_Exhaust_Temp"} 75
dell_hw_chassis_temps_max_warning{component="System_Board_Inlet_Temp"} 43
# HELP dell_hw_chassis_temps_min_failure Overall temperatures and status of system temperature readings.
# TYPE dell_hw_chassis_temps_min_failure gauge
dell_hw_chassis_temps_min_failure{component="CPU1_Temp"} 3
dell_hw_chassis_temps_min_failure{component="CPU2_Temp"} 3
dell_hw_chassis_temps_min_failure{component="System_Board_Exhaust_Temp"} 3
dell_hw_chassis_temps_min_failure{component="System_Board_Inlet_Temp"} -7
# HELP dell_hw_chassis_temps_min_warning Overall temperatures and status of system temperature readings.
# TYPE dell_hw_chassis_temps_min_warning gauge
dell_hw_chassis_temps_min_warning{component="System_Board_Exhaust_Temp"} 8
dell_hw_chassis_temps_min_warning{component="System_Board_Inlet_Temp"} 3
# HELP dell_hw_chassis_temps_reading Overall temperatures and status of system temperature readings.
# TYPE dell_hw_chassis_temps_reading gauge
dell_hw_chassis_temps_reading{component="CPU1_Temp"} 32
dell_hw_chassis_temps_reading{component="CPU2_Temp"} 34
dell_hw_chassis_temps_reading{component="System_Board_Exhaust_Temp"} 29
dell_hw_chassis_temps_reading{component="System_Board_Inlet_Temp"} 19
# HELP dell_hw_chassis_volts_reading Overall volts and status of power supply volt readings.
# TYPE dell_hw_chassis_volts_reading gauge
dell_hw_chassis_volts_reading{component="PS1_Voltage_1"} 228
dell_hw_chassis_volts_reading{component="PS2_Voltage_2"} 228
# HELP dell_hw_chassis_volts_status Overall volts and status of power supply volt readings.
# TYPE dell_hw_chassis_volts_status gauge
dell_hw_chassis_volts_status{component="CPU1_FIVR_PG"} 0
dell_hw_chassis_volts_status{component="CPU1_MEM012_VDDQ_PG"} 0
dell_hw_chassis_volts_status{component="CPU1_MEM012_VPP_PG"} 0
dell_hw_chassis_volts_status{component="CPU1_MEM012_VTT_PG"} 0
dell_hw_chassis_volts_status{component="CPU1_MEM345_VDDQ_PG"} 0
dell_hw_chassis_volts_status{component="CPU1_MEM345_VPP_PG"} 0
dell_hw_chassis_volts_status{component="CPU1_MEM345_VTT_PG"} 0
dell_hw_chassis_volts_status{component="CPU1_VCCIO_PG"} 0
dell_hw_chassis_volts_status{component="CPU1_VCORE_PG"} 0
dell_hw_chassis_volts_status{component="CPU1_VSA_PG"} 0
dell_hw_chassis_volts_status{component="CPU2_FIVR_PG"} 0
dell_hw_chassis_volts_status{component="CPU2_MEM012_VDDQ_PG"} 0
dell_hw_chassis_volts_status{component="CPU2_MEM012_VPP_PG"} 0
dell_hw_chassis_volts_status{component="CPU2_MEM012_VTT_PG"} 0
dell_hw_chassis_volts_status{component="CPU2_MEM345_VDDQ_PG"} 0
dell_hw_chassis_volts_status{component="CPU2_MEM345_VPP_PG"} 0
dell_hw_chassis_volts_status{component="CPU2_MEM345_VTT_PG"} 0
dell_hw_chassis_volts_status{component="CPU2_VCCIO_PG"} 0
dell_hw_chassis_volts_status{component="CPU2_VCORE_PG"} 0
dell_hw_chassis_volts_status{component="CPU2_VSA_PG"} 0
dell_hw_chassis_volts_status{component="PS1_Voltage_1"} 0
dell_hw_chassis_volts_status{component="PS2_Voltage_2"} 0
dell_hw_chassis_volts_status{component="System_Board_1.8V_SW_PG"} 0
dell_hw_chassis_volts_status{component="System_Board_2.5V_SW_PG"} 0
dell_hw_chassis_volts_status{component="System_Board_3.3V_A_PG"} 0
dell_hw_chassis_volts_status{component="System_Board_3.3V_B_PG"} 0
dell_hw_chassis_volts_status{component="System_Board_5V_SW_PG"} 0
dell_hw_chassis_volts_status{component="System_Board_BP0_PG"} 0
dell_hw_chassis_volts_status{component="System_Board_BP1_PG"} 0
dell_hw_chassis_volts_status{component="System_Board_BP2_PG"} 0
dell_hw_chassis_volts_status{component="System_Board_DIMM_PG"} 0
dell_hw_chassis_volts_status{component="System_Board_NDC_PG"} 0
dell_hw_chassis_volts_status{component="System_Board_PS1_PG_FAIL"} 0
dell_hw_chassis_volts_status{component="System_Board_PS2_PG_FAIL"} 0
dell_hw_chassis_volts_status{component="System_Board_PVNN_SW_PG"} 0
dell_hw_chassis_volts_status{component="System_Board_VSB11_SW_PG"} 0
dell_hw_chassis_volts_status{component="System_Board_VSBM_SW_PG"} 0
# HELP dell_hw_cmos_batteries_status Overall status of chassis batteries
# TYPE dell_hw_cmos_batteries_status gauge
dell_hw_cmos_batteries_status{id="0"} 0
# HELP dell_hw_firmware Version info of firmwares/bios.
# TYPE dell_hw_firmware gauge
dell_hw_firmware{idrac9="4.22.00.00 (build 20)"} 0
# HELP dell_hw_nic_status Connection status of network cards.
# TYPE dell_hw_nic_status gauge
dell_hw_nic_status{device="docker0",id="0"} 1
dell_hw_nic_status{device="eno1",id="2"} 0
dell_hw_nic_status{device="eno2",id="3"} 1
dell_hw_nic_status{device="eno3",id="0"} 1
dell_hw_nic_status{device="eno4",id="1"} 1
# HELP dell_hw_ps_rated_input_wattage Overall status of power supplies.
# TYPE dell_hw_ps_rated_input_wattage gauge
dell_hw_ps_rated_input_wattage{id="0"} 1260
dell_hw_ps_rated_input_wattage{id="1"} 1260
# HELP dell_hw_ps_rated_output_wattage Overall status of power supplies.
# TYPE dell_hw_ps_rated_output_wattage gauge
dell_hw_ps_rated_output_wattage{id="0"} 1100
dell_hw_ps_rated_output_wattage{id="1"} 1100
# HELP dell_hw_ps_status Overall status of power supplies.
# TYPE dell_hw_ps_status gauge
dell_hw_ps_status{id="0"} 0
dell_hw_ps_status{id="1"} 0
# HELP dell_hw_scrape_collector_duration_seconds dellhw_exporter: Duration of a collector scrape.
# TYPE dell_hw_scrape_collector_duration_seconds gauge
dell_hw_scrape_collector_duration_seconds{collector="chassis"} 2.516581654
dell_hw_scrape_collector_duration_seconds{collector="chassis_batteries"} 2.408411858
dell_hw_scrape_collector_duration_seconds{collector="fans"} 0.391451297
dell_hw_scrape_collector_duration_seconds{collector="firmwares"} 0.690688131
dell_hw_scrape_collector_duration_seconds{collector="memory"} 2.392257815
dell_hw_scrape_collector_duration_seconds{collector="nics"} 2.501337814
dell_hw_scrape_collector_duration_seconds{collector="processors"} 0.490092584
dell_hw_scrape_collector_duration_seconds{collector="ps"} 0.488725968
dell_hw_scrape_collector_duration_seconds{collector="ps_amps_sysboard_pwr"} 2.410220934
dell_hw_scrape_collector_duration_seconds{collector="storage_battery"} 0.596624274
dell_hw_scrape_collector_duration_seconds{collector="storage_controller"} 1.887127915
dell_hw_scrape_collector_duration_seconds{collector="storage_enclosure"} 0.590577123
dell_hw_scrape_collector_duration_seconds{collector="storage_pdisk"} 2.294178837
dell_hw_scrape_collector_duration_seconds{collector="storage_vdisk"} 0.69400828
dell_hw_scrape_collector_duration_seconds{collector="system"} 0.587728339
dell_hw_scrape_collector_duration_seconds{collector="temps"} 0.488827354
dell_hw_scrape_collector_duration_seconds{collector="volts"} 0.491565389
# HELP dell_hw_scrape_collector_success dellhw_exporter: Whether a collector succeeded.
# TYPE dell_hw_scrape_collector_success gauge
dell_hw_scrape_collector_success{collector="chassis"} 1
dell_hw_scrape_collector_success{collector="chassis_batteries"} 1
dell_hw_scrape_collector_success{collector="fans"} 1
dell_hw_scrape_collector_success{collector="firmwares"} 1
dell_hw_scrape_collector_success{collector="memory"} 1
dell_hw_scrape_collector_success{collector="nics"} 1
dell_hw_scrape_collector_success{collector="processors"} 1
dell_hw_scrape_collector_success{collector="ps"} 1
dell_hw_scrape_collector_success{collector="ps_amps_sysboard_pwr"} 1
dell_hw_scrape_collector_success{collector="storage_battery"} 1
dell_hw_scrape_collector_success{collector="storage_controller"} 1
dell_hw_scrape_collector_success{collector="storage_enclosure"} 1
dell_hw_scrape_collector_success{collector="storage_pdisk"} 1
dell_hw_scrape_collector_success{collector="storage_vdisk"} 1
dell_hw_scrape_collector_success{collector="system"} 1
dell_hw_scrape_collector_success{collector="temps"} 1
dell_hw_scrape_collector_success{collector="volts"} 1
# HELP dell_hw_storage_controller_status Overall status of storage controllers.
# TYPE dell_hw_storage_controller_status gauge
dell_hw_storage_controller_status{controller_name="Dell HBA330 Mini (Slot Embedded)",id="0"} 0
# HELP dell_hw_storage_enclosure_status Overall status of storage enclosures.
# TYPE dell_hw_storage_enclosure_status gauge
dell_hw_storage_enclosure_status{controller_name="Dell HBA330 Mini (Embedded)",enclosure="0_1"} 0
# HELP dell_hw_storage_pdisk_status Overall status of physical disks + failure prediction (if available).
# TYPE dell_hw_storage_pdisk_status gauge
dell_hw_storage_pdisk_status{controller="0",controller_name="Dell HBA330 Mini (Embedded)",disk="0_1_10"} 0
dell_hw_storage_pdisk_status{controller="0",controller_name="Dell HBA330 Mini (Embedded)",disk="0_1_11"} 0
# HELP dell_hw_storage_pdisk_failure_predicted Overall status of physical disks + failure prediction (if available).
# TYPE dell_hw_storage_pdisk_failure_predicted gauge
dell_hw_storage_pdisk_failure_predicted{controller="0",controller_name="Dell HBA330 Mini (Embedded)",disk="0_1_10"} 0
dell_hw_storage_pdisk_failure_predicted{controller="0",controller_name="Dell HBA330 Mini (Embedded)",disk="0_1_11"} 0
# HELP dell_hw_system_status Overall status of system components.
# TYPE dell_hw_system_status gauge
dell_hw_system_status{component="Main_System_Chassis"} 0
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 4.4303e-05
go_gc_duration_seconds{quantile="0.25"} 8.1471e-05
go_gc_duration_seconds{quantile="0.5"} 0.000220224
go_gc_duration_seconds{quantile="0.75"} 0.000391777
go_gc_duration_seconds{quantile="1"} 0.00080469
go_gc_duration_seconds_sum 0.692176499
go_gc_duration_seconds_count 1641
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 9
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.13.11"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 2.77076e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 2.049096176e+09
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.607338e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 9.131498e+06
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 1.1641209236887955e-05
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 2.394112e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 2.77076e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 5.9318272e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 4.612096e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 5914
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 5.7663488e+07
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 6.3930368e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.60153799910655e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 9.137412e+06
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 138880
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 147456
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 71944
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 98304
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 5.202448e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 5.38683e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 3.178496e+06
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 3.178496e+06
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 7.6742904e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 55
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 50.7
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 64
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 2.9724672e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.6014587121e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 1.22191872e+08
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes -1
```
