# Prometheus-collector
A plugin for Nightingale is used to collect metrics from prometheus exporters.

## Building and running 

    $ mkdir -p $GOPATH/src/github.com/n9e
    $ cd $GOPATH/src/github.com/n9e
    $ git clone https://github.com/n9e/prometheus-collector.git
    $ cd prometheus-collector
    $ go build
    $ cat plugin.test.json | ./prometheus-collector 


 ### Command Parameters
 Name                             |  type     | Description
 ---------------------------------|-----------|--------------------------------------------------------------------------------------------------
 exporter_urls                    | array     | Address to collect metric for prometheus exporter.
 append_tags                      | array     | Add tags for n9e metric default empty
 endpoint                         | string    | Field endpoint for n9e metric default empty
 username                         | string    | Not needed for now
 password                         | string    | Not needed for now
 ignore_metrics_prefix            | array     | Ignore metric prefix default empty
 timeout                          | int       | Timeout for access a exporter url default 500ms
 ###
=======
 
 ###
