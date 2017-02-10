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

	tmpl.Render(w, r, "account.tmpl", web.Model{
		"transactions": transactions,
		"categories":   getCategoryView(user.AccountId),
		"user":         user,
	})

	return
}}

var category = web.Route{"GET", "/category", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User

	// gets user and double checks that the user exists still
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}

	// gets all budgetGroups for an account
	/*var budgetGroups []BudgetGroup
	db.TestQuery("budgetGroup", &budgetGroups, adb.Eq("accountId", `"`+user.AccountId+`"`))*/

	// gets all budgetItems for an account
	/*var budgetItems []BudgetItem
	db.TestQuery("budgetItem", &budgetItems, adb.Eq("accountId", `"`+user.AccountId+`"`))*/

	tmpl.Render(w, r, "category.tmpl", web.Model{
		"categories": getCategoryView(user.AccountId),
		"user":       user,
	})

	return
}}

var transactionSave = web.Route{"POST", "/transaction", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User

	// gets user and double checks that the user exists still
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}

	// parses form and throws it into a variable
	var transaction Transaction
	r.ParseForm()
	if errs, ok := web.FormToStruct(&transaction, r.Form, "transaction"); !ok {
		web.SetFormErrors(w, errs)
		web.SetErrorRedirect(w, r, "/account", "Error saving transaction")
		return
	}
	// assign non parsed fields
	transaction.Id = genId()
	transaction.UserId = user.Id
	transaction.AccountId = user.AccountId
	transaction.Date = time.Now().Unix()

	// save to db with err check
	if !db.Add("transaction", transaction.Id, transaction) {
		web.SetErrorRedirect(w, r, "/account", "Error saving transcation")
		return
	}
	web.SetSuccessRedirect(w, r, "/account", "Transaction Added")

}}

var subcategorySave = web.Route{"POST", "/subcategory", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User

	// gets user and double checks that the user exists still
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}

	// err check for empty formvalue
	var subcategory Subcategory
	if r.FormValue("title") == "" {
		web.SetErrorRedirect(w, r, "/category", "Error saving budget item")
		return
	}
	// assign fields
	subcategory.Id = genId()
	subcategory.AccountId = user.AccountId
	subcategory.CategoryId = r.FormValue("categoryId")
	subcategory.Title = r.FormValue("title")

	// save to db with err check
	if !db.Add("subcategory", subcategory.Id, subcategory) {
		web.SetErrorRedirect(w, r, "/category", "Error saving budget item")
		return
	}
	web.SetSuccessRedirect(w, r, "/category", "Budget Item Added")

}}

var categorySave = web.Route{"POST", "/category", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User

	// gets user and double checks that the user exists still
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}

	// parses form and throws it into a variable
	var category Category
	r.ParseForm()
	if errs, ok := web.FormToStruct(&category, r.Form, "category"); !ok {
		web.SetFormErrors(w, errs)
		web.SetErrorRedirect(w, r, "/category", "Error saving budget group")
		return
	}
	// assign non parsed fields
	category.Id = genId()
	category.AccountId = user.AccountId

	// save to db with err check
	if !db.Add("category", category.Id, category) {
		web.SetErrorRedirect(w, r, "/category", "Error saving budget group")
		return
	}
	web.SetSuccessRedirect(w, r, "/category", "Budget Group Added")

}}

var categoryDel = web.Route{"POST", "/category/:id/del", func(w http.ResponseWriter, r *http.Request) {

	var category Category
	var subcategories []Subcategory

	id := web.GetId(r)
	var user User

	// gets user and double checks that the user exists still
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}

	db.Get("category", r.FormValue(":id"), &category)

	if user.AccountId != category.AccountId {
		web.SetErrorRedirect(w, r, "/category", "Error deleting, Please try again")
		return
	}
	db.TestQuery("subcategory", &subcategories, adb.Eq("categoryId", `"`+category.Id+`"`))

	for _, item := range subcategories {
		db.Del("subcategory", item.Id)
	}

	db.Del("category", category.Id)

	web.SetSuccessRedirect(w, r, "/category", "Budget Group Deleted")

	return

}}

var subcategoryDel = web.Route{"POST", "/subcategory/:id/del", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User
	var subcategory Subcategory

	// gets user and double checks that the user exists still
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}

	db.Get("subcateogry", r.FormValue(":id"), &subcategory)

	if user.AccountId != subcategory.AccountId {
		web.SetErrorRedirect(w, r, "/category", "Error deleting budget item, Please try again")
		return
	}

	db.Del("subcategory", subcategory.Id)

	web.SetSuccessRedirect(w, r, "/category", "Budget Item Deleted")

	return

}}
