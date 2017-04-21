package main

import (
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/Shaked/mobiledetect"
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

	detect := mobiledetect.NewMobileDetect(r, nil)
	page := "dashboard.tmpl"

	if detect.IsMobile() {
		page = "mobileDashboard.tmpl"
	}

	tmpl.Render(w, r, page, web.Model{
		"user":              user,
		"transactions":      transactions,
		"quickTransactions": quickTransactions,
	})

	return
}}

var account = web.Route{"GET", "/account", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}
	tmpl.Render(w, r, "account.tmpl", web.Model{
		"user": user,
	})
}}

var accountSave = web.Route{"POST", "/account", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}
	r.ParseForm()
	if r.FormValue("password") == "" {
		r.Form.Set("password", user.Password)
	}
	web.FormToStruct(&user, r.Form, "account")

	if errs, ok := web.FormToStruct(&user, r.Form, "account"); !ok {
		web.SetFormErrors(w, errs)
		web.SetErrorRedirect(w, r, "/account", "Error updating account information")
		return
	}

	// check for uniqueness
	var users []User
	db.TestQuery("user", &users, adb.Eq("email", user.Email), adb.Ne("id", `"`+user.Id+`"`))
	if len(users) > 0 {
		web.SetErrorRedirect(w, r, "/account", "Error updating account information.<br>Email is already in use.")
		return
	}

	db.Set("user", user.Id, user)
	web.SetSuccessRedirect(w, r, "/account", "Successfully updated information")
	return
}}

var categorySave = web.Route{"POST", "/category", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}
	oldCategory := NormalIzeString(r.FormValue("oldCategory"))
	newCategory := NormalIzeString(r.FormValue("newCategory"))

	if _, ok := user.Categories[oldCategory]; !ok {
		web.SetErrorRedirect(w, r, "/account", "Error finding category")
		return
	}

	if _, ok := user.Categories[newCategory]; ok {
		web.SetErrorRedirect(w, r, "/account", "Error changing category<br>A category with that title already exists")
		return
	}

	var transactions []Transaction

	db.TestQuery("transaction", &transactions, adb.Eq("category", oldCategory))
	for _, transaction := range transactions {
		transaction.Category = newCategory
		db.Set("transaction", transaction.Id, transaction)
	}

	db.TestQuery("transaction", &transactions, adb.Eq("secondaryCategory", oldCategory))
	for _, transaction := range transactions {
		transaction.SecondaryCategory = newCategory
		db.Set("transaction", transaction.Id, transaction)
	}

	db.TestQuery("quickTransaction", &transactions, adb.Eq("category", oldCategory))
	for _, transaction := range transactions {
		transaction.Category = newCategory
		db.Set("quickTransaction", transaction.Id, transaction)
	}

	db.TestQuery("quickTransaction", &transactions, adb.Eq("secondaryCategory", oldCategory))
	for _, transaction := range transactions {
		transaction.SecondaryCategory = newCategory
		db.Set("quickTransaction", transaction.Id, transaction)
	}

	user.Categories[newCategory] = struct{}{}
	delete(user.Categories, oldCategory)
	db.Set("user", user.Id, user)
	web.SetSuccessRedirect(w, r, "/account", "Successfully updated category")
	return

}}

var categoryDel = web.Route{"POST", "/category/:name/del", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}
	category := NormalIzeString(r.FormValue(":name"))

	if _, ok := user.Categories[category]; !ok {
		web.SetErrorRedirect(w, r, "/account", "Error finding category")
		return
	}

	var transactions []Transaction

	db.TestQuery("transaction", &transactions, adb.Eq("category", category))
	if len(transactions) > 0 {
		web.SetErrorRedirect(w, r, "/account", "Error deleting category<br>It is still in use by a transaction<br>You must change the category before you can delete it")
		return
	}

	db.TestQuery("transaction", &transactions, adb.Eq("secondaryCategory", category))
	if len(transactions) > 0 {
		web.SetErrorRedirect(w, r, "/account", "Error deleting category<br>It is still in use by a transaction<br>you must change the category before you can delete it")
		return
	}

	db.TestQuery("quickTransaction", &transactions, adb.Eq("category", category))
	if len(transactions) > 0 {
		web.SetErrorRedirect(w, r, "/account", "Error deleting category<br>It is still in use by a quick transaction<br>you must change the category before you can delete it")
		return
	}

	db.TestQuery("quickTransaction", &transactions, adb.Eq("secondaryCategory", category))
	if len(transactions) > 0 {
		web.SetErrorRedirect(w, r, "/account", "Error deleting category<br>It is still in use by a quick transaction<br>you must change the category before you can delete it")
		return
	}
	delete(user.Categories, category)
	db.Set("user", user.Id, user)
	web.SetSuccessRedirect(w, r, "/account", "Successfully deleted category")
	return

}}

