package omreport

import (
	"strconv"
	"strings"
)

type Options struct {
	OMReportExecutable string
}

type OMReport struct {
	Options *Options
}

type Value struct {
	Name   string
	Value  string
	Labels map[string]string
}

func New(opts *Options) *OMReport {

	return &OMReport{
		Options: opts,
	}
}

func (or *OMReport) readOmreport(f func([]string), args ...string) {
	args = append(args, "-fmt", "ssv")
	_ = readCommand(func(line string) error {
		sp := strings.Split(line, ";")
		for i, s := range sp {
			sp[i] = clean(s)
		}
		f(sp)
		return nil
	}, or.Options.OMReportExecutable, args...)
}

func (or *OMReport) Chassis() ([]Value, error) {
	values := []Value{}
	or.readOmreport(func(fields []string) {
		if len(fields) != 2 || fields[0] == "SEVERITY" {
			return
		}
		component := strings.Replace(fields[1], " ", "_", -1)
		values = append(values, Value{
			Name:   "chassis_status",
			Value:  severity(fields[0]),
			Labels: map[string]string{"component": component},
		})
	}, "chassis")
	return values, nil
}

func (or *OMReport) System() ([]Value, error) {
	values := []Value{}
	or.readOmreport(func(fields []string) {
		if len(fields) != 2 || fields[0] == "SEVERITY" {
			return
		}
		component := strings.Replace(fields[1], " ", "_", -1)
		values = append(values, Value{
			Name:   "system_status",
			Value:  severity(fields[0]),
			Labels: map[string]string{"component": component},
		})
	}, "system")
	return values, nil
}

func (or *OMReport) StorageEnclosure() ([]Value, error) {
	values := []Value{}
	or.readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		values = append(values, Value{
			Name:   "storage_enclosure_status",
			Value:  severity(fields[1]),
			Labels: map[string]string{"enclosure": id},
		})
	}, "storage", "enclosure")
	return values, nil
}

func (or *OMReport) StorageVdisk() ([]Value, error) {
	values := []Value{}
	or.readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		values = append(values, Value{
			Name:   "storage_vdisk_status",
			Value:  severity(fields[1]),
			Labels: map[string]string{"vdisk": id},
		})
	}, "storage", "vdisk")
	return values, nil
}

func (or *OMReport) Ps() ([]Value, error) {
	values := []Value{}
	or.readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "Index" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		ts := map[string]string{"id": id}
		values = append(values, Value{
			Name:   "ps_status",
			Value:  severity(fields[1]),
			Labels: ts,
		})
		if len(fields) < 6 {
			return
		}
		if fields[4] != "" {
			iWattage, err := extract(fields[4], "W")
			if err == nil {
				values = append(values, Value{
					Name:   "rated_input_wattage",
					Value:  iWattage,
					Labels: ts,
				})
			}
		}
		if fields[5] != "" {
			oWattage, err := extract(fields[5], "W")
			if err == nil {
				values = append(values, Value{
					Name:   "rated_output_wattage",
					Value:  oWattage,
					Labels: ts,
				})
			}
		}
	}, "chassis", "pwrsupplies")
	return values, nil
}

func (or *OMReport) PsAmpsSysboardPwr() ([]Value, error) {
	values := []Value{}
	or.readOmreport(func(fields []string) {
		if len(fields) == 2 && strings.Contains(fields[0], "Current") {
			iFields := strings.Split(fields[0], "Current")
			vFields := strings.Fields(fields[1])
			if len(iFields) < 2 && len(vFields) < 2 {
				return
			}
			id := strings.Replace(iFields[0], " ", "", -1)
			values = append(values, Value{
				Name:   "chassis_current_reading",
				Value:  vFields[0],
				Labels: map[string]string{"pwrsupply": id},
			})
		} else if len(fields) == 6 && (fields[2] == "System Board Pwr Consumption" || fields[2] == "System Board System Level") {
			vFields := strings.Fields(fields[3])
			warnFields := strings.Fields(fields[4])
			failFields := strings.Fields(fields[5])
			if len(vFields) < 2 || len(warnFields) < 2 || len(failFields) < 2 {
				return
			}
			values = append(values, Value{
				Name:   "chassis_power_reading",
				Value:  vFields[0],
				Labels: nil,
			})
			values = append(values, Value{
				Name:   "chassis_power_warn_level",
				Value:  warnFields[0],
				Labels: nil,
			})
			values = append(values, Value{
				Name:   "chassis_power_fail_level",
				Value:  failFields[0],
				Labels: nil,
			})
		}
	}, "chassis", "pwrmonitoring")
	return values, nil
}

