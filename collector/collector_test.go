package collector

import (
	"github.com/n9e/prometheus-collector/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

const validConfigParam = `{"target_urls":["http://192.168.2.30:9103/metrics?dns=192.168.1.7:3306"],"endpoint":"","service":"","step":10,"username":"","password":""}`

func TestConfigParse(t *testing.T) {
	err := config.Parse(validConfigParam)
	assert.NoError(t, err)
}

const validPromGaugeMetric = `# HELP mysql_global_status_buffer_pool_dirty_pages Innodb buffer pool dirty pages.
# TYPE mysql_global_status_buffer_pool_dirty_pages gauge
mysql_global_status_buffer_pool_dirty_pages 132
# HELP mysql_global_status_buffer_pool_pages Innodb buffer pool pages by state.
# TYPE mysql_global_status_buffer_pool_pages gauge
mysql_global_status_buffer_pool_pages{state="data"} 794548
mysql_global_status_buffer_pool_pages{state="free"} 243560
mysql_global_status_buffer_pool_pages{state="misc"} 10460
mysql_global_status_buffer_pool_pages{state="old"} 293137
`
const validPromCounterMetric = `# HELP mysql_global_status_buffer_pool_page_changes_total Innodb buffer pool page state changes.
# TYPE mysql_global_status_buffer_pool_page_changes_total counter
mysql_global_status_buffer_pool_page_changes_total{operation="flushed"} 1.4547024e+07
mysql_global_status_buffer_pool_page_changes_total{operation="lru_flushed"} 0
mysql_global_status_buffer_pool_page_changes_total{operation="made_not_young"} 1.7149202e+07
mysql_global_status_buffer_pool_page_changes_total{operation="made_young"} 979263
# HELP mysql_exporter_scrapes_total Total number of times MySQL was scraped for metrics.
# TYPE mysql_exporter_scrapes_total counter
mysql_exporter_scrapes_total 1
`
const validPromSummaryMetric = `# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
go_gc_duration_seconds{quantile="0.75"} 0
go_gc_duration_seconds{quantile="1"} 0
go_gc_duration_seconds_sum 0
go_gc_duration_seconds_count 0
`
const validPromHistogramMetric = `# HELP mysql_global_status_test histogram test
# TYPE mysql_global_status_test histogram
mysql_global_status_test_bucket{le="500"} 0
mysql_global_status_test_bucket{le="50"} 0
mysql_global_status_test_bucket{le="5"} 0
mysql_global_status_test_bucket{le="1"} 0
mysql_global_status_test_sum 0
mysql_global_status_test_count 0
`
const validPromUntypedMetric = `# HELP mysql_global_status_aborted_clients Generic metric from SHOW GLOBAL STATUS.
# TYPE mysql_global_status_aborted_clients untyped
mysql_global_status_aborted_clients 62539
`

func TestPromMetricParse(t *testing.T) {
	err := config.Parse(validConfigParam)
	assert.NoError(t, err)

	metrics, err := Parse([]byte(validPromUntypedMetric))
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)

	metrics, err = Parse([]byte(validPromGaugeMetric))
	assert.NoError(t, err)
	assert.Len(t, metrics, 5)

	metrics, err = Parse([]byte(validPromCounterMetric))
	assert.NoError(t, err)
	assert.Len(t, metrics, 5)

	metrics, err = Parse([]byte(validPromSummaryMetric))
	assert.NoError(t, err)
	assert.Len(t, metrics, 7)

	metrics, err = Parse([]byte(validPromHistogramMetric))
	assert.NoError(t, err)
	assert.Len(t, metrics, 6)
}
