package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type storageControllerCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["storage_controller"] = NewStorageControllerCollector
}

func NewStorageControllerCollector() (Collector, error) {
	return &storageControllerCollector{}, nil
}

func (c *storageControllerCollector) Update(ch chan<- prometheus.Metric) error {
	storageController, err := or.StorageController()
	if err != nil {
		return err
	}
	for _, value := range storageController {
		float, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return err
		}
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Overall status of storage controllers.",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}

	return nil
}
