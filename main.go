package main

import (
	"fmt"
	"net/http"

	account "./code"
)

func main() {

	account.DeleteCookie()

	fs := http.FileServer(http.Dir("assets"))
	img := http.FileServer(http.Dir("img"))
	css := http.FileServer(http.Dir("css"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.Handle("/img", http.StripPrefix("/img", img))
	http.Handle("/css", http.StripPrefix("/img", css))
	http.HandleFunc("/", account.IndexHandler)
	http.HandleFunc("/login", account.LoginHandler)
	http.HandleFunc("/welcome", account.WelcomeHandler)
	http.HandleFunc("/logout", account.LogoutHandler)
	http.HandleFunc("/post", account.PostHandler)
	http.HandleFunc("/user", account.UserHandler)
	http.HandleFunc("/showPost", account.ShowHandler)
	http.HandleFunc("/delete", account.DeleteHandler)
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
