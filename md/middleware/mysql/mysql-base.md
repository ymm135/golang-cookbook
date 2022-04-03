# mysql 基础知识  
[参考书籍《深入浅出MySQL：数据库开发、优化与管理维护(第2版)]()  

# 用于日常开发的知识  
## 存储引擎  
MySQL5.0支持的存储引擎包括 `MyISAM`、`InnoDB`、`BDB`、`MEMORY`、`MERGE`、`EXAMPLE`、`NDBCluster`、`ARCHIVE`、`CSV`、`BLACKHOLE`、`FEDERATED`等，其中I`nnoDB`和`BDB`提供事务安全表，其他存储引擎都是非事务安全表。  
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


## 数据类型 

## 字符集  

## 索引的设计及使用  

## 视图  

- ### [视图原理](mysql-viewer.md)  

## 存储过程和函数  

## 触发器  

## 事务控制和锁定语句  






