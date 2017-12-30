package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type voltsCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["volts"] = NewVoltsCollector
}

func NewVoltsCollector() (Collector, error) {
	return &voltsCollector{}, nil
}

func (c *voltsCollector) Update(ch chan<- prometheus.Metric) error {
	volts, err := or.Volts()
	if err != nil {
		return err
	}
	for _, value := range volts {
		float, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return err
		}
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Overall volts and status of power supply volt readings.",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}

	return nil
}
