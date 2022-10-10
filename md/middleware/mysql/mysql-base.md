# mysql 基础知识  
[参考书籍《深入浅出MySQL：数据库开发、优化与管理维护(第2版)]()  

# 用于日常开发的知识  
## 存储引擎  
MySQL5.0支持的存储引擎包括 `MyISAM`、`InnoDB`、`BDB`、`MEMORY`、`MERGE`、`EXAMPLE`、`NDBCluster`、`ARCHIVE`、`CSV`、`BLACKHOLE`、`FEDERATED`等，其中`InnoDB`和`BDB`提供事务安全表，其他存储引擎都是非事务安全表。  
创建新表时如果不指定存储引擎，那么系统就会使用默认存储引擎，MySQL5.5之前的默认存储引擎是MyISAM，5.5之后改为了`InnoDB`。如果要修改默认的存储引擎，可以在参数文件中设置default-table-type。查看当前的默认存储引擎，可以使用以下命令：  

```shell
mysql> show variables like '%storage_engine%';
+----------------------------------+--------+
| Variable_name                    | Value  |
+----------------------------------+--------+
| default_storage_engine           | InnoDB |
| default_tmp_storage_engine       | InnoDB |
| disabled_storage_engines         |        |
| internal_tmp_disk_storage_engine | InnoDB |
+----------------------------------+--------+
4 rows in set (0.01 sec)
```

查看所有支持的引擎及简介     
```
mysql> show engines;   
+--------------------+---------+----------------------------------------------------------------+--------------+------+------------+
| Engine             | Support | Comment                                                        | Transactions | XA   | Savepoints |
+--------------------+---------+----------------------------------------------------------------+--------------+------+------------+
| InnoDB             | DEFAULT | Supports transactions, row-level locking, and foreign keys     | YES          | YES  | YES        |
| MRG_MYISAM         | YES     | Collection of identical MyISAM tables                          | NO           | NO   | NO         |
| MEMORY             | YES     | Hash based, stored in memory, useful for temporary tables      | NO           | NO   | NO         |
| BLACKHOLE          | YES     | /dev/null storage engine (anything you write to it disappears) | NO           | NO   | NO         |
| MyISAM             | YES     | MyISAM storage engine                                          | NO           | NO   | NO         |
| CSV                | YES     | CSV storage engine                                             | NO           | NO   | NO         |
| ARCHIVE            | YES     | Archive storage engine                                         | NO           | NO   | NO         |
| PERFORMANCE_SCHEMA | YES     | Performance Schema                                             | NO           | NO   | NO         |
| FEDERATED          | NO      | Federated MySQL storage engine                                 | NULL         | NULL | NULL       |
+--------------------+---------+----------------------------------------------------------------+--------------+------+------------+
9 rows in set (0.00 sec)
```

<br>
<div align=center>
    <img src="../../../res/mysql-engine-feature.png" width="80%" height="80%" ></img>  
</div>
<br>

## [InnoDB](https://dev.mysql.com/doc/refman/5.7/en/innodb-storage-engine.html)  
### 整体架构
<br>
<div align=center>
    <img src="../../../res/innodb-architecture.png" width="100%" height="100%" ></img>  
</div>
<br>

### 内存中结构
#### 缓冲池(Buffer Pool) 
缓冲池是主内存中的一个区域，用于在 InnoDB`访问表`和`索引数据`时对其进行缓存。缓冲池允许直接从内存中访问经常使用的数据，从而加快处理速度。在专用服务器上，高达 80% 的物理内存通常分配给缓冲池。

为了提高大容量读取操作的效率，缓冲池被划分为可能包含多行的页面。为了缓存管理的效率，缓冲池被实现为`页链表`；很少使用的数据使用 LRU 算法的变体从缓存中老化。  

了解如何利用缓冲池将频繁访问的数据保存在内存中是 MySQL 调优的一个重要方面。  

```sql
SHOW ENGINE INNODB STATUS

BUFFER POOL AND MEMORY
----------------------
Total large memory allocated 137428992
Dictionary memory allocated 688795
Buffer pool size   8191
Free buffers       1024
Database pages     7077
Old database pages 2592
Modified db pages  0
Pending reads      0
Pending writes: LRU 0, flush list 0, single page 0
Pages made young 2886810, not young 270955789
0.00 youngs/s, 0.00 non-youngs/s
Pages read 23938015, created 176173, written 1046057
0.00 reads/s, 0.00 creates/s, 0.00 writes/s
Buffer pool hit rate 1000 / 1000, young-making rate 13 / 1000 not 0 / 1000
Pages read ahead 0.00/s, evicted without access 0.00/s, Random read ahead 0.00/s
LRU len: 7077, unzip_LRU len: 0
I/O sum[14]:cur[0], unzip sum[0]:cur[0]
```

