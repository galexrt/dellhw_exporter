package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
)

var (
	// RootCmd is the main command
	RootCmd = &cobra.Command{
		Use:   "hardware_exporter",
		Short: "Prometheus exporter for Dell Hardware components",
		Run: func(cmd *cobra.Command, args []string) {
			runMainCommand()
		},
	}

	HWEVersion string
	BuildDate  string
	logLevel   string

	enabledCollectors string
	pushGateway       string

	cache = newMetricStorage()
)

func init() {
	RootCmd.Flags().StringVarP(&logLevel, "loglevel", "L", "info", "Set log level")
	RootCmd.Flags().StringVarP(&enabledCollectors, "collect", "c", "chassis,fans,memory,processors,ps,ps_amps_sysboard_pwr,storage_battery,storage_enclosure,storage_controller,storage_vdisk,system,temps,volts", "Comma-separated list of collectors to use.")
	RootCmd.Flags().StringVarP(&pushGateway, "gateway", "g", "", "Push Gateway in the form of address:port")
	RootCmd.AddCommand(versionCmd)

}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of hardware_exporter",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("hardware_exporter\n\n")
		fmt.Printf("version    : %s\n", HWEVersion)
		if BuildDate != "" {
			fmt.Printf("build date : %s\n", BuildDate)
		}
	},
}

func runMainCommand() {

	if logLevel == "info" {
		log.SetLevel(log.InfoLevel)
	}
	if logLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	}
	if logLevel == "error" {
		log.SetLevel(log.ErrorLevel)
	}

	err := collect(collectors)
	if err != nil {
		log.Debug("Collect failed")
		os.Exit(1)
	}

	if err := prometheus.Push("dell_hardware", "", pushGateway); err != nil {
		log.Error("Failed to push to the prometheus gateway : ", err)
	}

}

func main() {

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}
