package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type storageVdiskCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["storage_vdisk"] = NewStorageVdiskCollector
}

func NewStorageVdiskCollector() (Collector, error) {
	return &storageVdiskCollector{}, nil
}

func (c *storageVdiskCollector) Update(ch chan<- prometheus.Metric) error {
	storageVdisk, err := or.StorageVdisk()
	if err != nil {
		return err
	}
	for _, value := range storageVdisk {
		float, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return err
		}
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Overall status of virtual disks.",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}

	return nil
}
