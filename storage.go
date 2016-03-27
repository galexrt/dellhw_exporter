package main

type metricStorage struct {
	metrics map[string]interface{}
}

func newMetricStorage() *metricStorage {
	ms := new(metricStorage)
	ms.metrics = make(map[string]interface{})
	return ms
}
