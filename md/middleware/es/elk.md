- # ELK

- [结构](#结构)
  - [官网](#官网)
- [安装logstash](#安装logstash)
  - [linux 安装](#linux-安装)
  - [docker 安装](#docker-安装)
- [logstash 指南](#logstash-指南)
  - [入门](#入门)
  - [原理](#原理)
  - [JDBC Input Plugin](#jdbc-input-plugin)
- [mysql同步工具canal](#mysql同步工具canal)
- [环境搭建](#环境搭建)
  - [filebeat](#filebeat)
  - [logstash](#logstash)
    - [输出到pg](#输出到pg)
  - [elasticsearch](#elasticsearch)
  - [kibana](#kibana)
- [以下连接es方式没有成功，可以忽略](#以下连接es方式没有成功可以忽略)


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

logstash规则
```sh
input {
  beats {
    port => 5044
  }
}

output {
  elasticsearch {
    hosts => ["http://localhost:9200"]
    index => "%{[@metadata][beat]}-%{[@metadata][version]}" 
    action => "create"
  }
}
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
path.config: /usr/share/logstash/config/conf.d/*.conf
path.logs: /usr/share/logstash/logs
```

```sh
# 启动
docker run -itd --net elastic --ip 172.20.1.5 \
  -p 5044:5044 -p 5000:5000/udp \
  --name logstash \
  docker.elastic.co/logstash/logstash:8.8.0
```

另外也可以使用自定义参数:
```sh
docker run -itd --net elastic --ip 172.20.1.5 \
  -p 5044:5044 -p 5000:5000 -p 5140:5140/udp \
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

查看`/usr/share/logstash/config/conf.d/logstash.conf`和自定义配置文件  
> 默认是`/usr/share/logstash/pipeline/logstash.conf`, 如果你配置了conf.d，pipeline中的conf就不生效了  
```sh
input {
  tcp {
        port => 5000
        mode => "server"
        ssl_enable => false
  }
}

output {
  stdout {
    codec => rubydebug
  }
}
```

启动日志:
```sh
2023-05-27 21:37:19 [2023-05-27T13:37:19,627][INFO ][logstash.javapipeline    ][main] Starting pipeline {:pipeline_id=>"main", "pipeline.workers"=>3, "pipeline.batch.size"=>125, "pipeline.batch.delay"=>50, "pipeline.max_inflight"=>375, "pipeline.sources"=>["/usr/share/logstash/config/conf.d/logstash.conf"], :thread=>"#<Thread:0x4d0d1e3b@/usr/share/logstash/logstash-core/lib/logstash/java_pipeline.rb:134 run>"}
2023-05-27 21:37:19 [2023-05-27T13:37:19,996][INFO ][logstash.javapipeline    ][main] Pipeline Java execution initialization time {"seconds"=>0.37}
2023-05-27 21:37:20 [2023-05-27T13:37:20,069][INFO ][logstash.javapipeline    ][main] Pipeline started {"pipeline.id"=>"main"}
2023-05-27 21:37:20 [2023-05-27T13:37:20,074][INFO ][logstash.inputs.tcp      ][main][5d3b15e941d2a5fd9052e87c49105661524bbbf5a664607aa8139aacaf6277b3] Starting tcp input listener {:address=>"0.0.0.0:5000", :ssl_enable=>false}
2023-05-27 21:37:20 [2023-05-27T13:37:20,081][INFO ][logstash.agent           ] Pipelines running {:count=>1, :running_pipelines=>[:main], :non_running_pipelines=>[]}
```

> 独立运行:`docker exec logstash bin/logstash -f /usr/share/logstash/pipeline/logstash.conf`   

使用宿主机测试:
```sh
echo "Hello Log" > olddata
nc 127.0.0.1 5000 < olddata
```

控制台输出的日志
```sh
2023-05-27 21:39:02 {
2023-05-27 21:39:02     "@timestamp" => 2023-05-27T13:39:02.257847877Z,
2023-05-27 21:39:02          "event" => {
2023-05-27 21:39:02         "original" => "Hello Log"
2023-05-27 21:39:02     },
2023-05-27 21:39:02       "@version" => "1",
2023-05-27 21:39:02        "message" => "Hello Log"
2023-05-27 21:39:02 }
```

syslog输入模块配置测试,新建一个配置文件`/usr/share/logstash/config/conf.d/syslog.conf`  
```sh
input {
  syslog {
    port => 5140
    codec => cef
    syslog_field => "syslog"
    grok_pattern => "<%{POSINT:priority}>%{SYSLOGTIMESTAMP:timestamp} CUSTOM GROK HERE"
  }
}

output {
  stdout {
    codec => rubydebug
  }
}
```

输出格式:
```sh
2023-05-27 22:16:48 {
2023-05-27 22:16:48            "log" => {
2023-05-27 22:16:48         "syslog" => {
2023-05-27 22:16:48             "facility" => {
2023-05-27 22:16:48                 "name" => "kernel",
2023-05-27 22:16:48                 "code" => 0
2023-05-27 22:16:48             },
2023-05-27 22:16:48             "priority" => 0,
2023-05-27 22:16:48             "severity" => {
2023-05-27 22:16:48                 "name" => "Emergency",
2023-05-27 22:16:48                 "code" => 0
2023-05-27 22:16:48             }
2023-05-27 22:16:48         }
2023-05-27 22:16:48     },
2023-05-27 22:16:48        "message" => "<34>Oct 11 22:14:15 mymachine myproc[10]: 'su root' failed for user",
2023-05-27 22:16:48           "tags" => [
2023-05-27 22:16:48         [0] "_cefparsefailure",
2023-05-27 22:16:48         [1] "_grokparsefailure_sysloginput"
2023-05-27 22:16:48     ],
2023-05-27 22:16:48           "host" => {
2023-05-27 22:16:48         "ip" => "172.20.1.1"
2023-05-27 22:16:48     },
2023-05-27 22:16:48        "service" => {
2023-05-27 22:16:48         "type" => "system"
2023-05-27 22:16:48     },
2023-05-27 22:16:48       "@version" => "1",
2023-05-27 22:16:48     "@timestamp" => 2023-05-27T14:16:48.630375811Z,
2023-05-27 22:16:48          "event" => {
2023-05-27 22:16:48         "original" => nil
2023-05-27 22:16:48     }
2023-05-27 22:16:48 }
```

或者
```sh
input {
  tcp {
    port => 5140
    type => syslog
  }

  udp {
    port => 5140
    type => syslog
  }
}

output {
  stdout {
    codec => rubydebug
  }
}
```

输出日志:
```sh
2023-05-27 22:14:26 {
2023-05-27 22:14:26           "type" => "syslog",
2023-05-27 22:14:26        "message" => "<34>Oct 11 22:14:15 mymachine myproc[10]: 'su root' failed for user",
2023-05-27 22:14:26       "@version" => "1",
2023-05-27 22:14:26     "@timestamp" => 2023-05-27T14:14:26.629098863Z,
2023-05-27 22:14:26          "event" => {
2023-05-27 22:14:26         "original" => "<34>Oct 11 22:14:15 mymachine myproc[10]: 'su root' failed for user"
2023-05-27 22:14:26     }
2023-05-27 22:14:26 }
```

启动日志
```sh
2023-05-27 21:50:37 [2023-05-27T13:50:37,161][INFO ][logstash.inputs.syslog   ][main][1d272c67008a0c809278b494f2080e76cb40a66167f9a0dc054c91c2ed734946] Starting syslog udp listener {:address=>"0.0.0.0:5140"}
2023-05-27 21:50:37 [2023-05-27T13:50:37,165][INFO ][logstash.inputs.syslog   ][main][1d272c67008a0c809278b494f2080e76cb40a66167f9a0dc054c91c2ed734946] Starting syslog tcp listener {:address=>"0.0.0.0:5140"}
```

#### 过滤  
https://www.elastic.co/guide/en/logstash/current/plugins-filters-grok.html  

发送的数据
```sh
echo '"0:09:10","wwm","/WAN2交换机组/","100.100.54.199","113.96.202.102","邮件","QQ邮箱[浏览]","未定义位置","/PC","记录","59693268","wp_zie_file.zip","128","压缩文件"' | nc -w1 -u 127.0.0.1 5140

echo '"0:00:01","zhangxue","/部门集群用户组/","10.22.58.155","115.236.118.54","访问网站","个人网站及博客","未定义位置","/未知类型","记录","api.money.126.net","api.money.126.net","示例标题","-"' | nc -w1 -u 127.0.0.1 5140

端口: 5140
数据库: t_behavior_log_url, 序列: nextval('t_behavior_log_url_seq'::regclass)
```

过滤规则:  
```sh
input {
  udp {
    port => 5140
  }
}

filter {
  grok {
    match => { "message" => "%{TIME:record_time},%{WORD:user},%{DATA:group},%{IP:host_ip},%{IP:dst_ip},%{DATA:serv},%{DATA:app},%{DATA:site},%{DATA:tm_type},%{DATA:net_action},%{DATA:url},%{DATA:DNS},%{DATA:title},%{DATA:snapshot}" }
  }
}
```

> 主要是字符串分割  


#### 输出到pg

docker安装pg
```sh
docker run -itd -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres --net elastic --ip 172.20.1.6 -p 5432:5432 --name postgresql postgres
```
>也可以添加自己的数据: -v /data:/var/lib/postgresql/data  

使用[pgadmin](https://www.pgadmin.org/download/)连接  

创建数据库和表
```sh
postgres=# CREATE DATABASE sysloddb;

CREATE TABLE t_behavior_log(
   ID INT PRIMARY KEY     NOT NULL,
   TIME           TEXT    NOT NULL,
   CONTENT        CHAR(500),
);

使用\d查看是否创建成功
```

最终发现还是使用navicat连接使用比较习惯  
```sh
INSERT INTO "public"."t_behavior_log" ("id", "time", "user", "group", "source_ip", "destination_ip", "action", "client", "location", "device", "log", "size1", "size2", "size_name") VALUES (1, '', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
```

安装jdbc模块:
```sh
# 容器内执行
logstash-plugin install --no-verify logstash-output-jdbc 

# 容器内创建
mkdir -p /usr/share/logstash/vendor/jar/jdbc

# 宿主机拷贝到容器
docker cp lpostgresql-42.5.1.jar logstash:/usr/share/logstash/vendor/jar/jdbc
```

输出模块配置
```sh
output {
  stdout {
    codec => rubydebug
  }

  jdbc {
    connection_string => "jdbc:postgresql://10.25.1.4:5432/xgxx_log?user=postgres&password=password"
    statement => ['INSERT INTO t_behavior_log_url (record_time, "user", "group", host_ip, dst_ip, serv, app, site, tm_type, net_action, url, dns, title, snapshot) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)',"record_time", "user", "group", "host_ip", "dst_ip", "serv", "app", "site", "tm_type", "net_action", "url", "DNS", "title", "snapshot"]
  }
}
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

