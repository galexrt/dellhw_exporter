package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type chassisBatteriesCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["chassis_batteries"] = NewChassisBatteriesCollector
}

// NewChassisBatteriesCollector returns a new chassisBatteriesCollector
func NewChassisBatteriesCollector() (Collector, error) {
	return &chassisBatteriesCollector{}, nil
}

// Update Prometheus metrics
func (c *chassisBatteriesCollector) Update(ch chan<- prometheus.Metric) error {
	chassisBatteries, err := or.ChassisBatteries()
	if err != nil {
		return err
	}
	for _, value := range chassisBatteries {
		float, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return err
		}
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Overall status of chassis batteries",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}

	return nil
}