```sql
mysql> SELECT * FROM information_schema.INNODB_BUFFER_POOL_STATS \G;
*************************** 1. row ***************************
                         POOL_ID: 0
                       POOL_SIZE: 8191
                    FREE_BUFFERS: 1024
                  DATABASE_PAGES: 7077
              OLD_DATABASE_PAGES: 2592
         MODIFIED_DATABASE_PAGES: 0
              PENDING_DECOMPRESS: 0
                   PENDING_READS: 0
               PENDING_FLUSH_LRU: 0
              PENDING_FLUSH_LIST: 0
                PAGES_MADE_YOUNG: 2886808
            PAGES_NOT_MADE_YOUNG: 270955789
           PAGES_MADE_YOUNG_RATE: 0
       PAGES_MADE_NOT_YOUNG_RATE: 0
               NUMBER_PAGES_READ: 23938015
            NUMBER_PAGES_CREATED: 176173
            NUMBER_PAGES_WRITTEN: 1046043
                 PAGES_READ_RATE: 0
               PAGES_CREATE_RATE: 0
              PAGES_WRITTEN_RATE: 0
                NUMBER_PAGES_GET: 424101556
                        HIT_RATE: 0
    YOUNG_MAKE_PER_THOUSAND_GETS: 0
NOT_YOUNG_MAKE_PER_THOUSAND_GETS: 0
         NUMBER_PAGES_READ_AHEAD: 2916964
       NUMBER_READ_AHEAD_EVICTED: 297961
                 READ_AHEAD_RATE: 0
         READ_AHEAD_EVICTED_RATE: 0
                    LRU_IO_TOTAL: 0
                  LRU_IO_CURRENT: 0
                UNCOMPRESS_TOTAL: 0
              UNCOMPRESS_CURRENT: 0
1 row in set (0.00 sec)
```

#### Change Buffer


```sql
mysql> SHOW ENGINE INNODB STATUS\G

-------------------------------------
INSERT BUFFER AND ADAPTIVE HASH INDEX
-------------------------------------
Ibuf: size 1, free list len 9, seg size 11, 93 merges
merged operations:
 insert 535, delete mark 806, delete 0
discarded operations:
 insert 0, delete mark 0, delete 0
Hash table size 34673, node heap has 1 buffer(s)
Hash table size 34673, node heap has 1 buffer(s)
Hash table size 34673, node heap has 2 buffer(s)
Hash table size 34673, node heap has 81 buffer(s)
Hash table size 34673, node heap has 2 buffer(s)
Hash table size 34673, node heap has 1 buffer(s)
Hash table size 34673, node heap has 1 buffer(s)
Hash table size 34673, node heap has 1 buffer(s)
0.00 hash searches/s, 0.00 non-hash searches/s
```

### 硬盘中结构
#### Table 
InnoDB使用 `CREATE TABLE` 语句创建表；例如： 

```
CREATE TABLE t1 (a INT, b CHAR (20), PRIMARY KEY (a)) ENGINE=InnoDB;
```

`.frm` 文件
MySQL 将表的数据字典信息存储 在数据库目录中的.frm 文件中。与其他 MySQL 存储引擎不同， InnoDB它还将有关表的信息编码在系统表空间内自己的内部数据字典中。当 MySQL 删除表或数据库时，它会删除一个或多个.frm文件以及InnoDB数据字典中的相应条目。您不能InnoDB仅通过移动.frm 文件来在数据库之间移动表。

查看`test`数据库的`t1`表的信息
```sql
mysql> SHOW TABLE STATUS FROM test LIKE 't%' \G;
*************************** 1. row ***************************
           Name: t1
         Engine: InnoDB
        Version: 10
     Row_format: Dynamic
           Rows: 0
 Avg_row_length: 0
    Data_length: 16384
Max_data_length: 0
   Index_length: 0
      Data_free: 0
 Auto_increment: NULL
    Create_time: 2022-05-16 15:20:31
    Update_time: NULL
     Check_time: NULL
      Collation: utf8mb4_croatian_ci
       Checksum: NULL
 Create_options: 
        Comment: 
```

