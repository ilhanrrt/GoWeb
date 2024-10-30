package controllers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"goweb_7/admin/helpers"
	logger "goweb_7/admin/log"
	"goweb_7/admin/models"
	"net/http"
	"strconv"

	"github.com/gosimple/slug"
	"github.com/julienschmidt/httprouter"
)

type Postman struct {
	ID          int64  `gorm:"primarykey"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Slug        string `json:"slug"`
	Picture_url string `json:"picture_url"`
	Status      string `json:"status"`
	CategoryID  int    `json:"categoryid"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

// Özelleştirilmiş "Method Not Allowed" handler'ı
func (postman Postman) MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {

	helpers.ErrorResponse(w, "Bu metod desteklenmiyor. Lütfen uygun bir HTTP metodu kullanın.", "405001", http.StatusBadRequest)
	currentUser := helpers.GetUserFromSessionPanel(r)
	logger.LogActionPanel(int(currentUser.ID), currentUser.Username, "Bu metod desteklenmiyor. Lütfen uygun bir HTTP metodu kullanın.", false)
}

func (postman Postman) HandlePostAdmin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer func() {
		if err := recover(); err != nil {
			logger.LogActionPanel(0, "unknown", fmt.Sprintf("Panic occurred: %v", err), false)
			helpers.ErrorResponse(w, "Beklenmedik bir hata oluştu", "500002", http.StatusBadRequest)

		}
	}()

	currentUser := helpers.GetUserFromSessionPanel(r)
	if !helpers.CheckUserPanel(w, r) {
		return
	}

	if r.Method == http.MethodGet {

		posts, err := models.CustomPostQuery{}.GetAllPosts()
		if err != nil {
			logger.LogActionPanel(int(currentUser.ID), currentUser.Username, "Postman üzerinden veri tabanı hatası oluştu", false)
			// Beklenmedik bir hata oluştuğunda panic tetikliyoruz
			panic("Veriler alınamadı: " + err.Error())
		}
		logger.LogActionPanel(int(currentUser.ID), currentUser.Username, "Postman üzerinden veriler çağrıldı", true)

		helpers.SuccessResponse(w, "Veriler başarıyla alındı", posts, http.StatusOK)

	}
}

func (postman Postman) LoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer func() {
		if err := recover(); err != nil {
			logger.LogSaveMessagePanel(fmt.Sprintf("Panic occurred: %v", err), false)
			helpers.ErrorResponse(w, "Beklenmedik bir hata oluştu", "500003", http.StatusBadRequest)

		}
	}()

	if r.Method == http.MethodPost {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			helpers.ErrorResponse(w, "Lütfen giriş yapın", "400001", http.StatusBadRequest)
			logger.LogSaveMessagePanel("Giriş yapılırken hata oluştu", false)
			return
		}

		password := fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password)))
		userFromDB, err := models.GetUserByUsernameAndPassword(user.Username, password)
		if err != nil {
			helpers.ErrorResponse(w, "Veri işleme hatası", "400001", http.StatusBadRequest)
			// Beklenmedik bir hata oluştuğunda panic tetikliyoruz
			logger.LogSaveMessagePanel(fmt.Sprintf("Veri işleme hatası: "+err.Error()), false)
			panic("Veri işleme hatası: " + err.Error())
		}

		if userFromDB.ID == 0 {
			helpers.ErrorResponse(w, "Geçersiz kullanıcı adı veya şifre", "401001", http.StatusBadRequest)

			logger.LogSaveMessagePanel("Giriş yapılırken hata oluştu", false)
			return
		}

		err = helpers.SetUserPanel(w, r, userFromDB)
		if err != nil {
			// Oturum oluşturulamadığında panic teti
			helpers.ErrorResponse(w, "Oturum oluşturulamadı", "401001", http.StatusBadRequest)
			logger.LogSaveMessagePanel(fmt.Sprintf("Oturum oluşturulamadı: "+err.Error()), false)
			panic("Oturum oluşturulamadı: " + err.Error())
		}

		logger.LogActionPanel(int(user.ID), user.Username, "Giriş yaptı", true)
		helpers.SuccessResponse(w, "Giriş Başarılı", user.Username, http.StatusOK)

	}
}

