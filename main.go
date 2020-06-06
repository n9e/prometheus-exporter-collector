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
