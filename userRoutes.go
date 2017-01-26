package main

import (
	"net/http"
	"time"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var account = web.Route{"GET", "/account", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User

	//gets user and double checks that the user exists still
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}

	//gets all transactions for an account
	var transactions []Transaction
	db.TestQuery("transaction", &transactions, adb.Eq("accountId", `"`+user.AccountId+`"`))

	//gets all categories for an account
	var categories []Category
	db.TestQuery("category", &categories, adb.Eq("accountId", `"`+user.AccountId+`"`))

	tmpl.Render(w, r, "account.tmpl", web.Model{
		"transactions": transactions,
		"categories":   categories,
		"user":         user,
	})

	return
}}

var transactionSave = web.Route{"POST", "/transactionSave", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User

	//gets user and double checks that the user exists still
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}

	//parses form and throws it into a variable
	var transaction Transaction
	r.ParseForm()
	if errs, ok := web.FormToStruct(&transaction, r.Form, "transaction"); !ok {
		web.SetFormErrors(w, errs)
		web.SetErrorRedirect(w, r, "/account", "Error saving transaction")
		return
	}
	//assign non parsed fields
	transaction.Id = genId()
	transaction.UserId = user.Id
	transaction.AccountId = user.AccountId
	transaction.Date = time.Now().Unix()

	//save to db with err check
	if !db.Add("transaction", transaction.Id, transaction) {
		web.SetErrorRedirect(w, r, "/account", "Error saving transcation")
		return
	}
	web.SetSuccessRedirect(w, r, "/account", "Transaction Added")

}}
