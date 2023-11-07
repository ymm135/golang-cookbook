- # mysql 日常问题整理

- [Mysql 表修复](#mysql-表修复)
  - [如何手工复现`表损坏`](#如何手工复现表损坏)
- [centos7 mysql启动问题](#centos7-mysql启动问题)
- [mysql启动异常](#mysql启动异常)
- [表锁](#表锁)
  - [表锁异常场景](#表锁异常场景)
  - [mysql表锁测试](#mysql表锁测试)


## Mysql 表修复  

```sh
Error 145: Table './firewall/log_base_policy' is marked as crashed and should be repaired
```

检查表的情况:
```sh
mysql> CHECK TABLE log_base_policy;
+--------------------------+-------+----------+---------------------------------------------------------------------------------+
| Table                    | Op    | Msg_type | Msg_text                                                                        |
+--------------------------+-------+----------+---------------------------------------------------------------------------------+
| firewall.log_base_policy | check | warning  | Table is marked as crashed                                                      |
| firewall.log_base_policy | check | warning  | 2 clients are using or haven't closed the table properly                        |
| firewall.log_base_policy | check | error    | Can't read key from filepos: 5120                                               |
| firewall.log_base_policy | check | Error    | Incorrect key file for table './firewall/log_base_policy.MYI'; try to repair it |
| firewall.log_base_policy | check | error    | Corrupt                                                                         |
+--------------------------+-------+----------+---------------------------------------------------------------------------------+
5 rows in set (0.37 sec)
```

从检查的情况来看，MyISAM的索引表损坏了，位置` Can't read key from filepos: 5120 `  

修复表
```sh
REPAIR TABLE log_base_policy;
```

修复的时候，会创建临时表中转
```sh
root@ubuntu:~# ls -lh /var/lib/mysql/firewall/log_base_policy.*
-rw-r----- 1 mysql mysql  18K Mar 11 18:15 /var/lib/mysql/firewall/log_base_policy.frm
-rw-r----- 1 mysql mysql  21G Mar 29 11:45 /var/lib/mysql/firewall/log_base_policy.MYD
-rw-r----- 1 mysql mysql 2.5G Mar 29 15:21 /var/lib/mysql/firewall/log_base_policy.MYI
-rw-r----- 1 mysql mysql 6.1G Mar 29 15:38 /var/lib/mysql/firewall/log_base_policy.TMD
```

修复失败后再次检测:
```sh
mysql> CHECK TABLE log_base_policy;
+--------------------------+-------+----------+----------------------------------------------------------+
| Table                    | Op    | Msg_type | Msg_text                                                 |
+--------------------------+-------+----------+----------------------------------------------------------+
| firewall.log_base_policy | check | warning  | Table is marked as crashed and last repair failed        |
| firewall.log_base_policy | check | warning  | 2 clients are using or haven't closed the table properly |
| firewall.log_base_policy | check | warning  | Size of indexfile is: 2674672640      Should be: 1024    |
| firewall.log_base_policy | check | error    | Record-count is not ok; is 184227914   Should be: 0      |
| firewall.log_base_policy | check | warning  | Found 164198520 deleted space.   Should be 0             |
| firewall.log_base_policy | check | warning  | Found 1368321 deleted blocks       Should be: 0          |
| firewall.log_base_policy | check | warning  | Found 185596235 key parts. Should be: 0                  |
| firewall.log_base_policy | check | error    | Corrupt                                                  |
+--------------------------+-------+----------+----------------------------------------------------------+
8 rows in set (1 min 40.85 sec)
```

> 最后发现问题了，如果在优化时`optimizer`，为了减小表的大小，这时突然断电，那就会导致表无法使用了。  

### 如何手工复现`表损坏`  

```sh
dd if=/dev/random of=/path-to-mysql-data-directory/数据库名/表名.ibd bs=1 count=1024 seek=512
```

修改
```sh
$ dd if=/dev/random of=log_threat_1.frm bs=1 count=1024 seek=512
mysql> CHECK TABLE log_threat_1;
+-----------------------+-------+----------+--------------------------------------------------------------+
| Table                 | Op    | Msg_type | Msg_text                                                     |
+-----------------------+-------+----------+--------------------------------------------------------------+
| firewall.log_threat_1 | check | Error    | Incorrect information in file: './firewall/log_threat_1.frm' |
| firewall.log_threat_1 | check | error    | Corrupt                                                      |
+-----------------------+-------+----------+--------------------------------------------------------------+
2 rows in set (0.01 sec)

# 这种表结构的损坏，不能再使用了，也查不到
select * from information_schema.TABLES where table_schema = 'firewall' AND ENGINE = 'MyISAM'  


$ dd if=/dev/random of=log_base_policy_1.MYD bs=1 count=10240 seek=512
10240+0 records in
10240+0 records out
10240 bytes (10 kB, 10 KiB) copied, 0.0107179 s, 955 kB/s


mysql> CHECK TABLE firewall.log_base_policy_1;
+----------------------------+-------+----------+--------------------------------------------------------+
| Table                      | Op    | Msg_type | Msg_text                                               |
+----------------------------+-------+----------+--------------------------------------------------------+
| firewall.log_base_policy_1 | check | warning  | Table is marked as crashed                             |
| firewall.log_base_policy_1 | check | error    | Size of datafile is: 10752         Should be: 68424328 |
| firewall.log_base_policy_1 | check | error    | Corrupt                                                |
+----------------------------+-------+----------+--------------------------------------------------------+
3 rows in set (0.00 sec)

$ 查询语句
mysql> select * from log_base_policy_1 limit 1;
ERROR 145 (HY000): Table './firewall/log_base_policy_1' is marked as crashed and should be repaired
```

> 如果是innoDB引擎，可以修改`table_name.ibd`. 其中包含索引与数据。myisam索引与数据是分离的。  
> select * from information_schema.TABLES where table_schema = 'firewall' AND ENGINE = 'MyISAM' 在`navicat`中执行看不到损坏的表，但是通过mysql命令远程和本地连接都可以看到。   



## centos7 mysql启动问题

之前启动方式从系统更改为`supervisor`,但是发现启动总会异常，mysql起不来，可能是supersor启动机制的问题。现在准备修改为原始方式:  

```sh
$ systemctl enable mysqld.service
mysqld.service is not a native service, redirecting to /sbin/chkconfig.
Executing /sbin/chkconfig mysqld on

# 需要执行命令
/sbin/chkconfig mysqld on

# 查看
$ chkconfig --list

Note: This output shows SysV services only and does not include native
      systemd services. SysV configuration data might be overridden by native
      systemd configuration.

      If you want to list systemd services use 'systemctl list-unit-files'.
      To see services enabled on particular target use
      'systemctl list-dependencies [target]'.

mysqld          0:off   1:off   2:on    3:on    4:on    5:on    6:off
netconsole      0:off   1:off   2:off   3:off   4:off   5:off   6:off
```


## mysql启动异常  

```sh
2023-09-13T07:44:50.032335Z 0 [Warning] Changed limits: max_open_files: 5000 (requested 10240)
2023-09-13T07:44:50.032759Z 0 [Warning] Changed limits: table_open_cache: 1471 (requested 2000)
2023-09-13T07:44:50.180828Z 0 [Warning] TIMESTAMP with implicit DEFAULT value is deprecated. Please use --explicit_defaults_for_timestamp server option (see documentation for more details).
2023-09-13T07:44:50.182114Z 0 [Note] /usr/sbin/mysqld (mysqld 5.7.39) starting as process 12411 ...
2023-09-13T07:44:50.184853Z 0 [Note] InnoDB: PUNCH HOLE support available
2023-09-13T07:44:50.184878Z 0 [Note] InnoDB: Mutexes and rw_locks use GCC atomic builtins
2023-09-13T07:44:50.184882Z 0 [Note] InnoDB: Uses event mutexes
2023-09-13T07:44:50.184885Z 0 [Note] InnoDB: GCC builtin __atomic_thread_fence() is used for memory barrier
2023-09-13T07:44:50.184889Z 0 [Note] InnoDB: Compressed tables use zlib 1.2.12
2023-09-13T07:44:50.184892Z 0 [Note] InnoDB: Using Linux native AIO
2023-09-13T07:44:50.185661Z 0 [Note] InnoDB: Number of pools: 1
2023-09-13T07:44:50.185754Z 0 [Note] InnoDB: Using CPU crc32 instructions
2023-09-13T07:44:50.187492Z 0 [Note] InnoDB: Initializing buffer pool, total size = 3G, instances = 8, chunk size = 128M
2023-09-13T07:44:50.315890Z 0 [Note] InnoDB: Completed initialization of buffer pool
2023-09-13T07:44:50.326666Z 0 [Note] InnoDB: If the mysqld execution user is authorized, page cleaner thread priority can be changed. See the man page of setpriority().
2023-09-13T07:44:50.337993Z 0 [Note] InnoDB: Highest supported file format is Barracuda.
2023-09-13T07:44:50.344536Z 0 [Note] InnoDB: Log scan progressed past the checkpoint lsn 19949026883
2023-09-13T07:44:50.344568Z 0 [Note] InnoDB: Doing recovery: scanned up to log sequence number 19949026939
2023-09-13T07:44:50.344579Z 0 [Note] InnoDB: Database was not shutdown normally!
2023-09-13T07:44:50.344583Z 0 [Note] InnoDB: Starting crash recovery.
2023-09-13T07:44:50.532988Z 0 [Note] InnoDB: Removed temporary tablespace data file: "ibtmp1"
2023-09-13T07:44:50.533036Z 0 [Note] InnoDB: Creating shared tablespace for temporary tables
2023-09-13T07:44:50.533194Z 0 [Note] InnoDB: Setting file './ibtmp1' size to 12 MB. Physically writing the file full; Please wait ...
2023-09-13T07:44:50.682779Z 0 [Note] InnoDB: File './ibtmp1' size is now 12 MB.
2023-09-13T07:44:50.684750Z 0 [Note] InnoDB: 96 redo rollback segment(s) found. 96 redo rollback segment(s) are active.
2023-09-13T07:44:50.684785Z 0 [Note] InnoDB: 32 non-redo rollback segment(s) are active.
2023-09-13T07:44:50.685615Z 0 [Note] InnoDB: Waiting for purge to start
2023-09-13 15:44:50 0x7f1c3ccc0700  InnoDB: Assertion failure in thread 139759255815936 in file fut0lst.ic line 93
InnoDB: Failing assertion: addr.page == FIL_NULL || addr.boffset >= FIL_PAGE_DATA
InnoDB: We intentionally generate a memory trap.
InnoDB: Submit a detailed bug report to http://bugs.mysql.com.
InnoDB: If you get repeated assertion failures or crashes, even
InnoDB: immediately after the mysqld startup, there may be
InnoDB: corruption in the InnoDB tablespace. Please refer to
InnoDB: http://dev.mysql.com/doc/refman/5.7/en/forcing-innodb-recovery.html
InnoDB: about forcing recovery.
07:44:50 UTC - mysqld got signal 6 ;
This could be because you hit a bug. It is also possible that this binary
or one of the libraries it was linked against is corrupt, improperly built,
or misconfigured. This error can also be caused by malfunctioning hardware.
Attempting to collect some information that could help diagnose the problem.
As this is a crash and something is definitely wrong, the information
collection process might fail.

key_buffer_size=8388608
read_buffer_size=131072
max_used_connections=0
max_threads=2048
thread_count=0
connection_count=0
It is possible that mysqld could use up to
key_buffer_size + (read_buffer_size + sort_buffer_size)*max_threads = 822048 K  bytes of memory
Hope that's ok; if not, decrease some variables in the equation.

Thread pointer: 0x7f1c24000900
Attempting backtrace. You can use the following information to find out
where mysqld died. If you see no messages after this, something went
terribly wrong...
stack_bottom = 7f1c3ccbfdb0 thread_stack 0x40000
/usr/sbin/mysqld(my_print_stacktrace+0x3b)[0xf57adb]
/usr/sbin/mysqld(handle_fatal_signal+0x486)[0x7e7d06]
/lib64/libpthread.so.0(+0xf630)[0x7f1d12f58630]
/lib64/libc.so.6(gsignal+0x37)[0x7f1d11940387]
/lib64/libc.so.6(abort+0x148)[0x7f1d11941a78]
/usr/sbin/mysqld[0x7b7c04]
/usr/sbin/mysqld[0x7b7714]
/usr/sbin/mysqld[0x12bb495]
/usr/sbin/mysqld[0x12bde01]
/usr/sbin/mysqld(_Z9trx_purgemmb+0x3e9)[0x12c10e9]
/usr/sbin/mysqld(srv_purge_coordinator_thread+0xded)[0x12968ad]
/lib64/libpthread.so.0(+0x7ea5)[0x7f1d12f50ea5]
/lib64/libc.so.6(clone+0x6d)[0x7f1d11a08b0d]

Trying to get some variables.
Some pointers may be invalid and cause the dump to abort.
Query (0): Connection ID (thread ID): 0
Status: NOT_KILLED
```

## 表锁  
### 表锁异常场景  
出现异常时会锁表，无法释放
```sh
BEGIN;
LOCK TABLES `firewall`.`user_setting` WRITE;
INSERT INTO `firewall`.`user_setting` (`id`, `setting_code`, `setting_name`, `description`, `args`, `created_at`, `updated_at`, `deleted_at`) 
VALUES (25, 'ct_status', '状态检测', '基于状态检测的防火墙', '0', NULL, NULL, NULL);
INSERT INTO `firewall`.`user_setting` (`id`, `setting_code`, `setting_name`, `description`, `args`, `created_at`, `updated_at`, `deleted_at`) 
VALUES (26, 'ct_status', '状态检测', '基于状态检测的防火墙', '0', NULLd, NULL, NULL);
UNLOCK TABLES
```
> NULLd 故意出错,无法释放锁, 如果增加INSERT IGNORE , 即使重复也不会报错  


msyqdump导出时，提示有问题:  
```sh
mysqldump: [Warning] Using a password on the command line interface can be insecure.
mysqldump: Got error: 1146: Table 'audit.flow_datas_pop3' doesn't exist when using LOCK TABLES
```

该表确实不存你在:  
```sh
mysql> select * from audit.flow_datas_pop3;
ERROR 1146 (42S02): Table 'audit.flow_datas_pop3' doesn't exist
```

但是`show tables`中还存在
```sh
| proto_detail_relation         |
| proto_dnp3                    |
| proto_ethercat                |
| proto_fins                    |
```

数据库的文件:  权限也是正常的  
```sh
-rw-r----- 1 mysql mysql   13188 Jul 25  2022 flow_datas_pop3.frm
-rw-r----- 1 mysql mysql  114688 Jul 25  2022 flow_datas_pop3.ibd
```

也没有看到表损坏的信息
```sh
mysql> select * from information_schema.TABLES where table_name='flow_datas_pop3';
Empty set (0.00 sec)
```

### mysql表锁测试  

MySQL 的 MyISAM 和 InnoDB 是两种常用的存储引擎，它们在表锁的机制上存在显著的差异：

- ### MyISAM:
- **锁定机制**: MyISAM 只支持`表级锁`（`table-level locking`），不支持行级锁。
- **影响**: 当一个用户对 MyISAM 表执行写操作（如 `INSERT`, `UPDATE`, `DELETE`）时，其他用户不能向表中插入新的记录，直到第一个用户完成操作。
- **适用场景**: 由于这种锁定机制，MyISAM 更适合读取频繁但更新不太频繁的应用程序。

- ### InnoDB:
- **锁定机制**: InnoDB 支持`行级锁`（row-level locking）和外键约束。
- **影响**: 这意味着当某个用户正在写某一行时，其他用户仍然可以写其他行。
- **适用场景**: InnoDB 适合于需要频繁更新操作的应用。

- ### 测试用例:

1. **MyISAM 锁定测试**:
   ```sql
   -- 在 Session A:
   CREATE TABLE test_myisam (id INT) ENGINE=MyISAM;
   INSERT INTO test_myisam VALUES (1), (2), (3);
   
   LOCK TABLES table_name WRITE;
   DELETE FROM test_myisam WHERE id=2; -- 这里我们不提交事务

   -- 在 Session B:
   INSERT INTO test_myisam VALUES (4); -- 这里会被阻塞，直到 Session A 提交或回滚事务。

   -- 解锁 
   UNLOCK TABLES;
   ```

2. **InnoDB 锁定测试**:
   ```sql
   -- 在 Session A:
   DROP TABLE IF EXISTS test_innodb;
   CREATE TABLE test_innodb (id INT PRIMARY KEY, value VARCHAR(50)) ENGINE=InnoDB;
   INSERT INTO test_innodb VALUES (1, 'one'), (2, 'two'), (3, 'three');

   
   BEGIN;
   DELETE FROM test_innodb WHERE id=2; -- 这里我们不提交事务
   UPDATE test_innodb SET `value` = 'one_update' WHERE `id` = 1;

   -- 在 Session B:
   INSERT INTO test_innodb VALUES (4,'four'); -- 这里能够成功插入，因为 InnoDB 使用的是行级锁。
   
   UPDATE test_innodb SET `value` = 'one_update' WHERE `id` = 2;  -- 也是没问题的

  
   ```

- ### 结果:
- 在 MyISAM 的测试中，Session B 中的 `INSERT` 会被阻塞，直到 Session A 提交或回滚事务。
- 在 InnoDB 的测试中，Session B 中的 `INSERT` 能够立即执行并插入记录，因为它锁定的只是正在修改的行。

这些例子展示了 MyISAM 和 InnoDB 在并发环境下的行为差异，选择适当的存储引擎非常重要，这取决于应用的具体需求和使用情境。


查看锁表的结果:
```sh
show processlist;
+----+------+--------------------+----------+---------+-------+---------------------------------+------------------------------------+
| Id | User | Host               | db       | Command | Time  | State                           | Info                               |
+----+------+--------------------+----------+---------+-------+---------------------------------+------------------------------------+
|  5 | root | 10.25.17.211:59211 | firewall | Sleep   | 20676 |                                 | NULL                               |
|  6 | root | localhost:41456    | firewall | Sleep   | 20895 |                                 | NULL                               |
|  7 | root | 10.25.17.211:60718 | NULL     | Sleep   | 20558 |                                 | NULL                               |
|  8 | root | 10.25.17.211:60809 | firewall | Sleep   |    11 |                                 | NULL                               |
|  9 | root | localhost          | NULL     | Query   |     0 | starting                        | show processlist                   |
| 10 | root | 10.25.17.211:62953 | firewall | Query   |     5 | Waiting for table metadata lock | INSERT INTO test_myisam VALUES (7) |
| 11 | root | 10.25.17.211:62954 | firewall | Sleep   |    94 |                                 | NULL                               |
| 12 | root | 10.25.17.211:63086 | firewall | Sleep   |    94 |                                 | NULL                               |
+----+------+--------------------+----------+---------+-------+---------------------------------+------------------------------------+
8 rows in set (0.00 sec)
```


```sh
mysql> show processlist;
+----+------+--------------------+----------+---------+------+----------+--------------------------------------------------------------+
| Id | User | Host               | db       | Command | Time | State    | Info                                                         |
+----+------+--------------------+----------+---------+------+----------+--------------------------------------------------------------+
|  2 | root | localhost:43492    | firewall | Sleep   |    1 |          | NULL                                                         |
|  6 | root | 10.25.17.211:63894 | firewall | Sleep   |  173 |          | NULL                                                         |
|  7 | root | 10.25.17.211:63895 | firewall | Query   |    2 | updating | UPDATE test_innodb SET `value` = 'one_update' WHERE `id` = 1 |
|  8 | root | 10.25.17.211:63999 | firewall | Sleep   |  108 |          | NULL                                                         |
|  9 | root | 10.25.17.211:64000 | firewall | Sleep   |  108 |          | NULL                                                         |
| 10 | root | localhost          | NULL     | Query   |    0 | starting | show processlist                                             |
+----+------+--------------------+----------+---------+------+----------+--------------------------------------------------------------+
6 rows in set (0.00 sec)
```

这个示例确实展示了，尽管 InnoDB 使用行级锁，操作同一行记录会等待，时间长了，会超时。
> 1205 - Lock wait timeout exceeded; try restarting transaction, Time: 51.024000s    

## 没有磁盘空间
测试在磁盘没有空间时，mysql存储会有什么问题? 

要快速填充磁盘空间，你可以使用 `dd` 命令。

以下是使用 `dd` 命令快速填充磁盘的方法：

1. 使用 `/dev/zero` 作为输入文件，并使用 `of` 选项指定输出文件：

```bash
dd if=/dev/zero of=/path/to/outputfile bs=1M count=1024
```

这会创建一个 1GB 的文件。你可以增加 `count` 的值以创建更大的文件。

2. 如果你想持续地写入文件，直到磁盘满，你可以这样做：

```bash
dd if=/dev/zero of=/path/to/outputfile bs=1M
```

没有 `count` 选项，`dd` 会继续写入直到发生错误，通常是因为磁盘已满。


mysql 测试程序:
安装依赖
- `go get -u gorm.io/gorm`  
- `go get -u gorm.io/driver/mysql`  

```go
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dsn = "root:123456@tcp(10.25.17.233:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
)

type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:50"`
}

func setupDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Errorf("Open%s", err)
	}

	// Create table
	if err := db.AutoMigrate(&User{}); err != nil {
		fmt.Errorf("Open%s", err)
	}

	return db
}

func TestGormInsertAndSelect(size int) {
	db := setupDB()

	for i := 0; i < size; i++ {
		db.Create(&User{Name: "John" + strconv.Itoa(i)})
	}

}

func main() {
	TestGormInsertAndSelect(100_000)
}
```

通过`dd`命令已经写满磁盘  
```sh
dd if=/dev/zero of=/root/test bs=1M
dd: error writing 'test': No space left on device
85777+0 records in
85776+0 records out
89942695936 bytes (90 GB, 84 GiB) copied, 776.217 s, 116 MB/s
```

mysql错误日志
```sh
2023-10-28T04:12:49.171696Z 271 [ERROR] InnoDB: posix_fallocate(): Failed to preallocate data for file ./testdb/users.ibd, desired size 32768 bytes. Operating system error number 28. Check that the disk is not full or a disk quota exceeded. Make sure the file system supports this function. Some operating system error numbers are described at http://dev.mysql.com/doc/refman/5.7/en/operating-sys
```

sql语句报错
```sh
2023/10/28 12:14:09 main.go:37 Error 1114 (HY000): The table 'users' is full
[14.766ms] [rows:0] INSERT INTO `users` (`name`) VALUES ('John4899')

2023/10/28 12:14:09 main.go:37 Error 1114 (HY000): The table 'users' is full
[14.967ms] [rows:0] INSERT INTO `users` (`name`) VALUES ('John4900')

2023/10/28 12:14:09 main.go:37 Error 1114 (HY000): The table 'users' is full
[14.047ms] [rows:0] INSERT INTO `users` (`name`) VALUES ('John4901')
```



redis也已经无法新增,但是可以读取数据
```sh
 set ddd bbb
(error) MISCONF Redis is configured to save RDB snapshots, but it is currently not able to persist on disk. Commands that may modify the data set are disabled, because this instance is configured to report errors during writes if RDB snapshotting fails (stop-writes-on-bgsave-error option). Please check the Redis logs for details about the RDB error.
```

redis-server日志
```sh
220513:M 28 Oct 2023 12:26:42.012 * 1 changes in 900 seconds. Saving...
220513:M 28 Oct 2023 12:26:42.014 * Background saving started by pid 223374
223374:C 28 Oct 2023 12:26:42.015 # Write error saving DB on disk: No space left on device
220513:M 28 Oct 2023 12:26:42.115 # Background saving error
```

磁盘空间满时，redis ping失败
```sh
127.0.0.1:6379> ping
(error) MISCONF Redis is configured to save RDB snapshots, but it is currently not able to persist on disk. Commands that may modify the data set are disabled, because this instance is configured to report errors during writes if RDB snapshotting fails (stop-writes-on-bgsave-error option). Please check the Redis logs for details about the RDB error.
```
