/*
Copyright 2021 The dellhw_exporter Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package omreport

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// Prefixes
	omReportControllerNamePrefix = "Controller "

	// Labels
	controllerLabel     = "controller"
	controllerNameLabel = "controller_name"
)

// Options allow to set options for the OMReport package
type Options struct {
	OMReportExecutable string
}

// OMReport contains the Options and a Reader to mock outputs during development
type OMReport struct {
	Options *Options
	Reader  func(func([]string), string, ...string) error
}

// Value contains a metrics name, value and labels
type Value struct {
	Name   string
	Value  string
	Labels map[string]string
}

const (
	// DefaultOMReportExecutable the default path of the omreport binary
	DefaultOMReportExecutable = "/opt/dell/srvadmin/bin/omreport"

	indexField = "Index"
)

// New returns a new *OMReport
func New(opts *Options) *OMReport {
	if opts.OMReportExecutable == "" {
		opts.OMReportExecutable = DefaultOMReportExecutable
	}
	return &OMReport{
		Options: opts,
		Reader:  readOmreport,
	}
}

func readOmreport(f func([]string), omreportExecutable string, args ...string) error {
	args = append(args, "-fmt", "ssv")
	return readCommand(func(line string) error {
		sp := strings.Split(line, ";")
		for i, s := range sp {
			sp[i] = clean(s)
		}
		f(sp)
		return nil
	}, omreportExecutable, args...)
}

func (or *OMReport) getOMReportExecutable() string {
	if or.Options != nil {
		return or.Options.OMReportExecutable
	}
	return DefaultOMReportExecutable
}

func (or *OMReport) readReport(f func([]string), omreportExecutable string, args ...string) error {
	return or.Reader(f, omreportExecutable, args...)
}

// Chassis returns the chassis status
func (or *OMReport) Chassis() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(fields []string) {
		if len(fields) != 2 || fields[0] == "SEVERITY" {
			return
		}
		component := strings.Replace(fields[1], " ", "_", -1)
		values = append(values, Value{
			Name:   "chassis_status",
			Value:  severity(fields[0]),
			Labels: map[string]string{"component": component},
		})
	}, or.getOMReportExecutable(), "chassis")
	return values, err
}

// Fans returns the fan status and if supported RPM reading
func (or *OMReport) Fans() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(fields []string) {
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
	}, or.getOMReportExecutable(), "chassis", "fans")
	return values, err
}

// Memory returns the memory status
func (or *OMReport) Memory() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(fields []string) {
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
	}, or.getOMReportExecutable(), "chassis", "memory")
	return values, err
}

// System returns the system status
func (or *OMReport) System() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(fields []string) {
		if len(fields) != 2 || fields[0] == "SEVERITY" {
			return
		}
		component := strings.Replace(fields[1], " ", "_", -1)
		values = append(values, Value{
			Name:   "system_status",
			Value:  severity(fields[0]),
			Labels: map[string]string{"component": component},
		})
	}, or.getOMReportExecutable(), "system")
	return values, err
}

// StorageBattery returns the storage battery ("RAID batteries")
func (or *OMReport) StorageBattery() ([]Value, error) {
	values := []Value{}
	controllerName := "N/A"
	err := or.readReport(func(fields []string) {
		if len(fields) == 1 && strings.HasPrefix(fields[0], omReportControllerNamePrefix) {
			controllerName = strings.TrimPrefix(fields[0], omReportControllerNamePrefix)
			return
		} else if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		values = append(values, Value{
			Name:  "storage_battery_status",
			Value: severity(fields[1]),
			Labels: map[string]string{
				controllerLabel:     id,
				controllerNameLabel: controllerName,
			},
		})
	}, or.getOMReportExecutable(), "storage", "battery")
	return values, err
}

// StorageController returns the storage controller status
func (or *OMReport) StorageController() ([]Value, error) {
	values := []Value{}
	controllerName := "N/A"
	err := or.readReport(func(fields []string) {
		// Use the fields instead of the "single line with the controller name on it"
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		controllerName = fmt.Sprintf("%s (Slot %s)", fields[2], fields[3])
		or.StoragePdisk(fields[0])
		id := strings.Replace(fields[0], ":", "_", -1)
		values = append(values, Value{
			Name:  "storage_controller_status",
			Value: severity(fields[1]),
			Labels: map[string]string{
				"id":                id,
				controllerNameLabel: controllerName,
			},
		})
	}, or.getOMReportExecutable(), "storage", "controller")
	return values, err
}

// StorageEnclosure returns the storage enclosure status
func (or *OMReport) StorageEnclosure() ([]Value, error) {
	values := []Value{}
	controllerName := "N/A"
	err := or.readReport(func(fields []string) {
		if len(fields) == 1 && strings.HasPrefix(fields[0], "Enclosure(s) on Controller ") {
			controllerName = strings.TrimPrefix(fields[0], "Enclosure(s) on Controller ")
			return
		} else if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		values = append(values, Value{
			Name:  "storage_enclosure_status",
			Value: severity(fields[1]),
			Labels: map[string]string{
				"enclosure":         id,
				controllerNameLabel: controllerName,
			},
		})
	}, or.getOMReportExecutable(), "storage", "enclosure")
	return values, err
}

// StoragePdisk is called from the controller func, since it needs the encapsulating IDs.
func (or *OMReport) StoragePdisk(cid string) ([]Value, error) {
	values := []Value{}
	controllerName := "N/A"
	err := or.readReport(func(fields []string) {
		if len(fields) == 1 && strings.HasPrefix(fields[0], omReportControllerNamePrefix) {
			controllerName = strings.TrimPrefix(fields[0], omReportControllerNamePrefix)
			return
		} else if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		// Need to find out what the various ID formats might be
		id := strings.Replace(fields[0], ":", "_", -1)

		values = append(values, Value{
			Name:  "storage_pdisk_status",
			Value: severity(fields[1]),
			Labels: map[string]string{
				controllerLabel:     cid,
				"disk":              id,
				controllerNameLabel: controllerName,
			},
		})

		values = append(values, Value{
			Name:  "storage_pdisk_state",
			Value: pdiskState(fields[3]),
			Labels: map[string]string{
				controllerLabel:     cid,
				"disk":              id,
				controllerNameLabel: controllerName,
			},
		})

		if len(fields) > 8 {
			values = append(values, Value{
				Name:  "storage_pdisk_failure_predicted",
				Value: yesNoToBool(fields[9]),
				Labels: map[string]string{
					controllerLabel:     cid,
					"disk":              id,
					controllerNameLabel: controllerName,
				},
			})
			values = append(values, Value{
				Name:  "storage_pdisk_remaining_rated_write_endurance",
				Value: getNumberFromString(fields[8]),
				Labels: map[string]string{
					controllerLabel:     cid,
					"disk":              id,
					controllerNameLabel: controllerName,
				},
			})
			if fields[15] == "Yes" {
				values = append(values, Value{
					Name:  "storage_pdisk_storage_encrypted",
					Value: yesNoToBool(fields[16]),
					Labels: map[string]string{
						controllerLabel:     cid,
						"disk":              id,
						controllerNameLabel: controllerName,
					},
				})
			}
		}
	}, or.getOMReportExecutable(), "storage", "pdisk", "controller="+cid)
	return values, err
}

// StorageVdisk returns the storage vdisk status
func (or *OMReport) StorageVdisk() ([]Value, error) {
	values := []Value{}
	controllerName := "N/A"
	err := or.readReport(func(fields []string) {
		if len(fields) == 1 && strings.HasPrefix(fields[0], omReportControllerNamePrefix) {
			controllerName = strings.TrimPrefix(fields[0], omReportControllerNamePrefix)
			return
		} else if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		values = append(values, Value{
			Name:  "storage_vdisk_status",
			Value: severity(fields[1]),
			Labels: map[string]string{
				"vdisk":             id,
				"vdisk_name":        fields[2],
				controllerNameLabel: controllerName,
			},
		})

		values = append(values, Value{
			Name:  "storage_vdisk_state",
			Value: vdiskState(fields[3]),
			Labels: map[string]string{
				"vdisk":             id,
				"vdisk_name":        fields[2],
				controllerNameLabel: controllerName,
			},
		})

		if len(fields) > 17 {
			values = append(values, Value{
				Name:  "storage_vdisk_raidlevel",
				Value: getNumberFromString(fields[6]),
				Labels: map[string]string{
					"vdisk":             id,
					"vdisk_name":        fields[2],
					controllerNameLabel: controllerName,
				},
			})

			values = append(values, Value{
				Name:  "storage_vdisk_read_policy",
				Value: vdiskReadPolicy(fields[13]),
				Labels: map[string]string{
					"vdisk":             id,
					"vdisk_name":        fields[2],
					controllerNameLabel: controllerName,
				},
			})

			values = append(values, Value{
				Name:  "storage_vdisk_write_policy",
				Value: vdiskWritePolicy(fields[14]),
				Labels: map[string]string{
					"vdisk":             id,
					"vdisk_name":        fields[2],
					controllerNameLabel: controllerName,
				},
			})

			values = append(values, Value{
				Name:  "storage_vdisk_cache_policy",
				Value: vdiskCachePolicy(fields[15]),
				Labels: map[string]string{
					"vdisk":             id,
					"vdisk_name":        fields[2],
					controllerNameLabel: controllerName,
				},
			})
		} else {
			values = append(values, Value{
				Name:  "storage_vdisk_raidlevel",
				Value: getNumberFromString(fields[5]),
				Labels: map[string]string{
					"vdisk":             id,
					"vdisk_name":        fields[2],
					controllerNameLabel: controllerName,
				},
			})

			values = append(values, Value{
				Name:  "storage_vdisk_read_policy",
				Value: vdiskReadPolicy(fields[12]),
				Labels: map[string]string{
					"vdisk":             id,
					"vdisk_name":        fields[2],
					controllerNameLabel: controllerName,
				},
			})

			values = append(values, Value{
				Name:  "storage_vdisk_write_policy",
				Value: vdiskWritePolicy(fields[13]),
				Labels: map[string]string{
					"vdisk":             id,
					"vdisk_name":        fields[2],
					controllerNameLabel: controllerName,
				},
			})

			values = append(values, Value{
				Name:  "storage_vdisk_cache_policy",
				Value: vdiskCachePolicy(fields[14]),
				Labels: map[string]string{
					"vdisk":             id,
					"vdisk_name":        fields[2],
					controllerNameLabel: controllerName,
				},
			})
		}
	}, or.getOMReportExecutable(), "storage", "vdisk")
	return values, err
}

// Ps returns the power supply state and if supported input/output wattage
func (or *OMReport) Ps() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(fields []string) {
		if len(fields) < 3 || fields[0] == indexField {
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
					Name:   "ps_rated_input_wattage",
					Value:  iWattage,
					Labels: ts,
				})
			}
		}
		if fields[5] != "" {
			oWattage, err := extract(fields[5], "W")
			if err == nil {
				values = append(values, Value{
					Name:   "ps_rated_output_wattage",
					Value:  oWattage,
					Labels: ts,
				})
			}
		}
	}, or.getOMReportExecutable(), "chassis", "pwrsupplies")
	return values, err
}

// Nics returns the connection status of the NICs
func (or *OMReport) Nics(nicList ...string) ([]Value, error) {
	values := []Value{}
	monitorAllNics := false
	monitoredNics := make(map[string]bool)

	if len(nicList) == 0 {
		monitorAllNics = true
	}

	for _, nic := range nicList {
		monitoredNics[nic] = true
	}

	err := or.readReport(func(fields []string) {
		if len(fields) < 5 || fields[0] == indexField {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		ts := map[string]string{"id": id, "device": fields[1]}
		var ret string
		if fields[4] == "Connected" || fields[4] == "Full" || fields[4] == "Not Applicable" {
			ret = "0"
		} else {
			ret = "1"
		}
		dev := ts["device"]
		if monitorAllNics || monitoredNics[dev] {
			values = append(values, Value{
				Name:   "nic_status",
				Value:  ret,
				Labels: ts,
			})
		}
	}, or.getOMReportExecutable(), "chassis", "nics")
	return values, err
}

// PsAmpsSysboardPwr returns the power supply system board amps power consumption
func (or *OMReport) PsAmpsSysboardPwr() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(fields []string) {
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
	}, or.getOMReportExecutable(), "chassis", "pwrmonitoring")
	return values, err
}

// Processors returns the processors status
func (or *OMReport) Processors() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(fields []string) {
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
	}, or.getOMReportExecutable(), "chassis", "processors")
	return values, err
}

// Temps returns the temperatures for the chassis including the min and max,
// for the max value, warning and failure thresholds are returned
func (or *OMReport) Temps() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(fields []string) {
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
		minWarningThreshold := strings.Fields(fields[4])
		if len(minWarningThreshold) == 2 && minWarningThreshold[1] == "C" {
			values = append(values, Value{
				Name:   "chassis_temps_min_warning",
				Value:  minWarningThreshold[0],
				Labels: ts,
			})
		}
		maxWarningThreshold := strings.Fields(fields[5])
		if len(maxWarningThreshold) == 2 && maxWarningThreshold[1] == "C" {
			values = append(values, Value{
				Name:   "chassis_temps_max_warning",
				Value:  maxWarningThreshold[0],
				Labels: ts,
			})
		}
		minFailureThreshold := strings.Fields(fields[6])
		if len(minFailureThreshold) == 2 && minFailureThreshold[1] == "C" {
			values = append(values, Value{
				Name:   "chassis_temps_min_failure",
				Value:  minFailureThreshold[0],
				Labels: ts,
			})
		}
		maxFailureThreshold := strings.Fields(fields[7])
		if len(maxFailureThreshold) == 2 && maxFailureThreshold[1] == "C" {
			values = append(values, Value{
				Name:   "chassis_temps_max_failure",
				Value:  maxFailureThreshold[0],
				Labels: ts,
			})
		}
	}, or.getOMReportExecutable(), "chassis", "temps")
	return values, err
}

// Volts returns the chassis volts statud and if support reading
func (or *OMReport) Volts() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(fields []string) {
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
	}, or.getOMReportExecutable(), "chassis", "volts")
	return values, err
}

// ChassisBatteries retursn the chassis batteries status
func (or *OMReport) ChassisBatteries() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(fields []string) {
		if len(fields) < 4 || fields[0] == indexField {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		ts := map[string]string{"id": id}

		values = append(values, Value{
			Name:   "cmos_batteries_status",
			Value:  severity(fields[1]),
			Labels: ts,
		})
	}, or.getOMReportExecutable(), "chassis", "batteries")
	return values, err
}

// ChassisBios returns the bios version name
func (or *OMReport) ChassisBios() ([]Value, error) {
	value := Value{
		Name:   "bios",
		Value:  "0",
		Labels: map[string]string{},
	}

	err := or.readReport(func(fields []string) {

		if len(fields) != 2 {
			return
		}
		id := strings.Replace(strings.ToLower(fields[0]), " ", "_", -1)
		value.Labels[id] = strings.ToLower(fields[1])

	}, or.getOMReportExecutable(), "chassis", "bios")

	values := []Value{}
	values = append(values, value)

	return values, err
}

// ChassisFirmware returns the firmware revisions
func (or *OMReport) ChassisFirmware() ([]Value, error) {
	value := Value{
		Name:   "firmware",
		Value:  "0",
		Labels: map[string]string{},
	}

	err := or.readReport(func(fields []string) {

		if len(fields) != 2 {
			return
		}
		id := strings.Replace(strings.ToLower(fields[0]), " ", "_", -1)
		value.Labels[id] = strings.ToLower(fields[1])

	}, or.getOMReportExecutable(), "chassis", "firmware")

	values := []Value{}
	values = append(values, value)

	return values, err
}
