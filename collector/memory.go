package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type memoryCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["memory"] = NewMemoryCollector
}

func NewMemoryCollector() (Collector, error) {
	return &memoryCollector{}, nil
}

func (c *memoryCollector) Update(ch chan<- prometheus.Metric) error {
	memory, err := or.Memory()
	if err != nil {
		return err
	}
	for _, value := range memory {
		float, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return err
		}
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"System RAM DIMM status.",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}

	return nil
}
