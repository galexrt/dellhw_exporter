package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type chassisCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["chassis"] = NewChassisCollector
}

func NewChassisCollector() (Collector, error) {
	return &chassisCollector{}, nil
}

func (c *chassisCollector) Update(ch chan<- prometheus.Metric) error {
	chassis, err := or.Chassis()
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
			"Overall status of chassis components.",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}

	return nil
}
