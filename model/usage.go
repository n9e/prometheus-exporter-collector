package model

const Usage = `Prometheus-Collector, The nightingale plugin for collecting and reporting prometheus metrics.

Usage:

  prometheus-collector [commands|flags]

The commands & flags are:

  -p <DATA>   json format data :
               {
			"exporter_urls": ["http://xx:9104/metrics?dns=xxip:3306"],
			"endpoint": "",
			"service": "",
			"step": 10,
			"username": "",
			"password": "",
			"ignore_exporter_metric": false
		}

Examples:

  # generate a prometheus-collector param:
  ./prometheus-collector -p {"exporter_urls": ["http://127.0.0.1:9104/metrics?dns=xxip:3306"],"endpoint": "","service": "","step": 10,"username": "","password": ""}
`
