package model

import (
	"github.com/didi/nightingale/src/dataobj"
	"github.com/n9e/prometheus-exporter-collector/config"
)

func newMetricValue(metric string, val interface{}, dataType string, ts int64, tagsMap map[string]string) *dataobj.MetricValue {
	mv := dataobj.MetricValue{
		Metric:       metric,
		ValueUntyped: val,
		CounterType:  dataType,
		Timestamp:    ts,
		TagsMap:      tagsMap,
	}

	return &mv
}

func NewGaugeMetric(metric string, val interface{}, ts int64, tagsMap map[string]string) *dataobj.MetricValue {
	return newMetricValue(metric, val, "GAUGE", ts, tagsMap)
}

func NewCounterMetric(metric string, val interface{}, ts int64, tagsMap map[string]string) *dataobj.MetricValue {
	return newMetricValue(metric, val, "COUNTER", ts, tagsMap)
}

func NewSubtractMetric(metric string, val interface{}, ts int64, tagsMap map[string]string) *dataobj.MetricValue {
	return newMetricValue(metric, val, "SUBTRACT", ts, tagsMap)
}

func NewCumulativeMetric(metric string, val interface{}, ts int64, tagsMap map[string]string) *dataobj.MetricValue {
	if config.Get().CumulativeMetricType == "COUNTER" {
		return NewCounterMetric(metric, val, ts, tagsMap)
	}

	return NewSubtractMetric(metric, val, ts, tagsMap)
}
