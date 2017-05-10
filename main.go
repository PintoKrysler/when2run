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
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var tpl *template.Template

const (
	minDefaultTemperature float64 = 0
	maxDefaultTemperature float64 = 110
	dbUser                        = "postgres"
	dbPassword                    = "postgres"
	dbName                        = "when2run"
)

// Settings
type settings struct {
	MinTemp float64
	MaxTemp float64
}

func (r *settings) newSettings() settings {
	return settings{
		MinTemp: minDefaultTemperature,
		MaxTemp: maxDefaultTemperature,
	}
}

//ResponseMain ...
type responseMain struct {
	Temp    float64 `json:"temp"`
	TempMin float64 `json:"temp_min"`
	TempMax float64 `json:"temp_max"`
}

//ResponseElem ...
type responseElem struct {
	Ts          int          `json:"dt"`
	TempValues  responseMain `json:"main"`
	TsFormatted time.Time
	GoRun       bool
}

// responsetype ...
type responsetype struct {
	List []responseElem `json:"list"`
}

// tplData
type tplData struct {
	Title     string
	TabActive string
	Data      responsetype
}

type user struct {
	Email    string
	Password string
	Settings settings
}

type app struct {
	User        user
	UserLogged  bool
	CurrentView string
	Data        interface{}
	MsgError    string
}

var myapp = app{
	UserLogged: false,
}

var database *sql.DB

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	connectDB()
	getUsers()
}

func main() {

	http.HandleFunc("/account", accountHandler)
	http.HandleFunc("/settings", settingsHandler)
	http.HandleFunc("/createUser", createUserHandler)
	http.HandleFunc("/setSettings", setSettingsHandler)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/login", loginHandler)
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
	myapp.CurrentView = "index"

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
	myapp.CurrentView = "settings"
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
	myapp.CurrentView = "account"
	myapp.Data = templateData

	err := tpl.ExecuteTemplate(w, "account.gohtml", myapp)
	if err != nil {
		log.Println(err)
	}
}

