package main

import (
	"bytes"
	"fmt"
	"github.com/n9e/prometheus-exporter-collector/collector"
	"github.com/n9e/prometheus-exporter-collector/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

const validConfigParam = `{
  "exporter_urls":["http://127.0.0.1:8080/metrics"],
  "append_tags": ["region=bj", "dept=cloud"],
  "endpoint": "127.0.0.1",
  "ignore_metrics_prefix": [],
  "timeout": 500
}
`

func TestConfigParse(t *testing.T) {
	err := config.Parse([]byte(validConfigParam))
	assert.NoError(t, err)
}

const validPromGaugeMetric = `# HELP mysql_global_status_buffer_pool_dirty_pages Innodb buffer pool dirty pages.
# TYPE mysql_global_status_buffer_pool_dirty_pages gauge
mysql_global_status_buffer_pool_dirty_pages 132
# HELP mysql_global_status_buffer_pool_pages Innodb buffer pool pages by state.
# TYPE mysql_global_status_buffer_pool_pages gauge
mysql_global_status_buffer_pool_pages{state="data"} 794548
mysql_global_status_buffer_pool_pages{state="free"} +Inf
mysql_global_status_buffer_pool_pages{state="misc"} 10460
mysql_global_status_buffer_pool_pages{state="old"} 293137
`
const validPromCounterMetric = `# HELP mysql_global_status_buffer_pool_page_changes_total Innodb buffer pool page state changes.
# TYPE mysql_global_status_buffer_pool_page_changes_total counter
mysql_global_status_buffer_pool_page_changes_total{operation="flushed"} 1.4547024e+07
mysql_global_status_buffer_pool_page_changes_total{operation="lru_flushed"} -Inf
mysql_global_status_buffer_pool_page_changes_total{operation="made_not_young"} 1.7149202e+07
mysql_global_status_buffer_pool_page_changes_total{operation="made_young"} 979263
# HELP mysql_exporter_scrapes_total Total number of times MySQL was scraped for metrics.
# TYPE mysql_exporter_scrapes_total counter
mysql_exporter_scrapes_total 1
`
const validPromSummaryMetric = `# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 543
go_gc_duration_seconds{quantile="0.25"} -Inf
go_gc_duration_seconds{quantile="0.5"} +Inf
go_gc_duration_seconds{quantile="0.75"} Nan
go_gc_duration_seconds{quantile="1"} 0
go_gc_duration_seconds_sum 794548
go_gc_duration_seconds_count 1444
`
const validPromHistogramMetric = `# HELP mysql_global_status_test histogram test
# TYPE mysql_global_status_test histogram
mysql_global_status_test_bucket{le="500"} 111
mysql_global_status_test_bucket{le="50"} 85
mysql_global_status_test_bucket{le="5"} 378
mysql_global_status_test_bucket{le="2"} Nan
mysql_global_status_test_sum 794548
mysql_global_status_test_count 1422
`
const validPromUntypedMetric = `# HELP mysql_global_status_aborted_clients Generic metric from SHOW GLOBAL STATUS.
# TYPE mysql_global_status_aborted_clients untyped
mysql_global_status_aborted_clients 62539
`

func TestPromMetricParser(t *testing.T) {
	err := config.Parse([]byte(validConfigParam))
	assert.NoError(t, err)

	metrics, err := collector.Parse([]byte(validPromUntypedMetric))
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)

	metrics, err = collector.Parse([]byte(validPromGaugeMetric))
	assert.NoError(t, err)
	assert.Len(t, metrics, 5)

	metrics, err = collector.Parse([]byte(validPromCounterMetric))
	assert.NoError(t, err)
	assert.Len(t, metrics, 5)

	metrics, err = collector.Parse([]byte(validPromSummaryMetric))
	assert.NoError(t, err)
	assert.Len(t, metrics, 7)

	metrics, err = collector.Parse([]byte(validPromHistogramMetric))
	assert.NoError(t, err)
	assert.Len(t, metrics, 6)
}

const validStdin = `{
  "exporter_urls": [
    "http://xxxxx:9103/metrics?dns=xxxx:3306",
	"http://xxxxx:9103/metrics?dns=xxxx:3307"
  ],
  "endpoint": "xxxx",
  "append_tags": ["region=bj", "dept=cloud"],
  "ignore_metrics_prefix": ["gc_","go_"],
  "timeout": 500
}
`

func TestStdin(t *testing.T) {
	bs := bytes.NewReader([]byte(validStdin))
	parseStdin(bs)
	p := config.Get()
	fmt.Println(p)
	assert.NotNil(t, p)
}

func TestStdout(t *testing.T) {
	err := config.Parse([]byte(validConfigParam))
	assert.NoError(t, err)

	metrics, err := collector.Parse([]byte(validPromHistogramMetric))
	assert.NoError(t, err)
	assert.Len(t, metrics, 6)

	printMetrics(metrics)
}

func TestGather(t *testing.T) {
	err := config.Parse([]byte(validConfigParam))
	assert.NoError(t, err)

	metrics := collector.Gather()
	assert.NotNil(t, metrics)
	printMetrics(metrics)
}

const validIgnoreConfigParam = `{
  "exporter_urls":["http://127.0.0.1:8080/metrics"],
  "append_tags": ["region=bj", "dept=cloud"],
  "endpoint": "127.0.0.1",
  "ignore_metrics_prefix": ["go_"],
  "timeout": 500
}
`
const validPromIgnoreMetric = `# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 1.2099e-05
go_gc_duration_seconds{quantile="0.25"} 1.3161e-05
go_gc_duration_seconds{quantile="0.5"} 1.3841e-05
go_gc_duration_seconds{quantile="0.75"} 1.4729e-05
go_gc_duration_seconds{quantile="1"} 0.015064859
go_gc_duration_seconds_sum 266.170701672
go_gc_duration_seconds_count 84408
`

func TestIgnoreMetric(t *testing.T) {
	err := config.Parse([]byte(validIgnoreConfigParam))
	assert.NoError(t, err)

	config.Config.IgnoreMetricsPrefix = []string{"mem"}
	metrics, err := collector.Parse([]byte(validPromIgnoreMetric))
	assert.NoError(t, err)
	assert.Len(t, metrics, 7)

	config.Config.IgnoreMetricsPrefix = []string{"go_"}
	metrics, err = collector.Parse([]byte(validPromIgnoreMetric))
	assert.NoError(t, err)
	assert.Len(t, metrics, 0)
}
