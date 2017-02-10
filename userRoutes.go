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
		"transactions":  transactions,
		"categoryViews": getCategoryView(user.AccountId),
		"user":          user,
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
		"categoryViews": getCategoryView(user.AccountId),
		"user":          user,
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
		web.SetErrorRedirect(w, r, "/category", "Error saving subcategory")
		return
	}
	// assign fields
	subcategory.Id = genId()
	subcategory.AccountId = user.AccountId
	subcategory.CategoryId = r.FormValue("categoryId")
	subcategory.Title = r.FormValue("title")

	// save to db with err check
	if !db.Add("subcategory", subcategory.Id, subcategory) {
		web.SetErrorRedirect(w, r, "/category", "Error saving subcategory")
		return
	}
	web.SetSuccessRedirect(w, r, "/category", "Subcategory Added")

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
		web.SetErrorRedirect(w, r, "/category", "Error saving category")
		return
	}
	// assign non parsed fields
	category.Id = genId()
	category.AccountId = user.AccountId

	// save to db with err check
	if !db.Add("category", category.Id, category) {
		web.SetErrorRedirect(w, r, "/category", "Error saving category")
		return
	}
	web.SetSuccessRedirect(w, r, "/category", "Category Added")

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

	// gets category from id in url
	db.Get("category", r.FormValue(":id"), &category)

	// err checks for ownership of account
	if user.AccountId != category.AccountId {
		web.SetErrorRedirect(w, r, "/category", "Error deleting, Please try again")
		return
	}

	// queries for subcategories with a matching categoryid
	db.TestQuery("subcategory", &subcategories, adb.Eq("categoryId", `"`+category.Id+`"`))

	// loops to delete subcateogries
	for _, item := range subcategories {
		db.Del("subcategory", item.Id)
	}

	// deletes category
	db.Del("category", category.Id)

	web.SetSuccessRedirect(w, r, "/category", "Cateogry Deleted")

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

	// gets subcategory from id in url
	db.Get("subcategory", r.FormValue(":id"), &subcategory)

	// err checks for ownership of account
	if user.AccountId != subcategory.AccountId {
		web.SetErrorRedirect(w, r, "/category", "Error deleting subcategory, Please try again")
		return
	}

	// deletes subcategory
	db.Del("subcategory", subcategory.Id)

	web.SetSuccessRedirect(w, r, "/category", "Subcategory Deleted")

	return

}}

var subcategoryRename = web.Route{"POST", "/subcategory/:id/rename", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User
	var subcategory Subcategory

	// gets user and double checks that the user exists still
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}

	// gets subcategory from id in url
	db.Get("subcategory", r.FormValue(":id"), &subcategory)

	// err checks for ownership of account
	if user.AccountId != subcategory.AccountId {
		web.SetErrorRedirect(w, r, "/category", "Error renaming subcategory, Please try again")
		return
	}

	// assigns new subcategory title
	subcategory.Title = r.FormValue("title")

	// saves changes to the db
	db.Set("subcategory", subcategory.Id, subcategory)

	web.SetSuccessRedirect(w, r, "/category", "Subcategory Renamed")

	return
}}

var categoryRename = web.Route{"POST", "/category", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User
	var category Category

	// gets user and double checks that the user exists still
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}
	// gets category from id in url
	db.Get("category", r.FormValue(":id"), &category)

	// err checks for ownership of account
	if user.AccountId != category.AccountId {
		web.SetErrorRedirect(w, r, "/category", "Error renaming category, Please try again")
		return
	}

	// assigns new category title
	category.Title = r.FormValue("title")

	// saves changes to the db
	db.Set("category", category.Id, category)

	web.SetSuccessRedirect(w, r, "/category", "Category Renamed")

	return

}}

var subcategoryMove = web.Route{"POST", "/subcategory/:id/move", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User
	var subcategory Subcategory

	// gets user and double checks that the user exists still
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}

	// gets subcategory from id in url
	db.Get("subcategory", r.FormValue(":id"), &subcategory)

	// err checks for ownership of account
	if user.AccountId != subcategory.AccountId {
		web.SetErrorRedirect(w, r, "/category", "Error moving subcategory, Please try again")
		return
	}
	// assigns new category id
	subcategory.CategoryId = r.FormValue("categoryId")

	// saves changes to the db
	db.Set("subcategory", subcategory.Id, subcategory)

	web.SetSuccessRedirect(w, r, "/category", "Subcategory Moved")

	return
}}
