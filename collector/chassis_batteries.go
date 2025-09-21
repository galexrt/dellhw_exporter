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
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type chassisBatteriesCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["chassis_batteries"] = NewChassisBatteriesCollector
}

// NewChassisBatteriesCollector returns a new chassisBatteriesCollector
func NewChassisBatteriesCollector(cfg *Config) (Collector, error) {
	return &chassisBatteriesCollector{}, nil
}

// Update Prometheus metrics
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
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Overall status of chassis batteries",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}

	return nil
}

// IsAvailable if the collector is available
func (c *chassisBatteriesCollector) IsAvailable() bool {
	_, err := or.ChassisBatteries()
	if err == nil {
		return true
	}

	e := strings.ToLower(err.Error())
	return strings.Contains(e, "no battery probes found on this system")
}
