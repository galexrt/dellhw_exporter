package omreport

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testResultOMReport struct {
	Input  string
	Values []Value
}

func getOMReport(input *string) *OMReport {
	return &OMReport{
		Reader: func(f func([]string), _ string, args ...string) error {
			for _, line := range strings.Split(*input, "\n") {
				sp := strings.Split(line, ";")
				for i, s := range sp {
					sp[i] = clean(s)
				}
				f(sp)
			}
			return nil
		},
	}
}

var chassisTests = []testResultOMReport{
	{
		Input: `Health

Main System Chassis

SEVERITY;COMPONENT
Ok;Fans
Ok;Intrusion

For further help, type the command followed by -?
`,
		Values: []Value{
			{
				Name:  "chassis_status",
				Value: "0",
				Labels: map[string]string{
					"component": "Fans",
				},
			},
			{
				Name:  "chassis_status",
				Value: "0",
				Labels: map[string]string{
					"component": "Intrusion",
				},
			},
		},
	},
}

func TestChassis(t *testing.T) {
	input := ""
	report := getOMReport(&input)
	for _, result := range chassisTests {
		input = result.Input
		values, _ := report.Chassis()
		assert.Equal(t, result.Values, values)
	}
}

var fansTests = []testResultOMReport{
	{
		Input: `Fan Probes Information

Fan Redundancy
Redundancy Status;Full

Probe List

Index;Status;Probe Name;Reading;Minimum Warning Threshold;Maximum Warning Threshold;Minimum Failure Threshold;Maximum Failure Threshold
0;Ok;System Board Fan1A;5040 RPM;840 RPM;[N/A];600 RPM;[N/A]
1;Ok;System Board Fan2A;5160 RPM;840 RPM;[N/A];600 RPM;[N/A]
`,
		Values: []Value{
			{
				Name:  "chassis_fan_status",
				Value: "0",
				Labels: map[string]string{
					"fan": "System_Board_Fan1A",
				},
			},
			{
				Name:  "chassis_fan_reading",
				Value: "5040",
				Labels: map[string]string{
					"fan": "System_Board_Fan1A",
				},
			},
			{
				Name:  "chassis_fan_status",
				Value: "0",
				Labels: map[string]string{
					"fan": "System_Board_Fan2A",
				},
			},
			{
				Name:  "chassis_fan_reading",
				Value: "5160",
				Labels: map[string]string{
					"fan": "System_Board_Fan2A",
				},
			},
		},
	},
}

func TestFans(t *testing.T) {
	input := ""
	report := getOMReport(&input)
	for _, result := range fansTests {
		input = result.Input
		values, _ := report.Fans()
		assert.Equal(t, result.Values, values)
	}
}

var memoryTests = []testResultOMReport{
	{
		Input: `Memory Information

Health;Ok

Attributes of Memory Array(s)

Attributes of Memory Array(s)
Location;System Board or Motherboard
Use;System Memory
Installed Capacity;131072  MB
Maximum Capacity;3145728  MB
Slots Available;24
Slots Used;8
Error Correction;Multibit ECC

Total of Memory Array(s)
Total Installed Capacity;131072  MB
Total Installed Capacity Available to the OS;128853  MB
Total Maximum Capacity;3145728  MB

Details of Memory Array 1

Index;Status;Connector Name;Type;Size
0;Ok;A1;DDR4 - Synchronous Registered (Buffered);16384  MB
;Unknown;A9;[Not Occupied];
`,
		Values: []Value{
			{
				Name:  "chassis_memory_status",
				Value: "0",
				Labels: map[string]string{
					"memory": "A1",
				},
			},
		},
	},
}

func TestMemory(t *testing.T) {
	input := ""
	report := getOMReport(&input)
	for _, result := range memoryTests {
		input = result.Input
		values, _ := report.Memory()
		assert.Equal(t, result.Values, values)
	}
}

var systemTests = []testResultOMReport{
	{
		Input: `Health

SEVERITY;COMPONENT
Ok;Main System Chassis

For further help, type the command followed by -?
`,
		Values: []Value{
			{
				Name:  "system_status",
				Value: "0",
				Labels: map[string]string{
					"component": "Main_System_Chassis",
				},
			},
		},
	},
}

