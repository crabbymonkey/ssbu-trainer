package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

//Compile templates on start
var templates = template.Must(template.ParseFiles(
	"templates/addCourse.html",
	"templates/basicTraining.html",
	"templates/footer.html",
	"templates/header.html",
	"templates/index.html",
	"templates/lesson.html",
	"templates/todo.html"))

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := Page{
		PageTitle: "Home",
	}
	display(w, "index", data)
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	data := TodoPageData{
		PageTitle: "TODO List",
		ListTitle: "My TODO List",
		Todos: []Todo{
			{Title: "Project Setup", Done: true},
			{Title: "Setup CI", Done: false},
			{Title: "Setup CD", Done: true},
			{Title: "Make Test Cases", Done: false},
			{Title: "Flushout Navbar", Done: false},
			{Title: "Create 404", Done: true},
			{Title: "Add How to Use Page", Done: false},
			{Title: "Make a logo", Done: false},
		},
	}
	display(w, "todo", data)
}

// getCharater(db *sql.DB, name string) { // &Character{
// age := 27
// row, err := db.Query("SELECT id FROM characters WHERE name=?", name)
// if err != nil {
// 	log.Fatal(err)
// }
// defer rows.Close()

// for rows.Next() {
// 	var char Character
// 	if err := rows.Scan(&name); err != nil {
// 		// Check for a scan error.
// 		// Query rows will be closed with defer.
// 		log.Fatal(err)
// 	}
// 	names = append(names, name)
// }
// // If the database is being written to ensure to check for Close
// // errors that may be returned from the driver. The query may
// // encounter an auto-commit error and be forced to rollback changes.
// rerr := rows.Close()
// if rerr != nil {
// 	log.Fatal(err)
// }

// // Rows.Err will report the last error encountered by Rows.Scan.
// if err := rows.Err(); err != nil {
// 	log.Fatal(err)
// }
// fmt.Printf("%s are %d years old", strings.Join(names, ", "), age)
// }

func lessonHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("IN LESSON: ")
	fmt.Println("Base of the Path = " + path.Base(r.URL.Path))
	lesson := Lesson{
		PageTitle:   "Short Hop",
		Name:        "The Short Hop",
		Number:      1,
		Gif:         "https://ftp.crabbymonkey.org/smash/smash_gifs/smash_examples/example_short_hop.gif",
		Description: "This is where I type explenations, Lorem ipsum dolor sit amet, fabulas nusquam facilisi per cu, ex ius voluptua principes. Quo te simul nullam. Illud aperiam accusamus mel no. Ex oporteat perfecto petentium qui, meis solum utamur sit te, per reque eligendi appellantur ei. Posse dictas laoreet pri ut, vide tamquam quaeque at his. Eu his bonorum dolorum, est vidisse discere verterem cu. Vim an veritus adipisci. An quaeque alienum electram vis, possim diceret efficiendi ex vis. Id offendit moderatius intellegam pro, ne usu atqui verterem philosophia, sit eu feugiat gloriatur expetendis. Vix ei aperiri scripserit.",
	}
	display(w, "lesson", lesson)
}

func basicTrainingHandler(w http.ResponseWriter, r *http.Request) {
	data := Page{
		PageTitle: "Basic Training",
	}
	display(w, "basicTraining", data)
}

func randomPageHandler(w http.ResponseWriter, r *http.Request) {
	// If empty show the home page
	// If a number simulate a dice with that many sides
	// If static page show the static package
	// Else show the 404 page
	if r.URL.Path == "/" {
		homeHandler(w, r)
	} else if r.URL.Path == "/lesson/lesson1" {
		lessonHandler(w, r)
	} else if strings.HasSuffix(r.URL.Path[1:], ".html") {
		http.ServeFile(w, r, "static/html/"+r.URL.Path[1:])
	} else {
		fmt.Println("Sorry but it seems this page does not exist...")
		errorHandler(w, r, http.StatusNotFound)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		display(w, "404", &Page{PageTitle: "404"})
	} else {
		http.ServeFile(w, r, "static/html/issue.html")
	}
}

// Gets a random value from the low to high values. This will include the low and high values.
func randomValue(low int, high int) int {
	scaledInt := high - low + 1 // The +1 is to offset the values so it can be the high value.
	randSource := rand.NewSource(time.Now().UnixNano() / int64(time.Millisecond))
	randomRand := rand.New(randSource)
	return randomRand.Intn(scaledInt) + low
}

var funcMap = template.FuncMap{
	"randomValue": randomValue,
}

func getPort() string {
	if value, ok := os.LookupEnv("PORT"); ok {
		return ":" + value
	}
	return ":8080"
}

func dbHandler(w http.ResponseWriter, r *http.Request) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS character (character_id SERIAL PRIMARY KEY, name VARCHAR(50) UNIQUE NOT NULL)")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error creating database table character: %q", err))
		return
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS lesson (lesson_id SERIAL PRIMARY KEY, character_id INTEGER REFERENCES character(character_id) NOT NULL, name VARCHAR(50) UNIQUE NOT NULL)")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error creating database table lesson: %q", err))
		return
	}

	data := Page{
		PageTitle: "Database",
	}
	display(w, "basicTraining", data)
}

var (
	db *sql.DB
)

func main() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL")+" sslmode=disable")
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/basic/", basicTrainingHandler)
	http.HandleFunc("/todo/", todoHandler)
	http.HandleFunc("/lesson", lessonHandler)
	http.HandleFunc("/db", dbHandler)

	http.HandleFunc("/", randomPageHandler)

	port := getPort()
	fmt.Println("Now listening to port " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// Page structure
type Page struct {
	PageTitle string
}

// Todo struct
type Todo struct {
	Title string
	Done  bool
}

// TodoPageData with titles and a list of Todos
type TodoPageData struct {
	PageTitle string
	ListTitle string
	Todos     []Todo
}

//Lesson is used to dynamicly create a lesson page
type Lesson struct {
	PageTitle   string
	Name        string
	Character   string
	Number      int
	Gif         string
	Description string
}
