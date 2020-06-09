package collector

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/didi/nightingale/src/dataobj"
	"github.com/n9e/prometheus-collector/config"
)

func Gather() []*dataobj.MetricValue {
	var wg sync.WaitGroup
	var res []*dataobj.MetricValue

	cfg := config.Get()
	metricChan := make(chan *dataobj.MetricValue)
	done := make(chan struct{}, 1)

	go func() {
		defer func() { done <- struct{}{} }()
		for m := range metricChan {
			res = append(res, m)
		}
	}()

	for _, url := range cfg.ExporterUrls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if metrics, err := gatherExporter(url); err == nil {
				for _, m := range metrics {
					metricChan <- m
				}
			}
		}()
	}

	wg.Wait()
	close(metricChan)

	<-done

	return res
}

func gatherExporter(url string) ([]*dataobj.MetricValue, error) {
	body, err := gatherExporterUrl(url)
	if err != nil {
		log.Printf("gather metrics from exporter error, url :[%s] ,error :%v", url, err)
		return nil, err
	}

	metrics, err := Parse(body)
	if err != nil {
		log.Printf("parse metrics error, url :[%s] ,error :%v", url, err)
		return nil, err
	}

	return metrics, nil
}

func gatherExporterUrl(url string) ([]byte, error) {
	var buf []byte
	var req *http.Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return buf, err
	}

	client := &http.Client{
		Timeout: time.Duration(config.Get().TimeOut) * time.Millisecond,
	}

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return buf, fmt.Errorf("error making HTTP request to %s: %s", url, err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return buf, fmt.Errorf("%s returned HTTP status %s", url, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return buf, fmt.Errorf("error reading body: %s", err)
	}

	return body, nil
}
