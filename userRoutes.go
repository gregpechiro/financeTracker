package main

import (
	"net/http"

	"github.com/cagnosolutions/web"
)

var account = web.Route{"GET", "/account", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "account.tmpl", nil)
}}
