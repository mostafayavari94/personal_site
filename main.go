package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"runtime"
	"server/connection"
	"strings"
	"text/template"
	"time"
)

type Page struct {
	Title  string
	Header string
	Body   string
	Data   string
}

type Event struct {
	Id           int
	StartTime    time.Time
	EndTime      time.Time
	ActivityDate time.Time
}

var Mux = http.NewServeMux()

func FunctionName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get context info for logger!")
	}

	funcName := runtime.FuncForPC(pc).Name()
	return funcName[strings.LastIndex(funcName, ".")+1:] + ".html"

}

func index(w http.ResponseWriter, R *http.Request) {

	//mux := http.NewServeMux()
	if _, pattern := Mux.Handler(R); pattern == "" {
		notFoundHandler(w, R)
		fmt.Println("offd")
		return
	}

	var file = FunctionName()

	funcsMap := map[string]interface{}{
		"upper": strings.ToUpper,
		"formatTime": func(t time.Time, layout string) string {
			return t.Format(layout)
		},
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

	t, err := template.New(page).ParseFiles(page)
	if err != nil {
		fmt.Println(err)
	}

	err = t.ExecuteTemplate(w, page, Page{Title: "Add Time"})
	if err != nil {
		fmt.Println(err)
	}
}

func storeTime(w http.ResponseWriter, r *http.Request) {

	// Establish a connection to the PostgreSQL database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		connection.Host, connection.Port, connection.User, connection.Password, connection.DBname)

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
		start_time TIME NOT NULL,
		end_time TIME NOT NULL,
		activity_date Date NOT NULL
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	var validation []string

	if r.PostFormValue("start_time") == "" {
		validation = append(validation, "enter your start time")
	}

	if r.PostFormValue("end_time") == "" {
		validation = append(validation, "enter your end time")
	}

	fmt.Println(validation)

	if len(validation) != 0 {
		w.Header().Set("Content-Type", "application/json")
		temp, err := json.Marshal(validation)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}
		w.Write(temp)
		return
	}

	// Insert data into the table
	start_time := string(r.PostFormValue("start_time"))
	end_time := string(r.PostFormValue("end_time"))

	activity_date := string(r.PostFormValue("activity_date"))
	if activity_date == "" {
		activity_date = time.Now().Format("2006-01-02 15:04:05")
	}

	insertQuery := `INSERT INTO event (start_time, end_time, activity_date) VALUES ($1, $2, $3);`

	_, err = db.Exec(insertQuery, start_time, end_time, activity_date)
	if err != nil {
		log.Fatal(err)
	}

}

func timesList(w http.ResponseWriter, r *http.Request) {
	var page = FunctionName()

	funcsMap := map[string]interface{}{
		"upper": strings.ToUpper,
		"formatTime": func(t time.Time, layout string) string {
			return t.Format(layout)
		},
	}

	t, err := template.New(page).Funcs(funcsMap).ParseFiles(page)
	if err != nil {
		fmt.Println(err)
	}

	//////////////////////////////////////////////////
	// Establish a connection to the PostgreSQL database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		connection.Host, connection.Port, connection.User, connection.Password, connection.DBname)

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
	//createTableQuery := `
	//CREATE TABLE IF NOT EXISTS event (
	//	id SERIAL PRIMARY KEY,
	//	start_time TIME NOT NULL,
	//	end_time TIME NOT NULL,
	//	activity_date Date NOT NULL
	//);
	//`

	createTableQuery := `
	SELECT * FROM event;
	`

	result, err := db.Query(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer result.Close()

	var Events = map[int]Event{}
	for result.Next() {
		var (
			_id           int
			_startTime    time.Time
			_endTime      time.Time
			_activityDate time.Time
			// Add more variables for each column in your "event" table
		)

		if err := result.Scan(&_id, &_startTime, &_endTime, &_activityDate); err != nil {
			log.Fatal(err)
		}

		Events[_id] = Event{StartTime: _startTime, EndTime: _endTime, ActivityDate: _activityDate}

		// Use the values of the columns as needed
	}

	//fmt.Println(Events)
	//////////////////////////////////////////////////

	err = t.ExecuteTemplate(w, page, Events)
	if err != nil {
		fmt.Println(err)
	}
}

const public = "../frontend/public/"

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	Mux.Handle("/public", http.FileServer(http.Dir(public)))

	Mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("ridi")
	})
	Mux.HandleFunc("/", index)
	Mux.HandleFunc("/404", notFoundHandler)
	Mux.HandleFunc("/addTime", addTime)
	Mux.HandleFunc("/storeTime", storeTime)
	Mux.HandleFunc("/timesList", timesList)

	http.NotFoundHandler()

	http.ListenAndServe(":4000", Mux)
}
