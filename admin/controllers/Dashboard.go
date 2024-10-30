package controllers

import (
	"encoding/json"
	"fmt"
	"goweb_7/admin/helpers"
	logger "goweb_7/admin/log"
	"goweb_7/admin/models"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gosimple/slug"
	"github.com/julienschmidt/httprouter"
)

type Dashboard struct{}

func (dashboard Dashboard) Index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentUser := helpers.GetUserFromSession(r)
	if !helpers.CheckUser(w, r) {
		return
	}

	logger.LogAction(int(currentUser.ID), currentUser.Username, fmt.Sprintln("Blog yazıları formuna giriş yaptı"))

	// Filtreleme kriterlerini alıyoruz
	categoryIDStr := r.URL.Query().Get("category_id")
	status := r.URL.Query().Get("status")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	var categoryID int
	if categoryIDStr != "" {
		categoryID, _ = strconv.Atoi(categoryIDStr)
	}

	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	// Filtrelenmiş gönderileri al
	posts, err := models.Post{}.GetFilteredPosts(categoryID, status, startDate, endDate)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not retrieve posts", http.StatusInternalServerError)
		return
	}

	view, err := template.New("index").Funcs(template.FuncMap{
		"getCategory": func(categoryID int) string {
			return models.Category{}.Get(categoryID).Title
		},
	}).ParseFiles(helpers.Include("/dashboard/list")...)
	if err != nil {
		fmt.Println(err)
		return
	}

	data := make(map[string]interface{})
	data["Posts"] = posts
	data["Alert"] = helpers.GetAlert(w, r)
	data["Categories"] = models.Category{}.GetAll()
	data["SelectedCategoryID"] = categoryID
	data["SelectedStatus"] = status
	data["StartDate"] = startDate
	data["EndDate"] = endDate

	view.ExecuteTemplate(w, "index", data)
}

