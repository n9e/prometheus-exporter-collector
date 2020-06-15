# Prometheus-collector
作为Nightingale的插件，用于收集prometheus的指标

prometheus作为优秀的开源监控产品，本身不仅完整的指标体系，还拥有丰富的指标采集解决方案。通过各种exporter可以覆盖中间件，操作系统，开发语言等等方面的监控指标采集

Prometheus-collector以插件形式集成在collector中，通过Nightingale插件采集，collector采集目标exporter指标并上报

## 快速构建 

    $ mkdir -p $GOPATH/src/github.com/n9e
    $ cd $GOPATH/src/github.com/n9e
    $ git clone https://github.com/n9e/prometheus-collector.git
    $ cd prometheus-collector
    $ go build
    $ cat plugin.test.json | ./prometheus-collector 


 ### 配置参数
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
