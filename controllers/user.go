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
	fmt.Println("User Dispatch!")
	fmt.Println(action)
	switch action {
	case "login":
		uc.loginHandler(w, r)
		break
	case "create":
		uc.create(w, r)
		break
	case "logout":
		uc.logout(w, r)
		break
	}

}

// CreateUser creates a new user instance
func (uc UserController) create(w http.ResponseWriter, r *http.Request) {
	// u := models.User{}

}

// UpdateUser updates an existing user instance
func (uc UserController) update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// u := models.User{}

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
		session, _ := server.Server.Sess.Get(req, "when2runSess")

		// Check if user is authenticated using session
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			fmt.Println("not session auth")
			email := req.FormValue("email")
			password := req.FormValue("password")
			loggedIn, loggedUser, errMsg := login(email, password)
			if loggedIn {
				session.Values["authenticated"] = true
				session.Values["user"] = loggedUser

			} else {
				log.Println(errMsg)
				session.Values["authenticated"] = false
				session.Values["user"] = loggedUser
				myapp.MsgError = errMsg
			}
			session.Save(req, w)
		}
	}

	myapp.Data = templateData
	err := server.Server.Tpl.ExecuteTemplate(w, "account.gohtml", myapp)
	if err != nil {
		log.Println(err)
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
			//myapp.MsgError = "Authentication Error"
		}
	}
	return success, u, msg
}

func (uc UserController) logout(w http.ResponseWriter, r *http.Request) {
	// u := models.User{}

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
		fmt.Println(email, password)
	}
}

// getUser
func getUser(email string) (models.User, error) {
	u := models.User{}
	var pass, em string
	var minTemp, maxTemp float64

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
		u.Settings.MinTemp = minTemp
		u.Settings.MaxTemp = maxTemp
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
	//" + minDefaultTemperature + "','" + maxDefaultTemperature + "
	err := server.Server.Db.QueryRow("INSERT INTO USUARIO (email,password) VALUES('" + email + "','" + password + "') RETURNING email").Scan(&userid)
	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("userid created", userid)
	return true
}
