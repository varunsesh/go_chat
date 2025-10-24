package main

import (
	"net/http"

	"github.com/go_chat/auth"
	"github.com/go_chat/views"
)

func main() {
	http.HandleFunc("/test", auth.TestHandler)
	http.HandleFunc("/debug", auth.DebugHandler)
	http.HandleFunc("/view/", views.ViewHandler)
	http.HandleFunc("/", views.LandingPageHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
