package models

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func (post Post) InitDB() {
	var err error
	Db, err = gorm.Open(mysql.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
	}
}

type Post struct {
	ID                                             int64 `gorm:"primarykey"`
	CreatedAt                                      time.Time
	UpdatedAt                                      time.Time
	DeletedAt                                      gorm.DeletedAt `gorm:"index"`
	Title, Slug, Description, Content, Picture_url string
	CategoryID                                     int
	Status                                         string `gorm:"column=status" json:"status"`
}

func (post Post) Migrate() {
	if Db == nil {
		fmt.Println("Database connection is not initialized")
		return
	}
	err := Db.AutoMigrate(&post)
	if err != nil {
		fmt.Println("Failed to migrate database:", err)
	}

}

func (post *Post) Add() (int64, error) {
	if Db == nil {
		return 0, fmt.Errorf("database connection is not initialized")
	}
	err := Db.Create(&post).Error
	if err != nil {
		return 0, err
	}
	return post.ID, nil
}

func (post Post) Get(where ...interface{}) Post {
	if Db == nil {
		fmt.Println("Database connection is not initialized")
		return post
	}
	err := Db.First(&post, where...).Error
	if err != nil {
		fmt.Println("Failed to get post:", err)
	}
	return post
}

func (post Post) GetAll(where ...interface{}) []Post {
	if Db == nil {
		fmt.Println("Database connection is not initialized")
		return nil
	}
	var posts []Post
	err := Db.Find(&posts, where...).Error
	if err != nil {
		fmt.Println("Failed to get all posts:", err)
		return nil
	}
	// Tüm verileri yazdırarak status sütununun çekilip çekilmediğini kontrol et
	for _, p := range posts {
		fmt.Printf("Post: %+v\n", p)
	}
	return posts
}

func (post Post) Update(column string, value interface{}) {
	if Db == nil {
		fmt.Println("Database connection is not initialized")
		return
	}
	err := Db.Model(&post).Update(column, value).Error
	if err != nil {
		fmt.Println("Failed to update post:", err)
	}
}

func (post Post) Updates(data Post) {
	if Db == nil {
		fmt.Println("Database connection is not initialized")
		return
	}
	err := Db.Model(&post).Updates(data).Error
	if err != nil {
		fmt.Println("Failed to update post:", err)
	}
}

func (post Post) Delete() {
	if Db == nil {
		fmt.Println("Database connection is not initialized")
		return
	}
	Db.Model(post).Updates(Post{
		Status:    "delete",
		DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
	})

	//Db.Delete(&post, post.ID)
}

func (post Post) GetFilteredPosts(categoryID int, status string, startDate, endDate string) ([]Post, error) {
	if Db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var posts []Post
	query := Db.Model(&Post{})

	if categoryID != 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	err := query.Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (post Post) SearchByTitle(query string) []Post {
	if Db == nil {
		fmt.Println("Database connection is not initialized")
		return nil
	}
	var posts []Post
	err := Db.Where("title LIKE ?", "%"+query+"%").Find(&posts).Error
	if err != nil {
		fmt.Println("Failed to search posts by title:", err)
	}
	return posts
}
