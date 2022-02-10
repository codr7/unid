package main

import (
	"github.com/codr7/unid/lib"
	"github.com/codr7/unid/lib/ui"
	"log"
	"net/http"
)

func StartWorker() {
}

func main() {
	log.Printf("Welcome to Unid v%v", unid.Version())
	http.Handle("/", http.FileServer(http.Dir("www")))

	http.HandleFunc("/login.html", ui.LoginView)
	http.HandleFunc("/login", ui.LoginHandler)
	http.HandleFunc("/logout", ui.LogoutHandler)
	
	http.HandleFunc("/rcs.html", ui.RcsView)
	log.Fatal(http.ListenAndServe(":8080", nil))	
}
