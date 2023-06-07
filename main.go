package main

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Password string
}

func main() {
	db, err := gorm.Open(sqlite.Open("data/sqlite.db?mode=wal"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to migrate database")
	}

	user := &User{Username: "alice", Password: "123456"}
	db.Create(&user)

	var users []User
	db.Find(&users)

	fmt.Println("users: ", users)
}
