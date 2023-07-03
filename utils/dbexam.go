package utils

import (
	"gorm.io/gorm"
	"log"
	"os/user"
)

// db.go
//package db
//
//import (
//"log"
//
//"github.com/jinzhu/gorm"
//_ "github.com/jinzhu/gorm/dialects/mysql"
//)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "user:password@tcp(127.0.0.1:3306)/database")
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.AutoMigrate(&User{})
}

func GetDB() *gorm.DB {
	return db
}

// user.go
//package user
//
//import (
//"log"
//
//"path/to/db"
//
//"github.com/jinzhu/gorm"
//)

type User struct {
	gorm.Model
	Name  string
	Email string
}

func CreateUser(u *User) error {
	if err := db.GetDB().Create(u).Error; err != nil {
		log.Printf("Failed to create user: %v", err)
		return err
	}
	return nil
}

func CreateUser2(u *User) error {
	// 使用事务
	tx := db.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(u).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to create user: %v", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return err
	}

	return nil
}

// main.go
//package main
//
//import (
//"log"
//
//"path/to/db"
//"path/to/user"
//)

func main() {
	u := &user.User{Name: "Alice", Email: "alice@example.com"}
	err := user.CreateUser(u)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	log.Printf("User created: %+v", u)
}