func TestSystem(t *testing.T) {
	input := ""
	report := getOMReport(&input)
	for _, result := range systemTests {
		input = result.Input
		values, _ := report.System()
		assert.Equal(t, result.Values, values)
	}
}

var storageBatteryTests = []testResultOMReport{
	{
		Input: `List of Batteries in the System

Controller PERC H730 Mini (Slot Embedded)

ID;Status;Name;State;Recharge Count;Max Recharge Count;Learn State;Next Learn Time;Maximum Learn Delay
0;Ok;Battery ;Ready;Not Applicable;Not Applicable;Not Applicable;Not Applicable;Not Applicable
`,
		Values: []Value{
			{
				Name:  "storage_battery_status",
				Value: "0",
				Labels: map[string]string{
					"controller":        "0",
					controllerNameLabel: "PERC H730 Mini (Slot Embedded)",
				},
			},
		},
	},
	{
		Input: `List of Batteries in the System

Controller PERC H840 Adapter  (Slot 1)

ID;Status;Name;State;Recharge Count;Max Recharge Count;Learn State;Next Learn Time;Maximum Learn Delay
0;Ok;Battery ;Ready;Not Applicable;Not Applicable;Not Applicable;Not Applicable;Not Applicable

Controller PERC H740P Mini  (Slot Embedded)

ID;Status;Name;State;Recharge Count;Max Recharge Count;Learn State;Next Learn Time;Maximum Learn Delay
0;Ok;Battery ;Ready;Not Applicable;Not Applicable;Not Applicable;Not Applicable;Not Applicable
`,
		Values: []Value{
			{
				Name:  "storage_battery_status",
				Value: "0",
				Labels: map[string]string{
					"controller":        "0",
					controllerNameLabel: "PERC H840 Adapter (Slot 1)",
				},
			},
			{
				Name:  "storage_battery_status",
				Value: "0",
				Labels: map[string]string{
					"controller":        "0",
					controllerNameLabel: "PERC H740P Mini (Slot Embedded)",
				},
			},
		},
	},
}

func TestStorageBattery(t *testing.T) {
	input := ""
	report := getOMReport(&input)
	for _, result := range storageBatteryTests {
		input = result.Input
		values, _ := report.StorageBattery()
		assert.Equal(t, result.Values, values)
	}
}

