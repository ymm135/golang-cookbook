- # mysql 日常问题整理

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

> 最后发现问题了，如果在优化时`optimizer`，突然断点，那就会导致表无法使用了。  

