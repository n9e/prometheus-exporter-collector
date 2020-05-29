package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/n9e/prometheus-collector/collector"
	"github.com/n9e/prometheus-collector/config"
	"github.com/n9e/prometheus-collector/model"
)

var pparam = flag.String("p", "", "json format data")

func usageExit(rc int) {
	fmt.Println(model.Usage)
	os.Exit(rc)
}

/*todo list
指标名称：_是否转.
指标类型：_count,_sum是否counter类型
exporter gc指标是否上报
exporter：是否支持传参方式
prom untype类型处理
*/
func main() {
	flag.Usage = func() { usageExit(0) }
	flag.Parse()

	if err := config.Parse(*pparam); err != nil {
		fmt.Println("cannot parse param:", err)
		os.Exit(1)
	}

	metrics := collector.Gather()
	// stdout
	fmt.Println(metrics)
}