var storageControllerTests = []testResultOMReport{
	{
		Input: `Controller  PERC H730 Mini (Slot Embedded)

Controller

ID;Status;Name;Slot ID;State;Firmware Version;Minimum Required Firmware Version;Driver Version;Minimum Required Driver Version;Storport Driver Version;Minimum Required Storport Driver Version;Number of Connectors;Rebuild Rate;BGI Rate;Check Consistency Rate;Reconstruct Rate;Alarm State;Cluster Mode;SCSI Initiator ID;Cache Memory Size;Patrol Read Mode;Patrol Read State;Patrol Read Rate;Patrol Read Iterations;Abort Check Consistency on Error;Allow Revertible Hot Spare and Replace Member;Load Balance;Auto Replace Member on Predictive Failure;Redundant Path view;CacheCade Capable;Persistent Hot Spare;Encryption Capable;Encryption Key Present;Encryption Mode;Preserved Cache;Spin Down Unconfigured Drives;Spin Down Hot Spares;Spin Down Configured Drives;Automatic Disk Power Saving (Idle C);Time Interval for Spin Down (in Minutes);Start Time (HH:MM);Time Interval for Spin Up (in Hours);T10 Protection Information Capable;Non-RAID HDD Disk Cache Policy;Current Controller Mode
0;Ok;PERC H730 Mini;Embedded;Ready;25.5.0.0018;Not Applicable;06.811.02.00-rc1;Not Applicable;Not Applicable;Not Applicable;1;30%;30%;30%;30%;Not Applicable;Not Applicable;Not Applicable;1024 MB;Auto;Stopped;30%;0;Disabled;Disabled;Not Applicable;Disabled;Not Applicable;Not Applicable;Disabled;Yes;No;None;Not Applicable;Enabled;Disabled;Disabled;Disabled;30;Not Applicable;Not Applicable;Yes;Unchanged;RAID
`,
		Values: []Value{
			{
				Name:  "storage_controller_status",
				Value: "0",
				Labels: map[string]string{
					"id":                "0",
					controllerNameLabel: "PERC H730 Mini (Slot Embedded)",
				},
			},
		},
	},
	{
		Input: `List of Controllers in the system

		Controller
		
		ID;Status;Name;Slot ID;State;Firmware Version;Minimum Required Firmware Version;Driver Version;Minimum Required Driver Version;Storport Driver Version;Minimum Required Storport Driver Version;Number of Connectors;Rebuild Rate;BGI Rate;Check Consistency Rate;Reconstruct Rate;Alarm State;Cluster Mode;SCSI Initiator ID;Cache Memory Size;Patrol Read Mode;Patrol Read State;Patrol Read Rate;Patrol Read Iterations;Abort Check Consistency on Error;Allow Revertible Hot Spare and Replace Member;Load Balance;Auto Replace Member on Predictive Failure;Redundant Path view;CacheCade Capable;Persistent Hot Spare;Encryption Capable;Encryption Key Present;Encryption Mode;Preserved Cache;Spin Down Unconfigured Drives;Spin Down Hot Spares;Spin Down Configured Drives;Automatic Disk Power Saving (Idle C);Start Time (HH:MM);Time Interval for Spin Up (in Hours);T10 Protection Information Capable;Non-RAID HDD Disk Cache Policy;Current Controller Mode
		0;Ok;PERC H730P Mini;Embedded;Ready;25.5.5.0005;Not Applicable;06.810.09.00-rc1;Not Applicable;Not Applicable;Not Applicable;1;30%;30%;30%;30%;Not Applicable;Not Applicable;Not Applicable;2048 MB;Auto;Stopped;30%;159;Disabled;Enabled;Not Applicable;Disabled;Not Applicable;Not Applicable;Disabled;Yes;No;None;Not Applicable;Disabled;Disabled;Disabled;Disabled;Not Applicable;Not Applicable;No;Unchanged;RAID
		
		Controller
		
		ID;Status;Name;Slot ID;State;Firmware Version;Minimum Required Firmware Version;Driver Version;Minimum Required Driver Version;Storport Driver Version;Minimum Required Storport Driver Version;Number of Connectors;Rebuild Rate;BGI Rate;Check Consistency Rate;Reconstruct Rate;Alarm State;Cluster Mode;SCSI Initiator ID;Cache Memory Size;Patrol Read Mode;Patrol Read State;Patrol Read Rate;Patrol Read Iterations;Abort Check Consistency on Error;Allow Revertible Hot Spare and Replace Member;Load Balance;Auto Replace Member on Predictive Failure;Redundant Path view;CacheCade Capable;Persistent Hot Spare;Encryption Capable;Encryption Key Present;Encryption Mode;Preserved Cache;Spin Down Unconfigured Drives;Spin Down Hot Spares;Spin Down Configured Drives;Automatic Disk Power Saving (Idle C);Start Time (HH:MM);Time Interval for Spin Up (in Hours);T10 Protection Information Capable;Non-RAID HDD Disk Cache Policy
		1;Ok;PERC H810 Adapter;PCI Slot 6;Ready;21.3.5-0002;Not Applicable;06.810.09.00-rc1;Not Applicable;Not Applicable;Not Applicable;2;30%;30%;30%;30%;Not Applicable;Not Applicable;Not Applicable;1024 MB;Auto;Stopped;30%;122;Disabled;Enabled;Auto;Disabled;Detected;Yes;Disabled;Yes;No;None;Not Applicable;Disabled;Disabled;Disabled;Disabled;Not Applicable;Not Applicable;No;Not Applicable
		
		Controller
		
		ID;Status;Name;Slot ID;State;Firmware Version;Minimum Required Firmware Version;Driver Version;Minimum Required Driver Version;Storport Driver Version;Minimum Required Storport Driver Version;Number of Connectors;Rebuild Rate;BGI Rate;Check Consistency Rate;Reconstruct Rate;Alarm State;Cluster Mode;SCSI Initiator ID;Cache Memory Size;Patrol Read Mode;Patrol Read State;Patrol Read Rate;Patrol Read Iterations;Abort Check Consistency on Error;Allow Revertible Hot Spare and Replace Member;Load Balance;Auto Replace Member on Predictive Failure;Redundant Path view;CacheCade Capable;Persistent Hot Spare;Encryption Capable;Encryption Key Present;Encryption Mode;Preserved Cache;Spin Down Unconfigured Drives;Spin Down Hot Spares;Spin Down Configured Drives;Automatic Disk Power Saving (Idle C);Start Time (HH:MM);Time Interval for Spin Up (in Hours);T10 Protection Information Capable;Non-RAID HDD Disk Cache Policy
		2;Ok;PERC H810 Adapter;PCI Slot 4;Ready;21.3.5-0002;Not Applicable;06.810.09.00-rc1;Not Applicable;Not Applicable;Not Applicable;2;30%;30%;30%;30%;Not Applicable;Not Applicable;Not Applicable;1024 MB;Auto;Stopped;30%;96;Disabled;Enabled;Auto;Disabled;Detected;Yes;Disabled;Yes;No;None;Not Applicable;Disabled;Disabled;Disabled;Disabled;Not Applicable;Not Applicable;No;Not Applicable
`,
		Values: []Value{
			{
				Name:  "storage_controller_status",
				Value: "0",
				Labels: map[string]string{
					"id":                "0",
					controllerNameLabel: "PERC H730P Mini (Slot Embedded)",
				},
			},
			{
				Name:  "storage_controller_status",
				Value: "0",
				Labels: map[string]string{
					"id":                "1",
					controllerNameLabel: "PERC H810 Adapter (Slot PCI Slot 6)",
				},
			},
			{
				Name:  "storage_controller_status",
				Value: "0",
				Labels: map[string]string{
					"id":                "2",
					controllerNameLabel: "PERC H810 Adapter (Slot PCI Slot 4)",
				},
			},
		},
	},
}

