package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/pintokrysler/when2run/controllers"
	"github.com/pintokrysler/when2run/models"
	"github.com/pintokrysler/when2run/server"
)

const (
	minDefaultTemperature float64 = 0
	maxDefaultTemperature float64 = 110
)

//Myapp ...
var Myapp = models.App{
	UserLogged: false,
}

func init() {
	// Init server
	server.InitServer()
}

func main() {

	router := httprouter.New()

	// Create controllers
	userController := controllers.NewUserController()
	indexController := controllers.NewIndexController()

	router.GET("/", indexController.Dispatch)
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.GET("/user/:action", userController.Dispatch)
	router.POST("/user/:action", userController.Dispatch)
	// router.GET("/settings", settingsHandler)
	// router.POST("/settings", settingsHandler)
	//router.GET("/favicon.ico", http.NotFoundHandler())
	// http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", router)
}

// func settingsHandler(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
// 	var myapp = models.App{}
// 	if req.Method == http.MethodPost {
// 		setSettingsHandler(w, req)
// 	} else {
// 		templateData := models.TplData{
// 			Title:     "Settings",
// 			TabActive: "settings",
// 		}
// 		myapp.CurrentView = "settings"
// 		myapp.Data = templateData
// 		fmt.Println(myapp)
// 		err := server.Server.Tpl.ExecuteTemplate(w, "settings.gohtml", myapp)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}
// }

// func createUserHandler(w http.ResponseWriter, req *http.Request) {
// 	var myapp = models.App{}
// 	fmt.Println("createuserhandler")
// 	fmt.Println(req.URL.Path)
// 	templateData := models.TplData{
// 		Title:     "Create Account",
// 		TabActive: "account",
// 	}
//
// 	// post request, http.MethodPost is a constant
// 	if req.Method == http.MethodPost {
// 		email := req.FormValue("email")
// 		password := req.FormValue("password")
//
// 		// Check if there is a user with that email
// 		_, err := getUser(email)
// 		if err != nil {
// 			if err == sql.ErrNoRows {
// 				// User does not exists
// 				if inserted := insertUser(email, password); inserted {
// 					myapp.UserLogged = true
// 					userSettings := models.Settings{}
// 					userSettings = userSettings.New()
// 					myapp.User = models.User{Email: email, Password: password, Settings: userSettings}
//
// 				}
// 				myapp.Data = templateData
// 				err := server.Server.Tpl.ExecuteTemplate(w, "account.gohtml", myapp)
// 				if err != nil {
// 					log.Println(err)
// 				}
// 			}
// 		} else {
// 			// User already exists
// 			myapp.MsgError = "User account already exists"
// 			err := server.Server.Tpl.ExecuteTemplate(w, "createaccount.gohtml", myapp)
// 			if err != nil {
// 				log.Println(err)
// 			}
// 		}
//
// 	} else {
// 		// Is not a POST request! we only want to show the form to create a new account
// 		myapp.Data = templateData
// 		err := server.Server.Tpl.ExecuteTemplate(w, "createaccount.gohtml", myapp)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}
// }

// func setSettingsHandler(w http.ResponseWriter, req *http.Request) {
// 	var myapp = models.App{}
// 	templateData := models.TplData{
// 		Title:     "Your Running Times",
// 		TabActive: "times",
// 	}
// 	var maxTempVal float64
// 	var minTempVal float64
//
// 	s := models.Settings{MinTemp: minDefaultTemperature, MaxTemp: maxDefaultTemperature}
// 	myapp.CurrentView = "settings"
// 	minTemp := req.FormValue("minTemp")
// 	maxTemp := req.FormValue("maxTemp")
//
// 	// if maxTemp was passed on the form
// 	if maxTemp != "" {
// 		maxTempVal, _ = strconv.ParseFloat(maxTemp, 64)
// 		s.MaxTemp = maxTempVal
// 	}
//
// 	if minTemp != "" {
// 		minTempVal, _ = strconv.ParseFloat(minTemp, 64)
// 		s.MinTemp = minTempVal
// 	}
//
// 	if myapp.UserLogged && (minTemp != "" || maxTemp != "") {
// 		// Update database if user is logged in and temperatures passed on the form
// 		// are different than DB values
// 		changedValues := make(map[string]interface{})
// 		// Check if maxTemp changed
// 		if s.MaxTemp != myapp.User.Settings.MaxTemp {
// 			changedValues["MaxTemp"] = maxTempVal
// 			myapp.User.Settings.MaxTemp = maxTempVal
// 		}
//
// 		// Check if minTemp changed
// 		//fmt.Println("minChanges", minTempVal, myapp.User.Settings.MinTemp)
// 		if s.MinTemp != myapp.User.Settings.MinTemp {
// 			changedValues["MinTemp"] = minTempVal
// 			myapp.User.Settings.MinTemp = minTempVal
// 		}
//
// 		if len(changedValues) > 0 {
// 			fmt.Println("Update user information ", minTemp, maxTemp)
// 			updateSuccess := updateUser(myapp.User.Email, changedValues)
// 			if updateSuccess {
// 			}
// 		}
// 	}
//
// 	var data = makeWeatherAPIcall(s)
// 	templateData.Data = data
// 	myapp.Data = templateData
// 	fmt.Println(myapp)
// 	err := server.Server.Tpl.ExecuteTemplate(w, "times.gohtml", myapp)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }
