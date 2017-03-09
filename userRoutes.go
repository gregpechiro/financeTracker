package main

import (
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var dashboard = web.Route{"GET", "/dashboard", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}

	var transactions []Transaction
	db.TestQuery("transaction", &transactions, adb.Eq("accountId", `"`+user.AccountId+`"`))
	sort.Slice(transactions, func(i int, j int) bool {
		return transactions[i].Date > transactions[j].Date
	})

	if len(transactions) > 10 {
		transactions = transactions[:10]
	}

	var quickTransactions []Transaction
	db.TestQuery("quickTransaction", &quickTransactions, adb.Eq("accountId", `"`+user.AccountId+`"`))

	sort.Slice(quickTransactions, func(i int, j int) bool {
		return quickTransactions[i].Title > quickTransactions[j].Title
	})

	tmpl.Render(w, r, "dashboard.tmpl", web.Model{
		"user":              user,
		"transactions":      transactions,
		"quickTransactions": quickTransactions,
	})

	return
}}

var transaction = web.Route{"GET", "/transaction", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}

	var transactions []Transaction
	db.TestQuery("transaction", &transactions, adb.Eq("accountId", `"`+user.AccountId+`"`))

	sort.Slice(transactions, func(i int, j int) bool {
		return transactions[i].Date > transactions[j].Date
	})

	var quickTransactions []Transaction
	db.TestQuery("quickTransaction", &quickTransactions, adb.Eq("accountId", `"`+user.AccountId+`"`))

	sort.Slice(quickTransactions, func(i int, j int) bool {
		return quickTransactions[i].Title > quickTransactions[j].Title
	})

	tmpl.Render(w, r, "transaction.tmpl", web.Model{
		"user":              user,
		"transactions":      transactions,
		"quickTransactions": quickTransactions,
	})

	return
}}

var transactionSave = web.Route{"POST", "/transaction", func(w http.ResponseWriter, r *http.Request) {

	redirect := "/dashboard"
	rUrl, err := url.Parse(r.Referer())
	if err == nil {
		redirect = rUrl.Path
	}

	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, redirect, "Error retrieving user")
		return
	}

	var transaction Transaction
	db.Get("transaction", r.FormValue("id"), &transaction)

	r.ParseForm()
	if errs, ok := web.FormToStruct(&transaction, r.Form, ""); !ok {
		web.SetFormErrors(w, errs)
		web.SetErrorRedirect(w, r, redirect, "Error saving transaction")
		return
	}

	t, err := time.ParseInLocation("1/2/2006", r.FormValue("dateString"), time.Local)
	if err != nil {
		web.SetErrorRedirect(w, r, redirect, "Error getting date")
		return
	}
	transaction.Date = t.Unix()

	if transaction.Id == "" {
		transaction.Id = genId()
		transaction.UserId = user.Id
		transaction.AccountId = user.AccountId
		transaction.DateAdded = time.Now().Unix()
	}

	if user.Categories == nil {
		user.Categories = map[string]struct{}{}
	}

	if user.People == nil {
		user.People = map[string]struct{}{}
	}

	if transaction.Category != "" {
		transaction.Category = NormalIzeString(transaction.Category)
		user.Categories[transaction.Category] = struct{}{}
	}
	if transaction.SecondaryCategory != "" {
		transaction.SecondaryCategory = NormalIzeString(transaction.SecondaryCategory)
		user.Categories[transaction.SecondaryCategory] = struct{}{}
	}
	if transaction.Who != "" {
		transaction.Who = NormalIzeString(transaction.Who)
		user.People[transaction.Who] = struct{}{}
	}

	if !db.Set("transaction", transaction.Id, transaction) {
		web.SetErrorRedirect(w, r, redirect, "Error saving transaction")
		return
	}

	db.Set("user", user.Id, user)

	if r.FormValue("save") == "save" {
		var quickTransactions []Transaction
		db.TestQuery("quickTransaction", &quickTransactions, adb.Eq("title", transaction.Title))
		if len(quickTransactions) > 0 {
			web.SetErrorRedirect(w, r, redirect, "Successfully added transaction.\nFailed to save as quick transaction\nA quick transaction already exists with that title")
			return
		}
		quickTransaction := transaction
		quickTransaction.Id = genId()
		quickTransaction.DateAdded = 0
		quickTransaction.Date = 0
		if !db.Add("quickTransaction", quickTransaction.Id, quickTransaction) {
			web.SetErrorRedirect(w, r, redirect, "Successfully added transaction.\nFailed to save as quick transaction")
			return
		}
		transaction.QuickTransactionId = quickTransaction.Id
		db.Set("transaction", transaction.Id, transaction)
		web.SetSuccessRedirect(w, r, redirect, "Transaction Added and quick transaction created")
		return
	}
	web.SetSuccessRedirect(w, r, redirect, "Transaction Saved")
	return
}}

