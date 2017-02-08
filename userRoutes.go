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
	/*var budgetItems []BudgetItem
	db.TestQuery("budgetItems", &budgetItems, adb.Eq("accountId", `"`+user.AccountId+`"`))*/

	tmpl.Render(w, r, "account.tmpl", web.Model{
		"transactions": transactions,
		"categories":   getCategories(user.AccountId),
		"user":         user,
	})

	return
}}

var budget = web.Route{"GET", "/budget", func(w http.ResponseWriter, r *http.Request) {

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

	tmpl.Render(w, r, "budget.tmpl", web.Model{
		"categories": getCategories(user.AccountId),
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

var budgetItemSave = web.Route{"POST", "/budgetItem", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User

	// gets user and double checks that the user exists still
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}

	// err check for empty formvalue
	var budgetItem BudgetItem
	if r.FormValue("title") == "" {
		web.SetErrorRedirect(w, r, "/budget", "Error saving budget item")
		return
	}
	// assign fields
	budgetItem.Id = genId()
	budgetItem.AccountId = user.AccountId
	budgetItem.BudgetGroupId = r.FormValue("budgetGroupId")
	budgetItem.Title = r.FormValue("title")

	// save to db with err check
	if !db.Add("budgetItem", budgetItem.Id, budgetItem) {
		web.SetErrorRedirect(w, r, "/budget", "Error saving budget item")
		return
	}
	web.SetSuccessRedirect(w, r, "/budget", "Budget Item Added")

}}

var budgetGroupSave = web.Route{"POST", "/budgetGroup", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User

	// gets user and double checks that the user exists still
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}

	// parses form and throws it into a variable
	var budgetGroup BudgetGroup
	r.ParseForm()
	if errs, ok := web.FormToStruct(&budgetGroup, r.Form, "budgetGroup"); !ok {
		web.SetFormErrors(w, errs)
		web.SetErrorRedirect(w, r, "/budget", "Error saving budget group")
		return
	}
	// assign non parsed fields
	budgetGroup.Id = genId()
	budgetGroup.AccountId = user.AccountId

	// save to db with err check
	if !db.Add("budgetGroup", budgetGroup.Id, budgetGroup) {
		web.SetErrorRedirect(w, r, "/budget", "Error saving budget group")
		return
	}
	web.SetSuccessRedirect(w, r, "/budget", "Budget Group Added")

}}

var budgetGroupDel = web.Route{"POST", "/budgetGroup/:id/del", func(w http.ResponseWriter, r *http.Request) {

	var budgetGroup BudgetGroup
	var budgetItems []BudgetItem

	id := web.GetId(r)
	var user User

	// gets user and double checks that the user exists still
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}

	db.Get("budgetGroup", r.FormValue(":id"), &budgetGroup)

	if user.AccountId != budgetGroup.AccountId {
		web.SetErrorRedirect(w, r, "/budget", "Error deleting, Please try again")
		return
	}
	db.TestQuery("budgetItems", &budgetItems, adb.Eq("budgetGroupId", `"`+budgetGroup.Id+`"`))

	for _, item := range budgetItems {
		db.Del("budgetItem", item.Id)
	}

	db.Del("budgetGroup", budgetGroup.Id)

	web.SetSuccessRedirect(w, r, "/budget", "Budget Group Deleted")

	return

}}

var budgetItemDel = web.Route{"POST", "/budgetItem/:id/del", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User
	var budgetItem BudgetItem

	// gets user and double checks that the user exists still
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}

	db.Get("budgetItem", r.FormValue(":id"), &budgetItem)

	if user.AccountId != budgetItem.AccountId {
		web.SetErrorRedirect(w, r, "/budget", "Error deleting budget item, Please try again")
		return
	}

	db.Del("budgetItem", budgetItem.Id)

	web.SetSuccessRedirect(w, r, "/budget", "Budget Item Deleted")

	return

}}