func TestStorageController(t *testing.T) {
	input := ""
	report := getOMReport(&input)
	for _, result := range storageControllerTests {
		input = result.Input
		values, _ := report.StorageController()
		assert.Equal(t, result.Values, values)
	}
}

var storageEnclosureTests = []testResultOMReport{
	{
		Input: `List of Enclosures in the System

Enclosure(s) on Controller PERC H730 Mini (Slot Embedded)


ID;Status;Name;State;Connector;Target ID;Configuration;Firmware Version;Downstream Firmware Version;Service Tag;Express Service Code;Asset Tag;Asset Name;Backplane Part Number;Split Bus Part Number;Enclosure Part Number;SAS Address;Enclosure Alarm
0:1;Ok;Backplane;Ready;0;Not Applicable;Not Applicable;3.31;Not Applicable;Not Applicable;Not Applicable;Not Applicable;Not Applicable;Not Applicable;Not Applicable;Not Applicable;500056B3B43B8CFD;Not Applicable
`,
		Values: []Value{
			{
				Name:  "storage_enclosure_status",
				Value: "0",
				Labels: map[string]string{
					"enclosure":         "0_1",
					controllerNameLabel: "PERC H730 Mini (Slot Embedded)",
				},
			},
		},
	},
}

func TestStorageEnclosure(t *testing.T) {
	input := ""
	report := getOMReport(&input)
	for _, result := range storageEnclosureTests {
		input = result.Input
		values, _ := report.StorageEnclosure()
		assert.Equal(t, result.Values, values)
	}
}

