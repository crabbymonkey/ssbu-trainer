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
	"templates/basicContent.html",
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

func lessonHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("IN LESSON: ")
	fmt.Println("Base of the Path = " + path.Base(r.URL.Path))
	lesson := Lesson{
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
	var dbToString []string

	rows, err := db.Query("SELECT * FROM character")
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
		// return, so no else is needed
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id   int
			name string
		)
		if err := rows.Scan(&id, &name); err != nil {
			panic(err)
		}
		dbToString = append(dbToString, fmt.Sprintf("%d | %s", id, name))
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	dbToString = append(dbToString, "-------------------------------------------------------")

	rows, err = db.Query("SELECT * FROM lesson")
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
		// return, so no else is needed
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id                  int
			name                string
			characterID         int
			number              int
			gif                 string
			description         string
			learningTimeSeconds int
			trainingTimeSeconds int
			testTimeSeconds     int
		)
		if err := rows.Scan(&id, &characterID, &name, &number, &gif, &description, &learningTimeSeconds, &trainingTimeSeconds, &testTimeSeconds); err != nil {
			panic(err)
		}
		dbToString = append(dbToString, fmt.Sprintf("%d | %s | %d | %d | %s | %s | %d | %d | %d |\n", id, name, characterID, number, gif, description, learningTimeSeconds, trainingTimeSeconds, testTimeSeconds))
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	data := PageContent{
		PageTitle:   "Database",
		PageContent: dbToString,
	}
	display(w, "basicContent", data)
}

var (
	db *sql.DB
)

func main() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL")+"?sslmode=disable")
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
	fmt.Println("Connected to DB " + os.Getenv("DATABASE_URL"))
	fmt.Println("Now listening to port " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// Page structure
type Page struct {
	PageTitle string
}

// PageContent structure
type PageContent struct {
	PageTitle   string
	PageContent []string
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
	ID                  int
	Name                string
	CharacterID         int
	Number              int
	Gif                 string
	Description         string
	LearningTimeSeconds int
	TrainingTimeSeconds int
	TestTimeSeconds     int
}

// LessonPage is the page for the Lesson
type LessonPage struct {
	PageTitle string
	Lesson    Lesson
}
