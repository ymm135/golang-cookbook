# 大数据统计与分析Demo
## [SZT-bigdata](https://github.com/ymm135/SZT-bigdata)    
该项目主要分析深圳通刷卡数据，通过大数据技术角度来研究深圳地铁客运能力，探索深圳地铁优化服务的方向；  


## 架构图

<br>
<div align=center>
    <img src="../../res/SZT-bigdata-2+.png" width="100%"></img>  
</div>
<br>

```shell
数字标记不分先后顺序，对应代码：
1-cn.java666.sztcommon.util.SZTData
2-cn.java666.etlflink.app.Jsons2Redis
3-cn.java666.etlspringboot.controller.RedisController#get
4-cn.java666.etlflink.app.Redis2ES
5-cn.java666.etlflink.app.Redis2Csv
6-Hive sql 脚本（开发维护成本最低）
7-Saprk 程序（开发维护成本最高，但是功能更强）
8-HUE 方便查询和展示 Hive 数据
9-cn.java666.etlflink.app.Redis2HBase
10、14-cn.java666.szthbase.controller.KafkaListen#sink2Hbase
11-cn.java666.etlflink.app.Redis2HBase
12-CDH HDFS+HUE+Hbase+Hive 一站式查询
13-cn.java666.etlflink.app.Redis2Kafka
15-cn.java666.sztflink.realtime.Kafka2MyCH
16-cn.java666.sztflink.realtime.sink.MyClickhouseSinkFun
```

## 环境搭建
- Java-1.8/Scala-2.11  `brew install scala@2.11`  
- Flink-1.10
- Redis-3.2
- Kafka-2.1
- Zookeeper-3.4.5
- CDH-6.2
- Docker-19
- SpringBoot-2.13
- Elasticsearch-7
- Kibana-7.4
- ClickHouse
- MongoDB-4.0
- Spark-2.3
- Mysql-5.7
- Hadoop3.0  

通过docker搭建开发环境

### Flume [fluːm]()
Flume是Cloudera提供的一个高可用的，高可靠的，分布式的海量日志采集、聚合和传输的系统，Flume支持在日志系统中定制各类数据发送方，用于收集数据；同时，Flume提供对数据进行简单处理，并写到各种数据接受方（可定制）的能力。

### HBase `Hadoop Database`
HBase是一个开源的非关系型分布式数据库（NoSQL），它参考了谷歌的BigTable建模，实现的编程语言为 Java。它是Apache软件基金会的Hadoop项目的一部分，运行于HDFS文件系统之上，为 Hadoop 提供类似于BigTable 规模的服务。因此，它可以对稀疏文件提供极高的容错率。

### HDFS 
The Hadoop Distributed File System (HDFS) is a distributed file system designed to run on commodity hardware.   
Hadoop分布式文件系统(HDFS)是指被设计成适合运行在通用硬件(commodity hardware)上的分布式文件系统（Distributed File System）。它和现有的分布式文件系统有很多共同点。但同时，它和其他的分布式文件系统的区别也是很明显的。HDFS是一个高度容错性的系统，适合部署在廉价的机器上。  

### Hive [haɪv]()
hive是基于Hadoop的一个数据仓库工具，用来进行数据提取、转化、加载，这是一种可以存储、查询和分析存储在Hadoop中的大规模数据的机制。hive数据仓库工具能将结构化的数据文件映射为一张数据库表，并提供SQL查询功能，能将SQL语句转变成MapReduce任务来执行。  

### Hue

### Impala


### Kafka Zookeeper 

[kafka基础中docker安装](kafka-base.md)  

### Oozie

### Spark 
[bitnami/spark](https://hub.docker.com/r/bitnami/spark)  


Apache Spark is a high-performance engine for large-scale computing tasks, such as data processing, machine learning and real-time data streaming. It includes APIs for Java, Python, Scala and R.

spark 是一个大数据处理技术栈，广义的spark包括 spark sql，spark shell，HDFS 和 YARN。  

`docker-compose.yml`  ,执行`docker-compose up`  
```yml
version: '2'

services:
  spark:
    image: docker.io/bitnami/spark:2.4.3-r10  
    environment:
      - SPARK_MODE=master
      - SPARK_RPC_AUTHENTICATION_ENABLED=no
      - SPARK_RPC_ENCRYPTION_ENABLED=no
      - SPARK_LOCAL_STORAGE_ENCRYPTION_ENABLED=no
      - SPARK_SSL_ENABLED=no
    ports:
      - '8080:8080'
  spark-worker:
    image: docker.io/bitnami/spark:2.4.3-r10
    environment:
      - SPARK_MODE=worker
      - SPARK_MASTER_URL=spark://spark:7077
      - SPARK_WORKER_MEMORY=1G
      - SPARK_WORKER_CORES=1
      - SPARK_RPC_AUTHENTICATION_ENABLED=no
      - SPARK_RPC_ENCRYPTION_ENABLED=no
      - SPARK_LOCAL_STORAGE_ENCRYPTION_ENABLED=no
      - SPARK_SSL_ENABLED=no
```


### YARN 

## 配置及运行
IDEA需要安装`Scala`插件，方便开发。  
有些源码的根目录是`scala`,需要标记`SZT-common/src/main/scala`为源码目录，要不然找不到`ParseCardNo.java`类  










