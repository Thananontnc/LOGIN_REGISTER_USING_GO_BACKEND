package main

import (
	"database/sql"
	"fmt" // print on console
	"log"
	"net/http"      // http server
	"path/filepath" // file path

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

// CONNECT TO DATABASE
var db *sql.DB

func connect_to_db() {
	var err error

	connection := "root:root@tcp(localhost:3306)/user_db"

	db, err = sql.Open("mysql", connection)

	if err != nil {
		fmt.Println("Error opening Database : ", err)
	}
}

// ClOSE DATABASE
func closeDB() {
	if db != nil {
		db.Close()
	}
}

// QUERY SQL

// Insert value to users table
func insertUSER(username, email, password, phone string) error {
	query := "INSERT INTO users (username , email ,password,phone) VALUES (?,?,?,?)"

	_, err := db.Exec(query, username, email, password, phone)
	if err != nil {
		fmt.Println("Error Insert the data: ", err)
		return err
	}
	return nil
}

// checkPassword
func checkPassword(username, password string) (bool, error) {
	var passwordContainer string
	// Query the database for the stored password based on the username
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&passwordContainer)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("no user found with that username")
		}
		return false, err
	}

	// Compare the entered password with the stored password
	if passwordContainer == password {
		return true, nil
	}
	return false, nil
}

// HANDLE HTML FILE
// Handle register.html
func registerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		filePath := filepath.Join("web", "register.html")
		http.ServeFile(w, r, filePath)
	} else if r.Method == http.MethodPost {

		r.ParseForm()

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		phone := r.FormValue("phone")

		err := insertUSER(username, email, password, phone)
		if err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Username: %s, Email: %s , Password: %s ,Phone: %s",
			username, email, password, phone)
	}
}

// Handle login.html
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		filepath := filepath.Join("web", "login.html")
		http.ServeFile(w, r, filepath)
	} else if r.Method == http.MethodPost {

		r.ParseForm()
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		valid, err := checkPassword(username, password)

		if err != nil {
			log.Fatal("Login failed: ", err)
		}
		if valid {
			fmt.Println("Login Successful!")
			http.Redirect(w, r, "/home-page", http.StatusSeeOther)
		} else {
			fmt.Println("Invalid username or password")
		}

		fmt.Fprintf(w, "Username : %s , Email: %s , Password: %s", username, email, password)
	}

}

// Handle Home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	filepath := filepath.Join("web", "main.html")
	http.ServeFile(w, r, filepath)
}

// HANDLE CSS FILE
// Serve register.css
func serveRegisterCSS(w http.ResponseWriter, r *http.Request) {
	filepath := filepath.Join("style", "register.css")
	http.ServeFile(w, r, filepath)
}

// Serve login.css
func serveLoginCSS(w http.ResponseWriter, r *http.Request) {
	filepath := filepath.Join("style", "login.css")
	http.ServeFile(w, r, filepath)
}

// Serve home.css
func serveHomeCss(w http.ResponseWriter, r *http.Request) {
	filepath := filepath.Join("style", "home.css")
	http.ServeFile(w, r, filepath)
}

func main() {
	// Connect to database
	connect_to_db()
	defer closeDB()
	//  Register
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/register-style", serveRegisterCSS)
	// Login
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/login-style", serveLoginCSS)
	// Home
	http.HandleFunc("/home-page", homeHandler)
	http.HandleFunc("/home-style", serveHomeCss)

	// Close connection
	defer closeDB()
	// Run Port :8080
	http.ListenAndServe(":8080", nil)
}
