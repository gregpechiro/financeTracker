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

	mux.AddRoutes(home, login, logout, loginPost, register)
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
	var user User
	r.ParseForm()
	if errs, ok := web.FormToStruct(&user, r.Form, "login"); !ok {
		web.SetFormErrors(w, errs)
		web.SetErrorRedirect(w, r, "/login", "Error Loging In")
		return
	}
	if !db.Auth("user", user.Email, user.Password, &user) {
		web.SetErrorRedirect(w, r, "/login", "Incorrect email or password")
		return
	}

	sess := web.Login(w, r, user.Role)
	sess.PutId(w, user.Id)
	sess["email"] = user.Email
	web.PutMultiSess(w, r, sess)
	user.LastSeen = time.Now().Unix()
	db.Set("user", user.Id, user)

	web.SetSuccessRedirect(w, r, "/account", "Welcome "+user.FirstName)
	return

}}

var register = web.Route{"POST", "/register", func(w http.ResponseWriter, r *http.Request) {
	var user User
	r.ParseForm()
	if errs, ok := web.FormToStruct(&user, r.Form, "register"); !ok {
		web.SetFormErrors(w, errs)
		web.SetErrorRedirect(w, r, "/login", "Error Registering")
		return
	}

	var users []User
	db.TestQuery("user", &users, adb.Eq("email", user.Email), adb.Ne("id", `"`+user.Id+`"`))
	if len(users) > 0 {
		web.SetErrorRedirect(w, r, "/login", "Error Registering, Email already exists")
		return
	}

	user.Id = genId()
	user.Active = true
	user.Role = "USER"
	user.Created = time.Now().Unix()
	user.Primary = true

	if !db.Add("user", user.Id, user) {
		web.SetErrorRedirect(w, r, "/login", "Error Registering, Please try again")
		return
	}
	web.SetSuccessRedirect(w, r, "/login", "Registered, Please Login")
	return

}}

var logout = web.Route{"GET", "/logout", func(w http.ResponseWriter, r *http.Request) {
	web.Logout(w)
	web.SetSuccessRedirect(w, r, "/login", "Goodbye")
	return
}}
