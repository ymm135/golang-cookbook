## gorm基础  
[官网](https://gorm.io/index.html)  

## CRUD 
[code](../../../code/middleware/gorm/base/main.go)  
```go
package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Name string
	Age  int
	Id   int
}

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/one?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return
	}

	user := &User{}
	db.Raw("SELECT id, name, age FROM users WHERE id = ?", 1).Scan(user)
	fmt.Println(user)

	user1 := &User{}
	result1 := db.Where("id = ?", 11).First(user1)
	fmt.Println(result1.RowsAffected, user1)

	// check error ErrRecordNotFound
	errors.Is(result1.Error, gorm.ErrRecordNotFound)

	// Create
	userCreate := User{Name: "xiaoming222", Age: 28}
	result := db.Create(&userCreate)
	resultId := userCreate.Id // 获取创建成功的ID
	fmt.Println(result, resultId)

	// Read
	var user2 User
	// user2 的值会被修改为 id = 1的记录，所有的值都会被修改
	db.First(&user2, resultId) // find user2 with integer primary key，
	// user2已经对应一条记录，加入id为8，这时查询语句为: SELECT * FROM `users` WHERE age = 20 AND `users`.`id` = 8 ORDER BY `users`.`id` LIMIT 1
	db.First(&user2, "age = ?", 20) // find user2 with code D42
	fmt.Println(user2)

	// Update - update product's price to 200
	db.Model(&user2).Update("Age", 19) // user2也会被修改 age = 19
	// Update - update multiple fields
	db.Model(&user2).Updates(User{Name: "xiaohong2", Age: 20}) // non-zero fields
	db.Model(&user2).Updates(map[string]interface{}{"Name": "mingming", "Age": 26})

	// Delete - delete product
	db.Delete(&user2)

}

```

## 实现  

代码位置: 
```shell
$ ~/go/pkg/mod/gorm.io/gorm@v1.20.11/callbacks 
▶ ls -l
total 168
-r--r--r--  1 root  staff  10950 10 27 17:50 associations.go
-r--r--r--  1 root  staff   2334 10 27 17:50 callbacks.go
-r--r--r--  1 root  staff    634 10 27 17:50 callmethod.go
-r--r--r--  1 root  staff  11888 10 27 17:50 create.go
-r--r--r--  1 root  staff   5281 10 27 17:50 delete.go
-r--r--r--  1 root  staff   2262 10 27 17:50 helper.go
-r--r--r--  1 root  staff    667 10 27 17:50 interfaces.go
-r--r--r--  1 root  staff   5520 10 27 17:50 preload.go
-r--r--r--  1 root  staff   7179 10 27 17:50 query.go
-r--r--r--  1 root  staff    337 10 27 17:50 raw.go
-r--r--r--  1 root  staff    527 10 27 17:50 row.go
-r--r--r--  1 root  staff    579 10 27 17:50 transaction.go
-r--r--r--  1 root  staff   8191 10 27 17:50 update.go
```

- ### 数据填充  

```go
type User struct {
	Name string
	Age  int
	Id   int
}

user1 := &User{}
result1 := db.Where("id = ?", 11).First(user1)

```

> user1 如何被填充的呢？  

```go
// First find first record that match given conditions, order by primary key
func (db *DB) First(dest interface{}, conds ...interface{}) (tx *DB) {
	tx = db.Limit(1).Order(clause.OrderByColumn{
		Column: clause.Column{Table: clause.CurrentTable, Name: clause.PrimaryKey},
	})
	if len(conds) > 0 {
		if exprs := tx.Statement.BuildCondition(conds[0], conds[1:]...); len(exprs) > 0 {
			tx.Statement.AddClause(clause.Where{Exprs: exprs})
		}
	}
	tx.Statement.RaiseErrorOnNotFound = true
	tx.Statement.Dest = dest
	return tx.callbacks.Query().Execute(tx)
}
```

> 保存到 `tx.Statement.Dest = dest`  
> `db.Statement.SQL.buf` 的值为 `SELECT id, name, age FROM users WHERE id = ?` 

`Dest`内容  
```
Dest = {interface{} | *main.User} 
 Name = {string} ""
 Age = {int} 0
 Id = {int} 0
```

查看执行计划: 
```
func (p *processor) Execute(db *DB) *DB {
	...
	var (
		curTime           = time.Now()
		stmt              = db.Statement
		resetBuildClauses bool
	)

	// assign model values   // 如果没有声明Model，一般会以目标获取表名    
	if stmt.Model == nil {
		stmt.Model = stmt.Dest
	} else if stmt.Dest == nil {
		stmt.Dest = stmt.Model
	}

	// assign stmt.ReflectValue
	if stmt.Dest != nil {
		stmt.ReflectValue = reflect.ValueOf(stmt.Dest)             // ValueOf 
		for stmt.ReflectValue.Kind() == reflect.Ptr {              // 确认是指针，能够修改值的  
			if stmt.ReflectValue.IsNil() && stmt.ReflectValue.CanAddr() {
				stmt.ReflectValue.Set(reflect.New(stmt.ReflectValue.Type().Elem()))
			}

			stmt.ReflectValue = stmt.ReflectValue.Elem()           // stmt.ReflectValue  
		}
		if !stmt.ReflectValue.IsValid() {
			db.AddError(ErrInvalidValue)
		}
	}

	for _, f := range p.fns {
		f(db)
	}
```

最后要填充的内容 `stmt.ReflectValue = stmt.ReflectValue.Elem()           // stmt.ReflectValue`

find查询语句执行的方法,通过`f(db)`依次调用方法:
```
fns = {[]func(*gorm.DB)} len:3, cap:4
 0 = {func(*gorm.DB)} gorm.io/gorm/callbacks.Query
 1 = {func(*gorm.DB)} gorm.io/gorm/callbacks.Preload
 2 = {func(*gorm.DB)} gorm.io/gorm/callbacks.AfterQuery
```

query整个调用为`gorm@v1.23.4/callbacks/query.go`:  

```go
func Query(db *gorm.DB) {
	if db.Error == nil {
		BuildQuerySQL(db)   // 构建语句

		if !db.DryRun && db.Error == nil {
			// SELECT * FROM `users` ORDER BY `users`.`id` LIMIT 1
			rows, err := db.Statement.ConnPool.QueryContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)  // 执行
			if err != nil {
				db.AddError(err)
				return
			}
			defer func() {
				db.AddError(rows.Close())
			}()
			// 最终填充在Scan， Find也是需要调用Scan 
			gorm.Scan(rows, db, 0)
		}
	}
}
```

`gorm@v1.23.4/scan.go` 
```
// Scan scan rows into db statement
func Scan(rows Rows, db *DB, mode ScanMode) {
	var (
		columns, _          = rows.Columns()
		values              = make([]interface{}, len(columns))
		initialized         = mode&ScanInitialized != 0
		update              = mode&ScanUpdate != 0
		onConflictDonothing = mode&ScanOnConflictDoNothing != 0
	)

	switch dest := db.Statement.Dest.(type) {
		default:
		var (
			fields             = make([]*schema.Field, len(columns))
			selectedColumnsMap = make(map[string]int, len(columns))
			joinFields         [][2]*schema.Field
			sch                = db.Statement.Schema
			reflectValue       = db.Statement.ReflectValue
		)
		...
		// struct
		reflectValueType := reflectValue.Type() 
		...
		case reflect.Struct, reflect.Ptr:
			if initialized || rows.Next() {
				db.scanIntoStruct(rows, reflectValue, values, fields, joinFields)
			}		
```
`db.scanIntoStruct(rows, reflectValue, values, fields, joinFields)` 变量值为:  
```go
执行前为nil, 执行结束后为: 
values = {[]interface{}} len:3, cap:3
 0 = {interface{} | **string} "xiaohong"
 1 = {interface{} | **int64} 18
 2 = {interface{} | **int64} 2

执行前为:
fields = {[]*schema.Field} len:3, cap:3
 0 = {*schema.Field} 
  Name = {string} "Name"
  DBName = {string} "name"
  BindNames = {[]string} len:1, cap:1
  DataType = {schema.DataType} "string"
  GORMDataType = {schema.DataType} "string"
  PrimaryKey = {bool} false
  ...
 1 = {*schema.Field} 
 2 = {*schema.Field} 

joinFields nil 
```

```go
func (db *DB) scanIntoStruct(rows Rows, reflectValue reflect.Value, values []interface{}, fields []*schema.Field, joinFields [][2]*schema.Field) {
	// 把值存储到values
	for idx, field := range fields {
		if field != nil {
			values[idx] = field.NewValuePool.Get()
		} else if len(fields) == 1 {
			if reflectValue.CanAddr() {
				values[idx] = reflectValue.Addr().Interface()
			} else {
				values[idx] = reflectValue.Interface()
			}
		}
	}

	db.RowsAffected++
	db.AddError(rows.Scan(values...))

	for idx, field := range fields {
		if field != nil {
			if len(joinFields) == 0 || joinFields[idx][0] == nil {
				// 设置值  
				db.AddError(field.Set(db.Statement.Context, reflectValue, values[idx]))
			} else {
				relValue := joinFields[idx][0].ReflectValueOf(db.Statement.Context, reflectValue)
				if relValue.Kind() == reflect.Ptr && relValue.IsNil() {
					if value := reflect.ValueOf(values[idx]).Elem(); value.Kind() == reflect.Ptr && value.IsNil() {
						return
					}

					relValue.Set(reflect.New(relValue.Type().Elem()))
				}
				db.AddError(joinFields[idx][1].Set(db.Statement.Context, relValue, values[idx]))
			}

			// release data to pool
			field.NewValuePool.Put(values[idx])
		}
	}
}
```
`db.AddError(field.Set(db.Statement.Context, reflectValue, values[idx]))`  调用的是:  

```
// create valuer, setter when parse struct
func (field *Field) setupValuerAndSetter() {
	case reflect.String:
		field.Set = func(ctx context.Context, value reflect.Value, v interface{}) (err error) {
			switch data := v.(type) {
			case **string:
				if data != nil && *data != nil {
					// 设置值 
					field.ReflectValueOf(ctx, value).SetString(**data)
				}

# field.ReflectValueOf(ctx, value).SetString(**data) 调用如下:  
// ReflectValueOf returns field's reflect value
	switch {
	case len(field.StructField.Index) == 1 && fieldIndex > 0:
		field.ReflectValueOf = func(ctx context.Context, value reflect.Value) reflect.Value {
			return reflect.Indirect(value).Field(fieldIndex)
		}

```

整个sql语句的处理像一个流水线，一层一层处理，不同场景，调用流程是不一样的，回调函数注册:  
`gorm@v1.23.4/callbacks/callbacks.go`  

```go
func RegisterDefaultCallbacks(db *gorm.DB, config *Config) {
	enableTransaction := func(db *gorm.DB) bool {
		return !db.SkipDefaultTransaction
	}

	if len(config.CreateClauses) == 0 {
		config.CreateClauses = createClauses
	}
	if len(config.QueryClauses) == 0 {
		config.QueryClauses = queryClauses
	}
	if len(config.DeleteClauses) == 0 {
		config.DeleteClauses = deleteClauses
	}
	if len(config.UpdateClauses) == 0 {
		config.UpdateClauses = updateClauses
	}

	createCallback := db.Callback().Create()
	createCallback.Match(enableTransaction).Register("gorm:begin_transaction", BeginTransaction)
	createCallback.Register("gorm:before_create", BeforeCreate)
	createCallback.Register("gorm:save_before_associations", SaveBeforeAssociations(true))
	createCallback.Register("gorm:create", Create(config))
	createCallback.Register("gorm:save_after_associations", SaveAfterAssociations(true))
	createCallback.Register("gorm:after_create", AfterCreate)
	createCallback.Match(enableTransaction).Register("gorm:commit_or_rollback_transaction", CommitOrRollbackTransaction)
	createCallback.Clauses = config.CreateClauses

	queryCallback := db.Callback().Query()
	queryCallback.Register("gorm:query", Query)
	queryCallback.Register("gorm:preload", Preload)
	queryCallback.Register("gorm:after_query", AfterQuery)
	queryCallback.Clauses = config.QueryClauses

	deleteCallback := db.Callback().Delete()
	deleteCallback.Match(enableTransaction).Register("gorm:begin_transaction", BeginTransaction)
	deleteCallback.Register("gorm:before_delete", BeforeDelete)
	deleteCallback.Register("gorm:delete_before_associations", DeleteBeforeAssociations)
	deleteCallback.Register("gorm:delete", Delete(config))
	deleteCallback.Register("gorm:after_delete", AfterDelete)
	deleteCallback.Match(enableTransaction).Register("gorm:commit_or_rollback_transaction", CommitOrRollbackTransaction)
	deleteCallback.Clauses = config.DeleteClauses

	updateCallback := db.Callback().Update()
	updateCallback.Match(enableTransaction).Register("gorm:begin_transaction", BeginTransaction)
	updateCallback.Register("gorm:setup_reflect_value", SetupUpdateReflectValue)
	updateCallback.Register("gorm:before_update", BeforeUpdate)
	updateCallback.Register("gorm:save_before_associations", SaveBeforeAssociations(false))
	updateCallback.Register("gorm:update", Update(config))
	updateCallback.Register("gorm:save_after_associations", SaveAfterAssociations(false))
	updateCallback.Register("gorm:after_update", AfterUpdate)
	updateCallback.Match(enableTransaction).Register("gorm:commit_or_rollback_transaction", CommitOrRollbackTransaction)
	updateCallback.Clauses = config.UpdateClauses

	rowCallback := db.Callback().Row()
	rowCallback.Register("gorm:row", RowQuery)
	rowCallback.Clauses = config.QueryClauses

	rawCallback := db.Callback().Raw()
	rawCallback.Register("gorm:raw", RawExec)
	rawCallback.Clauses = config.QueryClauses
}
```

- ### db.Model  

目前Model的作用是获取表格名称，取结构体名称，如果表格名称与结构体名称不一致，需要增加函数: 
```go
func (p *DeviceName) TableName() string {
	return "table_name"
}
```

```go
// Model specify the model you would like to run db operations
//    // update all users's name to `hello`
//    db.Model(&User{}).Update("name", "hello")
//    // if user's primary key is non-blank, will use it as condition, then will only update the user's name to `hello`
//    db.Model(&user).Update("name", "hello")
func (db *DB) Model(value interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Model = value
	return
}
```

```go
				if tx.Statement.Parse(tx.Statement.Model) == nil {
					if f := tx.Statement.Schema.LookUpField(dbName); f != nil {
						dbName = f.DBName
					}
				}

// 解析表名 
func (stmt *Statement) Parse(value interface{}) (err error) {
	return stmt.ParseWithSpecialTableName(value, "")
}


```


- ### Find 与 Scan 区别  

```go
// Find find records that match given conditions
func (db *DB) Find(dest interface{}, conds ...interface{}) (tx *DB) {
	tx = db.getInstance()
	if len(conds) > 0 {
		tx.Statement.AddClause(clause.Where{Exprs: tx.Statement.BuildCondition(conds[0], conds[1:]...)})
	}
	tx.Statement.Dest = dest
	tx.callbacks.Query().Execute(tx)   // 走标准流程
	return
}

func (p *processor) Execute(db *DB) {
	...
	func (p *processor) Execute(db *DB) {
	...
}
```

find查询语句执行的方法,通过`f(db)`依次调用方法:
```
fns = {[]func(*gorm.DB)} len:3, cap:4
 0 = {func(*gorm.DB)} gorm.io/gorm/callbacks.Query
 1 = {func(*gorm.DB)} gorm.io/gorm/callbacks.Preload
 2 = {func(*gorm.DB)} gorm.io/gorm/callbacks.AfterQuery
```

`gorm@v1.20.11/callbacks/query.go`  
```go
func Query(db *gorm.DB) {
	if db.Error == nil {
		BuildQuerySQL(db)

		if !db.DryRun && db.Error == nil {
			rows, err := db.Statement.ConnPool.QueryContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)
			if err != nil {
				db.AddError(err)
				return
			}
			defer rows.Close()

			gorm.Scan(rows, db, false)  // 最后还是需要调用sacn,只能find增加需要其他“语法功能”
		}
	}
}
```


Scan实现: 
```go
// Scan scan value to a struct
func (db *DB) Scan(dest interface{}) (tx *DB) {
	config := *db.Config
	currentLogger, newLogger := config.Logger, logger.Recorder.New()
	config.Logger = newLogger

	tx = db.getInstance()
	tx.Config = &config

	if rows, err := tx.Rows(); err != nil {
		tx.AddError(err)
	} else {
		defer rows.Close()
		if rows.Next() {  // 执行到这里  
			tx.ScanRows(rows, dest)
		}
	}

	currentLogger.Trace(tx.Statement.Context, newLogger.BeginAt, func() (string, int64) {
		return newLogger.SQL, tx.RowsAffected
	}, tx.Error)
	tx.Logger = currentLogger
	return
}
```
调用到: `ScanRows(rows *sql.Rows, dest interface{})`  
```go
func (db *DB) ScanRows(rows *sql.Rows, dest interface{}) error {
	tx := db.getInstance()
	if err := tx.Statement.Parse(dest); !errors.Is(err, schema.ErrUnsupportedDataType) {
		tx.AddError(err)
	}
	tx.Statement.Dest = dest
	tx.Statement.ReflectValue = reflect.ValueOf(dest)
	for tx.Statement.ReflectValue.Kind() == reflect.Ptr {
		tx.Statement.ReflectValue = tx.Statement.ReflectValue.Elem()
	}
	Scan(rows, tx, true)  // 走到这里
	return tx.Error
}
```  
`gorm@v1.20.11/scan.go`  
```go
func Scan(rows *sql.Rows, db *DB, initialized bool) {
	columns, _ := rows.Columns()
	values := make([]interface{}, len(columns))
	db.RowsAffected = 0
	switch dest := db.Statement.Dest.(type) {
		...
		default:
			Schema := db.Statement.Schema

			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
			...

			if Schema != nil {
				... 
				// 匹配字段  以数据库查到的列顺序为主  idx 就是列位置  
				for idx, column := range columns {
					if field := Schema.LookUpField(column); field != nil && field.Readable {
						fields[idx] = field
					} else if names := strings.Split(column, "__"); len(names) > 1 {

					}
					...
	...
	db.Statement.ReflectValue.Set(reflect.MakeSlice(db.Statement.ReflectValue.Type(), 0, 20))
	...
	// 开始赋值给结构体  
	for initialized || rows.Next() {

		for idx, field := range fields {
			if len(joinFields) != 0 && joinFields[idx][0] != nil {		
				...
			else if field != nil {
				// 设置值    
				field.Set(elem, values[idx])
			}
``` 


<br>
<div align=center>
    <img src="../../../res/gorm filed.png" width="80%" height="80%"></img>  
</div>
<br>

scan填充到结构体，也是需要一行行扫描，然后把每一个的不同列填充到结构体成员中:  
`rows.rowsi.mysqlRows.rs.columns`包含一行数据的不同列  
```go
rows = {*sql.Rows} 
 dc = {*sql.driverConn} 
 releaseConn = {func(error)} database/sql.(*driverConn).releaseConn-fm
 rowsi = {driver.Rows | *mysql.binaryRows} 
  mysqlRows = {mysql.mysqlRows} 
   mc = {*mysql.mysqlConn} 
   rs = {mysql.resultSet} 
    columns = {[]mysql.mysqlField} len:6, cap:6
     0 = {mysql.mysqlField} 
      tableName = {string} ""
      name = {string} "threatName"
      length = {uint32} 1020
      flags = {mysql.fieldFlag} 0
      fieldType = {mysql.fieldType} fieldTypeVarString (253)
      decimals = {uint8} 0
      charSet = {uint8} 45
     1 = {mysql.mysqlField} 
     2 = {mysql.mysqlField} 
     3 = {mysql.mysqlField} 
     4 = {mysql.mysqlField} 
     5 = {mysql.mysqlField} 
    columnNames = {[]string} len:6, cap:6
    done = {bool} false
   finish = {func()} github.com/go-sql-driver/mysql.(*mysqlConn).finish-fm
 cancel = {func()} nil
 closeStmt = {*sql.driverStmt} 
```

Find与Scan都会走到这里:  
`gorm@v1.20.11/callbacks/query.go`   
```
func Query(db *gorm.DB) {
	if db.Error == nil {
		BuildQuerySQL(db)

		if !db.DryRun && db.Error == nil {
			rows, err := db.Statement.ConnPool.QueryContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)
			if err != nil {
				db.AddError(err)
				return
			}
			defer rows.Close()

			gorm.Scan(rows, db, false)
		}
	}
}
```

`gorm.Scan(rows, db, false)` 的实现就是上面介绍的，把`rows`的数据一行行的存储到`db.Statement.ReflectValue`  


- ### SubQuery 子查询

- ### Iteration 遍历

- ### 多表查询  Joins  

官方demo 
```go
type result struct {
  Name  string
  Email string
}

db.Model(&User{}).Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&result{})
// SELECT users.name, emails.email FROM `users` left join emails on emails.user_id = users.id

rows, err := db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Rows()
for rows.Next() {
  ...
}

db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&results)

// multiple joins with parameter
db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "jinzhu@example.org").Joins("JOIN credit_cards ON credit_cards.user_id = users.id").Where("credit_cards.number = ?", "411111111111").Find(&user)

```



- ### Raw SQL  

`callbacks/row.go`  
```
func RowQuery(db *gorm.DB) {
	if db.Error == nil {
		BuildQuerySQL(db)

		if !db.DryRun {
			if isRows, ok := db.InstanceGet("rows"); ok && isRows.(bool) {
				db.Statement.Dest, db.Error = db.Statement.ConnPool.QueryContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)
			} else {
				db.Statement.Dest = db.Statement.ConnPool.QueryRowContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)
			}

			db.RowsAffected = -1
		}
	}
}
```

`Raw`执行时没有走工作流程`fn(db)`，相当于直接把语句给mysql。  

