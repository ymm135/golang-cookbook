- # mysql 双机热备部署  

参考: http://www.br8dba.com/how-to-configure-mysql-master-slave-replication/  

## 环境搭建  
### 简介  
MySQL双机热备（也称为MySQL主从复制）是一种常用的数据库高可用性解决方案，用于实现数据的备份和故障切换。它的基本原理是将一个MySQL数据库（主数据库）的更改操作实时地复制到另一个MySQL数据库（从数据库），从而保持主从之间的数据一致性。当主数据库发生故障时，可以将从数据库提升为新的主数据库，从而实现快速的故障切换。

以下是MySQL双机热备的关键知识点和基本原理：

1. **主从复制架构**：
   - 主数据库（Master）：负责处理客户端的写操作，记录数据的更改，并将更改记录在二进制日志中。
   - 从数据库（Slave）：通过读取主数据库的二进制日志，实时地复制主数据库的更改操作，从而保持与主数据库的数据一致性。

2. **二进制日志**：
   - 主数据库的二进制日志（binary log）记录了所有的更改操作，包括插入、更新和删除操作。
   - 从数据库通过读取主数据库的二进制日志来获取更新操作，然后在自身执行相同的操作，从而保持数据同步。

3. **复制过程**：
   - 从数据库通过`CHANGE MASTER TO`命令指定主数据库的连接信息，并开始复制过程。
   - 从数据库从主数据库的二进制日志中读取并应用更新操作，保持与主数据库的数据一致。

4. **复制模式**：
   - 异步复制：从数据库在接收到主数据库的更新操作后，立即应用，不等待确认。这可以提供更好的性能，但可能会有一定程度的数据延迟。
   - 半同步复制：从数据库在接收到主数据库的更新操作后，需要等待至少一个从数据库确认收到并应用了该操作，然后才会继续。这可以提供更高的数据一致性，但会略微影响性能。
   - 同步复制：从数据库在接收到主数据库的更新操作后，需要等待所有从数据库都确认收到并应用了该操作，然后才会继续。这提供了最高的数据一致性，但会对性能产生较大影响。

5. **主从切换**：
   - 当主数据库发生故障时，可以将从数据库提升为新的主数据库，使应用程序可以继续访问数据。
   - 在切换过程中，需要确保应用程序连接到新的主数据库，并将之前的主数据库设置为新的从数据库。

6. **自动化管理**：
   - 为了更好地管理主从复制，通常会使用工具来自动化配置、监控和故障切换过程，如MySQL的官方工具或第三方解决方案。

总之，MySQL双机热备通过主从复制实现了数据的实时同步，从而提供了数据库的高可用性和容错性。它在应用程序高可用性的架构设计中扮演了重要角色。

### 配置  

测试表及数据:
```sh
-- 创建数据库
CREATE DATABASE my_database CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE my_database;

-- 创建表
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    age INT
);

-- 插入测试数据
INSERT INTO users (name, age) VALUES ('Alice', 25);
INSERT INTO users (name, age) VALUES ('Bob', 30);
INSERT INTO users (name, age) VALUES ('Charlie', 28);
INSERT INTO users (name, age) VALUES ('David', 22);

SELECT * FROM my_database.users;
```

#### 主服务器配置
`my.cnf`  
```sh
[mysqld]
# 双机热备配置  
server-id=10 # 并将数字设置为您想要的任何数字，只要它在所有主奴中都是独一无二的
log_bin=/var/log/mysql/mysql-bin.log
binlog_do_db=my_database # 设置为要复制的数据库
auto-increment-increment=2
auto-increment-offset=1
```

重启服务:`systemctl restart mysql.service`  

查看id
```sh
mysql> SHOW VARIABLES LIKE '%server_id%';
+----------------+-------+
| Variable_name  | Value |
+----------------+-------+
| server_id      | 0     |
| server_id_bits | 32    |
+----------------+-------+
2 rows in set (0.00 sec)
```

配置从服务器同步的用户权限
```sh
GRANT REPLICATION SLAVE ON *.* TO 'slave_user'@'%' IDENTIFIED BY 'password';
FLUSH PRIVILEGES;
```

导出表结构及数据:
```sh
mysqldump -u root -p --opt my_database > my_database.sql
```
> 可以导出数据，也可以使用sql生成  

查看服务器状态:  
```sh
mysql> SHOW MASTER STATUS;
+------------------+----------+--------------+------------------+-------------------+
| File             | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+------------------+----------+--------------+------------------+-------------------+
| mysql-bin.000002 |      589 | my_database  |                  |                   |
+------------------+----------+--------------+------------------+-------------------+
1 row in set (0.00 sec)

更换binlog存储位置后:  
+------------------+----------+--------------+------------------+-------------------+
| File             | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+------------------+----------+--------------+------------------+-------------------+
| mysql-bin.000001 |      154 | my_database  |                  |                   |
+------------------+----------+--------------+------------------+-------------------+
1 row in set (0.00 sec)  
```

#### 从服务器配置

