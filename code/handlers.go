package account

import (
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

/*handlers*/

//Index will launch the firstpage which is localhost://8080
func IndexHandler(response http.ResponseWriter, request *http.Request) {
	InitDatabase("post.db")
	InitDatabase("postStats.db")

	data_index := make(map[string]interface{})
	SetCookie(data_index, response, request)

	ShowPosts(data_index, "")
	tmp, _ := template.ParseFiles("html/welcome.html")
	tmp.Execute(response, data_index)
}

//WelcomeHandler is same as Index, but it won't run the database and also it used as a redirection
func WelcomeHandler(response http.ResponseWriter, request *http.Request) {
	data_welcome := make(map[string]interface{})
	SetCookie(data_welcome, response, request)

	ShowPosts(data_welcome, "") //showPosts will show every posts made

	fmt.Println(data_welcome)
	tmp, _ := template.ParseFiles("html/welcome.html")
	tmp.Execute(response, data_welcome)
}

//LoginHandler will verify everything about signup and login
func LoginHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	data["cookieExist"] = false

	pseudoSignup := request.Form.Get("signup-pseudo")
	emailSignup := request.Form.Get("signup-email")
	passwordSignup := request.Form.Get("signup-password")

	emailLogin := request.Form.Get("login-email")
	passwordLogin := request.Form.Get("login-password")

	if emailLogin != "" && passwordLogin != "" {
		fmt.Println("Email : ", emailLogin)
		fmt.Println("Password : ", passwordLogin)
		boolLogin := CompareInputAndDB(emailLogin, passwordLogin)
		if boolLogin {
			session, _ := store.Get(request, "mysession")
			session.Values["username"] = emailLogin
			session.Save(request, response)
			data["cookieExist"] = true

			http.Redirect(response, request, "/welcome", http.StatusSeeOther)
		} else {
			data["cookieExist"] = false
			http.Redirect(response, request, "/login", http.StatusSeeOther)
		}
	} else if pseudoSignup != "" && emailSignup != "" && passwordSignup != "" {

		fmt.Println("Pseudo : ", pseudoSignup)
		fmt.Println("Email : ", emailSignup)
		fmt.Println("Password : ", passwordSignup)
		data["cookieExist"] = true
		data["userConnected"] = pseudoSignup

		passwordEncrypted := HashPassword(passwordSignup)
		db := InitDatabase("dblogin.db")
		defer db.Close()

		InsertIntoTypes(db, pseudoSignup, emailSignup, passwordEncrypted)

		session, _ := store.Get(request, "mysession")
		session.Values["username"] = emailSignup
		session.Save(request, response)
		http.Redirect(response, request, "/welcome", http.StatusSeeOther)

	} else {
		data := map[string]interface{}{
			"err": "Email & Password didn't match",
		}
		tmp, _ := template.ParseFiles("html/index.html")
		tmp.Execute(response, data)
	}
}

//LogoutHandler will redirect to / and also deletes the cookies
func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "mysession")
	session.Options.MaxAge = -1
	data["cookieExist"] = false
	data["userConnected"] = "" // supprime l'utilisateur connecté
	data["already_liked"] = false

	session.Save(request, response)
	http.Redirect(response, request, "/", http.StatusSeeOther)
}

//UserHanler is used for showing the user profile.
func UserHandler(response http.ResponseWriter, request *http.Request) {
	data_user := make(map[string]interface{})
	SetCookie(data_user, response, request)
	if data["cookieExist"] == false {
		http.Redirect(response, request, "/login", http.StatusSeeOther)
	}

	fmt.Println("Va sur user.html")
	tmp, _ := template.ParseFiles("html/user.html")
	tmp.Execute(response, data_user)
}

//Posthandler when we create a new post
func PostHandler(response http.ResponseWriter, request *http.Request) {
	var data_page = make(map[string]interface{})
	SetCookie(data_page, response, request)

	// redirige si pas de cookie
	if data["cookieExist"] == false {
		http.Redirect(response, request, "/login", http.StatusSeeOther)
	}

	title := request.FormValue("title")
	content := request.FormValue("content")
	filter := request.FormValue("filter")
	//Reading everything and then push inside the database
	if title != "" && content != "" && filter != "" {
		dbPost := InitDatabase("post.db")
		InsertIntoTypesPost(dbPost, title, content, filter)
		http.Redirect(response, request, "/welcome", http.StatusSeeOther)
	}

	t, _ := template.ParseFiles("html/post.html")
	t.Execute(response, data_page)
}

//ShowHandlers will show one particuliar post by getting the id
func ShowHandler(response http.ResponseWriter, request *http.Request) {
	data_post := make(map[string]interface{})
	SetCookie(data_post, response, request)
	id_post := request.FormValue("id-post")

	if id_post != "" {
		ShowPosts(data_post, id_post)
		ShowComments(data_post, id_post)
	}

	likes := request.FormValue("like-post")
	search := SearchAllUserInPostDb(id_post, data_post)
	if likes != "" && search != true {
		AddLikes(likes)
		http.Redirect(response, request, request.Header.Get("Referer"), http.StatusFound)
	}

	comment := request.FormValue("comment")
	idPost := request.FormValue("idPost")
	if comment != "" {
		AddComment(comment, idPost, data_post)
		http.Redirect(response, request, request.Header.Get("Referer"), http.StatusFound)
	}

	fmt.Println(data_post)

	tmp, _ := template.ParseFiles("html/showPost.html")
	tmp.Execute(response, data_post)
}
