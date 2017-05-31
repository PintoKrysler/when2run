package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/pintokrysler/when2run/models"
	"github.com/pintokrysler/when2run/server"
	"github.com/pintokrysler/when2run/utils"
)

// UserController ...
type UserController struct{}

// NewUserController creates a new controller instance
func NewUserController() *UserController {
	newController := UserController{}
	return &newController
}

// Dispatch creates a new user instance
func (uc UserController) Dispatch(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	action := params.ByName("action")
	switch action {
	case "login":
		uc.loginHandler(w, r)
		break
	case "create":
		uc.create(w, r)
		break
	case "logout":
		uc.logoutHandler(w, r)
		break
	case "settings":
		uc.settingsHandler(w, r)
		break
	}

}

// CreateUser creates a new user instance
func (uc UserController) create(w http.ResponseWriter, req *http.Request) {
	var myapp = models.App{}
	templateData := models.TplData{
		Title:     "Create Account",
		TabActive: "account",
	}
	session, _ := server.Server.Sess.Get(req, "when2runSess")
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
					userSettings := models.Settings{}
					userSettings = userSettings.New()
					myapp.User = models.User{Email: email, Password: password, Settings: userSettings}
					session.Values["authenticated"] = true
					session.Values["user"] = myapp.User
				}
				myapp.Data = templateData
				err := server.Server.Tpl.ExecuteTemplate(w, "account.gohtml", myapp)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			// User already exists
			session.Values["authenticated"] = false
			session.Values["user"] = models.User{}
			myapp.MsgError = "User account already exists"
			err := server.Server.Tpl.ExecuteTemplate(w, "createaccount.gohtml", myapp)
			if err != nil {
				log.Println(err)
			}
		}
		session.Save(req, w)
	} else {
		// Is not a POST request! we only want to show the form to create a new account
		myapp.Data = templateData
		err := server.Server.Tpl.ExecuteTemplate(w, "createaccount.gohtml", myapp)
		if err != nil {
			log.Println(err)
		}
	}
}

// UpdateUser updates an existing user instance
func (uc UserController) updateHandler(w http.ResponseWriter, req *http.Request) {
	var myapp = models.App{}
	templateData := models.TplData{
		Title:     "Your Running Times",
		TabActive: "times",
	}
	var maxTempVal float64
	var minTempVal float64

	session, err := server.Server.Sess.Get(req, "when2runSess")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Check if user is authenticated using session
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		// User is authenticated
		myapp.User = session.Values["user"].(models.User)
		myapp.UserLogged = true
	}

	s := models.Settings{MinTemp: models.MinDefaultTemperature, MaxTemp: models.MaxDefaultTemperature}
	myapp.CurrentView = "settings"
	minTemp := req.FormValue("minTemp")
	maxTemp := req.FormValue("maxTemp")

	// if maxTemp was passed on the form
	if maxTemp != "" {
		maxTempVal, _ = strconv.ParseFloat(maxTemp, 64)
		s.MaxTemp = maxTempVal
	}

	if minTemp != "" {
		minTempVal, _ = strconv.ParseFloat(minTemp, 64)
		s.MinTemp = minTempVal
	}

	if myapp.UserLogged && (minTemp != "" || maxTemp != "") {
		// Update database if user is logged in and temperatures passed on the form
		// are different than DB values
		changedValues := make(map[string]interface{})
		// Check if maxTemp changed
		if s.MaxTemp != myapp.User.Settings.MaxTemp {
			changedValues["MaxTemp"] = maxTempVal
			myapp.User.Settings.MaxTemp = maxTempVal
		}

		// Check if minTemp changed
		//fmt.Println("minChanges", minTempVal, myapp.User.Settings.MinTemp)
		if s.MinTemp != myapp.User.Settings.MinTemp {
			changedValues["MinTemp"] = minTempVal
			myapp.User.Settings.MinTemp = minTempVal
		}

		if len(changedValues) > 0 {
			//fmt.Println("Update user information ", minTemp, maxTemp)
			updateSuccess := updateUser(myapp.User.Email, changedValues)
			if updateSuccess {
				// Change value of user
				newUserInfo, err := getUser(myapp.User.Email)
				if err == nil {
					session, err2 := server.Server.Sess.Get(req, "when2runSess")
					if err2 != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					// Check if user is authenticated using session
					if auth, ok := session.Values["authenticated"].(bool); ok && auth {
						session.Values["user"] = newUserInfo
						session.Save(req, w)
					}
				}
			}
		}
	}

	var data = utils.MakeWeatherAPIcall(s)
	templateData.Data = data
	myapp.Data = templateData
	err = server.Server.Tpl.ExecuteTemplate(w, "times.gohtml", myapp)
	if err != nil {
		log.Println(err)
	}
}

// DeleteUser deletes an existing user
func (uc UserController) delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// u := models.User{}

}

// GetUser gets user information
func (uc UserController) get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// u := models.User{}

}

