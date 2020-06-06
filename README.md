# Prometheus-collector
A plugin for Nightingale is used to collect metrics from prometheus exporters.

## Building and running 

    $ mkdir -p $GOPATH/src/github.com/n9e
    $ cd $GOPATH/src/github.com/n9e
    $ git clone https://github.com/n9e/prometheus-collector.git
    $ cd prometheus-collector
    $ go build
    $ ./prometheus-collector -p "{\"exporter_urls\": [\"http://127.0.0.1:9104/metrics?dns=$ip:3306\"],\"endpoint\":\"xx\",\"service\":\"xxx\",\"step\":10,\"username\":\"\",\"password\":\"\"}"


 ### Command Parameters
 Name                             |  type     | Description
 ---------------------------------|-----------|--------------------------------------------------------------------------------------------------
 exporter_urls                    | array     | Address to collect metric for prometheus exporter.
 endpoint                         | string    | Field endpoint for n9e metric.
 service                          | string    | Add a service tag for n9e metric.
 step                             | int       | Step for report metrics 
 username                         | string    | Not needed for now
 password                         | string    | Not needed for now
 ignore_exportermetric            | bool      | Ignore exporter metrics default false
 ###
=======
 
 ###
