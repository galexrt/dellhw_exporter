package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type storageEnclosureCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["storage_enclosure"] = NewStorageEnclosureCollector
}

func NewStorageEnclosureCollector() (Collector, error) {
	return &storageEnclosureCollector{}, nil
}

func (c *storageEnclosureCollector) Update(ch chan<- prometheus.Metric) error {
	storageEnclosure, err := or.StorageEnclosure()
	if err != nil {
		return err
	}
	for _, value := range storageEnclosure {
		float, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return err
		}
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Overall status of storage enclosures.",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}

	return nil
}