func (or *OMReport) StorageBattery() ([]Value, error) {
	values := []Value{}
	or.readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		values = append(values, Value{
			Name:   "storage_battery_status",
			Value:  severity(fields[1]),
			Labels: map[string]string{"controller": id},
		})
	}, "storage", "battery")
	return values, nil
}

func (or *OMReport) StorageController() ([]Value, error) {
	values := []Value{}
	or.readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		or.StoragePdisk(fields[0])
		id := strings.Replace(fields[0], ":", "_", -1)
		values = append(values, Value{
			Name:   "storage_controller_status",
			Value:  severity(fields[1]),
			Labels: map[string]string{"id": id},
		})
	}, "storage", "controller")
	return values, nil
}

// omreportStoragePdisk is called from the controller func, since it needs the encapsulating id.
func (or *OMReport) StoragePdisk(id string) {
	values := []Value{}
	or.readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		//Need to find out what the various ID formats might be
		id := strings.Replace(fields[0], ":", "_", -1)
		ts := map[string]string{"disk": id}
		values = append(values, Value{
			Name:   "storage_pdisk_status",
			Value:  severity(fields[1]),
			Labels: ts,
		})
	}, "storage", "pdisk", "controller="+id)
}

func (or *OMReport) Processors() ([]Value, error) {
	values := []Value{}
	or.readOmreport(func(fields []string) {
		if len(fields) != 8 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		values = append(values, Value{
			Name:   "chassis_processor_status",
			Value:  severity(fields[1]),
			Labels: map[string]string{"processor": replace(fields[2])},
		})
	}, "chassis", "processors")
	return values, nil
}

func (or *OMReport) Fans() ([]Value, error) {
	values := []Value{}
	or.readOmreport(func(fields []string) {
		if len(fields) != 8 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		ts := map[string]string{"fan": replace(fields[2])}
		values = append(values, Value{
			Name:   "chassis_fan_status",
			Value:  severity(fields[1]),
			Labels: ts,
		})
		fs := strings.Fields(fields[3])
		if len(fs) == 2 && fs[1] == "RPM" {
			values = append(values, Value{
				Name:   "chassis_fan_reading",
				Value:  fs[0],
				Labels: ts,
			})
		}
	}, "chassis", "fans")
	return values, nil
}

func (or *OMReport) Memory() ([]Value, error) {
	values := []Value{}
	or.readOmreport(func(fields []string) {
		if len(fields) != 5 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		values = append(values, Value{
			Name:   "chassis_memory_status",
			Value:  severity(fields[1]),
			Labels: map[string]string{"memory": replace(fields[2])},
		})
	}, "chassis", "memory")
	return values, nil
}

func (or *OMReport) Temps() ([]Value, error) {
	values := []Value{}
	or.readOmreport(func(fields []string) {
		if len(fields) != 8 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		ts := map[string]string{"component": replace(fields[2])}
		values = append(values, Value{
			Name:   "chassis_temps",
			Value:  severity(fields[1]),
			Labels: ts,
		})
		fs := strings.Fields(fields[3])
		if len(fs) == 2 && fs[1] == "C" {
			values = append(values, Value{
				Name:   "chassis_temps_reading",
				Value:  fs[0],
				Labels: ts,
			})
		}
	}, "chassis", "temps")
	return values, nil
}

func (or *OMReport) Volts() ([]Value, error) {
	values := []Value{}
	or.readOmreport(func(fields []string) {
		if len(fields) != 8 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		ts := map[string]string{"component": replace(fields[2])}
		values = append(values, Value{
			Name:   "chassis_volts_status",
			Value:  severity(fields[1]),
			Labels: ts,
		})
		if i, err := extract(fields[3], "V"); err == nil {
			values = append(values, Value{
				Name:   "chassis_volts_reading",
				Value:  i,
				Labels: ts,
			})
		}
	}, "chassis", "volts")
	return values, nil
}
