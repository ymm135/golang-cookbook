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

```
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