var storagePdiskTests = []testResultOMReport{
	{
		Input: `List of Physical Disks on Controller PERC H730 Mini (Slot Embedded)

Controller PERC H730 Mini (Slot Embedded)

ID;Status;Name;State;Power Status;Bus Protocol;Media;Part of Cache Pool;Remaining Rated Write Endurance;Failure Predicted;Revision;Driver Version;Model Number;T10 PI Capable;Certified;Encryption Capable;Encrypted;Progress;Mirror Set ID;Capacity;Used RAID Disk Space;Available RAID Disk Space;Hot Spare;Vendor ID;Product ID;Serial No.;Part Number;Negotiated Speed;Capable Speed;PCIe Negotiated Link Width;PCIe Maximum Link Width;Sector Size;Device Write Cache;Manufacture Day;Manufacture Week;Manufacture Year;SAS Address;Non-RAID HDD Disk Cache Policy;Disk Cache Policy;Form Factor ;Sub Vendor;ISE Capable
0:1:0;Ok;Physical Disk 0:1:0;Ready;Not Applicable;SATA;SSD;Not Applicable;100%;No;G201DL2B;Not Applicable;Not Applicable;No;Yes;No;Not Applicable;Not Applicable;Not Applicable;185.75 GB (199447543808 bytes);185.75 GB (199447543808 bytes);0.00 GB (0 bytes);Dedicated;DELL(tm);INTEL SSDSC2BX200G4R;BTHC643503A2200TGN;CN03481GIT2006AT00P3A0;6.00 Gbps;6.00 Gbps;Not Applicable;Not Applicable;512B;Not Applicable;Not Available;Not Available;Not Available;500056B3B43B8CC0;Not Applicable;Not Applicable;Not Available;Not Available;No
0:1:1;Ok;Physical Disk 0:1:1;Online;Not Applicable;SATA;SSD;Not Applicable;100%;No;G201DL2B;Not Applicable;Not Applicable;No;Yes;No;Not Applicable;Not Applicable;Not Applicable;185.75 GB (199447543808 bytes);185.75 GB (199447543808 bytes);0.00 GB (0 bytes);No;DELL(tm);INTEL SSDSC2BX200G4R;BTHC643503BX200TGN;CN03481GIT2006AT00PGA0;6.00 Gbps;6.00 Gbps;Not Applicable;Not Applicable;512B;Not Applicable;Not Available;Not Available;Not Available;500056B3B43B8CC1;Not Applicable;Not Applicable;Not Available;Not Available;No
`,
		Values: []Value{
			{
				Name:  "storage_pdisk_status",
				Value: "0",
				Labels: map[string]string{
					"controller":        "0",
					"disk":              "0_1_0",
					controllerNameLabel: "PERC H730 Mini (Slot Embedded)",
				},
			},
			{
				Name:  "storage_pdisk_status",
				Value: "0",
				Labels: map[string]string{
					"controller":        "0",
					"disk":              "0_1_1",
					controllerNameLabel: "PERC H730 Mini (Slot Embedded)",
				},
			},
		},
	},
}

func TestStoragePdisk(t *testing.T) {
	input := ""
	report := getOMReport(&input)
	for _, result := range storagePdiskTests {
		input = result.Input
		values, _ := report.StoragePdisk("0")
		assert.Equal(t, result.Values, values)
	}
}

var storageVdiskTests = []testResultOMReport{
	{
		Input: `List of Virtual Disks in the System

Controller PERC H730 Mini (Slot Embedded)

ID;Status;Name;State;Hot Spare Policy violated;Encrypted;Layout;Size;T10 Protection Information Status;Associated Fluid Cache State ;Device Name;Bus Protocol;Media;Read Policy;Write Policy;Cache Policy;Stripe Element Size;Disk Cache Policy
0;Ok;GenericR5_0;Ready;Not Assigned;No;RAID-5;743.00 GB (797790175232 bytes);No;Not Applicable;/dev/sda;SATA;SSD;No Read Ahead;Write Through;Not Applicable;64 KB;Unchanged
`,
		Values: []Value{
			{
				Name:  "storage_vdisk_status",
				Value: "0",
				Labels: map[string]string{
					"vdisk":             "0",
					controllerNameLabel: "PERC H730 Mini (Slot Embedded)",
				},
			},
		},
	},
}

func TestStorageVdisk(t *testing.T) {
	input := ""
	report := getOMReport(&input)
	for _, result := range storageVdiskTests {
		input = result.Input
		values, _ := report.StorageVdisk()
		assert.Equal(t, result.Values, values)
	}
}

var psTests = []testResultOMReport{
	{
		Input: `Power Supplies Information

Power Supply Redundancy
Redundancy Status;Full

Individual Power Supply Elements

Index;Status;Location;Type;Rated Input Wattage;Maximum Output Wattage;Firmware Version;Online Status;Power Monitoring Capable
0;Ok;PS1 Status;AC;900 W;750 W;00.14.4B;Presence Detected;Yes
1;Ok;PS2 Status;AC;900 W;750 W;00.14.4B;Presence Detected;Yes
`,
		Values: []Value{
			{
				Name:  "ps_status",
				Value: "0",
				Labels: map[string]string{
					"id": "0",
				},
			},
			{
				Name:  "ps_rated_input_wattage",
				Value: "900",
				Labels: map[string]string{
					"id": "0",
				},
			},
			{
				Name:  "ps_rated_output_wattage",
				Value: "750",
				Labels: map[string]string{
					"id": "0",
				},
			},
			{
				Name:  "ps_status",
				Value: "0",
				Labels: map[string]string{
					"id": "1",
				},
			},
			{
				Name:  "ps_rated_input_wattage",
				Value: "900",
				Labels: map[string]string{
					"id": "1",
				},
			},
			{
				Name:  "ps_rated_output_wattage",
				Value: "750",
				Labels: map[string]string{
					"id": "1",
				},
			},
		},
	},
}

