package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var login = web.Route{"GET", "/login", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "login.tmpl", nil)
	return
}}

var loginPost = web.Route{"POST", "/login", func(w http.ResponseWriter, r *http.Request) {
	var user User
	//parses form and throws it into a variable
	r.ParseForm()
	if errs, ok := web.FormToStruct(&user, r.Form, "login"); !ok {
		fmt.Println(errs)
		web.SetFormErrors(w, errs)
		web.SetErrorRedirect(w, r, "/login", "Error Loging In")
		return
	}

	if !db.Auth("user", user.Email, user.Password, &user) {
		fmt.Println(">>> ", user.Email, user.Password)
		web.SetErrorRedirect(w, r, "/login", "Incorrect email or password")
		return
	}

	//assigns non parsed fields
	sess := web.Login(w, r, user.Role)
	sess.PutId(w, user.Id)
	sess["email"] = user.Email
	web.PutMultiSess(w, r, sess)

	//updates lastseen
	user.LastSeen = time.Now().Unix()

	//saves to db
	db.Set("user", user.Id, user)

	web.SetSuccessRedirect(w, r, "/dashboard", "Welcome "+user.FirstName)
	return

}}

var register = web.Route{"POST", "/register", func(w http.ResponseWriter, r *http.Request) {
	var user User

	//parses form and throws it into a variable
	r.ParseForm()
	if errs, ok := web.FormToStruct(&user, r.Form, "register"); !ok {
		web.SetFormErrors(w, errs)
		web.SetErrorRedirect(w, r, "/login", "Error Registering")
		return
	}

	//error checks for a user that already exists
	var users []User
	db.TestQuery("user", &users, adb.Eq("email", user.Email), adb.Ne("id", `"`+user.Id+`"`))
	if len(users) > 0 {
		web.SetErrorRedirect(w, r, "/login", "Error Registering, Email already exists")
		return
	}

	//assign non parsed fields
	user.Id = genId()
	user.Active = true
	user.Role = "USER"
	user.Created = time.Now().Unix()
	user.Primary = true
	user.AccountId = genId()

	//save to db with err check
	if !db.Add("user", user.Id, user) {
		web.SetErrorRedirect(w, r, "/login", "Error Registering, Please try again")
		return
	}
	web.SetSuccessRedirect(w, r, "/login", "Registered, Please Login")
	return

}}

var logout = web.Route{"GET", "/logout", func(w http.ResponseWriter, r *http.Request) {

	web.Logout(w)
	http.Redirect(w, r, "/login", 303)
	return
}}
