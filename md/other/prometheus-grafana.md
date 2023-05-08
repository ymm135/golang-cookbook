- # prometheus-grafana  

- [环境搭建](#环境搭建)
  - [prometheus](#prometheus)
  - [prometheus-docker](#prometheus-docker)
  - [grafana](#grafana)
  - [grafana-docker](#grafana-docker)
- [应用](#应用)
  - [prometheus+grafana](#prometheusgrafana)
  - [prometheus](#prometheus-1)
  - [mysql](#mysql)
  - [logging  日志统计](#logging--日志统计)



## 环境搭建  
### prometheus
https://prometheus.io/  

https://prometheus.io/docs/prometheus/latest/getting_started/  

```sh
wget https://github.com/prometheus/prometheus/releases/download/v2.44.0-rc.0/prometheus-2.44.0-rc.0.linux-amd64.tar.gz  
tar xvfz prometheus-*.tar.gz
cd prometheus-*
```

配置`prometheus.yml`  
```yml
global:
  scrape_interval:     15s # By default, scrape targets every 15 seconds.

  # Attach these labels to any time series or alerts when communicating with
  # external systems (federation, remote storage, Alertmanager).
  external_labels:
    monitor: 'codelab-monitor'

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: 'prometheus'

    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s

    static_configs:
      - targets: ['localhost:9090']
```

启动  
```sh
./prometheus --config.file=prometheus.yml
```

`http://localhost:9090/metrics`查看所有的数据  


<br>
<div align=center>
    <img src="../../res/other/prometheus-1.png" width="90%"></img>  
</div>
<br>


https://github.com/prometheus/node_exporter  
https://github.com/prometheus/mysqld_exporter  

```sh
tar -xzvf node_exporter-*.*.tar.gz
cd node_exporter-*.*

# Start 3 example targets in separate terminals:
./node_exporter --web.listen-address 0.0.0.0:8080
./node_exporter --web.listen-address 0.0.0.0:8081
./node_exporter --web.listen-address 0.0.0.0:8082
```

可以查看配置  
http://localhost:8080/metrics  
http://localhost:8081/metrics  
http://localhost:8082/metrics  

配置Prometheus监听这些数据:  

```sh
scrape_configs:
  - job_name:       'node'

    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s

    static_configs:
      - targets: ['localhost:8080', 'localhost:8081']
        labels:
          group: 'production'

      - targets: ['localhost:8082']
        labels:
          group: 'canary'
```

<br>
<div align=center>
    <img src="../../res/other/prometheus-2.png" width="90%"></img>  
</div>
<br>

`node_cpu_seconds_total{instance="localhost:8080"}` 查到的数据:  
```sh
node_cpu_seconds_total{cpu="0", group="production", instance="localhost:8080", job="node", mode="idle"}
20986.74
node_cpu_seconds_total{cpu="0", group="production", instance="localhost:8080", job="node", mode="iowait"}
117.21
node_cpu_seconds_total{cpu="0", group="production", instance="localhost:8080", job="node", mode="irq"}
0
node_cpu_seconds_total{cpu="0", group="production", instance="localhost:8080", job="node", mode="nice"}
4.12
node_cpu_seconds_total{cpu="0", group="production", instance="localhost:8080", job="node", mode="softirq"}
2.53
node_cpu_seconds_total{cpu="0", group="production", instance="localhost:8080", job="node", mode="steal"}
0
node_cpu_seconds_total{cpu="0", group="production", instance="localhost:8080", job="node", mode="system"}
118.47
node_cpu_seconds_total{cpu="0", group="production", instance="localhost:8080", job="node", mode="user"}
283.84
node_cpu_seconds_total{cpu="1", group="production", instance="localhost:8080", job="node", mode="idle"}
20941.04
node_cpu_seconds_total{cpu="1", group="production", instance="localhost:8080", job="node", mode="iowait"}
153.79
node_cpu_seconds_total{cpu="1", group="production", instance="localhost:8080", job="node", mode="irq"}
0
node_cpu_seconds_total{cpu="1", group="production", instance="localhost:8080", job="node", mode="nice"}
4.38
node_cpu_seconds_total{cpu="1", group="production", instance="localhost:8080", job="node", mode="softirq"}
3.65
node_cpu_seconds_total{cpu="1", group="production", instance="localhost:8080", job="node", mode="steal"}
0
node_cpu_seconds_total{cpu="1", group="production", instance="localhost:8080", job="node", mode="system"}
118.49
node_cpu_seconds_total{cpu="1", group="production", instance="localhost:8080", job="node", mode="user"}
286.06
```

节点状态:  
<br>
<div align=center>
    <img src="../../res/other/prometheus-4.png" width="90%"></img>  
</div>
<br>

系统配置:  
<br>
<div align=center>
    <img src="../../res/other/prometheus-5.png" width="90%"></img>  
</div>
<br>

### prometheus-docker
https://hub.docker.com/r/prom/prometheus/  

独立安装:
```sh
docker run \
    -p 9090:9090 \
    --name prometheus \
    -itd \
    prom/prometheus
```

```sh
docker run \
    -p 9090:9090 \
    --name prometheus \
    -itd \
    -v /path/to/prometheus.yml:/etc/prometheus/prometheus.yml \
    prom/prometheus
```

http://localhost:9090/  

使用docker启动一个ubuntu容器，`docker run --privileged=true -itd --name ssh-ubuntu -p 9100:9100 -p 222:22 atxiaoming/ssh-ubuntu:20.04`, 连接使用`ssh -p 222 root@localhost`, 密码123456  
安装一个node  
https://github.com/prometheus/node_exporter  

```sh
wget https://github.com/prometheus/node_exporter/releases/download/v1.5.0/node_exporter-1.5.0.linux-amd64.tar.gz  

tar -zxvf node_exporter-*linux-amd64.tar.gz
```

查看帮助指令:
```sh
usage: node_exporter [<flags>]

Flags:
  -h, --help                     Show context-sensitive help (also try --help-long and --help-man).
      --collector.arp.device-include=COLLECTOR.ARP.DEVICE-INCLUDE  
                                 Regexp of arp devices to include (mutually exclusive to device-exclude).
      --collector.cpu.info       Enables metric cpu_info
      --collector.diskstats.device-include=COLLECTOR.DISKSTATS.DEVICE-INCLUDE  
                                 Regexp of diskstats devices to include (mutually exclusive to device-exclude).
      --collector.ethtool.device-include=COLLECTOR.ETHTOOL.DEVICE-INCLUDE 
```

运行客户端:`./node_exporter --web.listen-address 0.0.0.0:9100`  

静态配置的局限性，可以使用自动发现:  
- 基于文件的服务发现
- 基于API的服务发现  

这里使用基于`Consul`服务发现  

<br>
<div align=center>
    <img src="../../res/other/prometheus-9.png" width="75%"></img>  
</div>
<br>

启动单节点consul服务:  
https://hub.docker.com/_/consul  
```sh
$ docker run --name consul -d -p 8500:8500 consul
```

访问: http://localhost:8500  , ip地址为:`172.17.0.5`  

通过api注册node_exporter, 在node节点运行:  
```sh
$ curl -X PUT -d '{"id": "node-exporter","name": "node-exporter-172.17.0.4","address": "172.17.0.4","port": 9100,"tags": ["test"],"checks": [{"http": "http://172.17.0.4:9100/metrics", "interval": "5s"}]}'  http://172.17.0.5:8500/v1/agent/service/register
```

这是可以在`consul`看见`node-exporter-172.17.0.4`节点  

注销服务:`$ curl -X PUT http://172.17.0.5:8500/v1/agent/service/deregister/node-exporter`  


配置 Prometheus 实现自动服务发现`prometheus.yml`  
```sh
...
- job_name: 'consul-prometheus'
  consul_sd_configs:
  - server: '172.17.0.5:8500'
    services: []  
```

使用`http:localhost:9090/targets`
```
可以看到设备: http://172.17.0.4:9100/metrics UP instance="172.17.0.4:9100" job="consul-prometheus"  
```

通过`http://localhost:9090/service-discovery?search=`可以看到`consul-prometheus`注册了多少服务  

http://localhost:9100/metrics  可以看到数据  

<br>
<div align=center>
    <img src="../../res/other/prometheus-10.png" width="85%"></img>  
</div>
<br>

### grafana

https://grafana.com/docs/grafana/latest/  

https://grafana.com/docs/grafana/latest/setup-grafana/installation/debian/  
```sh
sudo apt-get install -y adduser libfontconfig1
wget https://dl.grafana.com/enterprise/release/grafana-enterprise_9.5.1_amd64.deb
sudo dpkg -i grafana-enterprise_9.5.1_amd64.deb
```

启动
```sh
sudo systemctl daemon-reload
sudo systemctl start grafana-server
sudo systemctl status grafana-server
```

启动参数:  
```sh
/usr/share/grafana/bin/grafana server --config=/etc/grafana/grafana.ini --pidfile=/run/grafana/grafana-server.pid
```

访问:`http://localhost:3000`, 用户名和密码: `admin`  


<br>
<div align=center>
    <img src="../../res/other/prometheus-3.png" width="90%"></img>  
</div>
<br>

grafana配置  

https://grafana.com/docs/grafana/latest/setup-grafana/configure-grafana/  

> 也需要配置时区  

### grafana-docker
```sh
docker run -itd \
  -p 3000:3000 \
  --name=grafana \
  grafana/grafana-enterprise
```

>   -e "GF_INSTALL_PLUGINS=grafana-clock-panel,grafana-simple-json-datasource" \  

## 应用  
### prometheus+grafana 

启动prometheus之后，在启动一个节点，监听`./node_exporter --web.listen-address 0.0.0.0:8082`   

grafana首页，点击`Add data source`，选择`Prometheus`, 

<br>
<div align=center>
    <img src="../../res/other/prometheus-6.png" width="90%"></img>  
</div>
<br>


### prometheus  

<br>
<div align=center>
    <img src="../../res/other/prometheus-8.png" width="90%"></img>  
</div>
<br>

>  prometheus 时序数据库  

如果使用docker搭建，直接填写容器ip地址`http://172.17.0.2:9090`  




### mysql  

在数据源配置mysql  

<br>
<div align=center>
    <img src="../../res/other/prometheus-7.png" width="90%"></img>  
</div>
<br>

mysql后台更新数据,granfana其实已经同步了数据,只是视图范围没有更新`zoom`


### logging  日志统计  

- ELK ：用Logstash收集和处理日志；将数据存储到ElasticSearch进行索引；用Kibana进行可视化显示；  
- Loki:   Loki 操作更简单，运行成本更低，promtail收集日志，Grafana不用多说，界面酷炫！  
- Graylog: Elasticsearch用来持久化存储和检索日志文件数据(IO 密集), MongoDb用来存储关于 Graylog 的相关配置, Graylog来提供 Web 界面和对外接口的(CPU 密集)。 

> 还有一个日志中间件备选:  Fluentd  

可以把不同设备、系统、应用的日志收集在一起，包含coredumps文件等。统一存储及管理，方便分析及定位问题。  

<br>
<div align=center>
    <img src="../../res/other/fluentd.png" width="80%"></img>  
</div>
<br>

- Installing Loki
- Installing Promtail

https://github.com/grafana/loki   

```sh
$ curl -O -L "https://github.com/grafana/loki/releases/download/v2.8.1/loki-linux-amd64.zip"
# extract the binary
$ unzip "loki-linux-amd64.zip"
# make sure it is executable
$ chmod a+x "loki-linux-amd64"
```

<br>
<div align=center>
    <img src="../../res/other/simple-scalable-test-environment.png" width="80%"></img>  
</div>
<br>