// CreateUser creates a new user instance
func (uc UserController) loginHandler(w http.ResponseWriter, req *http.Request) {
	templateData := models.TplData{
		Title:     "Login",
		TabActive: "account",
	}
	var myapp = models.App{}
	if req.Method == http.MethodPost {
		// Get cookie information
		session, err := server.Server.Sess.Get(req, "when2runSess")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Check if user is authenticated using session
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			email := req.FormValue("email")
			password := req.FormValue("password")
			loggedIn, loggedUser, errMsg := login(email, password)
			if loggedIn {
				session.Values["authenticated"] = true
				session.Values["user"] = loggedUser
				myapp.User = loggedUser
				myapp.UserLogged = true
			} else {
				log.Println(errMsg)
				session.Values["authenticated"] = false
				session.Values["user"] = loggedUser
				myapp.MsgError = errMsg
			}
			session.Save(req, w)
		} else {
			fmt.Println("Already logged in")
			if u, ok := session.Values["user"].(models.User); !ok {
				// Handle the case that it's not an expected type
				if u.IsEmpty() {
					// user is empty
					session.Values["authenticated"] = false
					session.Save(req, w)
				} else {
					fmt.Println("current user is ", u)
				}
			}
		}
	}

	myapp.Data = templateData
	err := server.Server.Tpl.ExecuteTemplate(w, "account.gohtml", myapp)
	if err != nil {
		log.Println(err)
	}
}

func (uc UserController) logoutHandler(w http.ResponseWriter, req *http.Request) {
	var myapp = models.App{}
	// Get cookie information
	session, err := server.Server.Sess.Get(req, "when2runSess")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Check if user is authenticated using session
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		//fmt.Println("logout- authenthicated user")
		noUser := models.User{}
		myapp.UserLogged = false
		myapp.User = noUser
		templateData := models.TplData{
			Title:     "Account",
			TabActive: "account",
		}
		session.Values["authenticated"] = false
		session.Values["user"] = noUser
		session.Save(req, w)
		myapp.Data = templateData
		err = server.Server.Tpl.ExecuteTemplate(w, "account.gohtml", myapp)
		if err != nil {
			log.Println(err)
		}
	} else {
		//fmt.Println("logout un authenticated user")
	}
}

//login
func login(email string, password string) (bool, models.User, string) {
	var success = false
	var msg = ""
	var u = models.User{}
	// Find a user with email and password match
	foundUser, err := getUser(email)
	if err != nil {
		if err == sql.ErrNoRows {
			// User does not exists
			msg = "User account does not exists"
		} else {
			log.Println(err)
		}
	} else {
		if strings.Trim(foundUser.Password, " ") == strings.Trim(password, " ") {
			// Authentication passed
			u = foundUser
			success = true
		} else {
			// Authentication error
			msg = "Authentication Error"
		}
	}
	return success, u, msg
}

// getUsers
func getUsers() {
	rows, err := server.Server.Db.Query(`SELECT email,password FROM usuario`)
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
		//fmt.Println(email, password)
	}
}

// getUser
func getUser(email string) (models.User, error) {
	u := models.User{}
	var pass, em string
	var minTemp, maxTemp sql.NullFloat64

	sqlStatement := `SELECT email,password,mintemp,maxtemp FROM usuario WHERE email=$1;`
	// Replace 3 with an ID from your database or another random
	// value to test the no rows use case.

	err := server.Server.Db.QueryRow(sqlStatement, email).Scan(&em, &pass, &minTemp, &maxTemp)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		log.Fatal(err)
	default:
		u.Email = em
		u.Password = pass
		if minTemp.Valid {
			u.Settings.MinTemp = minTemp.Float64
		}
		if maxTemp.Valid {
			u.Settings.MaxTemp = maxTemp.Float64
		}
	}
	return u, err
}

// updateUser
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
	_, err := server.Server.Db.Exec(sqlStatement, setValues...)
	if err != nil {
		log.Fatal(err)
	}
	updated = true

	return updated
}

//insertUser
func insertUser(email string, password string) bool {
	var userid string
	//fmt.Println("insertUser", email, password)
	//" + minDefaultTemperature + "','" + maxDefaultTemperature + "
	err := server.Server.Db.QueryRow("INSERT INTO USUARIO (email,password) VALUES('" + email + "','" + password + "') RETURNING email").Scan(&userid)
	if err != nil {
		log.Fatal(err)
		return false
	}

	//fmt.Println("userid created", userid)
	return true
}

func (uc UserController) settingsHandler(w http.ResponseWriter, req *http.Request) {
	var myapp = models.App{}
	if req.Method == http.MethodPost {
		uc.updateHandler(w, req)
	} else {
		templateData := models.TplData{
			Title:     "Settings",
			TabActive: "settings",
		}
		myapp.CurrentView = "settings"
		myapp.Data = templateData
		session, _ := server.Server.Sess.Get(req, "when2runSess")
		// Check if user is authenticated using session
		if auth, ok := session.Values["authenticated"].(bool); ok && auth {
			myapp.UserLogged = session.Values["authenticated"].(bool)
			myapp.User = session.Values["user"].(models.User)
		}
		err := server.Server.Tpl.ExecuteTemplate(w, "settings.gohtml", myapp)
		if err != nil {
			log.Println(err)
		}
	}
}
