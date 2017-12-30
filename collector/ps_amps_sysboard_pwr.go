package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type psampssysboardpwrCollector struct {
	current *prometheus.Desc
}

func init() {
	Factories["ps_amps_sysboard_pwr"] = NewPsAmpsSysboardPwrCollector
}

func NewPsAmpsSysboardPwrCollector() (Collector, error) {
	return &psampssysboardpwrCollector{}, nil
}

func (c *psampssysboardpwrCollector) Update(ch chan<- prometheus.Metric) error {
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