func (postman Postman) LogoutHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer func() {
		if err := recover(); err != nil {
			logger.LogSaveMessagePanel(fmt.Sprintf("Panic occurred: %v", err), false)
			helpers.ErrorResponse(w, "Beklenmedik bir hata oluştu", "500003", http.StatusBadRequest)

		}
	}()
	// Oturumdan kullanıcı bilgilerini al
	currentUser := helpers.GetUserFromSessionPanel(r)
	if r.Method == http.MethodPost {
		if currentUser.ID == 0 {
			helpers.ErrorResponse(w, "Oturumdan kullanıcı bilgisi alınamadı", "500003", http.StatusBadRequest)
			// Kullanıcı bilgisi bulunamazsa panic tetikliyoruz
			logger.LogSaveMessagePanel("Oturumdan kullanıcı bilgisi alınamadı", false)
			panic("Oturumdan kullanıcı bilgisi alınamadı")
		}
		// Kullanıcıyı oturumdan çıkar
		err := helpers.RemoveUserPanel(w, r)
		if err != nil {
			helpers.ErrorResponse(w, "Kullanıcı oturumunu kapatırken hata oluştu", "500003", http.StatusBadRequest)
			// Beklenmedik bir hata oluştuğunda panic tetikliyoruz
			logger.LogSaveMessagePanel(fmt.Sprintf("Kullanıcı oturumunu kapatırken hata oluştu: "+err.Error()), false)
			panic("Kullanıcı oturumunu kapatırken hata oluştu: " + err.Error())
		}

		// Kullanıcıya başarı mesajı döndür
		helpers.SuccessResponse(w, "Güle güle yine bekleriz...", currentUser.Username, http.StatusOK)

		// Çıkış yapıldığına dair log kaydı
		logger.LogActionPanel(int(currentUser.ID), currentUser.Username, "Kullanıcı çıkış yaptı", true)
	}
}

