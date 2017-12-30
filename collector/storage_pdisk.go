package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type storagePdiskCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["storage_pdisk"] = NewStoragePdiskCollector
}

func NewStoragePdiskCollector() (Collector, error) {
	return &storagePdiskCollector{}, nil
}

func (c *storagePdiskCollector) Update(ch chan<- prometheus.Metric) error {
	controllers, err := or.StorageController()
	if err != nil {
		return err
	}
	for cid := range controllers {
		storagePdisk, err := or.StoragePdisk(strconv.Itoa(cid))
		if err != nil {
			return err
		}
		for _, value := range storagePdisk {
			float, err := strconv.ParseFloat(value.Value, 64)
			if err != nil {
				return err
			}
			c.current = prometheus.NewDesc(
				prometheus.BuildFQName(Namespace, "", value.Name),
				"Overall status of physical disks.",
				nil, value.Labels)
			ch <- prometheus.MustNewConstMetric(
				c.current, prometheus.GaugeValue, float)
		}
	}

	return nil
}
