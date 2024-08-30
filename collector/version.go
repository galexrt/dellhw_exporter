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
	"github.com/prometheus/common/version"
)

type versionCollector struct {
	metric *prometheus.GaugeVec
}

func init() {
	Factories["version"] = NewVersionCollector
}

// NewVersionCollector returns a new VersionCollector
func NewVersionCollector(cfg *Config) (Collector, error) {
	metric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dell_hw_exporter_version",
			Help: "Constant '1' value with version, revision, and branch labels from the dellhw_exporter version info.",
		},
		[]string{"version", "revision", "branch"},
	)
	metric.WithLabelValues(version.Version, version.Revision, version.Branch).Set(1)

	return &versionCollector{
		metric: metric,
	}, nil
}

// Update Prometheus metrics
func (c *versionCollector) Update(ch chan<- prometheus.Metric) error {
	c.metric.Collect(ch)

	return nil
}
