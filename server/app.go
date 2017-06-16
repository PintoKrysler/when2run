package server

import (
	"database/sql"
	"encoding/gob"
	"html/template"
	"log"

	"github.com/gorilla/sessions"
	"github.com/pintokrysler/when2run/models"
)

//Server ..
var Server models.Server

func init() {
	gob.Register(models.User{})
}

func main() {

}

// InitServer ..
func InitServer(conf models.AppConfiguration) {
	// Init session
	Server.Sess = sessions.NewCookieStore([]byte("something-very-secret"))
	Server.Tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	connectDB(conf.DB)

}

// connectDB ...
func connectDB(dbconf models.DBConfiguration) {
	databaseConn, err := sql.Open("postgres", "user="+dbconf.User+" password="+dbconf.Password+" dbname="+dbconf.Name+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	Server.Db = databaseConn
}
