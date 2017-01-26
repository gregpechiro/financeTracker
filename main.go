package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var tmpl *web.TmplCache
var mux *web.Mux
var db *adb.DB = adb.NewDB()

func init() {
	db.AddStore("user")
	db.AddStore("transaction")
	db.AddStore("category")

	web.SESSDUR = 15 * time.Minute
	mux = web.NewMux()

	mux.AddRoutes(home, login, logout, loginPost)
	mux.AddSecureRoutes(USER, account)

	tmpl = web.NewTmplCache()
}

func main() {
	fmt.Println("--------------------did you register all your routes you dumbass")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

var home = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "home.tmpl", nil)
	return
}}

var login = web.Route{"GET", "/login", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "login.tmpl", nil)
	return
}}

var loginPost = web.Route{"POST", "/login", func(w http.ResponseWriter, r *http.Request) {
	email, pass := r.FormValue("email"), r.FormValue("password")
	var user User
	if !db.Auth("user", email, pass, &user) {
		web.SetErrorRedirect(w, r, "/login", "Incorrect email or password")
		return
	}

	sess := web.Login(w, r, user.Role)
	sess.PutId(w, user.Id)
	sess["email"] = user.Email
	web.PutMultiSess(w, r, sess)
	web.SetSuccessRedirect(w, r, "/account", "Welcome "+user.FirstName)
	return

}}

var logout = web.Route{"GET", "/logout", func(w http.ResponseWriter, r *http.Request) {
	web.Logout(w)
	web.SetSuccessRedirect(w, r, "/login", "Goodbye")
	return
}}
