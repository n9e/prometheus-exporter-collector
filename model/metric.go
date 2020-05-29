package model

import "github.com/didi/nightingale/src/dataobj"

func newMetricValue(metric string, val interface{}, dataType string, ts int64, tags map[string]string) *dataobj.MetricValue {
	mv := dataobj.MetricValue{
		Metric:       metric,
		ValueUntyped: val,
		CounterType:  dataType,
		Timestamp:    ts,
		TagsMap:      tags,
	}

	return &mv
}

func NewGaugeMetric(metric string, val interface{}, ts int64, tags map[string]string) *dataobj.MetricValue {
	return newMetricValue(metric, val, "GAUGE", ts, tags)
}

func NewCounterMetric(metric string, val interface{}, ts int64, tags map[string]string) *dataobj.MetricValue {
	return newMetricValue(metric, val, "COUNTER", ts, tags)
}
