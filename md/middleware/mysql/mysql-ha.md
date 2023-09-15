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
- [环境搭建](#环境搭建)
  - [简介](#简介)
  - [配置](#配置)
    - [主服务器配置](#主服务器配置)
    - [从服务器配置](#从服务器配置)
- [测试](#测试)
- [疑问拓展](#疑问拓展)
  - [redis双机热备](#redis双机热备)


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

查看redis状态:
```sh
127.0.0.1:6379> info
# Server
redis_version:5.0.7
redis_git_sha1:00000000
redis_git_dirty:0
redis_build_id:66bd629f924ac924
redis_mode:standalone
os:Linux 5.4.0-155-generic x86_64
arch_bits:64
multiplexing_api:epoll
atomicvar_api:atomic-builtin
gcc_version:9.3.0
process_id:32138
run_id:0b07d3e6576b260648772ab132e5fc735f8ec601
tcp_port:6379
uptime_in_seconds:654
uptime_in_days:0
hz:10
configured_hz:10
lru_clock:16303477
executable:/usr/bin/redis-server
config_file:/etc/redis/redis.conf

# Clients
connected_clients:2
client_recent_max_input_buffer:2
client_recent_max_output_buffer:0
blocked_clients:0

# Memory
used_memory:1090168
used_memory_human:1.04M
used_memory_rss:7233536
used_memory_rss_human:6.90M
used_memory_peak:1090168
used_memory_peak_human:1.04M
used_memory_peak_perc:100.09%
used_memory_overhead:863624
used_memory_startup:796264
used_memory_dataset:226544
used_memory_dataset_perc:77.08%
allocator_allocated:1247104
allocator_active:1585152
allocator_resident:4288512
total_system_memory:10431987712
total_system_memory_human:9.72G
used_memory_lua:41984
used_memory_lua_human:41.00K
used_memory_scripts:0
used_memory_scripts_human:0B
number_of_cached_scripts:0
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
allocator_frag_ratio:1.27
allocator_frag_bytes:338048
allocator_rss_ratio:2.71
allocator_rss_bytes:2703360
rss_overhead_ratio:1.69
rss_overhead_bytes:2945024
mem_fragmentation_ratio:6.90
mem_fragmentation_bytes:6185352
mem_not_counted_for_evict:0
mem_replication_backlog:0
mem_clients_slaves:0
mem_clients_normal:66616
mem_aof_buffer:0
mem_allocator:jemalloc-5.2.1
active_defrag_running:0
lazyfree_pending_objects:0

# Persistence
loading:0
rdb_changes_since_last_save:0
rdb_bgsave_in_progress:0
rdb_last_save_time:1694024423
rdb_last_bgsave_status:ok
rdb_last_bgsave_time_sec:-1
rdb_current_bgsave_time_sec:-1
rdb_last_cow_size:0
aof_enabled:0
aof_rewrite_in_progress:0
aof_rewrite_scheduled:0
aof_last_rewrite_time_sec:-1
aof_current_rewrite_time_sec:-1
aof_last_bgrewrite_status:ok
aof_last_write_status:ok
aof_last_cow_size:0

# Stats
total_connections_received:439
total_commands_processed:471
instantaneous_ops_per_sec:0
total_net_input_bytes:348011
total_net_output_bytes:30530
instantaneous_input_kbps:0.00
instantaneous_output_kbps:0.00
rejected_connections:0
sync_full:0
sync_partial_ok:0
sync_partial_err:0
expired_keys:0
expired_stale_perc:0.00
expired_time_cap_reached_count:0
evicted_keys:0
keyspace_hits:5
keyspace_misses:21
pubsub_channels:0
pubsub_patterns:0
latest_fork_usec:0
migrate_cached_sockets:0
slave_expires_tracked_keys:0
active_defrag_hits:0
active_defrag_misses:0
active_defrag_key_hits:0
active_defrag_key_misses:0

# Replication
role:slave
master_host:10.25.10.52
master_port:6739
master_link_status:down
master_last_io_seconds_ago:-1
master_sync_in_progress:0
slave_repl_offset:1
master_link_down_since_seconds:1694025077
slave_priority:100
slave_read_only:1
connected_slaves:0
master_replid:c2d88dc49278b196d987ef91b3d3aa2a2739e41a
master_replid2:0000000000000000000000000000000000000000
master_repl_offset:0
second_repl_offset:-1
repl_backlog_active:0
repl_backlog_size:1048576
repl_backlog_first_byte_offset:0
repl_backlog_histlen:0

# CPU
used_cpu_sys:0.459283
used_cpu_user:0.485538
used_cpu_sys_children:0.000000
used_cpu_user_children:0.000000

# Cluster
cluster_enabled:0

# Keyspace
db0:keys=14,expires=1,avg_ttl=0
```

主从相关的配置
```sh
role:slave
master_host:10.25.10.52
master_port:6739
master_link_status:down
master_last_io_seconds_ago:-1
master_sync_in_progress:0
slave_repl_offset:1
master_link_down_since_seconds:1694025077
slave_priority:100
slave_read_only:1
```

在Redis 6.0之前的版本中，从服务器（slave）是始终只读的，无法通过内置选项来配置从服务器为读写模式。这是因为Redis早期的设计中，从服务器的主要目的是用于数据备份和读取负载均衡，而不是用于写入操作。  

### Duplicate entry '17' for key 'PRIMARY',
Could not execute Write_rows event on table firewall.security_policy_manager; Duplicate entry '17' for key 'PRIMARY', Error_code: 1062; handler error HA_ERR_FOUND_DUPP_KEY; the event's master log mysql-bin.000004, end_log_pos 20953  

可以不用增加binlog的同步位置:  
```sql
CHANGE MASTER TO MASTER_HOST='192.168.100.1',MASTER_USER='root',MASTER_PASSWORD='pass';
```

> 如果删除没有同步之前的数据，从数据库会找不到该id，也会报错。所以要先确保数据是一致的。  

### master and slave have equal MySQL server ids  
 Fatal error: The slave I/O thread stops because master and slave have equal MySQL server ids; these ids must be different for replication to work (or the --replicate-same-server-id option must be used on slave but this does not always make sense; please check the manual before using it).

首先已经确保配置文件总的id不相同了，sql语句查询到也是不同的，但是仍然有该问题？  
```sql
CHANGE MASTER TO MASTER_HOST='192.168.100.1',MASTER_USER='root',MASTER_PASSWORD='pass';
```

配置文件出错了，配置的`MASTER_HOST`是自己，所以一直报错.  






