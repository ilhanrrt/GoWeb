package models

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                           uint
	Username, Password, UserType string
}

func (user User) Migrate() {
	err := Db.AutoMigrate(&user) // Global veritabanı bağlantısını kullanıyoruz
	if err != nil {
		fmt.Println(err)
	}
}

func (user User) Add() {
	err := Db.Create(&user).Error // Global veritabanı bağlantısını kullanıyoruz
	if err != nil {
		fmt.Println(err)
	}
}

func (user User) Get(where ...interface{}) User {
	err := Db.First(&user, where...).Error // Global veritabanı bağlantısını kullanıyoruz
	if err != nil {
		fmt.Println(err)
	}
	return user
}

func (user User) GetAll(where ...interface{}) []User {
	var users []User
	err := Db.Find(&users, where...).Error // Global veritabanı bağlantısını kullanıyoruz
	if err != nil {
		fmt.Println(err)
	}
	return users
}

func (user User) Update(column string, value interface{}) {
	err := Db.Model(&user).Update(column, value).Error // Global veritabanı bağlantısını kullanıyoruz
	if err != nil {
		fmt.Println(err)
	}
}

func (user User) Updates(data User) {
	err := Db.Model(&user).Updates(data).Error // Global veritabanı bağlantısını kullanıyoruz
	if err != nil {
		fmt.Println(err)
	}
}

func (user User) Delete() {
	err := Db.Delete(&user, user.ID).Error // Global veritabanı bağlantısını kullanıyoruz
	if err != nil {
		fmt.Println(err)
	}
}
