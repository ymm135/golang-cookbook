- # mysql8 性能优化

- [ubuntu基础环境搭建](#ubuntu基础环境搭建)
  - [本地安装](#本地安装)
  - [docker](#docker)
    - [数据准备](#数据准备)
- [慢查询发现与分析](#慢查询发现与分析)
- [数据库调优原理](#数据库调优原理)
  - [B-Tree与B+Tree](#b-tree与btree)
  - [MyISAM与InnoDB](#myisam与innodb)
- [特定语句的原理与调优](#特定语句的原理与调优)


## ubuntu基础环境搭建
### 本地安装
下载apt源

`wget https://dev.mysql.com/get/mysql-apt-config_0.8.24-1_all.deb`  

```sh
sudo dpkg -i mysql-apt-config_0.8.24-1_all.deb
```

> 直接选择OK  

```sh
sudo apt update
sudo apt-get install mysql-server
```

安装过程中会要求输入`root`密码, 这里输入`root`  
另外选择密码加密方式，选择8.0的加密方式，而不是5.x  

`mysql -uroot -p`, 输入密码进入说明已经安装成功了。  


查看版本
```sh
mysql> select version();
+-----------+
| version() |
+-----------+
| 8.0.31    |
+-----------+
1 row in set (0.00 sec)
```

root远程连接
```sh
# 查看用户权限
mysql> select User,authentication_string,Host from mysql.user;
+------------------+------------------------------------------------------------------------+-----------+
| User             | authentication_string                                                  | Host      |
+------------------+------------------------------------------------------------------------+-----------+
| mysql.infoschema | $A$005$THISISACOMBINATIONOFINVALIDSALTANDPASSWORDTHATMUSTNEVERBRBEUSED | localhost |
| mysql.session    | $A$005$THISISACOMBINATIONOFINVALIDSALTANDPASSWORDTHATMUSTNEVERBRBEUSED | localhost |
| mysql.sys        | $A$005$THISISACOMBINATIONOFINVALIDSALTANDPASSWORDTHATMUSTNEVERBRBEUSED | localhost |
| root             | $A$005$rZ/PN\<N,6,IPG	&#O6CytVBkox0f0a1O/txLEEmCn6UuOGuj3U1mo/zu6P9 | localhost |
+------------------+------------------------------------------------------------------------+-----------+
4 rows in set (0.00 sec)

# 修改用户密码
alter user 'root'@'localhost' identified with mysql_native_password by 'root';

# 创建用户
create user 'root'@'%' identified by 'root'; 

# 授权
grant all privileges on *.* to 'root'@'%' with grant option;
flush privileges; 
```


身份验证插件
```sh
mysql> select user,host,plugin from mysql.user;    
+------------------+-----------+-----------------------+
| user             | host      | plugin                |
+------------------+-----------+-----------------------+
| root             | %         | caching_sha2_password |
| mysql.infoschema | localhost | caching_sha2_password |
| mysql.session    | localhost | caching_sha2_password |
| mysql.sys        | localhost | caching_sha2_password |
| root             | localhost | mysql_native_password |
+------------------+-----------+-----------------------+
5 rows in set (0.00 sec)
```

### docker  

#### 数据准备  

测试数据:  https://github.com/ymm135/mysql-test-db    

```sql
CREATE TABLE employees (
    emp_no      INT             NOT NULL,
    birth_date  DATE            NOT NULL,
    first_name  VARCHAR(14)     NOT NULL,
    last_name   VARCHAR(16)     NOT NULL,
    gender      ENUM ('M','F')  NOT NULL,    
    hire_date   DATE            NOT NULL,
    PRIMARY KEY (emp_no)
);

CREATE TABLE departments (
    dept_no     CHAR(4)         NOT NULL,
    dept_name   VARCHAR(40)     NOT NULL,
    PRIMARY KEY (dept_no),
    UNIQUE  KEY (dept_name)
);

CREATE TABLE dept_manager (
   emp_no       INT             NOT NULL,
   dept_no      CHAR(4)         NOT NULL,
   from_date    DATE            NOT NULL,
   to_date      DATE            NOT NULL,
   FOREIGN KEY (emp_no)  REFERENCES employees (emp_no)    ON DELETE CASCADE,
   FOREIGN KEY (dept_no) REFERENCES departments (dept_no) ON DELETE CASCADE,
   PRIMARY KEY (emp_no,dept_no)
); 

CREATE TABLE dept_emp (
    emp_no      INT             NOT NULL,
    dept_no     CHAR(4)         NOT NULL,
    from_date   DATE            NOT NULL,
    to_date     DATE            NOT NULL,
    FOREIGN KEY (emp_no)  REFERENCES employees   (emp_no)  ON DELETE CASCADE,
    FOREIGN KEY (dept_no) REFERENCES departments (dept_no) ON DELETE CASCADE,
    PRIMARY KEY (emp_no,dept_no)
);

CREATE TABLE titles (
    emp_no      INT             NOT NULL,
    title       VARCHAR(50)     NOT NULL,
    from_date   DATE            NOT NULL,
    to_date     DATE,
    FOREIGN KEY (emp_no) REFERENCES employees (emp_no) ON DELETE CASCADE,
    PRIMARY KEY (emp_no,title, from_date)
);

CREATE TABLE salaries (
    emp_no      INT             NOT NULL,
    salary      INT             NOT NULL,
    from_date   DATE            NOT NULL,
    to_date     DATE            NOT NULL,
    FOREIGN KEY (emp_no) REFERENCES employees (emp_no) ON DELETE CASCADE,
    PRIMARY KEY (emp_no, from_date)
);
```

关系图

<br>
<div align=center>
    <img src="../../../res/test-mysql-tables.png" width="80%"></img>  
</div>
<br>

## 慢查询发现与分析

查看慢查询状态
```sql
mysql> show variables like 'slow_query_log'; 
+----------------+-------+
| Variable_name  | Value |
+----------------+-------+
| slow_query_log | OFF   |
+----------------+-------+
1 row in set (0.25 sec)
```

临时开启

```sql
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL slow_query_log_file = '/var/log/mysql/mysql-slow.log';
SET GLOBAL log_queries_not_using_indexes = 'ON';
SET SESSION long_query_time = 1;
SET SESSION min_examined_row_limit = 100;
```

结果查询
```sql
mysql> show variables like '%quer%';
+----------------------------------------+-------------------------------+
| Variable_name                          | Value                         |
+----------------------------------------+-------------------------------+
| binlog_rows_query_log_events           | OFF                           |
| ft_query_expansion_limit               | 20                            |
| have_query_cache                       | NO                            |
| log_queries_not_using_indexes          | ON                            |
| log_throttle_queries_not_using_indexes | 0                             |
| long_query_time                        | 2.000000                      |
| query_alloc_block_size                 | 8192                          |
| query_prealloc_size                    | 8192                          |
| slow_query_log                         | ON                            |
| slow_query_log_file                    | /var/log/mysql/mysql-slow.log |
+----------------------------------------+-------------------------------+
```

永久开启


## 数据库调优原理
### explain与耗时分析
打开耗时统计
```sql
set profiling=1;
```

查看结果
```sql
mysql> select @@profiling;
+-------------+
| @@profiling |
+-------------+
|           1 |
+-------------+
```

查看sql耗时
```sql
mysql> show profiles; 
+----------+------------+------------------------------+
| Query_ID | Duration   | Query                        |
+----------+------------+------------------------------+
|        1 | 0.27903625 | select * from employees      |
|        2 | 0.10113750 | select emp_no from employees |
+----------+------------+------------------------------+
```

explain分析
```sql
mysql> explain select * from employees ;
+----+-------------+-----------+------------+------+---------------+------+---------+------+--------+----------+-------+
| id | select_type | table     | partitions | type | possible_keys | key  | key_len | ref  | rows   | filtered | Extra |
+----+-------------+-----------+------------+------+---------------+------+---------+------+--------+----------+-------+
|  1 | SIMPLE      | employees | NULL       | ALL  | NULL          | NULL | NULL    | NULL | 292025 |   100.00 | NULL  |
+----+-------------+-----------+------------+------+---------------+------+---------+------+--------+----------+-------+
1 row in set, 1 warning (0.01 sec)

mysql> explain select emp_no from employees;
+----+-------------+-----------+------------+-------+---------------+---------+---------+------+--------+----------+-------------+
| id | select_type | table     | partitions | type  | possible_keys | key     | key_len | ref  | rows   | filtered | Extra       |
+----+-------------+-----------+------------+-------+---------------+---------+---------+------+--------+----------+-------------+
|  1 | SIMPLE      | employees | NULL       | index | NULL          | PRIMARY | 4       | NULL | 292025 |   100.00 | Using index |
+----+-------------+-----------+------------+-------+---------------+---------+---------+------+--------+----------+-------------+
1 row in set, 1 warning (0.00 sec)
```

| 字段 | 含义  | 说明 | 
| ---- | ---- | ---- |
| id | select查询的序列号，包含一组数字，表示查询中执行select子句或操作表的顺序 |  | 
| select_type | 查询类型 或者是 其他操作类型 | `SIMPLE`简单查询,`PRIMARY`主查询,`SUBQUERY`子查询，`UNION`连接查询 |
| table | 正在访问哪个表 |  | 
| partitions | 匹配的分区 | | 
| type| 访问的类型 | NULL > system > const > eq_ref > ref > ref_or_null > index_merge > range > index > ALL | 
| possible_keys| 显示可能应用在这张表中的索引，一个或多个，但不一定实际使用到 | | 
| key | 实际使用到的索引，如果为NULL，则没有使用索引 | 
| key_len | 表示索引中使用的字节数，可通过该列计算查询中使用的索引的长度 | | 
| ref | 显示索引的哪一列被使用了，如果可能的话，是一个常数，哪些列或常量被用于查找索引列上的值 | | 
| rows | 根据表统计信息及索引选用情况，大致估算出找到所需的记录所需读取的行数
| filtered | 查询的表行占表的百分比 | | 
| Extra | 包含不适合在其它列中显示但十分重要的额外信息| | 


### B-Tree与B+Tree
一般说MySQL的索引，都清楚其索引主要以B+树为主，此外还有Hash、RTree、FullText。  

mysql innodb b+树索引  


<br>
<div align=center>
    <img src="../../../res/B+Tree-Structure.png" width="100%"></img>  
</div>
<br>

> 传统数据库使用 B+树方式存储.B+树 ==> 排序树（索引树）==> 每个节点存储文件块。 每个`磁盘块`存储多个数据。多个`磁盘块`构成 B+树。连续读写能力强,随机读写能力弱。  

### MyISAM与InnoDB

## 特定语句的原理与调优
### JOIN
```sql
select emp.emp_no,emp.first_name,emp.last_name 
from employees as emp
left join dept_emp on emp.emp_no = dept_emp.emp_no
where dept_emp.dept_no='d005'
limit 10,10

+--------+------------+-----------+
| emp_no | first_name | last_name |
+--------+------------+-----------+
|  10027 | Divier     | Reistad   |
|  10028 | Domenick   | Tempesti  |
|  10031 | Karsten    | Joslin    |
|  10037 | Pradeep    | Makrucki  |
|  10040 | Weiyi      | Meriste   |
|  10043 | Yishay     | Tzvieli   |
|  10048 | Florian    | Syrotiuk  |
|  10056 | Brendon    | Bernini   |
|  10057 | Ebbe       | Callaway  |
|  10062 | Anoosh     | Peyn      |
+--------+------------+-----------+
```

explain
```sql
mysql> explain select emp.emp_no,emp.first_name,emp.last_name 
    -> from employees as emp
    -> left join dept_emp on emp.emp_no = dept_emp.emp_no
    -> where dept_emp.dept_no='d005'
    -> limit 10,10;
+----+-------------+----------+------------+--------+-----------------+---------+---------+----------------------+--------+----------+--------------------------+
| id | select_type | table    | partitions | type   | possible_keys   | key     | key_len | ref                  | rows   | filtered | Extra                    |
+----+-------------+----------+------------+--------+-----------------+---------+---------+----------------------+--------+----------+--------------------------+
|  1 | SIMPLE      | dept_emp | NULL       | ref    | PRIMARY,dept_no | dept_no | 12      | const                | 148054 |   100.00 | Using where; Using index |
|  1 | SIMPLE      | emp      | NULL       | eq_ref | PRIMARY         | PRIMARY | 4       | test.dept_emp.emp_no |      1 |   100.00 | NULL                     |
+----+-------------+----------+------------+--------+-----------------+---------+---------+----------------------+--------+----------+--------------------------+
```


