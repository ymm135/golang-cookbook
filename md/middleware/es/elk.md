# ELK
## 结构
### [官网](https://www.elastic.co/cn/what-is/elk-stack)  
![elk-structure](../../../res/elk-structure.png)  

## 安装logstash
[官网](https://www.elastic.co/guide/en/logstash/6.8/installing-logstash.html)  

### linux 安装
```
sudo apt-get update && sudo apt-get install logstash  
```

### docker 安装


## logstash 指南
### [入门](https://www.elastic.co/guide/en/logstash/6.8/getting-started-with-logstash.html)  
### [原理](https://www.elastic.co/guide/en/logstash/6.8/pipeline.html)  
### JDBC Input Plugin

[原文](https://www.elastic.co/guide/en/logstash/6.8/plugins-inputs-jdbc.html)  

```shell
input {
  jdbc {
    jdbc_driver_library => "mysql-connector-java-5.1.36-bin.jar"
    jdbc_driver_class => "com.mysql.jdbc.Driver"
    jdbc_connection_string => "jdbc:mysql://localhost:3306/mydb"
    jdbc_user => "mysql"
    parameters => { "favorite_artist" => "Beethoven" }
    schedule => "* * * * *"
    statement => "SELECT * from songs where artist = :favorite_artist"
  }
}
```


输入日志文件，转发到`http rest` api接口  
```shell
input {
    jdbc {
      # mysql jdbc connection string to our backup databse
      jdbc_connection_string => "jdbc:mysql://ip:port/database?zeroDateTimeBehavior=convertToNull"

      # the user we wish to excute our statement as
      jdbc_user => "xxxx"
      jdbc_password => "xxxx"
      # the path to our downloaded jdbc driver
      jdbc_driver_library => "/home/admin/data/mysql-connector-java-5.1.36/mysql-connector-java-5.1.36-bin.jar"
      # the name of the driver class for mysql
      jdbc_driver_class => "com.mysql.jdbc.Driver"
      jdbc_paging_enabled => "true"
      jdbc_page_size => "50000"
      #statement_filepath => "jdbc.sql"
      statement => "SELECT * from mytable WHERE field = xx"
      type => "jdbc"
    }
}


output {
    http {
        url => "http://ip:port/xxxx"
        http_method => "post"
        format => "form"
        mapping => {"uid"=>"%{follwer_id}" "following_uid"=>"%{following_id}" "follow_source"=>1000}
    }
    stdout {
        codec => json_lines
    }
}
```

## mysql同步工具[canal](https://github.com/alibaba/canal)   

原理  
![canal.png](../../../res/canal.png)  

[docker安装](https://github.com/alibaba/canal/wiki/Docker-QuickStart)  


## 环境搭建
### filebeat  
https://www.elastic.co/guide/en/beats/filebeat/current/filebeat-installation-configuration.html  

```sh
curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-8.8.0-linux-x86_64.tar.gz
tar xzvf filebeat-8.8.0-linux-x86_64.tar.gz
```

连接ES
```sh
output.elasticsearch:
  hosts: ["https://myEShost:9200"]
  username: "filebeat_internal"
  password: "YOUR_PASSWORD" 
  ssl:
    enabled: true
    ca_trusted_fingerprint: "b9a10bbe64ee9826abeda6546fc988c8bf798b41957c33d05db736716513dc9c" 
```

收集日志
```sh
# 查看模块
./filebeat modules list
nginx
mysql
redis

# 开启
./filebeat modules enable nginx

# 配置 vim modules.d/nginx.yml.disabled
- module: nginx
  access:
    enabled: true
    var.paths: ["/var/log/nginx/access.log*"]


# 启动
sudo chown root filebeat.yml 
sudo chown root modules.d/nginx.yml 
sudo ./filebeat -e
```

接着就可以在kibana上查看了  

另外也可以配置输出到logstash, 配置文件`modules.d/logstash.yml.disabled`
```sh
output.logstash:
  hosts: ["127.0.0.1:5044"]
```

### logstash

https://www.elastic.co/guide/en/logstash/current/docker.html

```sh
docker pull docker.elastic.co/logstash/logstash:8.8.0

docker network create --subnet=172.20.1.0/24 elastic
```

配置文件`~/work/devops/elk/logstash/config/logstash.yml`
```yml
http.host: "0.0.0.0"
xpack.monitoring.enabled: true
xpack.monitoring.elasticsearch.hosts: "https://172.20.1.2:9200"  #es地址
xpack.monitoring.elasticsearch.username: "elastic"  #es xpack账号
xpack.monitoring.elasticsearch.password: "1zoajcY3=FP7uNsrdL*w"     #es xpack账号
path.config: /usr/share/logstash/config/conf.d/*.conf
path.logs: /usr/share/logstash/logs
```

```sh
# 启动
docker run --rm -itd --net elastic --ip 172.20.1.5 \
  -p 5044:5044 -p 5000:5000 \
  --name logstash \
  docker.elastic.co/logstash/logstash:8.8.0
```

另外也可以使用自定义参数:
```sh
docker run --rm -itd --net elastic --ip 172.20.1.5 \
  -p 5044:5044 -p 5000:5000 \
  --name logstash \
  -v ~/work/devops/elk/logstash/config/logstash.yml:/usr/share/logstash/config/logstash.yml  \
  -v ~/work/devops/elk/logstash/conf.d:/usr/share/logstash/config/conf.d  \
  -v ~/work/devops/elk/logstash/logs:/usr/share/logstash/logs  \
  docker.elastic.co/logstash/logstash:8.8.0
```

> http://172.20.1.2:9200/_xpack   172.20.1.2:9200 failed to respond  


修改容器配置文件:`/usr/share/logstash/config/logstash.yml`  
```sh
xpack.monitoring.elasticsearch.hosts: [ "https://172.20.1.2:9200" ]
```

### elasticsearch  
https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html

```sh
# 
docker pull docker.elastic.co/elasticsearch/elasticsearch:8.8.0

# 创建私有网络 
docker network create --subnet=172.20.1.0/24 elastic

# 启动
docker run --name es01 --net elastic --ip 172.20.1.2 -p 9200:9200 -p 9300:9300 -itd docker.elastic.co/elasticsearch/elasticsearch:8.8.0

# 登录 user: elastic ,密码: gH-jVifqyBh6rWAwjIX7, 从日志中看
https://localhost:9200/


# 重置密码  
docker exec -it es01 /usr/share/elasticsearch/bin/elasticsearch-reset-password -u elastic

# 节点信息
{
  "name" : "8b29f47ce5e6",
  "cluster_name" : "docker-cluster",
  "cluster_uuid" : "2vUMkNBUTja9Y1ng1p65Ng",
  "version" : {
    "number" : "8.8.0",
    "build_flavor" : "default",
    "build_type" : "docker",
    "build_hash" : "c01029875a091076ed42cdb3a41c10b1a9a5a20f",
    "build_date" : "2023-05-23T17:16:07.179039820Z",
    "build_snapshot" : false,
    "lucene_version" : "9.6.0",
    "minimum_wire_compatibility_version" : "7.17.0",
    "minimum_index_compatibility_version" : "7.0.0"
  },
  "tagline" : "You Know, for Search"
}
```

> 9300 展示的接口  

### kibana  
https://www.elastic.co/guide/en/kibana/current/docker.html  

```sh
docker pull docker.elastic.co/kibana/kibana:8.8.0

# 启动
docker run -itd --name kib01 --net elastic --ip 172.20.1.3 -p 5601:5601 docker.elastic.co/kibana/kibana:8.8.0
```

在es启动时，通过日志获取`enrollment token`,并使用es用户及密码登录:
```sh
ℹ️  Configure Kibana to use this cluster:
2023-05-26 23:09:18 • Run Kibana and click the configuration link in the terminal when Kibana starts.
2023-05-26 23:09:18 • Copy the following enrollment token and paste it into Kibana in your browser (valid for the next 30 minutes):
2023-05-26 23:09:18   eyJ2ZXIiOiI4LjguMCIsImFkciI6WyIxNzIuMjAuMS4yOjkyMDAiXSwiZmdyIjoiZTg4NjQyYThkYWQ3NTk4YzE1ZjQyZmE2ZGZjMmYxODQxMDYwZDhmYTQ1ZjE5MGRkY2UxZDlkOGZlM2Q2ZDM1MyIsImtleSI6IlM5bWJXSWdCeFdFcEFjVTdZRzF1OmZ3WEk4SWYyUzJlN0hmemJyMzQ3cVEifQ==

2023-05-26 23:09:18 ℹ️ Configure other nodes to join this cluster:
2023-05-26 23:09:18 • Copy the following enrollment token and start new Elasticsearch nodes with `bin/elasticsearch --enrollment-token <token>` (valid for the next 30 minutes):
2023-05-26 23:09:18   eyJ2ZXIiOiI4LjguMCIsImFkciI6WyIxNzIuMjAuMS4yOjkyMDAiXSwiZmdyIjoiZTg4NjQyYThkYWQ3NTk4YzE1ZjQyZmE2ZGZjMmYxODQxMDYwZDhmYTQ1ZjE5MGRkY2UxZDlkOGZlM2Q2ZDM1MyIsImtleSI6IlRkbWJXSWdCeFdFcEFjVTdZRzJCOk4wNm9WM1gxUUphbEhGTklqQmxleEEifQ==

2023-05-26 23:09:18 ℹ️  Password for the elastic user (reset with `bin/elasticsearch-reset-password -u elastic`):
2023-05-26 23:09:18   1zoajcY3=FP7uNsrdL*w
2023-05-26 23:09:18 
2023-05-26 23:09:18 ℹ️  HTTP CA certificate SHA-256 fingerprint:
2023-05-26 23:09:18   e88642a8dad7598c15f42fa6dfc2f1841060d8fa45f190ddce1d9d8fe3d6d353
```

再输入验证码:
```sh
docker exec -it kib01 bin/kibana-verification-code

```

以下连接es方式没有成功，可以忽略  
---


配置`http://localhost:5601/`, 配置文件位置: `/usr/share/kibana/config/kibana.yml`  

https://www.elastic.co/guide/en/kibana/current/settings.html  

```yaml
#
# ** THIS IS AN AUTO-GENERATED FILE **
#

# Default Kibana configuration for docker target
# Default Kibana configuration for docker target
server.host: "0.0.0.0"
server.shutdownTimeout: "5s"
elasticsearch.hosts: [ "http://172.20.1.2:9200" ]
elasticsearch.username: "elastic"
elasticsearch.password: "sODJ0yas9I0tDCl3M9xN"
monitoring.ui.container.elasticsearch.enabled: true
```

超级管理员用户不行
```sh
 FATAL  Error: [config validation of [elasticsearch].username]: value of "elastic" is forbidden. This is a superuser account that cannot write to system indices that Kibana needs to function. Use a service account token instead. Learn more: https://www.elastic.co/guide/en/elasticsearch/reference/8.0/service-accounts.html
```

也可以通过界面手动配置
```sh
docker exec -it es01 bin/elasticsearch-create-enrollment-token -s node
eyJ2ZXIiOiI4LjguMCIsImFkciI6WyIxNzIuMjAuMS4yOjkyMDAiXSwiZmdyIjoiZWU3NzM5ZWNmMzE4NTA2NmZiODRjNTQ0M2VmM2YxZjJmZTdlYWVjM2ZlYjIwODg4N2YwMzIzYzlmMGQ5MjFkMyIsImtleSI6ImdSQ0VXSWdCdnk5MDd0UWZqR3ViOnRjbmlkT2t0UnYyMEJWV0Y3VUgxV3cifQ==

# kibana
docker exec -it kib-01 bin/kibana-verification-code
Your verification code is:  451 393
```

---

