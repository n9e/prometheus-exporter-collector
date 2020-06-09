package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/didi/nightingale/src/dataobj"
	"github.com/n9e/prometheus-collector/collector"
	"github.com/n9e/prometheus-collector/config"
)

func main() {
	// stdin
	parseStdin(os.Stdin)
	// gather exporter metrics
	metrics := collector.Gather()
	// stdout
	printMetrics(metrics)
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
