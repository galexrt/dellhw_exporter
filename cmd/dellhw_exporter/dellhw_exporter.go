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

package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/galexrt/dellhw_exporter/collector"
	"github.com/galexrt/dellhw_exporter/pkg/omreport"
	"github.com/kardianos/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/sirupsen/logrus"
)

const (
	defaultCollectors = "chassis,chassis_batteries,fans,firmwares,memory,nics,processors,ps,ps_amps_sysboard_pwr,storage_battery,storage_controller,storage_enclosure,storage_pdisk,storage_vdisk,system,temps,volts"
)

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(collector.Namespace, "scrape", "collector_duration_seconds"),
		"dellhw_exporter: Duration of a collector scrape.",
		[]string{"collector"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(collector.Namespace, "scrape", "collector_success"),
		"dellhw_exporter: Whether a collector succeeded.",
		[]string{"collector"},
		nil,
	)
)

type program struct{}

// CmdLineOpts holds possible command line options/flags
type CmdLineOpts struct {
	version        bool
	showCollectors bool
	logLevel       string

	metricsAddr        string
	metricsPath        string
	enabledCollectors  string
	monitoredNics      []string
	omReportExecutable string
	cmdTimeout         int64

	cachingEnabled bool
	cacheDuration  int64
}

var (
	log   = logrus.New()
	opts  CmdLineOpts
	flags = flag.NewFlagSet("dellhw_exporter", flag.ExitOnError)
)

// DellHWCollector contains the collectors to be used
type DellHWCollector struct {
	lastCollectTime time.Time
	collectors      map[string]collector.Collector

	// Cache related
	cachingEnabled bool
	cacheDuration  time.Duration
	cache          []prometheus.Metric
	cacheMutex     sync.Mutex
}

func main() {
	// Service setup
	svcConfig := &service.Config{
		Name:        "DellOMSAExporter",
		DisplayName: "Dell OMSA Exporter",
		Description: "Prometheus exporter for Dell Hardware components using OMSA",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		log.Error(err)
	}
}

