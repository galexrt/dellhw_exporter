package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type nicsCollector struct {
	current []*prometheus.Desc
}

func init() {
	Factories["nics"] = NewNicsCollector
}

func NewNicsCollector() (Collector, error) {
	return &nicsCollector{}, nil
}

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
		current := prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Status of network card",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			current, prometheus.GaugeValue, float)
	}

	return nil
}
