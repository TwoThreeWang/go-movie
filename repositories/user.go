package repositories

import (
	"github.com/jinzhu/gorm"
	"log"
)

type User struct {
	gorm.Model
	Name  string
	Email string
}

func CreateUser(u *User) error {
	if err := GetDB().Create(u).Error; err != nil {
		log.Printf("Failed to create user: %v", err)
		return err
	}
	return nil
}

func CreateUser2(u *User) error {
	// 使用事务
	tx := GetDB().Begin()
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
