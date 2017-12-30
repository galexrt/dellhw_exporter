package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type processorsCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["processors"] = NewProcessorsCollector
}

func NewProcessorsCollector() (Collector, error) {
	return &processorsCollector{}, nil
}

func (c *processorsCollector) Update(ch chan<- prometheus.Metric) error {
	chassis, err := or.Processors()
	if err != nil {
		return err
	}
	for _, value := range chassis {
		float, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return err
		}
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Overall status of CPUs.",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}

	return nil
}