func (p *program) Start(s service.Service) error {
	if err := parseFlagsAndEnvVars(); err != nil {
		log.Fatal(err)
	}

	if opts.version {
		fmt.Fprintln(os.Stdout, version.Print("dellhw_exporter"))
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

	// Set log level
	l, err := logrus.ParseLevel(opts.logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(l)

	log.Infoln("Starting dellhw_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	if opts.cmdTimeout > 0 {
		log.Infof("Setting command timeout to %d", opts.cmdTimeout)
		omreport.SetCommandTimeout(opts.cmdTimeout)
	} else {
		log.Warnf("Not setting command timeout because it is zero")
	}

	if opts.cachingEnabled {
		log.Infof("Caching enabled. Cache Duration: %ds", opts.cacheDuration)
	} else {
		log.Info("Caching is disabled by default")
	}

	omrOpts := &omreport.Options{
		OMReportExecutable: opts.omReportExecutable,
	}

	collector.SetLogger(log)
	collector.SetOMReport(omreport.New(omrOpts))

	collectors, err := loadCollectors(opts.enabledCollectors, opts.monitoredNics)
	if err != nil {
		log.Fatalf("Couldn't load collectors: %v", err)
	}
	log.Infof("Enabled collectors:")
	for n := range collectors {
		log.Infof(" - %s", n)
	}

	if err = prometheus.Register(NewDellHWCollector(collectors, opts.cachingEnabled, opts.cacheDuration)); err != nil {
		log.Fatalf("Couldn't register collector: %v", err)
	}

	// non-blocking start
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	// non-blocking stop
	return nil
}

func NewDellHWCollector(collectors map[string]collector.Collector, cachingEnabled bool, cacheDurationSeconds int64) *DellHWCollector {
	return &DellHWCollector{
		cache:           make([]prometheus.Metric, 0),
		lastCollectTime: time.Unix(0, 0),
		collectors:      collectors,
		cachingEnabled:  cachingEnabled,
		cacheDuration:   time.Duration(cacheDurationSeconds) * time.Second,
	}
}

func init() {
	flags.BoolVar(&opts.version, "version", false, "Show version information")
	flags.StringVar(&opts.logLevel, "log-level", "INFO", "Set log level")

	flags.BoolVar(&opts.showCollectors, "collectors-print", false, "If true, print available collectors and exit.")
	flags.StringVar(&opts.enabledCollectors, "collectors-enabled", defaultCollectors, "Comma separated list of active collectors")
	flags.StringSliceVar(&opts.monitoredNics, "monitored-nics", []string{}, "Comma separated list of nics to monitor (default, empty list, is to monitor all)")
	flags.StringVar(&opts.omReportExecutable, "collectors-omreport", getDefaultOmReportPath(), "Path to the omreport executable (based on the OS (linux or windows) default paths are used if unset)")
	flags.Int64Var(&opts.cmdTimeout, "collectors-cmd-timeout", 15, "Command execution timeout for omreport")

	flags.StringVar(&opts.metricsAddr, "web-listen-address", ":9137", "The address to listen on for HTTP requests")
	flags.StringVar(&opts.metricsPath, "web-telemetry-path", "/metrics", "Path the metrics will be exposed under")

	flags.BoolVar(&opts.cachingEnabled, "cache-enabled", false, "Enable metrics caching to reduce load")
	flags.Int64Var(&opts.cacheDuration, "cache-duration", 20, "Cache duration in seconds")

	flags.SetNormalizeFunc(normalizeFlags)
	flags.SortFlags = true
}

// normalizeFlags "normalize" / alias flags that have been deprcated / removed
func normalizeFlags(f *flag.FlagSet, name string) flag.NormalizedName {
	switch name {
	case "collectors.print":
		name = "collectors-print"
	case "web.listen-address":
		name = "web-listen-address"
	case "web.telemetry-path":
		name = "web-telemetry-path"
	case "collectors.enabled":
		name = "collectors-enabled"
	case "collectors.omr-report":
		name = "collectors-omreport"
	case "collectors.cmd-timeout":
		name = "collectors-cmd-timeout"
	}
	return flag.NormalizedName(name)
}

func flagNameFromEnvName(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

func parseFlagsAndEnvVars() error {
	for _, v := range os.Environ() {
		vals := strings.SplitN(v, "=", 2)

		if !strings.HasPrefix(vals[0], "DELLHW_EXPORTER_") {
			continue
		}
		flagName := flagNameFromEnvName(strings.ReplaceAll(vals[0], "DELLHW_EXPORTER_", ""))

		fn := flags.Lookup(flagName)
		if fn == nil || fn.Changed {
			continue
		}

		if err := fn.Value.Set(vals[1]); err != nil {
			return err
		}
	}

	return flags.Parse(os.Args[1:])
}

// Describe implements the prometheus.Collector interface.
func (n *DellHWCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

// Collect implements the prometheus.Collector interface.
func (n *DellHWCollector) Collect(outgoingCh chan<- prometheus.Metric) {
	if n.cachingEnabled {
		n.cacheMutex.Lock()
		defer n.cacheMutex.Unlock()

		expiry := n.lastCollectTime.Add(n.cacheDuration)
		if time.Now().Before(expiry) {
			log.Debugf("Using cache. Now: %s, Expiry: %s, LastCollect: %s", time.Now().String(), expiry.String(), n.lastCollectTime.String())
			for _, cachedMetric := range n.cache {
				log.Debugf("Pushing cached metric %s to outgoingCh", cachedMetric.Desc().String())
				outgoingCh <- cachedMetric
			}
			return
		}
		// Clear cache, but keep slice
		n.cache = n.cache[:0]
	}

	metricsCh := make(chan prometheus.Metric)

	// Wait to ensure outgoingCh is not closed before the goroutine is finished
	wgOutgoing := sync.WaitGroup{}
	wgOutgoing.Add(1)
	go func() {
		for metric := range metricsCh {
			outgoingCh <- metric
			if n.cachingEnabled {
				log.Debugf("Appending metric %s to cache", metric.Desc().String())
				n.cache = append(n.cache, metric)
			}
		}
		log.Debug("Finished pushing metrics from metricsCh to outgoingCh")
		wgOutgoing.Done()
	}()

	wgCollection := sync.WaitGroup{}
	wgCollection.Add(len(n.collectors))
	for name, coll := range n.collectors {
		go func(name string, coll collector.Collector) {
			execute(name, coll, metricsCh)
			wgCollection.Done()
		}(name, coll)
	}

	log.Debug("Waiting for collectors")
	wgCollection.Wait()
	log.Debug("Finished waiting for collectors")

	n.lastCollectTime = time.Now()
	log.Debugf("Updated lastCollectTime to %s", n.lastCollectTime.String())

	close(metricsCh)

	log.Debug("Waiting for outgoing Adapter")
	wgOutgoing.Wait()
	log.Debug("Finished waiting for outgoing Adapter")
}

func execute(name string, c collector.Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)
	var success float64

	if err != nil {
		log.Errorf("%s collector failed after %fs: %v", name, duration.Seconds(), err)
		success = 0
	} else {
		log.Debugf("%s collector succeeded after %fs.", name, duration.Seconds())
		success = 1
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
}

func getCollectorConfig() *collector.Config {
	return &collector.Config{
		MonitoredNICs: opts.monitoredNics,
	}
}

func loadCollectors(list string, nicList []string) (map[string]collector.Collector, error) {
	cfg := getCollectorConfig()

	collectors := map[string]collector.Collector{}
	var c collector.Collector
	var err error
	for _, name := range strings.Split(list, ",") {
		fn, ok := collector.Factories[name]
		if !ok {
			return nil, fmt.Errorf("collector '%s' not available", name)
		}

		c, err = fn(cfg)
		if err != nil {
			return nil, err
		}
		collectors[name] = c
	}

	return collectors, nil
}

func getDefaultOmReportPath() string {
	if runtime.GOOS == "windows" {
		return "C:\\Program Files\\Dell\\SysMgt\\oma\\bin\\omreport.exe"
	}

	return "/opt/dell/srvadmin/bin/omreport"
}

func (p *program) run() {
	// Background work
	handler := promhttp.HandlerFor(prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			ErrorLog:      log,
			ErrorHandling: promhttp.ContinueOnError,
		})

	http.HandleFunc(opts.metricsPath, func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<!DOCTYPE html>
<html>
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
