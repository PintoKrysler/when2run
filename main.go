package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

var tpl *template.Template

const (
	minDefaultTemperature float64 = 0
	maxDefaultTemperature float64 = 110
	DB_USER                       = "postgres"
	DB_PASSWORD                   = "postgres"
	DB_NAME                       = "when2run"
)

// Settings
type settings struct {
	MinTemp float64
	MaxTemp float64
}

//ResponseMain ...
type responseMain struct {
	Temp    float64 `json:"temp"`
	TempMin float64 `json:"temp_min"`
	TempMax float64 `json:"temp_max"`
}

//ResponseElem ...
type responseElem struct {
	Ts           int          `json:"dt"`
	TempValues   responseMain `json:"main"`
	Ts_formatted time.Time
	GoRun        bool
}

// Response ...
type responsetype struct {
	List []responseElem `json:"list"`
}

// tplData
type tplData struct {
	Title     string
	FirstName string
	TabActive string
	Data      responsetype
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	connectDB()
}

var userSettings = settings{
	MinTemp: minDefaultTemperature,
	MaxTemp: maxDefaultTemperature,
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
		minTemp := req.FormValue("minTemp")
		maxTemp := req.FormValue("maxTemp")
		fmt.Println("values passed", minTemp, maxTemp)
		if minTemp != "" {
			fmt.Println(minTemp)
			val, err := strconv.ParseFloat(minTemp, 64)
			if err != nil {
				log.Println(err)
			}
			userSettings.MinTemp = val
		}
		if maxTemp != "" {
			fmt.Println(maxTemp)
			val, err := strconv.ParseFloat(maxTemp, 64)
			if err != nil {
				log.Println(err)
			}
			userSettings.MaxTemp = val
		}

		var data = makeWeatherAPIcall()
		templateData.Data = data

		err := tpl.ExecuteTemplate(w, "times.gohtml", templateData)
		if err != nil {
			log.Println(err)
		}
	}
}

func makeWeatherAPIcall() responsetype {
	apiKey := "4793867f02934a10b3033be4d68f385c"
	baseURL := "http://api.openweathermap.org/data/2.5/forecast?q=lakewood,co&units=imperial"
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
	r = parseData(r)

	return r
}

// This function parses the Weather API data
// Transforms ts into readable data for the view
//
func parseData(data responsetype) responsetype {
	for i := 0; i < len(data.List); i++ {
		elem := data.List[i]
		ts_sting := strconv.Itoa(elem.Ts)
		ts_formatted, err := strconv.ParseInt(ts_sting, 10, 64)
		if err != nil {
			panic(err)
		}
		data.List[i].Ts_formatted = time.Unix(ts_formatted, 0)
		// fmt.Println(userSettings.MinTemp)
		// fmt.Println(userSettings.MaxTemp)
		// fmt.Println(elem.TempValues.TempMin)
		// fmt.Println(elem.TempValues.TempMax)
		if userSettings.MinTemp <= elem.TempValues.TempMin && userSettings.MaxTemp >= elem.TempValues.TempMax {
			// fmt.Println("Between Min and max")
			data.List[i].GoRun = true
		}

	}
	return data
}

func connectDB() {
	db, err := sql.Open("postgres", "user="+DB_USER+" password="+DB_PASSWORD+" dbname="+DB_NAME+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// insert in DB
	var userid int = 1
	//err = db.QueryRow(`INSERT INTO USUARIO(id,name) VALUES(2,'kevin') RETURNING id`).Scan(&userid)
	fmt.Println("userid ", userid)
	rows, err := db.Query(`SELECT * FROM usuario`)
	defer rows.Close()

	var (
		id   int
		name string
	)
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}

}
