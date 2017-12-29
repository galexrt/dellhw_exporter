package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type chassisBatteriesCollector struct {
	current []*prometheus.Desc
}

func init() {
	Factories["chassis_batteries"] = NewChassisBatteriesCollector
}

func NewChassisBatteriesCollector() (Collector, error) {
	return &chassisBatteriesCollector{}, nil
}

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
		current := prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Status of cmos battery",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			current, prometheus.GaugeValue, float)
	}

	return nil
}
