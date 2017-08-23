package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/galexrt/dellhw_exporter/collector"
	"github.com/galexrt/dellhw_exporter/pkg/omreport"
	rcon "github.com/galexrt/go-rcon"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

const (
	defaultCollectors = "chassis,fans,memory,processors,ps,ps_amps_sysboard_pwr,storage_battery,storage_enclosure,storage_controller,storage_vdisk,system,temps,volts"
)

var (
	showhelp           bool
	showVersion        bool
	showCollectors     bool
	debugMode          bool
	enabledCollectors  string
	metricsAddr        string
	metricsPath        string
	omReportExecutable string
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

var (
	log = logrus.New()
)

// DellHWCollector
type DellHWCollector struct {
	lastCollectTime time.Time
	collectors      map[string]collector.Collector
}

func init() {
	flag.BoolVar(&showhelp, "help", false, "Show help menu")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showCollectors, "collectors.print", false, "If true, print available collectors and exit.")
	flag.BoolVar(&debugMode, "debug", false, "Enable debug output")
	flag.StringVar(&metricsAddr, "web.listen-address", ":9137", "The address to listen on for HTTP requests")
	flag.StringVar(&metricsPath, "web.telemetry-path", "/metrics", "Path the metrics will be exposed under")
	flag.StringVar(&enabledCollectors, "collectors.enabled", defaultCollectors, "Comma separated list of active collectors")
	flag.StringVar(&omReportExecutable, "collectors.omr-report", "/opt/dell/srvadmin/bin/omreport", "Path to the omReport executable")
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

func main() {
	flag.Parse()
	if showhelp {
		fmt.Println(os.Args[0] + " [FLAGS]")
		flag.PrintDefaults()
		return
	}
	if showVersion {
		fmt.Fprintln(os.Stdout, version.Print("srcds_exporter"))
		return
	}
	if showCollectors {
		collectorNames := make(sort.StringSlice, 0, len(collector.Factories))
		for n := range collector.Factories {
			collectorNames = append(collectorNames, n)
		}
		collectorNames.Sort()
		fmt.Printf("Available collectors:\n")
		for _, n := range collectorNames {
			fmt.Printf(" - %s\n", n)
		}
		return
	}
	log.Out = os.Stdout
	if debugMode {
		log.Level = logrus.DebugLevel
	}
	rcon.SetLog(log)
	log.Infoln("Starting srcds_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	omrOpts := &omreport.Options{
		OMReportExecutable: omReportExecutable,
	}

	collector.SetOMReport(omreport.New(omrOpts))

	collectors, err := loadCollectors(enabledCollectors)
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

	http.HandleFunc(metricsPath, func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>DellHW Exporter</title></head>
			<body>
			<h1>DellHW Exporter</h1>
			<p><a href="` + metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})
	err = http.ListenAndServe(metricsAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
