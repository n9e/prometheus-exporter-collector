package model

const Usage = `Prometheus-Collector, The nightingale plugin for collecting and reporting prometheus metrics.

Usage:

  prometheus-collector [commands|flags]

The commands & flags are:

  -p <DATA>   json format data :
               {
			"target_urls": ["http://xx:9104/metrics?dns=xxip:3306"],
			"endpoint": "",
			"service": "",
			"step": 10,
			"username": "",
			"password": ""
		}

Examples:

  # generate a prometheus-collector param:
  ./prometheus-collector -p {"target_urls": ["http://xx:9104/metrics?dns=xxip:3306"],"endpoint": "","service": "","step": 10,"username": "","password": ""}
`
