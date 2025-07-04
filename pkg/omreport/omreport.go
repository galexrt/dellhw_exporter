/*
Copyright 2024 The dellhw_exporter Authors. All rights reserved.

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
	"log/slog"
	"strconv"
	"strings"
)

const (
	// DefaultOMReportExecutable the default path of the omreport binary
	DefaultOMReportExecutable = "/opt/dell/srvadmin/bin/omreport"

	// Prefixes
	storageControllerNamePrefix = "Controller "
	storageEnclosureNamePrefix  = "Enclosure(s) on Controller "

	// Labels
	controllerLabel     = "controller"
	controllerNameLabel = "controller_name"
)

type ReaderMode int

const (
	DynamicReaderMode ReaderMode = iota
	KeyValueReaderMode
	TableReaderMode
)

// Options allow to set options for the OMReport package
type Options struct {
	OMReportExecutable string
}

// OMReport contains the Options and a Reader to mock outputs during development
type OMReport struct {
	Options *Options
	Reader  func(f func(Output), mode ReaderMode, cmd string, args ...string) error
}

// Value contains a metrics name, value and labels
type Value struct {
	Name   string
	Value  string
	Labels map[string]string
}

func (v Value) String() string {
	labels := []string{}
	for k, v := range v.Labels {
		labels = append(labels, fmt.Sprintf("%s=%s", k, v))
	}

	return fmt.Sprintf("%q{%q}=%s", v.Name, labels, v.Value)
}

// SetLogger
func SetLogger(l *slog.Logger) {
	logger = l
}

// New returns a new OMReport struct
func New(opts *Options) *OMReport {
	if opts.OMReportExecutable == "" {
		opts.OMReportExecutable = DefaultOMReportExecutable
	}

	return &OMReport{
		Options: opts,
		Reader:  readOmreport,
	}
}

func readOmreport(f func(Output), mode ReaderMode, omreportExecutable string, args ...string) error {
	args = append(args, "-fmt", "ssv")
	return readCommand(func(input string) error {
		output := parseOutput(mode, input)

		f(output)

		return nil
	}, omreportExecutable, args...)
}

func (or *OMReport) getOMReportExecutable() string {
	if or.Options != nil {
		return or.Options.OMReportExecutable
	}

	return DefaultOMReportExecutable
}

func (or *OMReport) readReport(f func(Output), mode ReaderMode, omreportExecutable string, args ...string) error {
	return or.Reader(f, mode, omreportExecutable, args...)
}

// Chassis returns the chassis status
func (or *OMReport) Chassis() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if !hasKeys(fields, "severity", "component") {
					continue
				}

				component := strings.Replace(fields["component"], " ", "_", -1)
				values = append(values, Value{
					Name:   "chassis_status",
					Value:  severity(fields["severity"]),
					Labels: map[string]string{"component": component},
				})
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "chassis")
	return values, err
}

// ChassisInfo returns the chassis information
func (or *OMReport) ChassisInfo() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if !hasKeys(fields, "chassis_model") {
					continue
				}

				model := strings.Replace(fields["chassis_model"], " ", "_", -1)
				values = append(values, Value{
					Name:   "chassis_info",
					Value:  "0",
					Labels: map[string]string{"chassis_model": model},
				})
			}
		}
	}, KeyValueReaderMode, or.getOMReportExecutable(), "chassis", "information")
	return values, err
}

// Fans returns the fan status and if supported RPM reading
func (or *OMReport) Fans() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if _, err := strconv.Atoi(fields["index"]); err != nil {
					continue
				}

				ts := map[string]string{"fan": replace(fields["probe_name"])}
				values = append(values, Value{
					Name:   "chassis_fan_status",
					Value:  severity(fields["status"]),
					Labels: ts,
				})

				fs := strings.Fields(fields["reading"])
				if len(fs) == 2 && fs[1] == "RPM" {
					values = append(values, Value{
						Name:   "chassis_fan_reading",
						Value:  fs[0],
						Labels: ts,
					})
				}
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "chassis", "fans")
	return values, err
}

// Memory returns the memory status
func (or *OMReport) Memory() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if len(fields) < 5 {
					continue
				}

				fs := strings.Fields(fields["size"])
				if len(fs) != 2 {
					continue
				}

				values = append(values, Value{
					Name:   "chassis_memory_status",
					Value:  severity(fields["status"]),
					Labels: map[string]string{"memory": replace(fields["connector_name"])},
				})
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "chassis", "memory")
	return values, err
}

// System returns the system status
func (or *OMReport) System() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if len(fields) != 2 || fields["severity"] == "SEVERITY" {
					continue
				}
				component := replace(fields["component"])
				values = append(values, Value{
					Name:   "system_status",
					Value:  severity(fields["severity"]),
					Labels: map[string]string{"component": component},
				})
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "system")
	return values, err
}

// StorageBattery returns the storage battery ("RAID batteries")
func (or *OMReport) StorageBattery() ([]Value, error) {
	values := []Value{}
	controllerName := "N/A"
	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if strings.HasPrefix(output.Title, storageControllerNamePrefix) {
					controllerName = strings.TrimPrefix(output.Title, storageControllerNamePrefix)
				} else if strings.HasPrefix(output.Description, storageControllerNamePrefix) {
					controllerName = strings.TrimPrefix(output.Description, storageControllerNamePrefix)
				}

				if len(fields) < 3 {
					continue
				}
				id := strings.Replace(fields["id"], ":", "_", -1)
				values = append(values, Value{
					Name:  "storage_battery_status",
					Value: severity(fields["status"]),
					Labels: map[string]string{
						controllerLabel:     id,
						controllerNameLabel: controllerName,
					},
				})
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "storage", "battery")
	return values, err
}

// StorageController returns the storage controller status
func (or *OMReport) StorageController() ([]Value, error) {
	values := []Value{}
	controllerName := "N/A"
	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if len(fields) < 3 {
					continue
				}

				controllerName = fmt.Sprintf("%s (Slot %s)", fields["name"], fields["slot_id"])
				or.StoragePdisk(fields["id"])
				id := strings.Replace(fields["id"], ":", "_", -1)
				values = append(values, Value{
					Name:  "storage_controller_status",
					Value: severity(fields["status"]),
					Labels: map[string]string{
						"id":                id,
						controllerNameLabel: controllerName,
					},
				})
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "storage", "controller")
	return values, err
}

// StorageEnclosure returns the storage enclosure status
func (or *OMReport) StorageEnclosure() ([]Value, error) {
	values := []Value{}
	controllerName := "N/A"
	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if strings.HasPrefix(output.Title, storageEnclosureNamePrefix) {
					controllerName = strings.TrimPrefix(output.Title, storageEnclosureNamePrefix)
				} else if strings.HasPrefix(output.Description, storageEnclosureNamePrefix) {
					controllerName = strings.TrimPrefix(output.Description, storageEnclosureNamePrefix)
				}

				if len(fields) < 3 {
					continue
				}

				id := strings.Replace(fields["id"], ":", "_", -1)
				values = append(values, Value{
					Name:  "storage_enclosure_status",
					Value: severity(fields["status"]),
					Labels: map[string]string{
						"enclosure":         id,
						controllerNameLabel: controllerName,
					},
				})
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "storage", "enclosure")
	return values, err
}

// StoragePdisk is called from the controller func, since it needs the encapsulating IDs.
func (or *OMReport) StoragePdisk(cid string) ([]Value, error) {
	values := []Value{}
	controllerName := "N/A"
	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if strings.HasPrefix(output.Title, storageControllerNamePrefix) {
					controllerName = strings.TrimPrefix(output.Title, storageControllerNamePrefix)
				} else if strings.HasPrefix(output.Description, storageControllerNamePrefix) {
					controllerName = strings.TrimPrefix(output.Description, storageControllerNamePrefix)
				}

				if len(fields) < 3 {
					continue
				}
				// Need to find out what the various ID formats might be
				id := strings.Replace(fields["id"], ":", "_", -1)

				values = append(values, Value{
					Name:  "storage_pdisk_status",
					Value: severity(fields["status"]),
					Labels: map[string]string{
						controllerLabel:     cid,
						"disk":              id,
						controllerNameLabel: controllerName,
					},
				})

				values = append(values, Value{
					Name:  "storage_pdisk_state",
					Value: pdiskState(fields["state"]),
					Labels: map[string]string{
						controllerLabel:     cid,
						"disk":              id,
						controllerNameLabel: controllerName,
					},
				})

				if hasKeys(fields, "Failure Predicted", "Remaining Rated Write Endurance") {
					values = append(values, Value{
						Name:  "storage_pdisk_failure_predicted",
						Value: yesNoToBool(fields["failure_predicted"]),
						Labels: map[string]string{
							controllerLabel:     cid,
							"disk":              id,
							controllerNameLabel: controllerName,
						},
					})

					values = append(values, Value{
						Name:  "storage_pdisk_remaining_rated_write_endurance",
						Value: getNumberFromString(fields["remaining_rated_write_endurance"]),
						Labels: map[string]string{
							controllerLabel:     cid,
							"disk":              id,
							controllerNameLabel: controllerName,
						},
					})
				}

				if hasKeys(fields, "cryptographic_erase_capable") {
					values = append(values, Value{
						Name:  "storage_pdisk_storage_encrypted",
						Value: yesNoToBool(fields["cryptographic_erase_capable"]),
						Labels: map[string]string{
							controllerLabel:     cid,
							"disk":              id,
							controllerNameLabel: controllerName,
						},
					})
				}
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "storage", "pdisk", "controller="+cid)
	return values, err
}

// StorageVdisk returns the storage vdisk status
func (or *OMReport) StorageVdisk() ([]Value, error) {
	values := []Value{}
	controllerName := "N/A"
	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if strings.HasPrefix(output.Title, storageControllerNamePrefix) {
					controllerName = strings.TrimPrefix(output.Title, storageControllerNamePrefix)
				} else if strings.HasPrefix(output.Description, storageControllerNamePrefix) {
					controllerName = strings.TrimPrefix(output.Description, storageControllerNamePrefix)
				}

				if len(fields) < 3 {
					continue
				}

				id := strings.Replace(fields["id"], ":", "_", -1)
				values = append(values, Value{
					Name:  "storage_vdisk_status",
					Value: severity(fields["status"]),
					Labels: map[string]string{
						"vdisk":             id,
						"vdisk_name":        fields["name"],
						controllerNameLabel: controllerName,
					},
				})

				values = append(values, Value{
					Name:  "storage_vdisk_state",
					Value: vdiskState(fields["state"]),
					Labels: map[string]string{
						"vdisk":             id,
						"vdisk_name":        fields["name"],
						controllerNameLabel: controllerName,
					},
				})

				values = append(values, Value{
					Name:  "storage_vdisk_raidlevel",
					Value: getNumberFromString(fields["layout"]),
					Labels: map[string]string{
						"vdisk":             id,
						"vdisk_name":        fields["name"],
						controllerNameLabel: controllerName,
					},
				})

				if hasKeys(fields, "read_policy") {
					values = append(values, Value{
						Name:  "storage_vdisk_read_policy",
						Value: vdiskReadPolicy(fields["read_policy"]),
						Labels: map[string]string{
							"vdisk":             id,
							"vdisk_name":        fields["name"],
							controllerNameLabel: controllerName,
						},
					})
				}

				if hasKeys(fields, "write_policy") {
					values = append(values, Value{
						Name:  "storage_vdisk_write_policy",
						Value: vdiskWritePolicy(fields["write_policy"]),
						Labels: map[string]string{
							"vdisk":             id,
							"vdisk_name":        fields["name"],
							controllerNameLabel: controllerName,
						},
					})
				}

				if hasKeys(fields, "cache_policy") {
					values = append(values, Value{
						Name:  "storage_vdisk_cache_policy",
						Value: vdiskCachePolicy(fields["cache_policy"]),
						Labels: map[string]string{
							"vdisk":             id,
							"vdisk_name":        fields["name"],
							controllerNameLabel: controllerName,
						},
					})
				}
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "storage", "vdisk")
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

	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if len(fields) < 5 {
					continue
				}

				id := strings.Replace(fields["index"], ":", "_", -1)
				ts := map[string]string{"id": id, "device": fields["interface_name"]}

				var ret string
				connStatus := "Not Applicable"
				if hasKeys(fields, "connection_status") {
					connStatus = fields["connection_status"]
				} else if hasKeys(fields, "redundancy_status") {
					connStatus = fields["redundancy_status"]
				}

				switch connStatus {
				case "Connected":
					fallthrough
				case "Full":
					fallthrough
				case "Not Applicable":
					ret = "0"

				default:
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
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "chassis", "nics")
	return values, err
}

// Ps returns the power supply state and if supported input/output wattage
func (or *OMReport) Ps() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if len(fields) < 3 {
					continue
				}

				id := strings.Replace(fields["index"], ":", "_", -1)
				ts := map[string]string{"id": id}
				values = append(values, Value{
					Name:   "ps_status",
					Value:  severity(fields["status"]),
					Labels: ts,
				})
				if len(fields) < 6 {
					continue
				}

				if hasKeys(fields, "rated_input_wattage") {
					iWattage, err := extract(fields["rated_input_wattage"], "W")
					if err == nil {
						values = append(values, Value{
							Name:   "ps_rated_input_wattage",
							Value:  iWattage,
							Labels: ts,
						})
					}
				}
				if hasKeys(fields, "maximum_output_wattage") {
					oWattage, err := extract(fields["maximum_output_wattage"], "W")
					if err == nil {
						values = append(values, Value{
							Name:   "ps_rated_output_wattage",
							Value:  oWattage,
							Labels: ts,
						})
					}
				}
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "chassis", "pwrsupplies")
	return values, err
}

// PsAmpsSysboardPwr returns the power supply system board amps power consumption
func (or *OMReport) PsAmpsSysboardPwr() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if len(fields) == 2 && strings.Contains(fields["psu"], "Current") {
					iFields := strings.Split(fields["psu"], "Current")
					vFields := strings.Fields(fields["amperage"])
					if len(iFields) < 2 && len(vFields) < 2 {
						continue
					}

					id := strings.Replace(iFields[0], " ", "", -1)
					values = append(values, Value{
						Name:   "chassis_current_reading",
						Value:  vFields[0],
						Labels: map[string]string{"pwrsupply": id},
					})
				} else if len(fields) == 6 && (fields["probe_name"] == "System Board Pwr Consumption" || fields["probe_name"] == "System Board System Level") {
					vFields := strings.Fields(fields["reading"])
					warnFields := strings.Fields(fields["warning_threshold"])
					failFields := strings.Fields(fields["failure_threshold"])
					if len(vFields) < 2 || len(warnFields) < 2 || len(failFields) < 2 {
						continue
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
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "chassis", "pwrmonitoring")
	return values, err
}

// Processors returns the processors status
func (or *OMReport) Processors() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if len(fields) != 8 {
					continue
				}

				if _, err := strconv.Atoi(fields["index"]); err != nil {
					continue
				}

				values = append(values, Value{
					Name:   "chassis_processor_status",
					Value:  severity(fields["status"]),
					Labels: map[string]string{"processor": replace(fields["connector_name"])},
				})
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "chassis", "processors")
	return values, err
}

// Temps returns the temperatures for the chassis including the min and max,
// for the max value, warning and failure thresholds are returned
func (or *OMReport) Temps() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if len(fields) != 8 {
					continue
				}

				if _, err := strconv.Atoi(fields["index"]); err != nil {
					continue
				}

				ts := map[string]string{"component": replace(fields["probe_name"])}
				values = append(values, Value{
					Name:   "chassis_temps",
					Value:  severity(fields["status"]),
					Labels: ts,
				})

				fs := strings.Fields(fields["reading"])
				if len(fs) == 2 && fs[1] == "C" {
					values = append(values, Value{
						Name:   "chassis_temps_reading",
						Value:  fs[0],
						Labels: ts,
					})
				}

				minWarningThreshold := strings.Fields(fields["minimum_warning_threshold"])
				if len(minWarningThreshold) == 2 && minWarningThreshold[1] == "C" {
					values = append(values, Value{
						Name:   "chassis_temps_min_warning",
						Value:  minWarningThreshold[0],
						Labels: ts,
					})
				}
				maxWarningThreshold := strings.Fields(fields["maximum_warning_threshold"])
				if len(maxWarningThreshold) == 2 && maxWarningThreshold[1] == "C" {
					values = append(values, Value{
						Name:   "chassis_temps_max_warning",
						Value:  maxWarningThreshold[0],
						Labels: ts,
					})
				}
				minFailureThreshold := strings.Fields(fields["minimum_failure_threshold"])
				if len(minFailureThreshold) == 2 && minFailureThreshold[1] == "C" {
					values = append(values, Value{
						Name:   "chassis_temps_min_failure",
						Value:  minFailureThreshold[0],
						Labels: ts,
					})
				}
				maxFailureThreshold := strings.Fields(fields["maximum_failure_threshold"])
				if len(maxFailureThreshold) == 2 && maxFailureThreshold[1] == "C" {
					values = append(values, Value{
						Name:   "chassis_temps_max_failure",
						Value:  maxFailureThreshold[0],
						Labels: ts,
					})
				}
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "chassis", "temps")
	return values, err
}

// Volts returns the chassis volts statud and if support reading
func (or *OMReport) Volts() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if len(fields) != 8 {
					continue
				}

				if _, err := strconv.Atoi(fields["index"]); err != nil {
					continue
				}

				ts := map[string]string{"component": replace(fields["probe_name"])}
				values = append(values, Value{
					Name:   "chassis_volts_status",
					Value:  severity(fields["status"]),
					Labels: ts,
				})
				if i, err := extract(fields["reading"], "V"); err == nil {
					values = append(values, Value{
						Name:   "chassis_volts_reading",
						Value:  i,
						Labels: ts,
					})
				}
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "chassis", "volts")
	return values, err
}

// ChassisBatteries returns the chassis batteries status
func (or *OMReport) ChassisBatteries() ([]Value, error) {
	values := []Value{}
	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if len(fields) < 4 {
					continue
				}

				index := strings.Replace(fields["index"], ":", "_", -1)
				ts := map[string]string{"index": index}

				values = append(values, Value{
					Name:   "cmos_batteries_status",
					Value:  severity(fields["status"]),
					Labels: ts,
				})
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "chassis", "batteries")
	return values, err
}

// ChassisBios returns the bios version name
func (or *OMReport) ChassisBios() ([]Value, error) {
	value := Value{
		Name:   "bios",
		Value:  "0",
		Labels: map[string]string{},
	}

	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if len(fields) != 1 {
					continue
				}

				for k, v := range fields {
					id := normalizeName(k)
					value.Labels[id] = strings.ToLower(v)
					break
				}
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "chassis", "bios")

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

	err := or.readReport(func(outputs Output) {
		for _, output := range outputs {
			for _, fields := range output.Lines {
				if len(fields) != 1 {
					continue
				}

				for k, v := range fields {
					id := normalizeName(k)
					value.Labels[id] = strings.ToLower(v)
					break
				}
			}
		}
	}, DynamicReaderMode, or.getOMReportExecutable(), "chassis", "firmware")

	values := []Value{}
	values = append(values, value)

	return values, err
}
