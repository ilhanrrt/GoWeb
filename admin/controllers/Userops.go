package controllers

import (
	"crypto/sha256"
	"fmt"
	"goweb_7/admin/helpers"
	logger "goweb_7/admin/log"
	"goweb_7/admin/models"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Userops struct{}

func (userops Userops) Index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	view, err := template.ParseFiles(helpers.Include("userops/login")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["Alert"] = helpers.GetAlert(w, r)
	view.ExecuteTemplate(w, "index", data)
}

func (userops Userops) Login(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	username := r.FormValue("username")
	password := fmt.Sprintf("%x", sha256.Sum256([]byte(r.FormValue("password"))))

	user := models.User{}.Get("username = ? AND password = ?", username, password)

	if user.ID != 0 { // Kullanıcı bulundu mu kontrolü
		err := helpers.SetUser(w, r, user)
		if err != nil {
			helpers.SetAlert(w, r, "Oturum başlatılamadı.")
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		if user.Username == username && user.Password == password {
			fmt.Println(user.UserType)
			helpers.SetUser(w, r, user)
			switch user.UserType {
			case "admin":
				logger.LogAction(int(user.ID), user.UserType, "Giriş yaptı")
				helpers.SetAlert(w, r, "Hoşgeldiniz, Admin!")
			case "editor":
				logger.LogAction(int(user.ID), user.UserType, "Giriş yaptı")
				helpers.SetAlert(w, r, "Hoşgeldiniz, Editör!")
			case "viewer":
				logger.LogAction(int(user.ID), user.UserType, "Giriş yaptı")
				helpers.SetAlert(w, r, "Hoşgeldiniz, İzleyici!")
			default:
				helpers.SetAlert(w, r, "Hoşgeldiniz!")
			}
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		}

	} else {
		helpers.SetAlert(w, r, "Yanlış Kullanıcı Adı veya Şifre")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}
}

func (userops Userops) Logout(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentUser := helpers.GetUserFromSession(r) // Oturumdaki kullanıcıyı al
	if currentUser.Username != "" {
		logger.LogAction(int(currentUser.ID), currentUser.Username, "Çıkış yaptı") // Çıkış yapanı logla
	}
	helpers.RemoveUser(w, r)
	helpers.SetAlert(w, r, "Hoşçakalın")
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)

}
