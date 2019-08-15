package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//Compile templates on start
var templates = template.Must(template.ParseFiles(
	"templates/notFound.html",
	"templates/header.html",
	"templates/footer.html",
	"templates/index.html",
	"templates/commonDice.html",
	"templates/customDice.html",
	"templates/dedicatedDice.html",
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
			{Title: "Basic Dice Rolls Based on URL", Done: true},
			{Title: "Flushout Navbar", Done: false},
			{Title: "Create 404", Done: true},
			{Title: "Create Common Dice Set Page", Done: false},
			{Title: "Create Custom Dice Set Page", Done: false},
			{Title: "Add Dice Graphics", Done: true},
			{Title: "Add How to Use Page", Done: false},
			{Title: "Flushout Dedicated Dice Page with Reroll button", Done: true},
			{Title: "Make a logo", Done: false},
			{Title: "Make a AJAX callback to server to get dice roll values", Done: false},
		},
	}
	display(w, "todo", data)
}

func dedicatedDiceHandler(w http.ResponseWriter, r *http.Request) {
	inputInt, err := strconv.Atoi(r.URL.Path[1:])
	if err != nil {
		// handle error
		http.ServeFile(w, r, "static/html/issue.html")
	}

	data := DiePage{
		PageTitle: "D" + r.URL.Path[1:],
		Die:       Dice{High: inputInt, Low: 1},
	}

	display(w, "dedicated", data)
}

func commonSetHandler(w http.ResponseWriter, r *http.Request) {
	data := DiceSetPage{
		PageTitle: "Common Dice Set",
		Dice: []Dice{
			{High: 20, Low: 1},
			{High: 12, Low: 1},
			{High: 10, Low: 1},
			{High: 8, Low: 1},
			{High: 6, Low: 1},
			{High: 4, Low: 1},
		},
	}
	display(w, "common", data)
}

func customSetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Make your own set of dice!")
}

func randomPageHandler(w http.ResponseWriter, r *http.Request) {
	// If empty show the home page
	// If a number simulate a dice with that many sides
	// If static page show the static package
	// Else show the 404 page
	if r.URL.Path == "/" {
		homeHandler(w, r)
	} else if _, err := strconv.Atoi(r.URL.Path[1:]); err == nil {
		dedicatedDiceHandler(w, r)
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

func (d Dice) roll() int {
	return randomValue(d.Low, d.High)
}

func getPort() string {
	if value, ok := os.LookupEnv("PORT"); ok {
		return ":" + value
	}
	return ":8080"
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/common/", commonSetHandler)
	http.HandleFunc("/custom/", customSetHandler)
	http.HandleFunc("/todo/", todoHandler)

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

// DiceSetPage with a page title and a list of Dice
type DiceSetPage struct {
	PageTitle string
	Dice      []Dice
}

// DiePage with a page title and a Dice
type DiePage struct {
	PageTitle string
	Die       Dice
}

// Dice object e.g. D20 would have a High of 20 and Low of 1.
type Dice struct {
	High int
	Low  int
}

// RolledDice with page title, high value, low value and the rolled value
// Note that you can add a function as a part of a struct that you can define when you create the object
// Roll func() int
type RolledDice struct {
	PageTitle string
	High      int
	Low       int
	Value     int
}
