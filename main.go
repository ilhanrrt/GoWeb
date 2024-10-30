package main

import (
	"fmt"
	"goweb_7/admin/auth"
	admin_models "goweb_7/admin/models"
	"goweb_7/config"
	site_models "goweb_7/site/models"

	//"goweb_7/config"
	"net/http"
)

func main() {
	// Kullanıcı giriş işlemi ve JWT token alma
	token, err := auth.GenerateJWT("exampleUser", "admin")
	if err != nil {
		fmt.Println("JWT oluşturma hatası:", err)
		return
	}

	// Token ile kullanıcı işlemlerini kontrol et
	HandleUserActions(token)

	// Veritabanı işlemleri
	admin_models.Post{}.InitDB()
	admin_models.Post{}.Migrate()
	admin_models.User{}.Migrate()
	admin_models.Category{}.InitDb()
	admin_models.Category{}.Migrate()
	site_models.Category{}.InitDB()
	site_models.Category{}.Migrate()
	// Statik dosyaları sun
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Sunucuyu başlat
	fmt.Println("Sunucu başlatılıyor... http://localhost:8080")
	if err := http.ListenAndServe(":8080", config.Routers()); err != nil {
		fmt.Println("Sunucu başlatma hatası:", err)
	}
}

func HandleUserActions(tokenString string) {
	_, err := auth.ValidateJWT(tokenString)
	if err != nil {
		fmt.Println("Token doğrulama hatası:", err)
		return
	}

}