var whoSave = web.Route{"POST", "/who", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}
	oldWho := NormalIzeString(r.FormValue("oldWho"))
	newWho := NormalIzeString(r.FormValue("newWho"))

	if _, ok := user.People[oldWho]; !ok {
		web.SetErrorRedirect(w, r, "/account", "Error finding Who")
		return
	}

	if _, ok := user.People[newWho]; ok {
		web.SetErrorRedirect(w, r, "/account", "Error changing category<br>A who with that title already exists")
		return
	}

	var transactions []Transaction

	db.TestQuery("transaction", &transactions, adb.Eq("who", oldWho))
	for _, transaction := range transactions {
		transaction.Who = newWho
		db.Set("transaction", transaction.Id, transaction)
	}

	db.TestQuery("quickTransaction", &transactions, adb.Eq("who", oldWho))
	for _, transaction := range transactions {
		transaction.Who = newWho
		db.Set("quickTransaction", transaction.Id, transaction)
	}

	user.People[newWho] = struct{}{}
	delete(user.People, oldWho)
	db.Set("user", user.Id, user)
	web.SetSuccessRedirect(w, r, "/account", "Successfully updated who")
	return

}}

var whoDel = web.Route{"POST", "/who/:name/del", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}
	who := NormalIzeString(r.FormValue(":name"))
	if _, ok := user.People[who]; !ok {
		web.SetErrorRedirect(w, r, "/account", "Error finding who")
		return
	}

	var transactions []Transaction

	db.TestQuery("transaction", &transactions, adb.Eq("who", who))
	if len(transactions) > 0 {
		web.SetErrorRedirect(w, r, "/account", "Error deleting who<br>It is still in use by a transaction<br>You must change the who before you can delete it")
		return
	}

	db.TestQuery("quickTransaction", &transactions, adb.Eq("who", who))
	if len(transactions) > 0 {
		web.SetErrorRedirect(w, r, "/account", "Error deleting who<br>It is still in use by a quick transaction<br>You must change the who before you can delete it")
		return
	}

	delete(user.People, who)
	db.Set("user", user.Id, user)
	web.SetSuccessRedirect(w, r, "/account", "Successfully deleted who")
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

	redirect := r.FormValue("redirect")
	if redirect == "" {
		redirect = "/dashboard"
		rUrl, err := url.Parse(r.Referer())
		if err == nil {
			redirect = rUrl.Path
		}
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
			web.SetErrorRedirect(w, r, redirect, "Successfully added transaction.<br>Failed to save as quick transaction<br>A quick transaction already exists with that title")
			return
		}
		quickTransaction := transaction
		quickTransaction.Id = genId()
		quickTransaction.DateAdded = 0
		quickTransaction.Date = 0
		if !db.Add("quickTransaction", quickTransaction.Id, quickTransaction) {
			web.SetErrorRedirect(w, r, redirect, "Successfully added transaction.<br>Failed to save as quick transaction")
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
var transactionDel = web.Route{"POST", "/transaction/del/:id", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/transaction", "Error retrieving user")
		return
	}

	db.Del("transaction", r.FormValue(":id"))
	web.SetSuccessRedirect(w, r, "/transaction", "Successfully deleted transaction")
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
		web.SetErrorRedirect(w, r, "/quickTransaction", "Failed to save quick quickTransaction<br>A quick quickTransaction already exists with that title")
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
