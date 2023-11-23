package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"text/template"
	// "database/sql"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "pguser"
	password = "pguser"
	dbname   = "pgdb"
)

type Page struct {
	Title  string
	Header string
	Body   string
}

func FunctionName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get context info for logger!")
	}

	funcName := runtime.FuncForPC(pc).Name()
	return funcName[strings.LastIndex(funcName, ".")+1:] + ".html"

}

func index(w http.ResponseWriter, R *http.Request) {

	var file = FunctionName()

	funcsMap := map[string]interface{}{
		"upper": strings.ToUpper,
	}

	t, err := template.New(file).Funcs(funcsMap).ParseFiles(file)
	if err != nil {
		fmt.Println(err)
	}

	err = t.ExecuteTemplate(w, file, Page{Title: "Mostafa yavari", Header: "It's a a personal website", Body: "This place should be created"})
	if err != nil {
		fmt.Println(err)
	}
}

func addTime(w http.ResponseWriter, R *http.Request) {

	var page = FunctionName()

	// funcsMap := map[string]interface{}{
	// 	"upper" : strings.ToUpper,
	// }

	t, err := template.New(page).ParseFiles(page)
	if err != nil {
		fmt.Println(err)
	}

	err = t.ExecuteTemplate(w, page, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func storeTime(w http.ResponseWriter, r *http.Request) {
	// Establish a connection to the PostgreSQL database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Open doesn't open a connection. Ping verifies that the database connection is alive.
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create a table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS event (
		id SERIAL PRIMARY KEY,
		start_time TIMESTAMP NOT NULL,
		end_time TIMESTAMP NOT NULL
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Insert data into the table
	startTime := time.Now()
	endTime := startTime.Add(2 * time.Hour) // Assuming the event duration is 2 hours

	insertQuery := `
	INSERT INTO event (start_time, end_time) VALUES ($1, $2);
	`
	_, err = db.Exec(insertQuery, startTime, endTime)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data inserted successfully.")

	fmt.Printf(r.PostFormValue("start_time"))
}

// var userDB map[string]string

// func add(w http.ResponseWriter, R *http.Request) {
// 	userName := R.FormValue("userName")
// 	password := R.FormValue("password")

// 	_, ok := userDB["bob"]
// 	if ok {
// 		userDB[userName] = password
// 	}
// 	fmt.Println(userDB)
// }

func main() {
	http.HandleFunc("/", index)
	// http.HandleFunc("/times", addTime)
	http.HandleFunc("/addTime", addTime)
	http.HandleFunc("/storeTime", storeTime)

	http.ListenAndServe(":4000", nil)
}

// func handleRequest(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		addTime(w, r)
// 	case http.MethodPost:
// 		handlePost(w, r)
// 	default:
// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 	}
// }

// http.Handle("/", http.FileServer(http.Dir("css/")))
// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
