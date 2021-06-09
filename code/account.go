package account

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

/*models*/
type dbLogin struct { //structure of the database of the login
	Id       int
	Pseudo   string
	Email    string
	Password string
}

type Post struct { //structure of the database of the Post
	Id      int
	Title   string
	Content string
	Author  string
	Filters string
}
type PostStats struct { //structure of the database of the PostStats
	PostId  int
	Pseudo  string
	Like    bool
	Comment string
}

var u dbLogin
var p Post
var store = sessions.NewCookieStore([]byte("mysession"))

var data = make(map[string]interface{})

/*database*/

//InitDatabase will create the database when we run the main.go we use it in the func Index
func InitDatabase(database string) *sql.DB {
	db, err := sql.Open("sqlite3", database)

	if err != nil {
		log.Fatal(err)
	}
	statement := ``
	if database == "dblogin.db" {
		statement = `CREATE TABLE IF NOT EXISTS dblogin (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, pseudo TEXT NOT NULL, email TEXT NOT NULL, password TEXT NOT NULL)`
	} else if database == "post.db" {
		statement = `CREATE TABLE IF NOT EXISTS post (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, title TEXT NOT NULL, content TEXT NOT NULL, author TEXT NOT NULL, filter TEXT NOT NULL)`
	} else if database == "postStats.db" {
		statement = `CREATE TABLE IF NOT EXISTS postStats (id INTEGER NOT NULL, pseudo TEXT NOT NULL, like INTEGER NOT NULL, comment TEXT NOT NULL)`
	}
	_, err = db.Exec(statement)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

//CompareInputAndDB will compare the input email and the email inside the databse, it's also do it for the password
func CompareInputAndDB(email string, password string) bool {
	// Open the database
	database, err := sql.Open("sqlite3", "./dblogin.db")
	CheckError(err)
	defer database.Close()

	// browse the database
	rows, err := database.Query("SELECT id, pseudo, email, password FROM dblogin")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	//create a loop which will scan the content of the database
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Pseudo, &u.Email, &u.Password)
		CheckError(err)

		if u.Email == email {
			fmt.Println("Email bien trouvé dans la BDD")
			if CheckPasswordHash(u.Password, []byte(password)) {
				fmt.Println("Bon MDP")
				data["userConnected"] = u.Pseudo // Register the username of the mail
				fmt.Println(data["userConnected"])
				return true

			} else {
				fmt.Println("Pas le bon MDP")
			}
			break
		} else {
			fmt.Println(u.Email, "ne correspond pas à", email)
		}
	}
	return false
}

