package models

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB // Global olarak bir veritabanı değişkeni tanımlıyoruz

// Veritabanı bağlantısını başlatan fonksiyon
func (category Category) InitDB() error {
	var err error
	Db, err = gorm.Open(mysql.Open(Dns), &gorm.Config{}) // Veritabanı bağlantısını açıyoruz
	if err != nil {
		return fmt.Errorf("veritabanı bağlantısı açılamadı: %v", err)
	}
	return nil
}

type Category struct {
	gorm.Model
	Title, Slug string
}

func (category Category) Migrate() {
	err := Db.AutoMigrate(&category) // Global veritabanı bağlantısını kullanıyoruz
	if err != nil {
		fmt.Println(err)
	}
}

func (category Category) Add() {
	err := Db.Create(&category).Error // Global veritabanı bağlantısını kullanıyoruz
	if err != nil {
		fmt.Println(err)
	}
}

func (category Category) Get(where ...interface{}) Category {
	err := Db.First(&category, where...).Error // Global veritabanı bağlantısını kullanıyoruz
	if err != nil {
		fmt.Println(err)
	}
	return category
}

func (category Category) GetAll(where ...interface{}) []Category {
	var categories []Category
	err := Db.Find(&categories, where...).Error // Global veritabanı bağlantısını kullanıyoruz
	if err != nil {
		fmt.Println(err)
	}
	return categories
}

func (category Category) Update(column string, value interface{}) {
	err := Db.Model(&category).Update(column, value).Error // Global veritabanı bağlantısını kullanıyoruz
	if err != nil {
		fmt.Println(err)
	}
}

func (category Category) Updates(data Category) {
	err := Db.Model(&category).Updates(data).Error // Global veritabanı bağlantısını kullanıyoruz
	if err != nil {
		fmt.Println(err)
	}
}

func (category Category) Delete() {
	err := Db.Delete(&category, category.ID).Error // Global veritabanı bağlantısını kullanıyoruz
	if err != nil {
		fmt.Println(err)
	}
}
