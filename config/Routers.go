package config

import (
	admin "goweb_7/admin/controllers"
	site "goweb_7/site/controllers"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Routers() *httprouter.Router {
	r := httprouter.New()
	//ADMIN

	//Postman
	// Özel Method Not Allowed Handler'ı ayarlayın
	r.MethodNotAllowed = http.HandlerFunc(admin.Postman{}.MethodNotAllowedHandler)

	r.Handle("POST", "/panel/giriş", admin.Postman{}.LoginHandler)
	r.Handle("POST", "/panel/çıkış", admin.Postman{}.LogoutHandler)
	r.Handle("GET", "/panel/postlar", admin.Postman{}.HandlePostAdmin)
	r.Handle("GET", "/panel/post/:id", admin.Postman{}.HandlePostAdminID)
	r.Handle("DELETE", "/panel/delete/:id", admin.Postman{}.HandlePostDelete)
	r.Handle("POST", "/panel/post_ekle", admin.Postman{}.HandlePostCreate)
	r.Handle("PUT", "/panel/posts/:id/status", admin.Postman{}.HandlePostStatusUpdate)
	r.Handle("GET", "/panel/kategori", admin.Postman{}.HandleGetCategory)
	r.Handle("DELETE", "/panel/kategori/sil/:id", admin.Postman{}.HandlePostCategoriesDelete)

	//Blog post
	r.GET("/admin", admin.Dashboard{}.Index)
	r.GET("/admin/yeni-ekle", admin.Dashboard{}.NewItem)
	r.POST("/admin/add", admin.Dashboard{}.Add)
	r.GET("/admin/delete/:id", admin.Dashboard{}.Delete)
	r.GET("/admin/edit/:id", admin.Dashboard{}.Edit)
	r.POST("/admin/update/:id", admin.Dashboard{}.Update)
	r.GET("/admin/search", admin.Dashboard{}.Search)

	//Categories
	r.GET("/admin/kategoriler", admin.Categories{}.Index)
	r.POST("/admin/kategoriler/add", admin.Categories{}.Add)
	r.GET("/admin/kategoriler/delete/:id", admin.Categories{}.Delete)

	//Onay Formu
	r.GET("/admin/onayformu", admin.Approval{}.Index)
	r.GET("/admin/review/:id", admin.Approval{}.Show)
	r.POST("/admin/update-review/:id", admin.Approval{}.Update)
	r.GET("/admin/update-review/:id", admin.Approval{}.Update)

	//Userops
	r.GET("/admin/login", admin.Userops{}.Index)
	r.POST("/admin/do_login", admin.Userops{}.Login)
	r.GET("/admin/logout", admin.Userops{}.Logout)

	//SITE
	//Homepage
	r.GET("/", site.Homepage{}.Index)
	r.GET("/yazilar/:slug", site.Homepage{}.Detail)

	//SERVE FILES
	r.ServeFiles("/admin/assests/*filepath", http.Dir("admin/assests"))
	r.ServeFiles("/uploads/*filepath", http.Dir("uploads/"))

	return r

}