func TestPs(t *testing.T) {
	input := ""
	report := getOMReport(&input)
	for _, result := range psTests {
		input = result.Input
		values, _ := report.Ps()
		assert.Equal(t, result.Values, values)
	}
}

var psAmpsSysboardPwrTests = []testResultOMReport{
	{
		Input: `Power Consumption Information

Power Consumption

Index;Status;Probe Name;Reading;Warning Threshold;Failure Threshold
2;Ok;System Board Pwr Consumption;84 W;896 W;980 W

Amperage
PS1 Current 1;0.2 A
PS2 Current 2;0.2 A

Power Headroom
System Instantaneous Headroom;811 W
System Peak Headroom;0 W

Power Tracking Statistics

Statistic;Measurement Start Time;Measurement Finish Time;Reading
Energy Consumption;Wed Dec 14 21:57:40 2016;Wed Aug 23 18:46:28 2017;584.5 kWh

Statistic;Measurement Start Time;Peak Time;Peak Reading
System Peak Power;Wed Dec 14 21:57:41 2016;Wed Dec 28 08:41:13 2016;1023 W
System Peak Amperage;Wed Dec 14 21:57:41 2016;Wed Dec 28 08:41:13 2016;1.3 A
`,
		Values: []Value{
			//{Name:chassis_power_reading Value:84 Labels:map[]}
			{
				Name:   "chassis_power_reading",
				Value:  "84",
				Labels: nil,
			},
			{
				Name:   "chassis_power_warn_level",
				Value:  "896",
				Labels: nil,
			},
			{
				Name:   "chassis_power_fail_level",
				Value:  "980",
				Labels: nil,
			},
			{
				Name:  "chassis_current_reading",
				Value: "0.2",
				Labels: map[string]string{
					"pwrsupply": "PS1",
				},
			},
			{
				Name:  "chassis_current_reading",
				Value: "0.2",
				Labels: map[string]string{
					"pwrsupply": "PS2",
				},
			},
		},
	},
}

func TestPsAmpsSysboardPwr(t *testing.T) {
	input := ""
	report := getOMReport(&input)
	for _, result := range psAmpsSysboardPwrTests {
		input = result.Input
		values, _ := report.PsAmpsSysboardPwr()
		assert.Equal(t, result.Values, values)
	}
}

var processorsTests = []testResultOMReport{
	{
		Input: `Processors Information

Health;Ok

Index;Status;Connector Name;Processor Brand;Processor Version;Current Speed;State;Core Count
0;Ok;CPU1;Intel(R) Xeon(R) CPU E5-2630 v3 @ 2.40GHz;Model 63 Stepping 2;2400  MHz;Present;8
1;Unknown;CPU2;[Not Occupied];NA;NA;NA;NA;
`,
		Values: []Value{
			{
				Name:  "chassis_processor_status",
				Value: "0",
				Labels: map[string]string{
					"processor": "CPU1",
				},
			},
		},
	},
}

func TestProcessors(t *testing.T) {
	input := ""
	report := getOMReport(&input)
	for _, result := range processorsTests {
		input = result.Input
		values, _ := report.Processors()
		assert.Equal(t, result.Values, values)
	}
}

