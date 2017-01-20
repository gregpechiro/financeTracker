package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var tmpl *web.TmplCache
var mux *web.Mux
var db *adb.DB = adb.NewDB()

func init() {
	db.AddStore("user")
	db.AddStore("transaction")
	db.AddStore("category")

	web.SESSDUR = 15 * time.Minute
	mux = web.NewMux()

	mux.AddRoutes()

	tmpl = web.NewTmplCache()
}

func main() {
	fmt.Println("--------------------did you register all your routes you dumbass")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