func createUserHandler(w http.ResponseWriter, req *http.Request) {
	templateData := tplData{
		Title:     "Create Account",
		TabActive: "account",
	}

	// post request, http.MethodPost is a constant
	if req.Method == http.MethodPost {
		email := req.FormValue("email")
		password := req.FormValue("password")

		// Check if there is a user with that email
		_, err := getUser(email)
		if err != nil {
			if err == sql.ErrNoRows {

				// User does not exists
				if inserted := insertUser(email, password); inserted {
					myapp.UserLogged = true
					userSettings := settings{}
					userSettings = userSettings.newSettings()
					myapp.User = user{Email: email, Password: password, Settings: userSettings}

				}
				myapp.Data = templateData
				err := tpl.ExecuteTemplate(w, "account.gohtml", myapp)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			// User already exists
			myapp.MsgError = "User account already exists"
			err := tpl.ExecuteTemplate(w, "createaccount.gohtml", myapp)
			if err != nil {
				log.Println(err)
			}
		}

	} else {
		// Is not a POST request! we only want to show the form to create a new account
		myapp.Data = templateData
		err := tpl.ExecuteTemplate(w, "createaccount.gohtml", myapp)
		if err != nil {
			log.Println(err)
		}
	}
}

func logoutHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		if myapp.UserLogged {
			myapp.UserLogged = false
			myapp.User = user{}
			templateData := tplData{
				Title:     "Account",
				TabActive: "account",
			}
			myapp.Data = templateData
			err := tpl.ExecuteTemplate(w, "account.gohtml", myapp)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	templateData := tplData{
		Title:     "Login",
		TabActive: "account",
	}
	if req.Method == http.MethodPost {
		email := req.FormValue("email")
		password := req.FormValue("password")
		loggedIn := login(email, password)
		if !loggedIn {
			log.Println("Authentication error")
		}
	}
	myapp.Data = templateData
	err := tpl.ExecuteTemplate(w, "account.gohtml", myapp)
	if err != nil {
		log.Println(err)
	}
}

func login(email string, password string) bool {
	var success = false
	if !myapp.UserLogged {
		// Find a user with email and password match
		foundUser, err := getUser(email)
		if err != nil {
			if err == sql.ErrNoRows {
				// User does not exists
				myapp.MsgError = "User account does not exists"
			} else {
				log.Println(err)
			}
		} else {
			if strings.Trim(foundUser.Password, " ") == strings.Trim(password, " ") {
				// Authentication passed
				myapp.UserLogged = true
				myapp.User = user{Email: foundUser.Email, Password: foundUser.Password}
				success = true
			} else {
				// Authentication error
				myapp.MsgError = "Authentication Error"
			}
		}

	}
	return success
}

func setSettingsHandler(w http.ResponseWriter, req *http.Request) {
	templateData := tplData{
		Title:     "Your Running Times",
		TabActive: "times",
	}
	myapp.CurrentView = "settings"
	// post request, http.MethodPost is a constant
	if req.Method == http.MethodPost {
		minTemp := req.FormValue("minTemp")
		maxTemp := req.FormValue("maxTemp")

		if myapp.UserLogged && (minTemp != "" || maxTemp != "") {
			// Update database if user is logged in and temperatures passed on the form
			// are different than DB values
			fmt.Println("Update user information ", minTemp, maxTemp)

			changedValues := make(map[string]interface{})

			if maxTemp != "" {
				maxTempVal, err := strconv.ParseFloat(maxTemp, 64)
				if err != nil {
					log.Println(err)
				}
				// Check if maxTemp changed
				//fmt.Println("maxChanges", maxTempVal, myapp.User.Settings.MaxTemp)
				if maxTempVal != myapp.User.Settings.MaxTemp {
					changedValues["MaxTemp"] = maxTempVal
				}
			}

			if minTemp != "" {
				minTempVal, err := strconv.ParseFloat(minTemp, 64)
				if err != nil {
					log.Println(err)
				}
				// Check if minTemp changed
				//fmt.Println("minChanges", minTempVal, myapp.User.Settings.MinTemp)
				if minTempVal != myapp.User.Settings.MinTemp {
					changedValues["MinTemp"] = minTempVal
				}
			}
			updateUser(myapp.User.Email, changedValues)
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
		tsString := strconv.Itoa(elem.Ts)
		tsFormatted, err := strconv.ParseInt(tsString, 10, 64)
		if err != nil {
			panic(err)
		}
		data.List[i].TsFormatted = time.Unix(tsFormatted, 0)
		if myapp.User.Settings.MinTemp <= elem.TempValues.TempMin && myapp.User.Settings.MaxTemp >= elem.TempValues.TempMax {
			data.List[i].GoRun = true
		}

	}
	return data
}

func connectDB() {
	databaseConn, err := sql.Open("postgres", "user="+dbUser+" password="+dbPassword+" dbname="+dbName+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	database = databaseConn
}
func getUsers() {
	rows, err := database.Query(`SELECT email,password FROM usuario`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		email    string
		password string
	)
	for rows.Next() {
		err := rows.Scan(&email, &password)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(email, password)
	}
}

func getUser(email string) (user, error) {
	u := user{}
	var pass, em string
	var tempMin, tempMax interface{}

	sqlStatement := `SELECT email,password,mintemp,maxtemp FROM usuario WHERE email=$1;`
	// Replace 3 with an ID from your database or another random
	// value to test the no rows use case.

	err := database.QueryRow(sqlStatement, email).Scan(&em, &pass, &tempMin, &tempMax)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		log.Fatal(err)
	default:
		u.Email = em
		u.Password = pass
	}
	return u, err
}

func updateUser(userID string, values map[string]interface{}) bool {
	updated := false
	setFields := ""
	var setValues []interface{}
	cnt := 0
	paramIndex := 2
	// First parameter is user id
	setValues = append(setValues, userID)
	for key, val := range values {
		if cnt > 0 {
			setFields += " , "
		}
		setFields += key + "= $" + strconv.Itoa(paramIndex)
		setValues = append(setValues, val)
		cnt++
		paramIndex++
	}
	//fmt.Println("setvalues", setFields)
	// Update user information new values
	sqlStatement := `
	UPDATE USUARIO
	SET ` + setFields + `
	WHERE email = $1;`
	_, err := database.Exec(sqlStatement, setValues...)
	if err != nil {
		log.Fatal(err)
	}

	return updated
}

func insertUser(email string, password string) bool {
	var userid string
	//" + minDefaultTemperature + "','" + maxDefaultTemperature + "
	err := database.QueryRow("INSERT INTO USUARIO (email,password) VALUES('" + email + "','" + password + "') RETURNING email").Scan(&userid)
	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("userid created", userid)
	return true
}
