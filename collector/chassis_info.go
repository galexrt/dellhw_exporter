/*
Copyright 2024 The dellhw_exporter Authors. All rights reserved.

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
	"github.com/prometheus/client_golang/prometheus"
)

type chassisInfoCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["chassis_info"] = NewChassisInfoCollector
}

// NewChassisCollector returns a new chassisInfoCollector
func NewChassisInfoCollector(cfg *Config) (Collector, error) {
	return &chassisInfoCollector{}, nil
}

// Update Prometheus metrics
func (c *chassisInfoCollector) Update(ch chan<- prometheus.Metric) error {
	chassisInfo, err := or.ChassisInfo()
	if err != nil {
		return err
	}
	for _, value := range chassisInfo {
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"Chassis info details in labels.",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, 0)
	}

	return nil
}
