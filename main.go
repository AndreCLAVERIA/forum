package main

import (
	"fmt"
	"net/http"

	account "./code"
)

func main() {

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", account.IndexHandler)
	http.HandleFunc("/login", account.LoginHandler)
	http.HandleFunc("/welcome", account.WelcomeHandler)
	http.HandleFunc("/logout", account.LogoutHandler)
	http.HandleFunc("/post", account.PostHandler)
	http.HandleFunc("/user", account.UserHandler)
	http.HandleFunc("/showPost", account.ShowHandler)
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