//insertIntoTypes will insert data in the table dblogin which is in dbLogin.db
func insertIntoTypes(db *sql.DB, pseudo string, email string, password string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO dblogin (pseudo, email, password) VALUES (?, ?, ?)`, pseudo, email, password)

	return result.LastInsertId()
}

//insertIntoTypesPost will insert data in the post which is in post.db
func insertIntoTypesPost(db *sql.DB, title string, content string, filter string) {
	// Open the database
	database, err := sql.Open("sqlite3", "./post.db")
	CheckError(err)
	defer database.Close()

	fmt.Println(title)
	fmt.Println(content)
	fmt.Println(filter)
	fmt.Println(data["userConnected"])

	// démarre une transaction
	tx, err := database.Begin()
	CheckError(err)
	// Prépare la transaction
	stmt, err := tx.Prepare("INSERT INTO post (title, content, author, filter) VALUES (?, ?, ?, ?)")
	CheckError(err)
	// Execute la transaction
	_, err = stmt.Exec(title, content, data["userConnected"], filter)
	CheckError(err)
	// Commit la transaction
	tx.Commit()
}

//showPosts will browse the database post.db and taking the values inside it for using it later
func showPosts(data_post map[string]interface{}, id_post string) {

	var arrPost []Post
	var query string
	// Open the database
	database, err := sql.Open("sqlite3", "./post.db")
	CheckError(err)
	defer database.Close()

	if id_post == "" {
		query = "SELECT id, title, content, author, filter FROM post" //this query will show every posts
	} else {
		query = "SELECT id, title, content, author, filter FROM post WHERE id = " + id_post //this query will show one particular post
	}

	rows, err := database.Query(query)
	CheckError(err)
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&p.Id, &p.Title, &p.Content, &p.Author, &p.Filters)
		p.Content = strings.ReplaceAll(p.Content, "\r", "<br>")
		p.Content = strings.ReplaceAll(p.Content, "\n", "<br>")
		arrPost = append(arrPost, p)

	}

	data_post["showPost"] = arrPost

}

// func showDBLogin(data_login map[string]interface{}, id_login string) {

// 	var arrLogin []dbLogin
// 	var query string
// 	// Open the database
// 	database, err := sql.Open("sqlite3", "./dblogin.db")
// 	CheckError(err)
// 	defer database.Close()

// 	// Parcourir la BDD
// 	query = "SELECT id, pseudo, email, password FROM post WHERE id = " + id_login

// 	rows, err := database.Query(query)
// 	CheckError(err)
// 	defer rows.Close()

// 	for rows.Next() {
// 		rows.Scan(&u.Id, &u.Pseudo, &u.Email, &u.Password)

// 		arrLogin = append(arrLogin, u)

// 	}
// 	data_login["showLogin"] = arrLogin
// 	fmt.Println(data_login)
// }

/*handlers*/

//Index will launch the firstpage which is localhost://8080
func IndexHandler(response http.ResponseWriter, request *http.Request) {
	InitDatabase("dbLogin.db")
	InitDatabase("post.db")
	InitDatabase("postStats.db")

	session, _ := store.Get(request, "mysession")
	session.Options.MaxAge = -1
	username := session.Values["username"]
	data_index := map[string]interface{}{
		"username": username,
	}
	data_index["cookieExist"] = data["cookieExist"]

	showPosts(data_index, "")
	tmp, _ := template.ParseFiles("html/welcome.html")
	tmp.Execute(response, data_index)
}

//WelcomeHandler is same as Index, but it won't run the database and also it used as a redirection
func WelcomeHandler(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "mysession")
	username := session.Values["username"]
	data_welcome := map[string]interface{}{
		"username": username,
	}
	data_welcome["cookieExist"] = data["cookieExist"]

	showPosts(data_welcome, "") //showPosts will show every posts made

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
		passwordEncrypted := HashPassword(passwordSignup)
		db := InitDatabase("dblogin.db")
		defer db.Close()

		insertIntoTypes(db, pseudoSignup, emailSignup, passwordEncrypted)

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

	session.Save(request, response)
	http.Redirect(response, request, "/", http.StatusSeeOther)
}

//UserHanler is used for showing the user profile.
func UserHandler(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "mysession")
	username := session.Values["username"]
	data_user := map[string]interface{}{
		"username": username,
	}
	data_user["cookieExist"] = data["cookieExist"]
	if data["cookieExist"] == false {
		http.Redirect(response, request, "/login", http.StatusSeeOther)
	}

	fmt.Println("Va sur user.html")
	tmp, _ := template.ParseFiles("html/user.html")
	tmp.Execute(response, data_user)

}

//Posthandler when we create a new post
func PostHandler(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "mysession")
	title := request.FormValue("title")
	content := request.FormValue("content")
	filter := request.FormValue("filter")
	username := session.Values["username"]
	data_page := map[string]interface{}{
		"username": username,
	}
	data_page["cookieExist"] = data["cookieExist"]
	if data["cookieExist"] == false {
		http.Redirect(response, request, "/login", http.StatusSeeOther)
	}

	//Reading everything and then push inside the database
	if title != "" && content != "" && filter != "" {
		dbPost := InitDatabase("post.db")

		insertIntoTypesPost(dbPost, title, content, filter)
		http.Redirect(response, request, "/welcome", http.StatusSeeOther)
	}

	t, _ := template.ParseFiles("html/post.html")
	t.Execute(response, data_page)
}

//ShowHandlers will show one particuliar post by getting the id
func ShowHandler(response http.ResponseWriter, request *http.Request) {

	session, _ := store.Get(request, "mysession")
	username := session.Values["username"]

	data_post := map[string]interface{}{
		"username": username,
	}
	data_post["cookieExist"] = data["cookieExist"]

	id_post := request.FormValue("id-post")
	fmt.Println(id_post)
	if id_post != "" {

		showPosts(data_post, id_post)

	}

	tmp, _ := template.ParseFiles("html/showPost.html")
	tmp.Execute(response, data_post)
}

/*Encrypt*/

//CheckPasswordhash will compare a hashed password to the plain password
func CheckPasswordHash(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	CheckError(err)

	return true
}

//HashPassword will encrypt the password
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println(err)
	}
	return string(bytes)
}

/*utils*/

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
