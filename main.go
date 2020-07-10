package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/didi/nightingale/src/dataobj"
	"github.com/n9e/prometheus-exporter-collector/collector"
	"github.com/n9e/prometheus-exporter-collector/config"
	"github.com/n9e/prometheus-exporter-collector/model"
	fmodel "github.com/open-falcon/falcon-plus/common/model"
)

var (
	h bool

	backend string
	step    int64
)

func init() {
	flag.BoolVar(&h, "h", false, "help")
	flag.StringVar(&backend, "b", "n9e", "send metrics to backend: n9e, falcon")
	flag.Int64Var(&step, "s", 60, "set default step of falcon metrics")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: ./prometheus-exporter-collector [-h] [-b backend] [-s step]

Options: 
`)
	flag.PrintDefaults()
}

func parseStdin(r io.Reader) {
	// stdin一般是个json
	stdinData, err := ioutil.ReadAll(r)
	if err != nil {
		// 错误信息使用log库，默认输出到stderr
		log.Printf("cannot read stdin, error: %v", err)
	}

	if err := config.Parse(stdinData); err != nil {
		log.Printf("cannot parse param: %v", err)
		os.Exit(1)
	}
}

func printMetrics(metrics []*dataobj.MetricValue) {
	metricStr, err := json.Marshal(metrics)
	if err != nil {
		log.Printf("parse metrics result error: %v", err)
	}
	// stdout
	fmt.Println(string(metricStr))
}

func printFalconMetrics(metrics []*fmodel.MetricValue) {
	metricStr, err := json.Marshal(metrics)
	if err != nil {
		log.Printf("parse metrics result error: %v", err)
	}
	// stdout
	fmt.Println(string(metricStr))
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return
	}

	// stdin
	parseStdin(os.Stdin)
	// gather exporter metrics
	metrics := collector.Gather()

	if backend == "falcon" {
		// falcon support interval above 10 seconds
		if step < 10 {
			step = 10
		}
		vs := model.FmtFalconMetricValue(metrics, step)
		// stdout to open-falcon
		printFalconMetrics(vs)
	} else {
		// stdout to nightingale
		printMetrics(metrics)
	}
}
