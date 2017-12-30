package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type systemCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["system"] = NewSystemCollector
}

func NewSystemCollector() (Collector, error) {
	return &systemCollector{}, nil
}

func (c *systemCollector) Update(ch chan<- prometheus.Metric) error {
	system, err := or.System()
	if err != nil {
		return err
	}
	for _, value := range system {
		float, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return err
		}
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Overall status of system components.",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}

	return nil
}
