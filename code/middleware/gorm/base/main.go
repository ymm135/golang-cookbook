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
