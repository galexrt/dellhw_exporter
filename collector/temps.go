package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type tempsCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["temps"] = NewTempsCollector
}

func NewTempsCollector() (Collector, error) {
	return &tempsCollector{}, nil
}

func (c *tempsCollector) Update(ch chan<- prometheus.Metric) error {
	temps, err := or.Temps()
	if err != nil {
		return err
	}
	for _, value := range temps {
		float, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return err
		}
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Overall temperatures and status of system temperature readings.",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}

	return nil
}
