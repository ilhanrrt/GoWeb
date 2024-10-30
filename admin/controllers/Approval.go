package controllers

import (
	"fmt"
	"goweb_7/admin/helpers"
	logger "goweb_7/admin/log"
	"goweb_7/admin/models"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Approval struct{}

func (approval Approval) Index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentUser := helpers.GetUserFromSession(r)
	if !helpers.CheckUser(w, r) {
		return
	}
	switch currentUser.UserType {
	case "admin":
		logger.LogAction(int(currentUser.ID), currentUser.Username, fmt.Sprintln("Blog onay formuna giriş yaptı"))
	case "editor":
		helpers.SetAlert(w, r, "Editör bu işlemi yapamaz.")
		logger.LogAction(int(currentUser.ID), currentUser.Username, fmt.Sprintln("Blog onay formuna giriş yapmayı denedi"))
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	case "viewer":
		helpers.SetAlert(w, r, "Görüntüleyici kullanıcılar bu işlemi yapamaz.")
		logger.LogAction(int(currentUser.ID), currentUser.Username, fmt.Sprintln("Blog onay formuna giriş yapmayı denedi"))
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
	//veri tabanı bağlantısı ve filtreleme işlemi sql kodu
	db, err := gorm.Open(mysql.Open(models.Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	var posts []models.Post
	db.Where("status = ?", "pending").Find(&posts)

	view, err := template.New("index").Funcs(template.FuncMap{
		"getCategory": func(categoryID int) string {
			return models.Category{}.Get(categoryID).Title
		},
	}).ParseFiles(helpers.Include("/dashboard/approval")...)
	if err != nil {
		fmt.Println(err)
		return
	}

	data := make(map[string]interface{})
	data["Posts"] = posts
	data["Alert"] = helpers.GetAlert(w, r)
	data["Categories"] = models.Category{}.GetAll()

	view.ExecuteTemplate(w, "index", data)
}

func (approval Approval) Show(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentUser := helpers.GetUserFromSession(r)
	if !helpers.CheckUser(w, r) {
		return
	}
	if currentUser.UserType == "viewer" {
		helpers.SetAlert(w, r, "Görüntüleyici kullanıcılar bu işlemi yapamaz.")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	if currentUser.UserType == "editor" {
		helpers.SetAlert(w, r, "Editör bu işlemi yapamaz.")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	view, err := template.ParseFiles(helpers.Include("/dashboard/review")...)
	if err != nil {
		fmt.Println()
		return
	}
	pos := models.Post{}.Get(params.ByName("id"))
	postid := pos.ID

	logger.LogAction(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Blog onay formunda ID %d veriyi görüntüledi.", postid))

	data := make(map[string]interface{})
	data["Post"] = models.Post{}.Get(params.ByName("id"))
	data["Categories"] = models.Category{}.GetAll()
	view.ExecuteTemplate(w, "index", data)
}

func (approval Approval) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentUser := helpers.GetUserFromSession(r)
	if !helpers.CheckUser(w, r) {
		return
	}
	if currentUser.UserType == "viewer" {
		helpers.SetAlert(w, r, "Görüntüleyici kullanıcılar bu işlemi yapamaz.")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	if currentUser.UserType == "editor" {
		helpers.SetAlert(w, r, "Editör bu işlemi yapamaz.")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	// ID ve action parametrelerini alıyoruz
	postID := params.ByName("id")
	// Post'u veritabanından çekiyoruz
	post := models.Post{}.Get(postID)
	action := r.URL.Query().Get("action")
	switch action {
	case "approve":
		post.Status = "approved"
		helpers.SetAlert(w, r, "Blog yazısı onaylandı.")
		logger.LogAction(int(currentUser.ID), currentUser.UserType, fmt.Sprintf("Admin tarafından onaylandı, onaylanan veri ID'si %d", post.ID))
	case "reject":
		post.Status = "rejected"
		helpers.SetAlert(w, r, "Blog yazısı reddedildi.")
		logger.LogAction(int(currentUser.ID), currentUser.UserType, fmt.Sprintf("Admin tarafından reddedildi, reddedilen veri ID'si %d", post.ID))
	}

	//Formdan gelen veriyi işleme
	action = r.FormValue("action")

	// Action'a göre post statüsünü güncelliyoruz
	switch action {
	case "approve":
		post.Status = "approved"
		helpers.SetAlert(w, r, "Blog yazısı onaylandı.")
		logger.LogAction(int(currentUser.ID), currentUser.UserType, fmt.Sprintf("Admin tarafından onaylandı, onaylanan veri ID'si %d", post.ID))
	case "reject":
		post.Status = "rejected"
		helpers.SetAlert(w, r, "Blog yazısı reddedildi.")
		logger.LogAction(int(currentUser.ID), currentUser.UserType, fmt.Sprintf("Admin tarafından reddedildi, reddedilen veri ID'si %d", post.ID))
	}

	// Diğer alanların güncellenmesi
	post.Title = r.FormValue("blog-title")
	post.Description = r.FormValue("blog-desc")
	post.CategoryID, _ = strconv.Atoi(r.FormValue("blog-category"))
	post.Content = r.FormValue("blog-content")

	// Görselin değiştirilmesi
	isSelected := r.FormValue("is_selected")
	if isSelected == "1" {
		r.ParseMultipartForm(10 << 20)
		file, header, err := r.FormFile("blog-picture")
		if err == nil {
			// Dosya kaydetme ve eski dosyayı silme işlemi
			f, err := os.OpenFile("uploads/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err == nil {
				io.Copy(f, file)
				post.Picture_url = "uploads/" + header.Filename
				os.Remove(post.Picture_url)
			}
		}
	}

	// Post'u veritabanında güncelliyoruz
	post.Updates(post)
	http.Redirect(w, r, "/admin/onayformu", http.StatusSeeOther)
}
