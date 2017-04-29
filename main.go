package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var tpl *template.Template

//"encoding/json"
// "fmt"
// "io/ioutil"
// "log"
// "net/http"

//ResponseMain ...
type responseMain struct {
	Temp    float64 `json:"temp"`
	TempMin float64 `json:"temp_min"`
	TempMax float64 `json:"temp_max"`
}

//ResponseElem ...
type responseElem struct {
	Ts         int          `json:"dt"`
	TempValues responseMain `json:"main"`
}

// Response ...
type responsetype struct {
	List []responseElem `json:"list"`
}

type tplData struct {
	Title     string
	FirstName string
	TabActive string
	Data      responsetype
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {

	http.HandleFunc("/account", accountHandler)
	http.HandleFunc("/settings", settingsHandler)
	http.HandleFunc("/createUser", createUserHandler)
	http.HandleFunc("/setSettings", setSettingsHandler)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", indexHandler)

	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	templateData := tplData{
		Title:     "Index",
		TabActive: "index",
	}

	err := tpl.ExecuteTemplate(w, "index.gohtml", templateData)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error - Resource not found", http.StatusInternalServerError)
	}

}

func settingsHandler(w http.ResponseWriter, req *http.Request) {
	templateData := tplData{
		Title:     "Settings",
		TabActive: "settings",
	}
	err := tpl.ExecuteTemplate(w, "settings.gohtml", templateData)
	if err != nil {
		log.Println(err)
	}
}

func accountHandler(w http.ResponseWriter, req *http.Request) {
	templateData := tplData{
		Title:     "Account",
		TabActive: "account",
	}
	err := tpl.ExecuteTemplate(w, "account.gohtml", templateData)
	if err != nil {
		log.Println(err)
	}
}

func createUserHandler(w http.ResponseWriter, req *http.Request) {
	templateData := tplData{
		Title:     "Account",
		TabActive: "account",
	}

	// post request, http.MethodPost is a constant
	if req.Method == http.MethodPost {
		userName := req.FormValue("uname")
		templateData.FirstName = userName
		err := tpl.ExecuteTemplate(w, "account.gohtml", templateData)
		if err != nil {
			log.Println(err)
		}
	}
}

func setSettingsHandler(w http.ResponseWriter, req *http.Request) {
	templateData := tplData{
		Title:     "Your Running Times",
		TabActive: "times",
	}

	// post request, http.MethodPost is a constant
	if req.Method == http.MethodPost {
		// minTemp := req.FormValue("minTemp")
		// maxTemp := req.FormValue("maxTemp")
		var data = makeWeatherAPIcall()
		fmt.Println("Data", data)
		templateData.Data = data

		err := tpl.ExecuteTemplate(w, "times.gohtml", templateData)
		if err != nil {
			log.Println(err)
		}
	}
}

func makeWeatherAPIcall() responsetype {
	apiKey := "4793867f02934a10b3033be4d68f385c"
	baseURL := "http://api.openweathermap.org/data/2.5/forecast?q=lakewood,co"
	query := baseURL + "&appid=" + apiKey + "&id=5427946"

	res, err := http.Get(query)
	if err != nil {
		log.Fatal(err)
	}
	response, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var r = responsetype{}
	json.Unmarshal(response, &r)

	return r
}