```sql
mysql> SELECT * FROM INFORMATION_SCHEMA.INNODB_SYS_TABLES WHERE NAME='test/t1' \G
*************************** 1. row ***************************
     TABLE_ID: 215
         NAME: test/t1
         FLAG: 33
       N_COLS: 5
        SPACE: 257
  FILE_FORMAT: Barracuda
   ROW_FORMAT: Dynamic
ZIP_PAGE_SIZE: 0
   SPACE_TYPE: Single
1 row in set (0.00 sec)
```

#### Index 
每个InnoDB表都有一个特殊的索引，称为`聚集索引`，用于存储行数据。通常，聚集索引与主键同义。为了从查询、插入和其他数据库操作中获得最佳性能，了解如何InnoDB使用聚集索引来优化常见的查找和 DML 操作非常重要。

`PRIMARY KEY`在表上 定义 a时，InnoDB将其用作聚集索引。应该为每个表定义一个主键。如果没有逻辑唯一且非空的列或列集来使用主键，请添加一个自动增量列。自动增量列值是唯一的，并在插入新行时自动添加。

如果您没有`PRIMARY KEY`为表定义 a，InnoDB则使用第一个 `UNIQUE` 索引，其中所有键列都定义为`NOT NULL`聚集索引。

如果表没有索引PRIMARY KEY或没有合适 的UNIQUE索引，则InnoDB 生成一个隐藏的聚集索引 ，该索引以`GEN_CLUST_INDEX`包含行 ID 值的合成列命名。行按InnoDB分配的行 ID 排序。行 ID 是一个 6 字节的字段，随着新行的插入而单调增加。因此，按行 ID 排序的行在物理上是按插入顺序排列的。

- 聚集索引如何加速查询  

通过聚集索引访问行很快，因为索引搜索直接指向包含`行数据的页面`。如果表很大，与使用与索引记录不同的页面存储行数据的存储组织相比，聚集索引架构通常会节省磁盘 I/O 操作。

- 二级索引与聚集索引的关系

聚集索引以外的索引称为`二级索引`。在InnoDB中，二级索引中的每条记录都包含该行的主键列，以及为二级索引指定的列。 InnoDB使用此主键值在聚集索引中搜索行。

如果主键长，二级索引占用的空间就更多，所以主键短是有利的。


- 大小

除空间索引外，InnoDB 索引都是`B-tree`数据结构。空间索引使用 `R-trees`，这是用于索引多维数据的专用数据结构。索引记录存储在其 `B` 树或 `R` 树数据结构的叶页中。索引页的默认大小为 16KB。页面大小由 innodb_page_sizeMySQL 实例初始化时的设置决定。  
  


## 数据类型 

## 字符集  

## 索引的设计及使用  

## 视图  

- ### [视图原理](mysql-viewer.md)  

## 存储过程和函数  

## 触发器  

## 事务控制和锁定语句  


## 安装mysql 

- ### CentOs7 安装 Mysql5.7 

下载mysql源安装包 
```shell
wget http://dev.mysql.com/get/mysql57-community-release-el7-8.noarch.rpm
```

安装mysql源
```shell
yum localinstall mysql57-community-release-el7-8.noarch.rpm
```

检查mysql源是否安装成功
```shell
yum repolist enabled | grep "mysql.*-community.*"
```

返回
```shell
mysql-connectors-community/x86_64 MySQL Connectors Community                 230
mysql-tools-community/x86_64      MySQL Tools Community                      138
mysql57-community/x86_64          MySQL 5.7 Community Server                 564
```

也可以修改 vim /etc/yum.repos.d/mysql-community.repo源，改变默认安装的mysql版本。比如要安装5.6版本，将5.7源的enabled=1改成enabled=0。然后再将5.6源的enabled=0改成enabled=1即可。改完之后的效果如下所示：  
```shell
..........
# Enable to use MySQL 5.5
[mysql55-community]
name=MySQL 5.5 Community Server
baseurl=http://repo.mysql.com/yum/mysql-5.5-community/el/7/$basearch/
enabled=0 # 这里 0表示不选
gpgcheck=1
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-mysql

# Enable to use MySQL 5.6
[mysql56-community]
name=MySQL 5.6 Community Server
baseurl=http://repo.mysql.com/yum/mysql-5.6-community/el/7/$basearch/
enabled=0 # 这里 0表示不选
gpgcheck=1
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-mysql

[mysql57-community]
name=MySQL 5.7 Community Server
baseurl=http://repo.mysql.com/yum/mysql-5.7-community/el/7/$basearch/
enabled=1 # 这里 1 表示 选中
gpgcheck=1 # 如果出现 GPG 失败，可以修改为gpgcheck = 0 
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-mysql
..........
```

