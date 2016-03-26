package main

import (
	"fmt"
	"net/http"
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

	exporterType      string
	listenAddress     string
	listenPort        string
	metricsPath       string
	enabledCollectors string

	cache = newMetricStorage()
)

func init() {
	RootCmd.Flags().StringVarP(&logLevel, "loglevel", "L", "info", "Set log level")
	RootCmd.Flags().StringVarP(&listenAddress, "web.listen", "l", "127.0.0.1", "Address on which to expose metrics and web interface.")
	RootCmd.Flags().StringVarP(&listenPort, "web.port", "p", "4242", "Port on which to expose metrics.")
	RootCmd.Flags().StringVarP(&metricsPath, "web.path", "m", "/metrics", "Path under which to expose metrics.")
	RootCmd.Flags().StringVarP(&enabledCollectors, "collect", "c", "chassis,fans,memory,processors,ps,ps_amps_sysboard_pwr,storage_battery,storage_enclosure,storage_controller,storage_vdisk,system,temps,volts", "Comma-separated list of collectors to use.")
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

	http.Handle(metricsPath, prometheus.Handler())
	log.Info("listening to ", listenAddress+":"+listenPort)
	log.Fatal(http.ListenAndServe(listenAddress+":"+listenPort, nil))

}

func main() {

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}
