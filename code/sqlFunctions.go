package account

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

//CompareInputAndDB will compare the input email and the email inside the databse, it's also do it for the password
func CompareInputAndDB(email string, password string) bool {
	// Open the database
	database, err := sql.Open("sqlite3", "./dblogin.db")
	CheckError(err)
	defer database.Close()

	// browse the database
	rows, err := database.Query("SELECT id, pseudo, email, password FROM dblogin")
	CheckError(err)
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
func InsertIntoTypes(db *sql.DB, pseudo string, email string, password string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO dblogin (pseudo, email, password) VALUES (?, ?, ?)`, pseudo, email, password)

	return result.LastInsertId()
}

//insertIntoTypesPost will insert data in the post which is in post.db
func InsertIntoTypesPost(db *sql.DB, title string, content string, filter string) {
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
	stmt, err := tx.Prepare("INSERT INTO post (title, content, author, filter, like, userlike) VALUES (?, ?, ?, ?, 0, '')")
	CheckError(err)
	// Execute la transaction
	_, err = stmt.Exec(title, content, data["userConnected"], filter)
	CheckError(err)
	// Commit la transaction
	tx.Commit()
}

//showPosts will browse the database post.db and taking the values inside it for using it later
func ShowPosts(data_post map[string]interface{}, id_post string) {

	var arrPost []Post
	var query string
	// Open the database
	database, err := sql.Open("sqlite3", "./post.db")
	CheckError(err)
	defer database.Close()

	if id_post == "" {
		query = "SELECT id, title, content, author, filter, like FROM post" //this query will show every posts
	} else {
		query = "SELECT id, title, content, author, filter, like FROM post WHERE id = " + id_post //this query will show one particular post
	}

	rows, err := database.Query(query)
	CheckError(err)
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&p.Id, &p.Title, &p.Content, &p.Author, &p.Filters, &p.Likes)

		arrPost = append([]Post{p}, arrPost...)

	}
	data_post["showPost"] = arrPost

}

func ShowComments(data_post map[string]interface{}, id_post string) {
	var all_comment []Comment
	database, err := sql.Open("sqlite3", "./post.db")
	CheckError(err)
	defer database.Close()
	rows, err := database.Query("SELECT content, author FROM comment WHERE idPost = ?", id_post)
	CheckError(err)
	for rows.Next() {
		err = rows.Scan(&c.Content, &c.Author)
		CheckError(err)
		all_comment = append(all_comment, c)
	}
	data_post["all_comment"] = all_comment
}

func AddComment(comment string, idPost string, data_post map[string]interface{}) {
	database, err := sql.Open("sqlite3", "./post.db")
	CheckError(err)
	defer database.Close()

	// démarre une transaction
	tx, err := database.Begin()
	CheckError(err)
	// Prépare la transaction
	stmt, err := tx.Prepare("INSERT INTO comment (idPost, content, author) VALUES (?, ?, ?)")
	CheckError(err)
	// Execute la transaction
	_, err = stmt.Exec(idPost, comment, data_post["userConnected"].(string))
	CheckError(err)
	// Commit la transaction
	tx.Commit()
}

func DeleteHandler(response http.ResponseWriter, request *http.Request) {
	id_post := request.FormValue("id-post-delete")

	fmt.Println(id_post)
	// Open the database
	database, err := sql.Open("sqlite3", "./post.db")
	CheckError(err)
	defer database.Close()

	// démarre une transaction
	tx, err := database.Begin()
	CheckError(err)
	// Prépare la transaction
	stmt, err := tx.Prepare("DELETE FROM post WHERE id = ?")
	CheckError(err)
	// Execute la transaction
	_, err = stmt.Exec(id_post)
	CheckError(err)
	// Commit la transaction
	tx.Commit()
	http.Redirect(response, request, "/welcome", http.StatusSeeOther)
}

func AddLikes(likes string) {
	// Open the database

	var nb_like int
	var query string
	var userLike string
	database, err := sql.Open("sqlite3", "./post.db")
	CheckError(err)
	defer database.Close()
	rows, err := database.Query("SELECT like, userlike FROM post WHERE id = ?", likes)
	CheckError(err)
	for rows.Next() {
		err = rows.Scan(&nb_like, &userLike)
		CheckError(err)
		fmt.Println(nb_like)
	}
	nb_like++
	fmt.Println(nb_like)

	// démarre une transaction
	tx, err := database.Begin()
	CheckError(err)
	// Prépare la transaction
	query = "UPDATE post SET like = " + strconv.Itoa(nb_like) + " WHERE id = " + likes
	stmt, err := tx.Prepare(query)
	CheckError(err)
	// Execute la transaction
	_, err = stmt.Exec()
	CheckError(err)
	if userLike == "" {
		userLike = data["userConnected"].(string)
	} else {
		userLike += " " + data["userConnected"].(string)
	}
	query = "UPDATE post SET userlike = ? WHERE id = " + likes
	stmt, err = tx.Prepare(query)
	CheckError(err)
	// Execute la transaction
	_, err = stmt.Exec(userLike)
	CheckError(err)
	// Commit la transaction
	tx.Commit()
}

func SearchAllUserInPostDb(likes string, data_post map[string]interface{}) bool {
	var userLike string
	var arrLike []string
	fmt.Println(likes)
	database, err := sql.Open("sqlite3", "./post.db")
	CheckError(err)
	defer database.Close()
	rows, err := database.Query("SELECT userlike FROM post WHERE id = ?", likes)
	CheckError(err)
	for rows.Next() {
		err = rows.Scan(&userLike)
		CheckError(err)

	}
	fmt.Println("test", userLike)
	if userLike != "" && data["userConnected"].(string) != "" {
		arrLike = strings.Split(userLike, " ")
		fmt.Println(arrLike)
		for _, val := range arrLike {
			if val == data["userConnected"].(string) {
				data_post["alreadyLiked"] = true
				return true
			}
		}
	}

	data_post["alreadyLiked"] = false
	return false
}

// func insertIntoTypesUserInfos(db *sql.DB, bio string, linkGit string, linkLinkedin string, username interface{}) {
// 	// Open the database
// 	database, err := sql.Open("sqlite3", "./dblogin.db")
// 	CheckError(err)
// 	defer database.Close()

// 	fmt.Println(bio)
// 	fmt.Println(linkGit)
// 	fmt.Println(linkLinkedin)
// 	fmt.Println(username)

// 	// démarre une transaction
// 	tx, err := database.Begin()
// 	CheckError(err)
// 	// Prépare la transaction
// 	stmt, err := tx.Prepare("UPDATE dblogin SET bio = ?, linkGit = ?, linkLinkedin = ? WHERE email = ?")
// 	CheckError(err)
// 	// Execute la transaction
// 	_, err = stmt.Exec(bio, linkGit, linkLinkedin, username)
// 	CheckError(err)
// 	// Commit la transaction
// 	tx.Commit()
// }
