package model

import "github.com/didi/nightingale/src/dataobj"

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

func NewMetricWithType(metric string, val interface{}, typ string, ts int64, tagsMap map[string]string) *dataobj.MetricValue {
	return newMetricValue(metric, val, typ, ts, tagsMap)
}
