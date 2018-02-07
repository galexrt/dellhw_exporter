package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type storageBatteryCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["storage_battery"] = NewStorageBatteryCollector
}

// NewStorageBatteryCollector returns a new storageBatteryCollector
func NewStorageBatteryCollector() (Collector, error) {
	return &storageBatteryCollector{}, nil
}

// Update Prometheus metrics
func (c *storageBatteryCollector) Update(ch chan<- prometheus.Metric) error {
	storageBattery, err := or.StorageBattery()
	if err != nil {
		return err
	}
	for _, value := range storageBattery {
		float, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return err
		}
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Status of storage controller backup batteries.",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}

	return nil
}
