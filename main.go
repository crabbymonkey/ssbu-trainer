package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

//Compile templates on start
var templates = template.Must(template.ParseFiles(
	"templates/addCourse.html",
	"templates/admin.html",
	"templates/basicContent.html",
	"templates/trainingList.html",
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
	character, err := getCharacterByNameFromDB("basic")
	if err != nil {
		panic(err)
	}
	lessons, err := getLessonsForCharacterIDFromDB(character.ID)
	if err != nil {
		panic(err)
	}

	data := CharacterTrainingPage{
		PageTitle: "Basic Training",
		Character: *character,
		Lessons:   *lessons,
	}
	display(w, "trainingList", data)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	data := Page{
		PageTitle: "Admin",
	}
	display(w, "admin", data)
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

func convertRowsToCharacters(rows *sql.Rows) (*[]Character, error) {
	var chars []Character
	for rows.Next() {
		var c Character
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			log.Fatalf(err.Error())
			return nil, err
		}
		chars = append(chars, c)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf(err.Error())
		return nil, err
	}
	return &chars, nil
}

func convertRowsToLessons(rows *sql.Rows) (*[]Lesson, error) {
	var less []Lesson
	for rows.Next() {
		var l Lesson
		if err := rows.Scan(&l.ID, &l.CharacterID, &l.Name, &l.Number, &l.Gif, &l.Description, &l.LearningTimeSeconds, &l.TrainingTimeSeconds, &l.TestTimeSeconds); err != nil {
			log.Fatalf(err.Error())
			return nil, err
		}
		less = append(less, l)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf(err.Error())
		return nil, err
	}
	return &less, nil
}

func getCharacterByNameFromDB(name string) (*Character, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM character WHERE name='%s'", name))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
		return nil, err
	}
	defer rows.Close()

	chars, err := convertRowsToCharacters(rows)
	if err != nil {
		log.Fatalf("Error reading rows: %q", err)
		return nil, err
	}

	if len(*chars) != 1 {
		return nil, nil
	}

	return &(*chars)[0], nil
}

func getLessonsForCharacterIDFromDB(characterID int) (*[]Lesson, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM lesson WHERE character_id=%d", characterID))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
		// return, so no else is needed
		return nil, err
	}
	defer rows.Close()

	lessons, err := convertRowsToLessons(rows)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
		// return, so no else is needed
		return nil, err
	}

	return lessons, nil
}

func signinHandeler(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into a new `Credentials` instance
	creds := &Password{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword(hashedAdminPassword, []byte(creds.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// If we reach this point, that means the users password was correct, and that they are authorized
	// The default 200 status is sent
}

func addCharHandeler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Trying to add Character")
	if req.Method == "POST" {
		req.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
		// attention: If you do not call ParseForm method, the following data can not be obtained form
		fmt.Println(req.Form) // print information on server side.
		fmt.Println("path", req.URL.Path)
		fmt.Println("scheme", req.URL.Scheme)
		fmt.Println(req.Form["url_long"])
		for k, v := range req.Form {
			fmt.Println("key:", k)
			fmt.Println("val:", strings.Join(v, ""))
		}
		password := (req.Form["password"])[0]
		name := (req.Form["name"])[0]

		// Compare the stored hashed password, with the hashed version of the password that was received
		if err := bcrypt.CompareHashAndPassword(hashedAdminPassword, []byte(password)); err != nil {
			// If the two passwords don't match, return a 401 status
			fmt.Println("Bad Password")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Println("Good Password")

		_, err := db.Exec(fmt.Sprintf("INSERT INTO character (name) VALUES ('%s')", strings.ToLower(name)))
		if err != nil {
			log.Fatalf("Error writing character to database: %q", err)
			res.WriteHeader(http.StatusInternalServerError)
		}

		fmt.Println("Added Character " + strings.ToLower(name))

		context := []string{"Successfully Saved the new character to the database!"}

		data := PageContent{
			PageTitle:   "Database",
			PageContent: context,
		}
		display(res, "basicContent", data)
	}
}

var (
	db                  *sql.DB
	hashedAdminPassword []byte
)

func main() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL")+"?sslmode=disable")
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	file, err := os.Open(".hashedpass")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	hashedAdminPassword, err = ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Can't get admin password: %q", err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/basic/", basicTrainingHandler)
	http.HandleFunc("/todo/", todoHandler)
	http.HandleFunc("/lesson", lessonHandler)
	http.HandleFunc("/db", dbHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/signin", signinHandeler)
	http.HandleFunc("/addchar", addCharHandeler)

	http.HandleFunc("/", randomPageHandler)

	// Comparing the password with the hash
	// err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	// fmt.Println(err) // nil means it is a match

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

// Character struct
type Character struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Icon        string `json:"icon"`
}

// CharacterCreation struct
type CharacterCreation struct {
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LessonPage is the page for the Lesson
type LessonPage struct {
	PageTitle string
	Lesson    Lesson
}

// CharacterTrainingPage is the page for the Lesson
type CharacterTrainingPage struct {
	PageTitle string
	Character Character
	Lessons   []Lesson
}

// Password Create a struct that models the structure of a user, both in the request body, and in the DB
type Password struct {
	Password string `json:"password"`
}
