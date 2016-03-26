package main

import "sync"

type metricStorage struct {
	Lock    sync.RWMutex
	metrics map[string]interface{}
}

func newMetricStorage() *metricStorage {
	ms := new(metricStorage)
	ms.metrics = make(map[string]interface{})
	return ms
}