var quickTransaction = web.Route{"GET", "/quickTransaction", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}

	var quickTransactions []Transaction
	db.TestQuery("quickTransaction", &quickTransactions, adb.Eq("accountId", `"`+user.AccountId+`"`))

	sort.Slice(quickTransactions, func(i int, j int) bool {
		return quickTransactions[i].Title > quickTransactions[j].Title
	})

	tmpl.Render(w, r, "quickTransaction.tmpl", web.Model{
		"user":              user,
		"quickTransactions": quickTransactions,
	})
}}

var quickTransacitonSave = web.Route{"POST", "/quickTransaction", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/quickTransaction", "Error retrieving user")
		return
	}

	var quickTransaction Transaction
	db.Get("quickTransaction", r.FormValue("id"), &quickTransaction)

	r.ParseForm()
	if errs, ok := web.FormToStruct(&quickTransaction, r.Form, ""); !ok {
		web.SetFormErrors(w, errs)
		web.SetErrorRedirect(w, r, "/quickTransaction", "Error saving quick quickTransaction")
		return
	}

	var quickTransactions []Transaction
	db.TestQuery("quickTransaction", &quickTransactions, adb.Eq("title", quickTransaction.Title), adb.Ne("id", `"`+quickTransaction.Id+`"`))
	if len(quickTransactions) > 0 {
		web.SetErrorRedirect(w, r, "/quickTransaction", "Failed to save quick quickTransaction\nA quick quickTransaction already exists with that title")
		return
	}

	if quickTransaction.Id == "" {
		quickTransaction.Id = genId()
		quickTransaction.UserId = user.Id
		quickTransaction.AccountId = user.AccountId
		quickTransaction.DateAdded = time.Now().Unix()
	}

	if user.Categories == nil {
		user.Categories = map[string]struct{}{}
	}

	if user.People == nil {
		user.People = map[string]struct{}{}
	}

	if quickTransaction.Category != "" {
		quickTransaction.Category = NormalIzeString(quickTransaction.Category)
		user.Categories[quickTransaction.Category] = struct{}{}
	}
	if quickTransaction.SecondaryCategory != "" {
		quickTransaction.SecondaryCategory = NormalIzeString(quickTransaction.SecondaryCategory)
		user.Categories[quickTransaction.SecondaryCategory] = struct{}{}
	}
	if quickTransaction.Who != "" {
		quickTransaction.Who = NormalIzeString(quickTransaction.Who)
		user.People[quickTransaction.Who] = struct{}{}
	}

	if !db.Set("quickTransaction", quickTransaction.Id, quickTransaction) {
		web.SetErrorRedirect(w, r, "/quickTransaction", "Error saving quick transaction")
		return
	}
	if quickTransaction.Id != "" {
		var transactions []Transaction
		db.TestQuery("transaction", &transactions, adb.Eq("quickTransactionId", `"`+quickTransaction.Id+`"`))
		for _, transaction := range transactions {
			transaction.Title = quickTransaction.Title
			transaction.Description = quickTransaction.Description
			transaction.Who = quickTransaction.Who
			transaction.Category = quickTransaction.Category
			transaction.SecondaryCategory = quickTransaction.SecondaryCategory

			db.Set("transaction", transaction.Id, transaction)
		}
	}
	db.Set("user", user.Id, user)

	web.SetSuccessRedirect(w, r, "/quickTransaction", "Successfully saved quick transaction")
	return

}}
