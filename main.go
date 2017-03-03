package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var tmpl *web.TmplCache
var mux *web.Mux
var db *adb.DB = adb.NewDB()

func init() {
	//stores
	db.AddStore("user")
	db.AddStore("transaction")
	db.AddStore("quickTransaction")
	// db.AddStore("category")

	//session timeout length
	web.SESSDUR = 15 * time.Minute
	mux = web.NewMux()

	//unsecured routes
	mux.AddRoutes(home, login, logout, loginPost, register, updateSession)

	//secured routes
	mux.AddSecureRoutes(USER, dashboard, transaction, transactionSave)

	web.Funcs["OrderCategories"] = OrderCategories
	web.Funcs["Title"] = strings.Title
	web.Funcs["PrettyDate"] = PrettyDate
	web.Funcs["isIncome"] = isIncome
	web.Funcs["toJson"] = toJson

	tmpl = web.NewTmplCache()
}

func main() {
	fmt.Println(">>> DID YOU REMEMBER TO REGISTER ANY NEW ROUTES? <<<")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

var home = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "home.tmpl", nil)
	return
}}

var updateSession = web.Route{"POST", "/updateSession", func(w http.ResponseWriter, r *http.Request) {
	return
}}
