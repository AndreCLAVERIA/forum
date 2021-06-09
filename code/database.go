package account

import (
	"database/sql"

	"github.com/gorilla/sessions"
)

/*models*/
type dbLogin struct { //structure of the database of the login
	Id       int
	Pseudo   string
	Email    string
	Password string
}

// type dbLogin struct { //structure of the database of the login
// 	Id           int
// 	Pseudo       string
// 	Email        string
// 	Password     string
// 	Bio          string
// 	LinkGit      string
// 	LinkLinkedin string
// }
type Post struct { //structure of the database of the Post
	Id       int
	Title    string
	Content  string
	Author   string
	Filters  string
	Likes    int
	UserLike string
}
type Comment struct { //structure of the database of the PostStats
	Id      int
	IdPost  string
	Content string
	Author  string
}

var u dbLogin
var p Post
var c Comment
var store = sessions.NewCookieStore([]byte("mysession"))

var data = make(map[string]interface{})

//InitDatabase will create the database when we run the main.go we use it in the func Index
func InitDatabase(database string) *sql.DB {
	db, err := sql.Open("sqlite3", database)
	CheckError(err)

	statement := ``
	if database == "dblogin.db" {
		statement = `CREATE TABLE IF NOT EXISTS dblogin (
						id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
						pseudo TEXT NOT NULL, 
						email TEXT NOT NULL,
			 			password TEXT NOT NULL
					)`
	} else if database == "post.db" {
		statement = `CREATE TABLE IF NOT EXISTS post (
						id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
						title TEXT NOT NULL, content TEXT NOT NULL, 
						author TEXT NOT NULL, filter TEXT NOT NULL, 
						like INTEGER NOT NULL, 
						userlike TEXT NOT NULL
					)`
		_, err = db.Exec(statement)
		CheckError(err)
		statement = `CREATE TABLE IF NOT EXISTS comment (
						id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
						idPost TEXT, 
						content TEXT,
			 			author TEXT
					)`
	}

	_, err = db.Exec(statement)

	CheckError(err)

	return db
}
