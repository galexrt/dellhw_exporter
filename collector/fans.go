package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type fansCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["fans"] = NewFansCollector
}

func NewFansCollector() (Collector, error) {
	return &fansCollector{}, nil
}

func (c *fansCollector) Update(ch chan<- prometheus.Metric) error {
	fans, err := or.Fans()
	if err != nil {
		return err
	}
	for _, value := range fans {
		float, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return err
		}
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Overall status of system fans.",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}

	return nil
}
