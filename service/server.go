package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"

	"github.com/semedi/epfiot/driver"
	"github.com/semedi/epfiot/model"
	"golang.org/x/crypto/bcrypt"
	"html/template"
)

// from const.go:
var dashboardTemplate = template.Must(template.New("").Parse(dashBoardPage))
var logUserTemplate = template.Must(template.New("").Parse(logUserPage))
var mainTemplate = template.Must(template.New("").Parse(mainPage))

type Server struct {
	db *model.DB
}

func New() *Server {
	s := new(Server)

	database, err := model.NewDB("./db.sqlite")
	if err != nil {
		panic(err)
	}

	s.db = database

	return s
}

func DashBoardPageHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		conditionsMap := map[string]interface{}{}

		islogged, user := Current(r)

		if islogged == true {
			log.Println("Username : ", user)
			conditionsMap["Username"] = user
		}

		if err := dashboardTemplate.Execute(w, conditionsMap); err != nil {
			log.Println(err)
		}
	})
}

func LoginPageHandler(res *model.Resolver) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conditionsMap := map[string]interface{}{}

		// check if session is active
		islogged, user := Current(r)

		if islogged == true {
			conditionsMap["Username"] = user
			log.Println("entro en logeado")
		}

		// verify username and password
		if r.FormValue("Login") != "" && r.FormValue("Username") != "" {
			username := r.FormValue("Username")
			password := r.FormValue("Password")

			// NOTE: here is where you want to query your database to retrieve the hashed password
			// for username.
			// For this tutorial and simplicity sake, we will simulate the retrieved hashed password
			// as $2a$10$4Yhs5bfGgp4vz7j6ScujKuhpRTA4l4OWg7oSukRbyRN7dc.C1pamu
			// the plain password is 'mynakedpassword'
			// see https://www.socketloop.com/tutorials/golang-bcrypting-password for more details
			// on how to generate bcrypted password
			hashedPasswordFromDatabase := []byte("$2a$10$4Yhs5bfGgp4vz7j6ScujKuhpRTA4l4OWg7oSukRbyRN7dc.C1pamu")
			if err := bcrypt.CompareHashAndPassword(hashedPasswordFromDatabase, []byte(password)); err != nil {
				log.Println("Either username or password is wrong")
				conditionsMap["LoginError"] = true
			} else {
                if (Authenticate(username, res, w, r)){
                    http.Redirect(w, r, "/dashboard", http.StatusFound)

				    conditionsMap["LoginError"] = false
				    conditionsMap["Username"]   = username
                } else {
                    conditionsMap["LoginError"] = true
                }
			}
		}

		if err := logUserTemplate.Execute(w, conditionsMap); err != nil {
			log.Println(err)
		}
	})
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	conditionsMap := map[string]interface{}{}

	//read from session
	islogged, user := Current(r)
	fmt.Println("putaaaa")

	if islogged == true {
		log.Println("Username : ", user)
		conditionsMap["Username"] = user
	}

	if err := mainTemplate.Execute(w, conditionsMap); err != nil {
		log.Println(err)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	Close(w, r)

	w.Write([]byte("Logged out!"))
}

func (s *Server) Run(drv *driver.Driver) {
	fileschema, err := ioutil.ReadFile("model/schema")

	if err != nil {
		log.Fatalf("failed read schema")
	}

	database := s.db
	r := &model.Resolver{Db: database, Drv: drv}

	schema := graphql.MustParseSchema(string(fileschema), r, graphql.UseStringDescriptions())

	mux := http.NewServeMux()

	mux.Handle("/", DashBoardPageHandler())
	mux.Handle("/login", LoginPageHandler(r))
	mux.Handle("/dashboard", http.HandlerFunc(MainHandler))
	mux.Handle("/logout", http.HandlerFunc(LogoutHandler))
	mux.Handle("/query", authenticated(&relay.Handler{Schema: schema}))

	log.WithFields(log.Fields{"time": time.Now()}).Info("starting server")

	log.Fatal(http.ListenAndServe("localhost:8080", logged(mux)))
}

func authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//islogged, user := Current(r)
        islogged, request := Retrieve_session(r)

		if islogged == true {
		    next.ServeHTTP(w, request)
        }else{
            http.Redirect(w, r, "/login", http.StatusForbidden)
        }
	})
}

// logging middleware
func logged(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()

		next.ServeHTTP(w, r)

		log.WithFields(log.Fields{
			"path":    r.RequestURI,
			"IP":      r.RemoteAddr,
			"elapsed": time.Now().UTC().Sub(start),
		}).Info()
	})
}
