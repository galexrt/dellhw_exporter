# hardware_exporter

## Prometheus exporter for Dell Hardware components

*Supports Dell OMSA 7.4*

This exporter wraps the "omreport" command from Dell OMSA. If you can't run omreport on your system, the exporter won't export any metrics.


	Usage:
	  hardware_exporter [flags]
	  hardware_exporter [command]
	
	Available Commands:
	  version     Print the version number of hardware_exporter
	  help        Help about any command
	
	Flags:
	  -c, --collect="chassis,fans,memory,processors,ps,ps_amps_sysboard_pwr,storage_battery,storage_enclosure,storage_controller,storage_vdisk,system,temps,volts": Comma-separated list of collectors to use.
	  -h, --help[=false]: help for hardware_exporter
	  -L, --loglevel="info": Set log level
	  -l, --web.listen="127.0.0.1": Address on which to expose metrics and web interface.
	  -m, --web.path="/metrics": Path under which to expose metrics.
	  -p, --web.port="4242": Port on which to expose metrics.
	
	
	Use "hardware_exporter [command] --help" for more information about a command.