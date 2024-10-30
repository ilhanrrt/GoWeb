package helpers

import (
	"encoding/json"
	"fmt"
	"goweb_7/admin/models"
	"net/http"
)

func GetUserFromSessionPanel(r *http.Request) models.User {
	session, err := store.Get(r, "blog-user")
	if err != nil {
		fmt.Println("Session error:", err)
		return models.User{}
	}

	user := models.User{}

	if id, ok := session.Values["user_id"].(uint); ok {
		user.ID = id
	} else {
		fmt.Println("User ID not found in session")
		return models.User{}
	}

	if username, ok := session.Values["username"].(string); ok {
		user.Username = username
	} else {
		fmt.Println("Username not found in session")
		return models.User{}
	}

	if userType, ok := session.Values["usertype"].(string); ok {
		user.UserType = userType
	} else {
		fmt.Println("UserType not found in session")
		return models.User{}
	}

	fmt.Printf("User fetched from session: %+v\n", user.UserType)
	return user
}

func SetUserPanel(w http.ResponseWriter, r *http.Request, user models.User) error {
	session, err := store.Get(r, "blog-user")
	if err != nil {
		return err
	}

	session.Values["user_id"] = user.ID
	session.Values["usertype"] = user.UserType
	session.Values["username"] = user.Username
	session.Values["password"] = user.Password
	return session.Save(r, w)
}

func CheckUserPanel(w http.ResponseWriter, r *http.Request) bool {
	session, err := store.Get(r, "blog-user")
	if err != nil {
		return false
	}
	username := session.Values["username"]
	password := session.Values["password"]
	user := models.User{}.Get("username=? AND password=?", username, password)

	if user.Username == username && user.Password == password {
		return true
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Lütfen giriş yapın")
	return false
}

func RemoveUserPanel(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, "blog-user")
	if err != nil {
		return err
	}

	session.Options.MaxAge = -1

	return session.Save(r, w)
}
