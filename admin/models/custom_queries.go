package models

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// CustomPostQuery is a struct for custom queries
type CustomPostQuery struct {
	ID                                             int64 `gorm:"primarykey"`
	CreatedAt                                      time.Time
	UpdatedAt                                      time.Time
	DeletedAt                                      gorm.DeletedAt `gorm:"index"`
	Title, Slug, Description, Content, Picture_url string
	CategoryID                                     int
	Status                                         string `gorm:"column=status" json:"status"`
}

func GetUserByUsernameAndPassword(username, password string) (User, error) {
	db, err := gorm.Open(mysql.Open(Dns), &gorm.Config{})
	if err != nil {
		return User{}, err
	}

	var user User
	db.First(&user, "username = ? AND password = ?", username, password)
	return user, nil
}

func (query CustomPostQuery) GetPosts(condition string, args ...interface{}) ([]Post, error) {
	var posts []Post
	if Db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	// Koşulu ve argümanları kullanarak veritabanından postları al
	err := Db.Where(condition, args...).Find(&posts).Error
	if err != nil {
		fmt.Println("Failed to get posts:", err)
		return nil, err
	}

	return posts, nil
}

// GetAllPosts fetches all posts from the database
func (query CustomPostQuery) GetAllPosts() ([]Post, error) {
	var posts []Post
	if Db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	// Veritabanındaki tüm post verilerini getir
	err := Db.Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (query CustomPostQuery) DeletePostByID(id int) error {
	if Db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	var post Post
	// Post kaydını bul
	result := Db.First(&post, id)
	if result.Error != nil {
		return result.Error
	}

	// Eğer herhangi bir satır etkilenmediyse hata döndür
	if result.RowsAffected == 0 {
		return fmt.Errorf("post bulunamadı")
	}

	// Statü ve silinme tarihini güncelle
	post.Status = "delete"                                         // Statü değişimi
	post.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true} // Silinme tarihi ekleniyor

	// Veritabanında güncelleme işlemi
	result = Db.Save(&post)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (query CustomPostQuery) GetAllCategories() ([]Category, error) {
	var categories []Category
	if Db == nil {
		return nil, fmt.Errorf("veritabanı bağlantısı yapılandırılmadı")
	}

	// Veritabanından kategorileri al
	err := Db.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// DeleteCategoryByID silme işlemi için veritabanına bağlanan fonksiyon
func (query CustomPostQuery) DeleteCategoryByID(id int) error {
	if Db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	var category Category

	result := Db.First(&category, id) // id'ye göre ilk satırı bul
	// Eğer herhangi bir satır etkilenmediyse hata döndür
	if result.RowsAffected == 0 {
		return fmt.Errorf("kategori bulunamadı")
	}
	category.Status = "delete"                                         // Statü değişimi
	category.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true} // Silinme tarihi ekleniyor

	// Veritabanında güncelleme işlemi
	result = Db.Save(&category)
	// Kategori statüsünü "delete" olarak güncelleme işlemi

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (query CustomPostQuery) CreatePost(post Post) (int64, error) {
	if Db == nil {
		return 0, fmt.Errorf("database connection is not initialized")
	}
	err := Db.Create(&post).Error
	if err != nil {
		return 0, err
	}
	return post.ID, nil
}

func (query CustomPostQuery) UpdatePostStatus(id int, newStatus string) error {
	if Db == nil {
		return fmt.Errorf("veritabanı bağlantısı yapılandırılmadı")
	}

	// Statüyü güncelle
	result := Db.Model(&Post{}).Where("id = ?", id).Update("status", newStatus)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("Post bulunamadı")
	}

	return nil
}
