package main

import (
	"io"
	"log/slog"
	"testing"

	"github.com/galexrt/dellhw_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
)

type availabilityTestCollector struct {
	available bool
}

func (c availabilityTestCollector) Update(chan<- prometheus.Metric) error {
	return nil
}

func (c availabilityTestCollector) IsAvailable() bool {
	return c.available
}

func TestLoadCollectorsAvailabilityCheck(t *testing.T) {
	previousFactory, hadPreviousFactory := collector.Factories["availability_test"]
	defer func() {
		if hadPreviousFactory {
			collector.Factories["availability_test"] = previousFactory
		} else {
			delete(collector.Factories, "availability_test")
		}
	}()

	logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	collector.Factories["availability_test"] = func(*collector.Config) (collector.Collector, error) {
		return availabilityTestCollector{available: true}, nil
	}

	loaded, err := loadCollectors([]string{"availability_test"}, []string{"availability_test"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := loaded["availability_test"]; !ok {
		t.Fatal("available collector was not loaded")
	}

	collector.Factories["availability_test"] = func(*collector.Config) (collector.Collector, error) {
		return availabilityTestCollector{available: false}, nil
	}

	loaded, err = loadCollectors([]string{"availability_test"}, []string{"availability_test"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := loaded["availability_test"]; ok {
		t.Fatal("unavailable collector was loaded")
	}
}