安装MySQL
```shell
yum install mysql-community-server
```

启动MySQL服务
```shell
systemctl start mysqld
```

查看MySQL的启动状态
```shell
# systemctl status mysqld
● mysqld.service - MySQL Server
   Loaded: loaded (/usr/lib/systemd/system/mysqld.service; enabled; vendor preset: disabled)
   Active: active (running) since 四 2018-08-23 15:27:28 CST; 1h 26min ago
     Docs: man:mysqld(8)
           http://dev.mysql.com/doc/refman/en/using-systemd.html
  Process: 21453 ExecStart=/usr/sbin/mysqld --daemonize --pid-file=/var/run/mysqld/mysqld.pid $MYSQLD_OPTS (code=exited, status=0/SUCCESS)
  Process: 21432 ExecStartPre=/usr/bin/mysqld_pre_systemd (code=exited, status=0/SUCCESS)
 Main PID: 21457 (mysqld)
   Memory: 202.1M
   CGroup: /system.slice/mysqld.service
           └─21457 /usr/sbin/mysqld --daemonize --pid-file=/var/run/mysqld/mysqld.pid
```

设置开机启动
```
systemctl enable mysqld
systemctl daemon-reload
```

获取root登陆密码
```
mysql安装完成之后，在/var/log/mysqld.log文件中给root生成了一个默认密码。通过下面的方式找到root默认密码，然后登录mysql进行修改：

#  grep 'temporary password' /var/log/mysqld.log
2022-04-13T08:17:37.040012Z 1 [Note] A temporary password is generated for root@localhost: #6we;/NBP3IP
```
ps:如果没有返回，找不到root密码，解决方案：
```
# 1删除原来安装过的mysql残留的数据（这一步非常重要，问题就出在这）
rm -rf /var/lib/mysql

# 2重启mysqld服务
systemctl restart mysqld

# 3再去找临时密码
grep 'temporary password' /var/log/mysqld.log
```
原因有可能是之前安装过一次，没有安装好。  

登陆 
```
[root@VM_18_105_centos ~]# mysql -uroot -p
---- 输入密码：thI/5wEl_chk
# 修改密码
mysql> ALTER USER 'root'@'localhost' IDENTIFIED BY '123456Aa!';   
```
MySql 默认密码级别一定要有大小写字母和特殊符号，所以比较麻烦。  

修改密码策略

在/etc/my.cnf文件添加validate_password_policy配置，指定密码策略  
```
# 0（LOW）：验证 Length
# 1（MEDIUM）：验证 Length; numeric, lowercase/uppercase, and special characters
# 2（STRONG）：验证 Length; numeric, lowercase/uppercase, and special characters; dictionary file
validate_password_policy=0
```

当然如果不需要密码策略，可以禁用：
在/etc/my.cnf文件添加
```
validate_password = off
```
重启生效：
```
systemctl restart mysqld
```

Mysql的root用户，只能本地访问，这里在创建一个远程可以访问的 用户。
```
GRANT ALL PRIVILEGES ON *.* TO 'user'@'%' IDENTIFIED BY '123456' WITH GRANT OPTION;
```

查看用户及权限: 
```
mysql> select User,authentication_string,Host from mysql.user;
+---------------+-------------------------------------------+-----------+
| User          | authentication_string                     | Host      |
+---------------+-------------------------------------------+-----------+
| root          | *5860629DEA0B537A9258E5B884D9A660887098D6 | localhost |
| mysql.session | *THISISNOTAVALIDPASSWORDTHATCANBEUSEDHERE | localhost |
| mysql.sys     | *THISISNOTAVALIDPASSWORDTHATCANBEUSEDHERE | localhost |
| root          | *81F5E21E35407D884A6CD4A731AEBFB6AF209E1B | %         |
+---------------+-------------------------------------------+-----------+
```






