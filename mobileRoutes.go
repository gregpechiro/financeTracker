package main

import (
	"net/http"
	"sort"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var addTransaction = web.Route{"GET", "/transaction/add", func(w http.ResponseWriter, r *http.Request) {

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

	tmpl.Render(w, r, "mobileAddTransaction.tmpl", web.Model{
		"user":              user,
		"quickTransactions": quickTransactions,
	})

	return
}}
