package controllers

import (
	"fmt"
	"goweb_7/admin/helpers"
	logger "goweb_7/admin/log"
	"goweb_7/admin/models"
	"html/template"
	"net/http"

	"github.com/gosimple/slug"
	"github.com/julienschmidt/httprouter"
)

type Categories struct{}

func (categories Categories) Index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentUser := helpers.GetUserFromSession(r)
	if !helpers.CheckUser(w, r) {

		return
	}
	switch currentUser.UserType {
	case "admin":
		logger.LogAction(int(currentUser.ID), currentUser.Username, fmt.Sprintln("Admin kategoriler formuna giriş yaptı"))
	case "editor":
		logger.LogAction(int(currentUser.ID), currentUser.Username, fmt.Sprintln("Editör kategoriler formuna giriş yapmayı denedi"))
		helpers.SetAlert(w, r, "Editör bu işlemi yapamaz.")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	case "viewer":
		logger.LogAction(int(currentUser.ID), currentUser.Username, fmt.Sprintln("Viewer kategoriler formuna giriş yapmayı denedi"))
		helpers.SetAlert(w, r, "Görüntüleyici kullanıcılar bu işlemi yapamaz.")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	view, err := template.ParseFiles(helpers.Include("categories/list")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["Categories"] = models.Category{}.GetAll()
	data["Alert"] = helpers.GetAlert(w, r)
	view.ExecuteTemplate(w, "index", data)

}

func (categories Categories) Add(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentUser := helpers.GetUserFromSession(r)
	if !helpers.CheckUser(w, r) {
		return
	}
	switch currentUser.UserType {
	case "editor":
		logger.LogAction(int(currentUser.ID), currentUser.UserType, fmt.Sprintln("Editör kategoriler formuna giriş yapmayı denedi"))
		helpers.SetAlert(w, r, "Editör bu işlemi yapamaz.")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	case "viewer":
		logger.LogAction(int(currentUser.ID), currentUser.UserType, fmt.Sprintln("Viewer kategoriler formuna giriş yapmayı denedi"))
		helpers.SetAlert(w, r, "Görüntüleyici kullanıcılar bu işlemi yapamaz.")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	categoryTitle := r.FormValue("category-title")
	categorySlug := slug.Make(categoryTitle)
	categortstatus := "approved"
	categor := models.Category{Title: categoryTitle, Slug: categorySlug, Status: categortstatus}
	categorid, err := categor.Add()
	if err != nil {
		fmt.Println("Failed to add category:", err)
		helpers.SetAlert(w, r, "Kayıt Eklenemedi")
		http.Redirect(w, r, "/admin/kategoriler", http.StatusSeeOther)
		return
	}

	helpers.SetAlert(w, r, "Kayıt Başarıyla Eklendi")

	logger.LogAction(int(currentUser.ID), currentUser.UserType, fmt.Sprintf("Kategori Eklendi: %d", categorid))

	http.Redirect(w, r, "/admin/kategoriler", http.StatusSeeOther)
}

func (categories Categories) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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
	category := models.Category{}.Get(params.ByName("id"))
	category.Delete()

	helpers.SetAlert(w, r, "Kayıt Başarıyla Silindi")

	logger.LogAction(int(currentUser.ID), currentUser.UserType, fmt.Sprintf("Kategori silindi: %d", category.ID))

	http.Redirect(w, r, "/admin/kategoriler", http.StatusSeeOther)
}
