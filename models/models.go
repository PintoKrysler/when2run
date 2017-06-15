package models

import (
	"database/sql"
	"html/template"
	"time"

	"github.com/gorilla/sessions"
)

// Constants temperatures
const (
	MinDefaultTemperature float64 = 0
	MaxDefaultTemperature float64 = 110
)

type (

	// Controller interface has a Dispatcher
	Controller interface {
		Dispatcher()
	}

	//Server struct
	Server struct {
		Sess *sessions.CookieStore
		Tpl  *template.Template
		Db   *sql.DB
	}

	//User type
	User struct {
		Email    string
		Password string
		Settings Settings
	}

	//Settings type
	Settings struct {
		MinTemp float64
		MaxTemp float64
		Days    map[int]bool
	}

	//ResponseMain ...
	responseMain struct {
		Temp    float64 `json:"temp"`
		TempMin float64 `json:"temp_min"`
		TempMax float64 `json:"temp_max"`
	}

	//ResponseElem ...
	ResponseElem struct {
		Ts            int          `json:"dt"`
		TempValues    responseMain `json:"main"`
		TsFormatted   time.Time
		GoRun         bool
		Weekday       time.Weekday
		Day           int
		Month         time.Month
		TimeFormatted string
	}

	// Responsetype ...
	Responsetype struct {
		List []ResponseElem `json:"list"`
	}

	// TplData ...
	TplData struct {
		Title     string
		TabActive string
		Data      Responsetype
	}

	// App ...
	App struct {
		User        User
		UserLogged  bool
		CurrentView string
		Data        interface{}
		MsgError    string
	}
)

//New creates new settings with defaults temperatures
func (r *Settings) New() Settings {
	return Settings{
		MinTemp: MinDefaultTemperature,
		MaxTemp: MaxDefaultTemperature,
	}
}

func (u User) IsEmpty() bool {
	if u.Email == "" && u.Password == "" {
		return true
	}
	return false
}