var tempsTests = []testResultOMReport{
	{
		Input: `Temperature Probes Information

Main System Chassis Temperatures : Ok

Index;Status;Probe Name;Reading;Minimum Warning Threshold;Maximum Warning Threshold;Minimum Failure Threshold;Maximum Failure Threshold
0;Ok;System Board Inlet Temp;17.0 C;3.0 C;42.0 C;-7.0 C;47.0 C
2;Ok;CPU1 Temp;34.0 C;8.0 C;82.0 C;3.0 C;87.0 C
`,
		Values: []Value{
			{
				Name:  "chassis_temps",
				Value: "0",
				Labels: map[string]string{
					"component": "System_Board_Inlet_Temp",
				},
			},
			{
				Name:  "chassis_temps_reading",
				Value: "17.0",
				Labels: map[string]string{
					"component": "System_Board_Inlet_Temp",
				},
			},
			{
				Name:  "chassis_temps_min_warning",
				Value: "3.0",
				Labels: map[string]string{
					"component": "System_Board_Inlet_Temp",
				},
			},
			{
				Name:  "chassis_temps_max_warning",
				Value: "42.0",
				Labels: map[string]string{
					"component": "System_Board_Inlet_Temp",
				},
			},
			{
				Name:  "chassis_temps_min_failure",
				Value: "-7.0",
				Labels: map[string]string{
					"component": "System_Board_Inlet_Temp",
				},
			},
			{
				Name:  "chassis_temps_max_failure",
				Value: "47.0",
				Labels: map[string]string{
					"component": "System_Board_Inlet_Temp",
				},
			},

			{
				Name:  "chassis_temps",
				Value: "0",
				Labels: map[string]string{
					"component": "CPU1_Temp",
				},
			},
			{
				Name:  "chassis_temps_reading",
				Value: "34.0",
				Labels: map[string]string{
					"component": "CPU1_Temp",
				},
			},
			{
				Name:  "chassis_temps_min_warning",
				Value: "8.0",
				Labels: map[string]string{
					"component": "CPU1_Temp",
				},
			},
			{
				Name:  "chassis_temps_max_warning",
				Value: "82.0",
				Labels: map[string]string{
					"component": "CPU1_Temp",
				},
			},
			{
				Name:  "chassis_temps_min_failure",
				Value: "3.0",
				Labels: map[string]string{
					"component": "CPU1_Temp",
				},
			},
			{
				Name:  "chassis_temps_max_failure",
				Value: "87.0",
				Labels: map[string]string{
					"component": "CPU1_Temp",
				},
			},
		},
	},
}

func TestTemps(t *testing.T) {
	input := ""
	report := getOMReport(&input)
	for _, result := range tempsTests {
		input = result.Input
		values, _ := report.Temps()
		assert.Equal(t, result.Values, values)
	}
}

var voltsTests = []testResultOMReport{
	{
		Input: `Voltage Probes Information

Health : Ok


Index;Status;Probe Name;Reading;Minimum Warning Threshold;Maximum Warning Threshold;Minimum Failure Threshold;Maximum Failure Threshold
0;Ok;CPU1 VCORE PG;Good;[N/A];[N/A];[N/A];[N/A]
1;Ok;System Board 3.3V PG;Good;[N/A];[N/A];[N/A];[N/A]
`,
		Values: []Value{
			{
				Name:  "chassis_volts_status",
				Value: "0",
				Labels: map[string]string{
					"component": "CPU1_VCORE_PG",
				},
			},
			{
				Name:  "chassis_volts_status",
				Value: "0",
				Labels: map[string]string{
					"component": "System_Board_3.3V_PG",
				},
			},
		},
	},
	{
		Input: `Voltage Probes Information

		Health : Ok
		
		
		Index;Status;Probe Name;Reading;Minimum Warning Threshold;Maximum Warning Threshold;Minimum Failure Threshold;Maximum Failure Threshold
		0;Ok;CPU1 VCORE PG;Good;[N/A];[N/A];[N/A];[N/A]
		2;Ok;System Board 3.3V PG;Good;[N/A];[N/A];[N/A];[N/A]
		31;Ok;System Board 2.5V AUX PG;Good;[N/A];[N/A];[N/A];[N/A]`,
		Values: []Value{
			{
				Name:  "chassis_volts_status",
				Value: "0",
				Labels: map[string]string{
					"component": "CPU1_VCORE_PG",
				},
			},
			{
				Name:  "chassis_volts_status",
				Value: "0",
				Labels: map[string]string{
					"component": "System_Board_3.3V_PG",
				},
			},
			{
				Name:  "chassis_volts_status",
				Value: "0",
				Labels: map[string]string{
					"component": "System_Board_2.5V_AUX_PG",
				},
			},
		},
	},
}

func TestVolts(t *testing.T) {
	input := ""
	report := getOMReport(&input)
	for _, result := range voltsTests {
		input = result.Input
		values, _ := report.Volts()
		assert.Equal(t, result.Values, values)
	}
}
