package collector

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/didi/nightingale/src/dataobj"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"

	"github.com/n9e/prometheus-exporter-collector/config"
	"github.com/n9e/prometheus-exporter-collector/model"
)

var now = time.Now().Unix()

func Parse(buf []byte) ([]*dataobj.MetricValue, error) {
	var metricList []*dataobj.MetricValue
	var parser expfmt.TextParser
	cfg := config.Get()
	// parse even if the buffer begins with a newline
	buf = bytes.TrimPrefix(buf, []byte("\n"))
	// Read raw data
	buffer := bytes.NewBuffer(buf)
	reader := bufio.NewReader(buffer)

	// Prepare output
	metricFamilies := make(map[string]*dto.MetricFamily)
	metricFamilies, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		return nil, fmt.Errorf("reading text format failed: %s", err)
	}

	now = time.Now().Unix()
	// read metrics
	for basename, mf := range metricFamilies {

		metrics := []*dataobj.MetricValue{}
		for _, m := range mf.Metric {

			// pass ignore metric
			if filterIgnoreMetric(basename) {
				continue
			}

			switch mf.GetType() {
			case dto.MetricType_GAUGE:
				// gauge metric
				metrics = makeCommon(basename, m)
			case dto.MetricType_COUNTER:
				// counter metric
				metrics = makeCommon(basename, m)
			case dto.MetricType_SUMMARY:
				// summary metric
				metrics = makeQuantiles(basename, m)
			case dto.MetricType_HISTOGRAM:
				// histogram metric
				metrics = makeBuckets(basename, m)
			case dto.MetricType_UNTYPED:
				// untyped as gauge
				metrics = makeCommon(basename, m)
			}

			// render endpoint info
			for _, metric := range metrics {
				metric.Endpoint = cfg.Endpoint
				metric.Tags = makeAppendTags(metric.TagsMap, config.AppendTagsMap)
				// set provided Time, ms to s
				if m.GetTimestampMs() > 0 {
					metric.Timestamp = m.GetTimestampMs() / 1000
				}
				metricList = append(metricList, metric)
			}
		}
	}

	return metricList, err
}

// Get Quantiles from summary metric
func makeQuantiles(basename string, m *dto.Metric) []*dataobj.MetricValue {
	metrics := []*dataobj.MetricValue{}
	tags := makeLabels(m)

	// 不能保证取到的metric一定是按照标准输出的,只能尽量去兼容

	if m.GetSummary().SampleCount != nil {
		if !math.IsNaN(float64(m.GetSummary().GetSampleCount())) && !math.IsInf(float64(m.GetSummary().GetSampleCount()), 0) {
			countName := fmt.Sprintf("%s_count", basename)
			metrics = append(metrics, model.NewCumulativeMetric(countName, m.GetSummary().SampleCount, now, tags))
		}
	}
	if m.GetSummary().SampleSum != nil {
		if !math.IsNaN(m.GetSummary().GetSampleSum()) && !math.IsInf(m.GetSummary().GetSampleSum(), 0) {
			sumName := fmt.Sprintf("%s_sum", basename)
			metrics = append(metrics, model.NewCumulativeMetric(sumName, m.GetSummary().SampleSum, now, tags))
		}
	}

	for _, q := range m.GetSummary().Quantile {
		if !math.IsNaN(q.GetValue()) && !math.IsInf(q.GetValue(), 0) {
			tagsNew := make(map[string]string)
			for tagKey, tagValue := range tags {
				tagsNew[tagKey] = tagValue
			}
			tagsNew["quantile"] = fmt.Sprint(q.GetQuantile())

			metrics = append(metrics, model.NewGaugeMetric(basename, q.GetValue(), now, tagsNew))
		}
	}

	return metrics
}

