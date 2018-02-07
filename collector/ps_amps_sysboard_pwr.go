package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type psAmpsSysboardPwrCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["ps_amps_sysboard_pwr"] = NewPsAmpsSysboardPwrCollector
}

// NewPsAmpsSysboardPwrCollector returns a new psAmpsSysboardPwrCollector
func NewPsAmpsSysboardPwrCollector() (Collector, error) {
	return &psAmpsSysboardPwrCollector{}, nil
}

// Update Prometheus metrics
func (c *psAmpsSysboardPwrCollector) Update(ch chan<- prometheus.Metric) error {
	psampssysboardpwr, err := or.PsAmpsSysboardPwr()
	if err != nil {
		return err
	}
	for _, value := range psampssysboardpwr {
		float, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return err
		}
		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", value.Name),
			"System board power usage.",
			nil, value.Labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float)
	}

	return nil
}
