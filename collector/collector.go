package collector

import (
	"github.com/galexrt/dellhw_exporter/pkg/omreport"
	"github.com/prometheus/client_golang/prometheus"
)

// Namespace holds the metrics namespace/first part
const Namespace = "dell_hw"

var or *omreport.OMReport

// Factories contains the list of all available collectors.
var Factories = make(map[string]func() (Collector, error))

// Collector is the interface a collector has to implement.
type Collector interface {
	// Get new metrics and expose them via prometheus registry.
	Update(ch chan<- prometheus.Metric) error
}

// SetOMReport a given OMReport for the collectors
func SetOMReport(omrep *omreport.OMReport) {
	or = omrep
}
