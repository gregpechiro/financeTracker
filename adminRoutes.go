package main

import (
	"net/http"
	"time"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var adminUsers = web.Route{"GET", "/admin/user", func(w http.ResponseWriter, r *http.Request) {
	var users []User
	db.All("user", &users)
	tmpl.Render(w, r, "adminUser.tmpl", web.Model{
		"users": users,
		"user":  User{},
	})
}}

var adminUser = web.Route{"GET", "/admin/user/:id", func(w http.ResponseWriter, r *http.Request) {
	var user User
	if !db.Get("user", r.FormValue(":id"), &user) {
		web.SetErrorRedirect(w, r, "/admin/user", "Error finding user")
		return
	}
	var users []User
	db.All("user", &users)
	tmpl.Render(w, r, "adminUser.tmpl", web.Model{
		"users": users,
		"user":  user,
	})
}}

var adminUserSave = web.Route{"POST", "/admin/user", func(w http.ResponseWriter, r *http.Request) {
	var user User
	db.Get("user", r.FormValue("id"), &user)

	r.ParseForm()
	if r.FormValue("password") == "" {
		r.Form.Set("password", user.Password)
	}
	web.FormToStruct(&user, r.Form, "adminUser")

	redirect := "/admin/user"
	if user.Id != "" {
		redirect = "/admin/user/" + user.Id
	}

	if errs, ok := web.FormToStruct(&user, r.Form, "adminUser"); !ok {
		if user.Id == "" && user.Password == "" {
			errs["adminUser.password"] = "Password is Required"
		}
		web.SetFormErrors(w, errs)
		web.SetErrorRedirect(w, r, redirect, "Error saving user")
		return
	}

	// check for uniqueness
	var users []User
	db.TestQuery("user", &users, adb.Eq("email", user.Email), adb.Ne("id", `"`+user.Id+`"`))
	if len(users) > 0 {
		web.SetErrorRedirect(w, r, redirect, "Error saving user.<br>Email is already in use.")
		return
	}

	if user.Id == "" {
		if user.Password == "" {
			web.SetFormErrors(w, map[string]string{"adminUser.password": "Password must not be empty"})
			web.SetErrorRedirect(w, r, redirect, "Error saving user")
			return
		}
		user.Id = genId()
		user.Created = time.Now().Unix()
		user.AccountId = genId()
	}

	db.Set("user", user.Id, user)
	web.SetSuccessRedirect(w, r, "/admin/user/"+user.Id, "Successfully saved user")
	return
}}
