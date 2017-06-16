package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/pintokrysler/when2run/controllers"
	"github.com/pintokrysler/when2run/models"
	"github.com/pintokrysler/when2run/server"
)

//Myapp ...
var Myapp = models.App{
	UserLogged: false,
}

var appConf models.AppConfiguration

func init() {
	// Read conf file
	configurationFile := "config.json"
	file, err := os.Open(configurationFile)
	if err != nil {
		log.Println(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&appConf)
	if err != nil {
		log.Println(err)
	}
	// Init server
	server.InitServer(appConf)
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
	//router.GET("/favicon.ico", http.NotFoundHandler())
	// http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":"+strconv.Itoa(appConf.Port), router)
}
