package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type psCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["ps"] = NewPsCollector
}

func NewPsCollector() (Collector, error) {
	return &psCollector{}, nil
}

func (c *psCollector) Update(ch chan<- prometheus.Metric) error {
	ps, err := or.Ps()
	if err != nil {
		return err
	}
	for _, value := range ps {
		float, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return err
		}
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Overall status of power supplies.",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}

	return nil
}
