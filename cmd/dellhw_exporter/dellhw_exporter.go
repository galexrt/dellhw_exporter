package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"flag"

	"github.com/Sirupsen/logrus"
	"github.com/galexrt/dellhw_exporter/collector"
	"github.com/galexrt/dellhw_exporter/pkg/omreport"
	"github.com/galexrt/pkg/flagutil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

const (
	defaultCollectors = "chassis,fans,memory,processors,ps,ps_amps_sysboard_pwr,storage_battery,storage_enclosure,storage_controller,storage_pdisk,storage_vdisk,system,temps,volts,nics,chassis_batteries"
)

var (
	connectionDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(collector.Namespace, "scrape", "connection_duration_seconds"),
		"srcds_exporter: Duration of the server connection.",
		[]string{"connection"},
		nil,
	)
	connectionSucessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(collector.Namespace, "scrape", "connection_success"),
		"srcds_exporter: Whether the server connection succeeded.",
		[]string{"connection"},
		nil,
	)
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(collector.Namespace, "scrape", "collector_duration_seconds"),
		"srcds_exporter: Duration of a collector scrape.",
		[]string{"collector"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(collector.Namespace, "scrape", "collector_success"),
		"srcds_exporter: Whether a collector succeeded.",
		[]string{"collector"},
		nil,
	)
)

// CmdLineOpts holds possible command line options/flags
type CmdLineOpts struct {
	version            bool
	help               bool
	container          bool
	showCollectors     bool
	debugMode          bool
	metricsAddr        string
	metricsPath        string
	enabledCollectors  string
	omReportExecutable string
}

var (
	log                 = logrus.New()
	opts                CmdLineOpts
	dellhwExporterFlags = flag.NewFlagSet("dellhw_exporter", flag.ExitOnError)
)

// DellHWCollector contains the collectors to be used
type DellHWCollector struct {
	lastCollectTime time.Time
	collectors      map[string]collector.Collector
}

func init() {
	dellhwExporterFlags.BoolVar(&opts.help, "help", false, "Show help menu")
	dellhwExporterFlags.BoolVar(&opts.version, "version", false, "Show version information")
	dellhwExporterFlags.BoolVar(&opts.container, "container", false, "Start the Dell OpenManage service")
	dellhwExporterFlags.BoolVar(&opts.showCollectors, "collectors.print", false, "If true, print available collectors and exit.")
	dellhwExporterFlags.BoolVar(&opts.debugMode, "debug", false, "Enable debug output")
	dellhwExporterFlags.StringVar(&opts.metricsAddr, "web.listen-address", ":9137", "The address to listen on for HTTP requests")
	dellhwExporterFlags.StringVar(&opts.metricsPath, "web.telemetry-path", "/metrics", "Path the metrics will be exposed under")
	dellhwExporterFlags.StringVar(&opts.enabledCollectors, "collectors.enabled", defaultCollectors, "Comma separated list of active collectors")
	dellhwExporterFlags.StringVar(&opts.omReportExecutable, "collectors.omr-report", "/opt/dell/srvadmin/bin/omreport", "Path to the omReport executable")

	// Define the usage function
	dellhwExporterFlags.Usage = usage

	dellhwExporterFlags.Parse(os.Args[1:])
}

// Describe implements the prometheus.Collector interface.
func (n DellHWCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

// Collect implements the prometheus.Collector interface.
func (n DellHWCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(n.collectors))
	for name, c := range n.collectors {
		go func(name string, c collector.Collector) {
			execute(name, c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
}

func filterAvailableCollectors(collectors string) string {
	var availableCollectors []string
	for _, c := range strings.Split(collectors, ",") {
		_, ok := collector.Factories[c]
		if ok {
			availableCollectors = append(availableCollectors, c)
		}
	}
	return strings.Join(availableCollectors, ",")
}

func execute(name string, c collector.Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)
	var success float64

	if err != nil {
		log.Errorf("ERROR: %s collector failed after %fs: %s", name, duration.Seconds(), err)
		success = 0
	} else {
		log.Debugf("OK: %s collector succeeded after %fs.", name, duration.Seconds())
		success = 1
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
}

func loadCollectors(list string) (map[string]collector.Collector, error) {
	collectors := map[string]collector.Collector{}
	for _, name := range strings.Split(list, ",") {
		fn, ok := collector.Factories[name]
		if !ok {
			return nil, fmt.Errorf("collector '%s' not available", name)
		}
		c, err := fn()
		if err != nil {
			return nil, err
		}
		collectors[name] = c
	}
	return collectors, nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]...\n", os.Args[0])
	dellhwExporterFlags.PrintDefaults()
	os.Exit(0)
}

func main() {
	flagutil.SetFlagsFromEnv(dellhwExporterFlags, "DELLHW_EXPORTER")
	if opts.version {
		fmt.Fprintln(os.Stdout, version.Print("srcds_exporter"))
		os.Exit(0)
	}
	if opts.showCollectors {
		collectorNames := make(sort.StringSlice, 0, len(collector.Factories))
		for n := range collector.Factories {
			collectorNames = append(collectorNames, n)
		}
		collectorNames.Sort()
		fmt.Printf("Available collectors:\n")
		for _, n := range collectorNames {
			fmt.Printf(" - %s\n", n)
		}
		os.Exit(0)
	}
	log.Out = os.Stdout
	if opts.debugMode {
		log.Level = logrus.DebugLevel
	}
	log.Infoln("Starting srcds_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	if opts.container {
		log.Infoln("Starting srvadmin-services ...")
		cmd := exec.Command("/opt/dell/srvadmin/sbin/srvadmin-services.sh", "start")
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		timer := time.AfterFunc(30*time.Second, func() {
			cmd.Process.Kill()
		})
		err := cmd.Wait()
		if err != nil {
			log.Fatal(err)
		}
		timer.Stop()
		log.Infoln("Started srvadmin-services.")
	}

	omrOpts := &omreport.Options{
		OMReportExecutable: opts.omReportExecutable,
	}

	collector.SetOMReport(omreport.New(omrOpts))

	collectors, err := loadCollectors(opts.enabledCollectors)
	if err != nil {
		log.Fatalf("Couldn't load collectors: %s", err)
	}
	log.Infof("Enabled collectors:")
	for n := range collectors {
		log.Infof(" - %s", n)
	}

	if err = prometheus.Register(DellHWCollector{lastCollectTime: time.Now(), collectors: collectors}); err != nil {
		log.Fatalf("Couldn't register collector: %s", err)
	}
	handler := promhttp.HandlerFor(prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			ErrorLog:      log,
			ErrorHandling: promhttp.ContinueOnError,
		})

	http.HandleFunc(opts.metricsPath, func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>DellHW Exporter</title></head>
			<body>
			<h1>DellHW Exporter</h1>
			<p><a href="` + opts.metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	if err := http.ListenAndServe(opts.metricsAddr, nil); err != nil {
		log.Fatal(err)
	}
}
