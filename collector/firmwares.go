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
