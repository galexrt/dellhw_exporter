package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type firmwaresCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["firmwares"] = NewFirmwaresCollector
}

// NewFirmwaresCollector returns a new firmwaresCollector
func NewFirmwaresCollector() (Collector, error) {
	return &firmwaresCollector{}, nil
}

// Update Prometheus metrics
func (c *firmwaresCollector) Update(ch chan<- prometheus.Metric) error {
	chassisBios, err := or.ChassisBios()
	if err != nil {
		return err
	}
	chassisFirmware, err := or.ChassisFirmware()
	if err != nil {
		return err
	}
	for _, value := range chassisBios {
		float, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return err
		}
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Version info of firmwares/bios.",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}
	for _, value := range chassisFirmware {
		float, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return err
		}
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Version info of firmwares/bios.",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}

	return nil
}