`my.cnf`  
```sh
[mysqld]
server-id = 2 # 并将数字设置为您想要的任何数字，只要它在所有主奴中都是独一无二的
log_bin = /var/log/mysql/mysql-bin.log
binlog_do_db = my_database # 设置为要复制的数据库
relay-log=/var/log/mysql/mysql-relay-bin.log
auto-increment-increment=2
auto-increment-offset=2
```

重启服务:`systemctl restart mysql.service`,查看状态`systemctl status mysql`  

增加主用户权限
```sh
CHANGE MASTER TO MASTER_HOST='master_ip',MASTER_USER='master_user',MASTER_PASSWORD='password',MASTER_LOG_FILE='mysql-bin.000001',MASTER_LOG_POS=154;

START SLAVE;

show SLAVE STATUS\G;

```

## 测试

插入数据:`INSERT INTO my_database.users (name, age) VALUES ('Test', 100);`  

```sh
mysql> INSERT INTO users (name, age) VALUES ('Test', 100);
ERROR 1046 (3D000): No database selected
mysql> INSERT INTO my_database.users (name, age) VALUES ('Test', 100);
Query OK, 1 row affected (0.00 sec)

mysql> SHOW MASTER STATUS;
+------------------+----------+--------------+------------------+-------------------+
| File             | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+------------------+----------+--------------+------------------+-------------------+
| mysql-bin.000001 |      578 | my_database  |                  |                   |
+------------------+----------+--------------+------------------+-------------------+
1 row in set (0.00 sec)

mysql> INSERT INTO my_database.users (name, age) VALUES ('Test1', 101);
Query OK, 1 row affected (0.00 sec)

mysql> SHOW MASTER STATUS;
+------------------+----------+--------------+------------------+-------------------+
| File             | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+------------------+----------+--------------+------------------+-------------------+
| mysql-bin.000001 |      851 | my_database  |                  |                   |
+------------------+----------+--------------+------------------+-------------------+
1 row in set (0.00 sec)
```

查看从服务器的配置:`SELECT * FROM my_database.users;`  


```sh
# 在从设备插入数据
INSERT INTO my_database.users (name, age) VALUES ('SlaveTest1', 111);
| 13 | Test5      |  105 |
| 14 | SlaveTest1 |  111 |
+----+------------+------+


# 在从设备插入数据
NSERT INTO my_database.users (name, age) VALUES ('MasterTest1', 111);;
| 13 | Test5       |  105 |
| 15 | MasterTest1 |  111 |
+----+-------------+------+

# 从同步后的数据
| 13 | Test5       |  105 |
| 14 | SlaveTest1  |  111 |
| 15 | MasterTest1 |  111 |
+----+-------------+------+
13 rows in set (0.00 sec)
```

## 疑问拓展
### redis双机热备  
``
```sh
# Master-Replica replication. Use replicaof to make a Redis instance a copy of
# another Redis server. A few things to understand ASAP about Redis replication.
#
#   +------------------+      +---------------+
#   |      Master      | ---> |    Replica    |
#   | (receive writes) |      |  (exact copy) |
#   +------------------+      +---------------+
#
# 1) Redis replication is asynchronous, but you can configure a master to
#    stop accepting writes if it appears to be not connected with at least
#    a given number of replicas.
# 2) Redis replicas are able to perform a partial resynchronization with the
#    master if the replication link is lost for a relatively small amount of
#    time. You may want to configure the replication backlog size (see the next
#    sections of this file) with a sensible value depending on your needs.
# 3) Replication is automatic and does not need user intervention. After a
#    network partition replicas automatically try to reconnect to masters
#    and resynchronize with them.
slaveof 10.25.10.52 6379
masterauth Netvine123#@!
```

> 有账户的需要配置密码  

主节点状态:
```sh
127.0.0.1:6379> info replication
# Replication
role:master
connected_slaves:0
master_replid:4178616bf87482a0b1384926dddda3a6a8e97a03
master_replid2:0000000000000000000000000000000000000000
master_repl_offset:0
second_repl_offset:-1
repl_backlog_active:0
repl_backlog_size:1048576
repl_backlog_first_byte_offset:0
repl_backlog_histlen:0
```

查看从节点状态
```sh
127.0.0.1:6379> info replication
# Replication
role:slave
master_host:10.25.10.52
master_port:6379
master_link_status:down
master_last_io_seconds_ago:-1
master_sync_in_progress:0
slave_repl_offset:1
master_link_down_since_seconds:1693822755
slave_priority:100
slave_read_only:1
connected_slaves:0
master_replid:c9b11a705a40e1a38dcec0d789e7bb8caae3e334
master_replid2:0000000000000000000000000000000000000000
master_repl_offset:0
second_repl_offset:-1
repl_backlog_active:0
repl_backlog_size:1048576
repl_backlog_first_byte_offset:0
repl_backlog_histlen:0
```

测试
```sh
# 主设备
127.0.0.1:6379> get slave-test
(nil)
127.0.0.1:6379> set slave-test true
OK
127.0.0.1:6379> get slave-test
"true"

# 从设备
"true"
127.0.0.1:6379> get slave-test
"true"
```