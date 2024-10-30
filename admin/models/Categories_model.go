package models

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// Initdb initializes the database connection and stores it in a global variable
func (category Category) InitDb() {
	var err error
	db, err = gorm.Open(mysql.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
	}
}

// Category model
type Category struct {
	ID                  int64 `gorm:"primarykey"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`
	Title, Slug, Status string
}

// Migrate runs the auto-migration for the Category model
func (category Category) Migrate() {
	if db == nil {
		fmt.Println("Database connection is not initialized - category migration")
		return
	}
	err := db.AutoMigrate(&category)
	if err != nil {
		fmt.Println("Failed to migrate:", err)
	}
}

// Add inserts a new category record in the database
func (category Category) Add() (int64, error) {
	if db == nil {
		return 0, fmt.Errorf("database connection is not initialized - category add")

	}
	err := db.Create(&category).Error
	if err != nil {
		return 0, err
	}
	return category.ID, nil
}

// Get retrieves a category record by the provided condition
func (category Category) Get(where ...interface{}) Category {
	if db == nil {
		fmt.Println("Database connection is not initialized")
		return category
	}
	err := db.First(&category, where...).Error
	if err != nil {
		fmt.Println("Failed to get category:", err)
	}
	return category
}

// GetAll retrieves all category records that match the provided condition
func (category Category) GetAll(where ...interface{}) []Category {
	if db == nil {
		fmt.Println("Database connection is not initialized")
		return nil
	}
	var categories []Category
	err := db.Find(&categories, where...).Error
	if err != nil {
		fmt.Println("Failed to get categories:", err)
	}
	return categories
}

// Update updates a specific column of a category record
func (category Category) Update(column string, value interface{}) {
	if db == nil {
		fmt.Println("Database connection is not initialized")
		return
	}
	err := db.Model(&category).Update(column, value).Error
	if err != nil {
		fmt.Println("Failed to update category:", err)
	}
}

// Updates updates multiple fields of a category record
func (category Category) Updates(data Category) {
	if db == nil {
		fmt.Println("Database connection is not initialized")
		return
	}
	err := db.Model(&category).Updates(data).Error
	if err != nil {
		fmt.Println("Failed to update category:", err)
	}
}

// Delete soft-deletes a category record (i.e., sets deleted_at timestamp)
func (category Category) Delete() {
	if db == nil {
		fmt.Println("Database connection is not initialized")
		return
	}
	db.Model(category).Updates(Category{
		Status:    "delete",
		DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
	})
	//err := db.Delete(&category, category.ID).Error
	//if err != nil {
	//	fmt.Println("Failed to delete category:", err)
	//}
}