func (dashboard Dashboard) NewItem(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentUser := helpers.GetUserFromSession(r)
	if !helpers.CheckUser(w, r) {
		return
	}
	if currentUser.UserType == "viewer" {
		helpers.SetAlert(w, r, "Görüntüleyici kullanıcılar bu işlemi yapamaz.")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	logger.LogAction(int(currentUser.ID), currentUser.Username, fmt.Sprintln("Yeni veri eklemek için giriş yaptı"))

	view, err := template.ParseFiles(helpers.Include("dashboard/add")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["Categories"] = models.Category{}.GetAll()
	view.ExecuteTemplate(w, "index", data)
}

func (dashboard Dashboard) Add(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentUser := helpers.GetUserFromSession(r)
	if !helpers.CheckUser(w, r) {
		return
	}
	if currentUser.UserType == "viewer" {
		helpers.SetAlert(w, r, "Görüntüleyici kullanıcılar bu işlemi yapamaz.")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	title := r.FormValue("blog-title")
	slug := slug.Make(title)
	description := r.FormValue("blog-desc")
	categoryID, _ := strconv.Atoi(r.FormValue("blog-category"))
	content := r.FormValue("blog-content")

	r.ParseMultipartForm(10 << 20)
	file, header, err := r.FormFile("blog-picture")
	if err != nil {
		fmt.Println(err)
		return
	}
	f, err := os.OpenFile("uploads/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Hangi durumun atanacağını belirleyin
	status := "approved" // Varsayılan olarak approved
	if currentUser.UserType == "editor" {
		status = "pending" // Eğer kullanıcı editörse pending
	}

	// Yeni gönderiyi veritabanına ekleyin
	pos := models.Post{
		Title:       title,
		Slug:        slug,
		Description: description,
		CategoryID:  categoryID,
		Content:     content,
		Picture_url: "uploads/" + header.Filename,
		Status:      status,
	}
	postID, err := pos.Add() // Post ID'sini ve hatayı alıyoruz
	if err != nil {
		fmt.Println("Failed to add post:", err)
		helpers.SetAlert(w, r, "Kayıt Eklenemedi")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	logger.LogAction(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Yeni veri eklendi: %s, ID: %d", title, postID))

	helpers.SetAlert(w, r, "Kayıt Başarıyla Eklendi")
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (dashboard Dashboard) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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
	fmt.Println(currentUser.UserType)
	post := models.Post{}.Get(params.ByName("id"))
	post.Delete()
	title := post.Title
	postID := post.ID

	logger.LogAction(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Kayıt Silindi: %s, ID: %d", title, postID))

	helpers.SetAlert(w, r, "Kayıt Başarıyla Silindi")
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (dashboard Dashboard) Edit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentUser := helpers.GetUserFromSession(r)
	if !helpers.CheckUser(w, r) {
		return
	}
	if currentUser.UserType == "viewer" {
		helpers.SetAlert(w, r, "Görüntüleyici kullanıcılar bu işlemi yapamaz.")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	view, err := template.ParseFiles(helpers.Include("/dashboard/edit")...)
	if err != nil {
		fmt.Println()
		return
	}
	id := params.ByName("id")

	logger.LogAction(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Veri Güncellemesi için detay bakıldı:ID %s", id))

	data := make(map[string]interface{})
	data["Post"] = models.Post{}.Get(params.ByName("id"))
	data["Categories"] = models.Category{}.GetAll()
	view.ExecuteTemplate(w, "index", data)
}

func (dashboard Dashboard) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentUser := helpers.GetUserFromSession(r)
	if !helpers.CheckUser(w, r) {
		return
	}
	if currentUser.UserType == "viewer" {
		helpers.SetAlert(w, r, "Görüntüleyici kullanıcılar bu işlemi yapamaz.")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	post := models.Post{}.Get(params.ByName("id"))
	title := r.FormValue("blog-title")
	slug := slug.Make(title)
	description := r.FormValue("blog-desc")
	categoryID, _ := strconv.Atoi(r.FormValue("blog-category"))
	content := r.FormValue("blog-content")
	is_selected := r.FormValue("is_selected")
	var picture_url string

	if is_selected == "1" {
		r.ParseMultipartForm(10 << 20)
		file, header, err := r.FormFile("blog-picture")
		if err != nil {
			fmt.Println(err)
			return
		}
		f, err := os.OpenFile("uploads/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		io.Copy(f, file)
		picture_url = "uploads/" + header.Filename
		os.Remove(post.Picture_url)
	} else {
		picture_url = post.Picture_url
	}

	// Durum güncellemesi
	status := post.Status // Varsayılan olarak mevcut durumu alıyoruz
	if currentUser.UserType == "editor" {
		status = "pending" // Eğer kullanıcı editörse, durum "pending" olur
	}

	// Gönderiyi güncelle
	post.Updates(models.Post{
		Title:       title,
		Slug:        slug,
		Description: description,
		CategoryID:  categoryID,
		Content:     content,
		Picture_url: picture_url,
		Status:      status, // Status alanını güncelliyoruz
	})
	id := params.ByName("id")

	logger.LogAction(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Veri Güncellendi:ID %s", id))

	helpers.SetAlert(w, r, "Kayıt Başarıyla Güncellendi")
	http.Redirect(w, r, "/admin/edit/"+params.ByName("id"), http.StatusSeeOther)
}

func (dashboard Dashboard) Search(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	currentUser := helpers.GetUserFromSession(r)
	if !helpers.CheckUser(w, r) {
		return
	}
	query := r.URL.Query().Get("query")

	var posts []models.Post
	if query != "" {
		posts = models.Post{}.SearchByTitle(query)
	} else {
		posts = models.Post{}.GetAll()
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"posts": posts,
	})

	logger.LogAction(int(currentUser.ID), currentUser.Username, fmt.Sprintf("Veri ismi arandı: %s", query))

}
