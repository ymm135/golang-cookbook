# Mysql 查询流程分析 
> 主要是了解mysql查询语句，通过主键索引查看ID后，查询表文件的具体流程是什么？耗时的主要原因？

[前提知识>>InnoDB整体结构](mysql-base.md)  

- [Mysql 查询流程分析](#mysql-查询流程分析)
  - [测试数据](#测试数据)
    - [Select执行步骤](#select执行步骤)
  - [\[1万\]测试ID不同的单条数据查询及耗时](#1万测试id不同的单条数据查询及耗时)
  - [\[500万\]测试ID不同的单条数据查询及耗时](#500万测试id不同的单条数据查询及耗时)


## 测试数据 

```sql
DROP TABLE IF EXISTS employees;
CREATE TABLE employees (
    emp_no      INT             NOT NULL,
    birth_date  DATE            NOT NULL,
    first_name  VARCHAR(14)     NOT NULL,
    last_name   VARCHAR(16)     NOT NULL,
    gender      ENUM ('M','F')  NOT NULL,    
    hire_date   DATE            NOT NULL,
    PRIMARY KEY (emp_no)
);
```

数据样本: 
```
mysql> select * from employees ;
+--------+------------+------------+-------------+--------+------------+
| emp_no | birth_date | first_name | last_name   | gender | hire_date  |
+--------+------------+------------+-------------+--------+------------+
|  10001 | 1953-09-02 | Georgi     | Facello     | M      | 1986-06-26 |
|  10002 | 1964-06-02 | Bezalel    | Simmel      | F      | 1985-11-21 |
|  10003 | 1959-12-03 | Parto      | Bamford     | M      | 1986-08-28 |
|  10004 | 1954-05-01 | Chirstian  | Koblick     | M      | 1986-12-01 |
|  10005 | 1955-01-21 | Kyoichi    | Maliniak    | M      | 1989-09-12 |
|  10006 | 1953-04-20 | Anneke     | Preusig     | F      | 1989-06-02 |
```

查询的日志信息:
```shell
do_command: info: Command on socket (52) = 3 (Query)
do_command: info: packet: ''; command: 3
dispatch_command: info: command: 3
dispatch_command: query: select * from employees where emp_no = 10006
gtid_pre_statement_checks: info: gtid_next->type=0 owned_gtid.{sidno,gno}={0,0}
mysql_execute_command: info: derived: 0  view: 0
column_bitmaps_signal: info: read_set: 0x7fffc800f348  write_set: 0x7fffc800f368
Field_iterator_table_ref::set_field_iterator: info: field_it for 'employees' is Field_iterator_table
SELECT_LEX::prepare: info: setup_ref_array this 0x7fffc8005930   45 :    0    0    6    1    2    0
setup_fields: info: thd->mark_used_columns: 1
setup_fields: info: thd->mark_used_columns: 1
SELECT_LEX::setup_conds: info: thd->mark_used_columns: 1
get_lock_data: info: count 1
get_lock_data: info: sql_lock->table_count 1 sql_lock->lock_count 0
mysql_lock_tables: info: thd->proc_info System lock
lock_external: info: count 1
THD::decide_logging_format: info: query: select * from employees where emp_no = 10006
THD::decide_logging_format: info: variables.binlog_format: 2
THD::decide_logging_format: info: lex->get_stmt_unsafe_flags(): 0x0
THD::decide_logging_format: info: decision: no logging since mysql_bin_log.is_open() = 0 and (options & OPTION_BIN_LOG) = 0x40000 and binlog_format = 2 and binlog_filter->db_ok(db) = 1
THD::is_classic_protocol: info: type=0

WHERE:(after const change) 0x7fffc80170f8 multiple equal(10006, `test`.`employees`.`emp_no`)
add_key_fields: info: add_key_field for field emp_no
get_lock_data: info: count 1
get_lock_data: info: sql_lock->table_count 1 sql_lock->lock_count 0

WHERE:(after substitute_best_equal) 0x7fffc8018138 1

WHERE:(constants) 0x7fffc8018138 1

Info about JOIN
employees         type: const    q_keys: 1  refs: 1  key: 0  len: 4
                  refs:  10006  
JOIN::make_tmp_tables_info: info: Using end_send
JOIN::exec: info: Sending data
Protocol_classic::start_result_metadata: info: num_cols 6, flags 5
Protocol_classic::end_result_metadata: info: num_cols 6, flags 5
do_select: info: Using end_send
do_select: info: 1 records output
ha_commit_trans: info: all=0 thd->in_sub_stmt=0 ha_info=0x7fffc80020d8 is_real_trans=1
close_thread_tables: info: thd->open_tables: 0x7fffc800f240
MDL_context::release_locks_stored_before: info: found lock to release ticket=0x7fffc800ed80
dispatch_command: info: query ready
net_send_ok: info: affected_rows: 0  id: 0  status: 2  warning_count: 0
net_send_ok: info: OK sent, so no more error sending allowed
```

### [Select执行步骤](https://dev.mysql.com/doc/internals/en/selects.html)   
每个选择都在以下基本步骤中执行：

- JOIN::prepare
  - Initialization and linking JOIN structure to `st_select_lex`.
  - fix_fields() for all items (after fix_fields(), we know everything about item).
  - Moving HAVING to WHERE if possible.
  - Initialization procedure if there is one.

- JOIN::optimize
  - Single select optimization.
  - Creation of first temporary table if needed.

- JOIN::exec
  - Performing select (a second temporary table may be created).

- JOIN::cleanup
  - Removing all temporary tables, other cleanup.

- JOIN::reinit
  - Prepare all structures for execution of SELECT (with JOIN::exec).


[官方查询结构说明](https://dev.mysql.com/doc/internals/en/select-structure.html) 
复杂的查询结构   
有两种描述选择的结构：

- st_select_lex ( SELECT_LEX) 代表 SELECT自己  
- st_select_lex_unit ( SELECT_LEX_UNIT) 用于将多个选择分组  

后一项表示`UNION`操作（没有`UNION`是只有一个的联合， `SELECT`并且在任何情况下都存在此结构）。将来，这种结构也将用于 `EXCEPT`和`INTERSECT`。

例如：  
```
(SELECT ...) UNION (SELECT ... (SELECT...)...(SELECT...UNION...SELECT))
   1           2      3           4             5        6       7
```

将表示为：  
```
------------------------------------------------------------------------
                                                                 level 1
SELECT_LEX_UNIT(2)
|
+---------------+
|               |
SELECT_LEX(1)   SELECT_LEX(3)
                |
--------------- | ------------------------------------------------------
                |                                                level 2
                +-------------------+
                |                   |
                SELECT_LEX_UNIT(4)  SELECT_LEX_UNIT(6)
                |                   |
                |                   +--------------+
                |                   |              |
                SELECT_LEX(4)       SELECT_LEX(5)  SELECT_LEX(7)

------------------------------------------------------------------------
```

> 注意：单个子查询 4 ​​有自己的 SELECT_LEX_UNIT.  




`sql/sql_optimizer.h`
```c++

class JOIN :public Sql_alloc
{
  JOIN(const JOIN &rhs);                        /**< not implemented */
  JOIN& operator=(const JOIN &rhs);             /**< not implemented */

  /// Query block that is optimized and executed using this JOIN
  SELECT_LEX *const select_lex;
  /// Query expression referring this query block
  SELECT_LEX_UNIT *const unit;
  /// Thread handler
  THD *const thd;

  int optimize();
  void reset();
  void exec();
  bool prepare_result();
  bool destroy();
  void restore_tmp();
  bool alloc_func_list();
  bool make_sum_func_list(List<Item> &all_fields,
                          List<Item> &send_fields,
                          bool before_group_by, bool recompute= FALSE);
}
```

`sql/sql_lex.h`
```c++
class st_select_lex: public Sql_alloc 
{
  void set_query_result(Query_result *result) { m_query_result= result; }
  Query_result *query_result() const { return m_query_result; }
  bool change_query_result(Query_result_interceptor *new_result,
                           Query_result_interceptor *old_result);
  /// Result of this query block
  Query_result *m_query_result;
}
```

调用栈:
```
end_send(JOIN * join, QEP_TAB * qep_tab, bool end_of_records) (mysql-server/sql/sql_executor.cc:2933)
do_select(JOIN * join) (mysql-server/sql/sql_executor.cc:902)
JOIN::exec(JOIN * const this) (mysql-server/sql/sql_executor.cc:206)     同步调用 st_select_lex::set_query_result
handle_query(THD * thd, LEX * lex, Query_result * result, ulonglong added_options, ulonglong removed_options) (mysql-server/sql/sql_select.cc:191)
execute_sqlcom_select(THD * thd, TABLE_LIST * all_tables) (mysql-server/sql/sql_parse.cc:5167)
mysql_execute_command(THD * thd, bool first_level) (mysql-server/sql/sql_parse.cc:2829)
mysql_parse(THD * thd, Parser_state * parser_state) (mysql-server/sql/sql_parse.cc:5600)
dispatch_command(THD * thd, const COM_DATA * com_data, enum_server_command command) (mysql-server/sql/sql_parse.cc:1493)
do_command(THD * thd) (mysql-server/sql/sql_parse.cc:1032)
handle_connection(void * arg) (mysql-server/sql/conn_handler/connection_handler_per_thread.cc:313)
pfs_spawn_thread(void * arg) (mysql-server/storage/perfschema/pfs.cc:2197)
libpthread.so.0!start_thread (未知源:0)
libc.so.6!clone (未知源:0)
```

目前主要想知道通过主键索引查询到数据位置，如何到磁盘中查询对应记录的？  



[测试数据](https://github.com/datacharmer/test_db/blob/master/load_employees.dump)  

## [1万]测试ID不同的单条数据查询及耗时  
`select * from employees where emp_no = 10006 `  



## [500万]测试ID不同的单条数据查询及耗时  
