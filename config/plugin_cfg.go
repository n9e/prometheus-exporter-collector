package config

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PluginCfg struct {
	ExporterUrls        []string          `json:"exporter_urls,omitempty"`
	AppendTags          []string          `json:"append_tags"`
	Endpoint            string            `json:"endpoint"`
	Timeout             int               `json:"timeout"`
	IgnoreMetricsPrefix []string          `json:"ignore_metrics_prefix"`
	MetricPrefix        string            `json:"metric_prefix"`
	MetricType          map[string]string `json:"metric_type"`
}

var (
	Config        *PluginCfg
	AppendTagsMap = make(map[string]string)
)

func Get() *PluginCfg {
	return Config
}

func Parse(bs []byte) error {
	Config = &PluginCfg{
		ExporterUrls:        []string{},
		AppendTags:          []string{},
		Endpoint:            "",
		Timeout:             500,
		IgnoreMetricsPrefix: []string{},
		MetricPrefix:        "",
		MetricType:          make(map[string]string),
	}

	if err := json.Unmarshal(bs, &Config); err != nil {
		return err
	}

	if len(Config.ExporterUrls) == 0 {
		return fmt.Errorf("exporter urls is nil")
	}

	if err := parseAppendTagsMap(); err != nil {
		return err
	}

	return nil
}

func parseAppendTagsMap() error {
	appendTags := Config.AppendTags
	if appendTags == nil {
		return nil
	}

	size := len(appendTags)
	if size == 0 {
		return nil
	}

	for _, tag := range appendTags {

		tag = strings.Replace(tag, " ", "", -1)
		if tag == "" {
			continue
		}

		tagPair := strings.SplitN(tag, "=", 2)
		if len(tagPair) == 2 {
			AppendTagsMap[tagPair[0]] = tagPair[1]
		} else {
			return fmt.Errorf("bad tag %s", tag)
		}
	}

	return nil
}
