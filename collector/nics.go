package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type nicsCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["nics"] = NewNicsCollector
}

// NewNicsCollector returns a new nicsCollector
func NewNicsCollector() (Collector, error) {
	return &nicsCollector{}, nil
}

// Update Prometheus metrics
func (c *nicsCollector) Update(ch chan<- prometheus.Metric) error {
	nics, err := or.Nics()
	if err != nil {
		return err
	}
	for _, value := range nics {
		float, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return err
		}
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Connection status of network cards.",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}

	return nil
}
