- # mysql 日常问题整理

- [Mysql 表修复](#mysql-表修复)
- [centos7 mysql启动问题](#centos7-mysql启动问题)
- [mysql启动异常](#mysql启动异常)


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