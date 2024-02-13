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
	"github.com/galexrt/dellhw_exporter/pkg/omreport"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// Namespace holds the metrics namespace/first part
const Namespace = "dell_hw"

var (
	or  *omreport.OMReport
	log = logrus.New()
)

type Config struct {
	MonitoredNICs []string
}

// Factories contains the list of all available collectors.
var Factories = make(map[string]func(*Config) (Collector, error))

// Collector is the interface a collector has to implement.
type Collector interface {
	// Get new metrics and expose them via prometheus registry.
	Update(ch chan<- prometheus.Metric) error
}

// SetOMReport a given OMReport for the collectors
func SetOMReport(omrep *omreport.OMReport) {
	or = omrep
}

// SetLogger
func SetLogger(logger *logrus.Logger) {
	log = logger
}

func init() {
	log.SetLevel(logrus.ErrorLevel)
}