// Get Buckets from histogram metric
func makeBuckets(basename string, m *dto.Metric) []*dataobj.MetricValue {
	metrics := []*dataobj.MetricValue{}
	tags := makeLabels(m)

	if m.GetHistogram().SampleCount != nil {
		if !math.IsNaN(float64(m.GetHistogram().GetSampleCount())) && !math.IsInf(float64(m.GetHistogram().GetSampleCount()), 0) {
			countName := fmt.Sprintf("%s_count", basename)
			metrics = append(metrics, model.NewCumulativeMetric(countName, m.GetHistogram().SampleCount, now, tags))
		}
	}

	if m.GetHistogram().SampleSum != nil {
		if !math.IsNaN(m.GetHistogram().GetSampleSum()) && !math.IsInf(m.GetHistogram().GetSampleSum(), 0) {
			sumName := fmt.Sprintf("%s_sum", basename)
			metrics = append(metrics, model.NewCumulativeMetric(sumName, m.GetHistogram().SampleSum, now, tags))
		}
	}

	for _, b := range m.GetHistogram().Bucket {
		//if !math.IsNaN(float64(b.GetCumulativeCount())) && !math.IsInf(float64(b.GetCumulativeCount()), 0) {
		tagsNew := make(map[string]string)
		for tagKey, tagValue := range tags {
			tagsNew[tagKey] = tagValue
		}
		tagsNew["le"] = fmt.Sprint(b.GetUpperBound())
		bucketName := fmt.Sprintf("%s_bucket", basename)
		metrics = append(metrics, model.NewGaugeMetric(bucketName, float64(b.GetCumulativeCount()), now, tagsNew))
	}
	//}

	return metrics
}

// Get gauge and counter from metric
func makeCommon(metricName string, m *dto.Metric) []*dataobj.MetricValue {
	var val float64
	metrics := []*dataobj.MetricValue{}
	tags := makeLabels(m)
	if m.Gauge != nil {
		if !math.IsNaN(m.GetGauge().GetValue()) && !math.IsInf(m.GetGauge().GetValue(), 0) {
			val = float64(m.GetGauge().GetValue())
			metrics = append(metrics, model.NewGaugeMetric(metricName, val, now, tags))
		}
	} else if m.Counter != nil {
		if !math.IsNaN(m.GetCounter().GetValue()) && !math.IsInf(m.GetCounter().GetValue(), 0) {
			val = float64(m.GetCounter().GetValue())
			metrics = append(metrics, model.NewCumulativeMetric(metricName, val, now, tags))
		}
	} else if m.Untyped != nil {
		// untyped as gauge
		if !math.IsNaN(m.GetUntyped().GetValue()) && !math.IsInf(m.GetUntyped().GetValue(), 0) {
			val = float64(m.GetUntyped().GetValue())
			metrics = append(metrics, model.NewGaugeMetric(metricName, val, now, tags))
		}
	}
	return metrics
}

// Get labels from metric
func makeLabels(m *dto.Metric) map[string]string {
	tags := map[string]string{}
	for _, lp := range m.Label {
		tags[lp.GetName()] = lp.GetValue()
	}
	return tags
}

// append tags
func makeAppendTags(tagsMap map[string]string, appendTagsMap map[string]string) string {
	if len(tagsMap) == 0 && len(appendTagsMap) == 0 {
		return ""
	}

	if len(tagsMap) == 0 {
		return dataobj.SortedTags(appendTagsMap)
	}

	if len(appendTagsMap) == 0 {
		return dataobj.SortedTags(tagsMap)
	}

	for k, v := range appendTagsMap {
		tagsMap[k] = v
	}

	return dataobj.SortedTags(tagsMap)
}

func filterIgnoreMetric(basename string) bool {
	ignorePrefix := config.Get().IgnoreMetricsPrefix
	if len(config.Get().IgnoreMetricsPrefix) == 0 {
		return false
	}

	for _, pre := range ignorePrefix {
		if strings.HasPrefix(basename, pre) {
			return true
		}
	}
	return false
}