func (postman Postman) HandlePostAdminID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	currentUser := helpers.GetUserFromSessionPanel(r)
	if r.Method == http.MethodGet {
		defer func() {
			if err := recover(); err != nil {
				logger.LogActionPanel(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Panic occurred: %v", err), false)
				helpers.ErrorResponse(w, "Beklenmedik bir hata oluştu", "500003", http.StatusBadRequest)

			}
		}()

		if !helpers.CheckUserPanel(w, r) {
			return
		}

		idStr := ps.ByName("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			helpers.ErrorResponse(w, "Geçersiz ID", "400001", http.StatusBadRequest)

			logger.LogActionPanel(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Geçersiz ID arandı Post ID'si %s", idStr), false)
			return
		}

		// Veritabanından sorgulama yapıyoruz
		posts, err := models.CustomPostQuery{}.GetPosts("id = ?", id)
		if err != nil {
			helpers.ErrorResponse(w, "Post bulunamadı", "400002", http.StatusBadRequest)
			// Verileri alırken hata olduğunda panic tetikliyoruz
			logger.LogSaveMessagePanel(fmt.Sprintf("Veriler alınamadı: "+err.Error()), false)
			panic("Veriler alınamadı: " + err.Error())
		}

		// ID'ye karşılık gelen post var mı kontrol ediyoruz
		if len(posts) == 0 {
			// Eğer sonuç boşsa, 404 Not Found döndürüyoruz
			helpers.ErrorResponse(w, "Belirtilen ID ile post bulunamadı", "404001", http.StatusBadRequest)

			logger.LogActionPanel(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Post bulunamadı ID: %s", idStr), false)
			return
		}
		helpers.SuccessResponse(w, "İstenilen id başarıyla getirildi", posts, http.StatusOK)

		logger.LogActionPanel(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Postman üzerinden veri ID'si arandı: %s", idStr), true)

	}
}

func (postman Postman) HandlePostDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer func() {
		if err := recover(); err != nil {
			logger.LogActionPanel(0, "unknown", fmt.Sprintf("Panic occurred: %v", err), false)
			helpers.ErrorResponse(w, "Beklenmedik bir hata oluştu", "500003", http.StatusBadRequest)

		}
	}()

	currentUser := helpers.GetUserFromSessionPanel(r)
	if !helpers.CheckUserPanel(w, r) {
		return
	}

	if r.Method == http.MethodDelete {
		// URL parametresinden ID'yi al (örneğin /api/category?id=8)
		idStr := ps.ByName("id") // URL'deki :id parametresini al
		id, err := strconv.Atoi(idStr)
		if err != nil {
			helpers.ErrorResponse(w, "Geçersiz ID", "400001", http.StatusBadRequest)

			logger.LogActionPanel(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Geçersiz ID, Post silinemedi: %s", idStr), false)
			return
		}

		// Veritabanından silme işlemi
		err = models.CustomPostQuery{}.DeletePostByID(id)
		if err != nil {
			// Hata oluşursa panic tetiklenir
			logger.LogSaveMessagePanel(fmt.Sprint("Post silinemedi: "+err.Error()), false)
			panic("Post silinemedi: " + err.Error())
		}

		logger.LogActionPanel(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Postman üzerinden veri silindi, ID: %s", idStr), true)
		responseData := map[string]interface{}{
			"Silinen id": id,
		}
		helpers.SuccessResponse(w, "Veri başarıyla  silindi", responseData, http.StatusOK)

	}
}

func (postman Postman) HandleGetCategory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer func() {
		if err := recover(); err != nil {
			logger.LogActionPanel(0, "unknown", fmt.Sprintf("Panic occurred: %v", err), false)
			helpers.ErrorResponse(w, "Beklenmedik bir hata oluştu", "500003", http.StatusBadRequest)

		}
	}()

	currentUser := helpers.GetUserFromSessionPanel(r)
	if !helpers.CheckUserPanel(w, r) {
		return
	}

	if r.Method == http.MethodGet {
		// Veritabanından tüm kategorileri al
		posts, err := models.CustomPostQuery{}.GetAllCategories()
		if err != nil {
			// Hata durumunda panic tetikleniyor
			logger.LogSaveMessagePanel(fmt.Sprint("Veriler alınamadı: "+err.Error()), false)
			panic("Veriler alınamadı: " + err.Error())
		}
		helpers.SuccessResponse(w, "Kategoriler başarıyla geldi", posts, http.StatusOK)
		// Başarılı bir durumda response'u güncelle

		logger.LogActionPanel(int(currentUser.ID), currentUser.Username, "Postman üzerinden kategoriler çağrıldı", true)

	}
}

func (postman Postman) HandlePostCategoriesDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer func() {
		if err := recover(); err != nil {
			logger.LogActionPanel(0, "unknown", fmt.Sprintf("Panic occurred: %v", err), false)
			helpers.ErrorResponse(w, "Beklenmedik bir hata oluştu", "500003", http.StatusBadRequest)

		}
	}()

	currentUser := helpers.GetUserFromSessionPanel(r)
	if !helpers.CheckUserPanel(w, r) {
		return
	}

	if r.Method == http.MethodDelete {
		// URL parametresinden ID'yi al
		idStr := ps.ByName("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			helpers.ErrorResponse(w, "Geçersiz ID", "400001", http.StatusBadRequest)

			logger.LogActionPanel(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Postman üzerinden kategori silme isteği başarısız oldu, ID: %d", id), false)
			return
		}

		// Veritabanından silme işlemi
		err = models.CustomPostQuery{}.DeleteCategoryByID(id)
		if err != nil {
			helpers.ErrorResponse(w, "Kategori silinemedi", "500100", http.StatusBadRequest)
			logger.LogSaveMessagePanel(fmt.Sprint("Kategori silinemedi: "+err.Error()), false)
			panic("Kategori silinemedi: " + err.Error())
		}

		// Başarılı silme işlemi
		logger.LogActionPanel(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Postman üzerinden kategori silindi ID'si: %s", idStr), true)
		// JSON formatında başarı mesajı döndür
		helpers.SuccessResponse(w, "Kategori başarıyla silindi. ID", id, http.StatusOK)

	}
}

func (postman Postman) HandlePostCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer func() {
		if err := recover(); err != nil {
			logger.LogActionPanel(0, "unknown", fmt.Sprintf("Panic occurred: %v", err), false)
			helpers.ErrorResponse(w, "Beklenmedik bir hata oluştu", "500003", http.StatusBadRequest)

		}
	}()

	currentUser := helpers.GetUserFromSessionPanel(r)
	if !helpers.CheckUserPanel(w, r) {
		return
	}

	if r.Method == http.MethodPost {
		var postman Postman
		err := json.NewDecoder(r.Body).Decode(&postman)
		if err != nil {
			fmt.Println(err)
			helpers.ErrorResponse(w, "Geçersiz veri formatı", "400001", http.StatusBadRequest)

			logger.LogActionPanel(int(currentUser.ID), currentUser.Username, "Geçersiz veri formatı", false)
			return
		}

		// Gerekli alanların boş olup olmadığını kontrol et
		if postman.Title == "" || postman.Description == "" || postman.Content == "" || postman.CategoryID == 0 {
			helpers.ErrorResponse(w, "Gerekli alanlar eksik: Title, Description, Content, CategoryID", "400002", http.StatusBadRequest)

			logger.LogActionPanel(int(currentUser.ID), currentUser.Username, "Postman üzerinden veri eklenemedi, gerekli alanlar eksik", false)
			return
		}

		// Postman verisini models.Post yapısına dönüştür
		post := models.Post{
			Title:       postman.Title,
			Slug:        slug.Make(postman.Title),
			Description: postman.Description,
			Content:     postman.Content,
			CategoryID:  postman.CategoryID,
			Picture_url: postman.Picture_url,
			Status:      "pending",
		}

		// Veritabanına ekleme işlemi
		idStr, err := post.Add()
		if err != nil {
			helpers.ErrorResponse(w, "Veri  ekleme hatası", "500001", http.StatusInternalServerError)

			logger.LogSaveMessagePanel("Veri eklenemedi: "+err.Error(), false)
			panic("Veri eklenemedi: " + err.Error())
		}

		// Log kaydı
		logger.LogActionPanel(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Postman üzerinden veri eklendi ID'si: %d", idStr), true)

		helpers.SuccessResponse(w, "Veri başarıyla eklendi", idStr, http.StatusOK)

	}
}

func (postman Postman) HandlePostStatusUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer func() {
		if err := recover(); err != nil {
			logger.LogActionPanel(0, "unknown", fmt.Sprintf("Panic occurred: %v", err), false)
			helpers.ErrorResponse(w, "Beklenmeyen bir hata oluştu", "500003", http.StatusInternalServerError)

		}
	}()

	currentUser := helpers.GetUserFromSessionPanel(r)
	if !helpers.CheckUserPanel(w, r) {
		return
	}

	if r.Method == http.MethodPut {

		// URL'den ID'yi al ve dönüştür
		idStr := ps.ByName("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			helpers.ErrorResponse(w, "Geçersiz ID", "400001", http.StatusBadRequest)

			logger.LogActionPanel(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Geçersiz ID: %s", idStr), false)
			return
		}

		// İstek gövdesinden status bilgisini al
		var reqBody struct {
			Status string `json:"status"`
		}
		err = json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil || reqBody.Status == "" {
			helpers.ErrorResponse(w, "Geçersiz veri", "400002", http.StatusBadRequest)

			logger.LogActionPanel(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Geçersiz veri, ID: %s", idStr), false)
			return
		}

		// Veritabanında post statüsünü güncelle
		err = models.CustomPostQuery{}.UpdatePostStatus(id, reqBody.Status)
		if err != nil {
			helpers.ErrorResponse(w, "Post statüsü güncellenemedi", "500002", http.StatusBadRequest)

			logger.LogActionPanel(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Post statüsü güncellenemedi, ID: %d, Hata: %s", id, err.Error()), false)
			return
		}

		// Başarılı yanıt
		responseData := map[string]interface{}{
			"id":     id,
			"status": reqBody.Status,
		}
		helpers.SuccessResponse(w, "Post statüsü başarıyla güncellendi", responseData, http.StatusOK)

		// Log ve yanıt
		logger.LogActionPanel(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Post statüsü başarıyla güncellendi, ID: %d", id), true)
	}
}
