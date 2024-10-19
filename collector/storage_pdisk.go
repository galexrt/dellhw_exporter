/*
Copyright 2021 The dellhw_exporter Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type storagePdiskCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["storage_pdisk"] = NewStoragePdiskCollector
}

// NewStoragePdiskCollector returns a new storagePdiskCollector
func NewStoragePdiskCollector(cfg *Config) (Collector, error) {
	return &storagePdiskCollector{}, nil
}

// Update Prometheus metrics
func (c *storagePdiskCollector) Update(ch chan<- prometheus.Metric) error {
	controllers, err := or.StorageController()
	if err != nil {
		return err
	}
	for cid := range controllers {
		logger := logger.With(zap.Int("controller", cid))
		logger.Debug("collecting pdisks from controller")

		storagePdisk, err := or.StoragePdisk(strconv.Itoa(cid))
		if err != nil {
			return err
		}

		for _, value := range storagePdisk {
			logger.Debug("iterating pdisks from controller", zap.Stringer("pdisk", value))
			float, err := strconv.ParseFloat(value.Value, 64)
			if err != nil {
				return err
			}

			c.current = prometheus.NewDesc(
				prometheus.BuildFQName(Namespace, "", value.Name),
				"Overall status of physical disks + failure prediction (if available).",
				nil, value.Labels)
			ch <- prometheus.MustNewConstMetric(
				c.current, prometheus.GaugeValue, float)
		}
	}

	return nil
}
