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
		"transactions":      transactions,
		"user":              user,
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

	tmpl.Render(w, r, "transaction.tmpl", web.Model{
		"transactions": transactions,
		"user":         user,
	})

	return
}}

var transactionSave = web.Route{"POST", "/transaction", func(w http.ResponseWriter, r *http.Request) {

	rUrl, err := url.Parse(r.Referer())
	redirect := "/dashboard"
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
	r.ParseForm()
	if errs, ok := web.FormToStruct(&transaction, r.Form, "transaction"); !ok {
		web.SetFormErrors(w, errs)
		web.SetErrorRedirect(w, r, redirect, "Error saving transaction")
		return
	}

	t, err := time.Parse("1/2/2006", r.FormValue("dateString"))
	if err != nil {
		web.SetErrorRedirect(w, r, redirect, "Error getting date")
		return
	}

	transaction.Date = t.UnixNano()
	transaction.Id = genId()
	transaction.UserId = user.Id
	transaction.AccountId = user.AccountId
	transaction.DateAdded = time.Now().UnixNano()

	/*if transaction.Category != "" {
		transaction.Category = NormalIzeString(transaction.Category)
		user.Categories[transaction.Category] = struct{}{}
		category := Category{
			Id:            genId(),
			AccountId:     user.AccountId,
			Title:         transaction.Category,
			TransactionId: transaction.Id,
		}

		db.Add("category", category.Id, category)
		transaction.CategoryId = category.Id
	}

	if transaction.SecondaryCategory != "" {
		transaction.SecondaryCategory = NormalIzeString(transaction.SecondaryCategory)
		user.Categories[transaction.SecondaryCategory] = struct{}{}
		secondaryCategory := Category{
			Id:            genId(),
			AccountId:     user.AccountId,
			Title:         transaction.SecondaryCategory,
			TransactionId: transaction.Id,
		}

		db.Add("category", secondaryCategory.Id, secondaryCategory)
		transaction.SecondaryCategoryId = secondaryCategory.Id
	}*/

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

	if !db.Add("transaction", transaction.Id, transaction) {
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
	web.SetSuccessRedirect(w, r, redirect, "Transaction Added")
	return
}}

/*var transactionFilter = web.Route{"POST", "/transaction/filter", func(w http.ResponseWriter, r *http.Request) {

	resp := map[string]interface{}{}

	tme := r.FormValue("time")
	aType := r.FormValue("aType")
	category1 := r.FormValue("category1")
	category2 := r.FormValue("category2")

	var fullQuery [][]byte

	fullQuery = append(fullQuery, adb.Eq("key", val))

	switch aType {
	case "income":
		fullQuery = append(fullQuery, adb.Gt("ammount", "0"))
	case "expense":
		fullQuery = append(fullQuery, adb.Lt("ammount", "0"))
	}

	var beg, end int64
	loc, _ := time.LoadLocation("Local")
	now := time.Now()
	end = now.Unix() + 1
	switch tme {
	case "week":

	case "":
		begDate := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, loc)
		beg = begDate.Unix() - 1
	case "3months":
		begDate := time.Date(now.Year(), now.Month()-3, 1, 0, 0, 0, 0, loc)
		beg = begDate.Unix() - 1
	case "6months":
		begDate := time.Date(now.Year(), now.Month()-6, 1, 0, 0, 0, 0, loc)
		beg = begDate.Unix() - 1
	case "12months":
		begDate := time.Date(now.Year()-1, now.Month(), 1, 0, 0, 0, 0, loc)
		beg = begDate.Unix() - 1
	case "all":
		beg = 0
	}

	var transactions []Transaction
	db.TestQuery("transaction", &transaction, adb.Gt("date", strconv.Itoa(beg)), adb.Lt("date", strconv.ItoA(end)))

}}*/
