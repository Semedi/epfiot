package service

import (
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"

    "golang.org/x/crypto/bcrypt"
    "html/template"
	"github.com/semedi/epfiot/driver"
	"github.com/semedi/epfiot/service/auth"
)


// from const.go:
var dashboardTemplate = template.Must(template.New("").Parse(dashBoardPage))
var logUserTemplate   = template.Must(template.New("").Parse(logUserPage))
var mainTemplate      = template.Must(template.New("").Parse(mainPage))


type Server struct{
    *DB
}

func New() *Server{
    s := new(Server)

	db, err := newDB("./db.sqlite")
	if err != nil {
		panic(err)
	}

    s.DB       = db

    return s
}

func DashBoardPageHandler() http.Handler {

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

        conditionsMap := map[string]interface{}{}


        username = auth.Current()
        //read from session
        session, err := loggedUserSession.Get(r, "authenticated-user-session")

        if err != nil {
                log.Println("Unable to retrieve session data!", err)
        }

        log.Println("Session name : ", session.Name())

        log.Println("Username : ", session.Values["username"])

        conditionsMap["Username"] = session.Values["username"]

        if err := dashboardTemplate.Execute(w, conditionsMap); err != nil {
                log.Println(err)
        }
    })
}

func LoginPageHandler() http.Handler {

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        conditionsMap := map[string]interface{}{}

        // check if session is active
        islogged, user := auth.Current()

        if islogged == true {
            conditionsMap["Username"] = user
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
                        log.Println("Logged in :", username)
                        conditionsMap["Username"] = username
                        conditionsMap["LoginError"] = false

                        auth.Authenticate(username, r)

                        http.Redirect(w, r, "/dashboard", http.StatusFound)
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
        session, err := loggedUserSession.Get(r, "authenticated-user-session")

        if err != nil {
                log.Println("Unable to retrieve session data!", err)
        }

        conditionsMap["Username"] = session.Values["username"]

        if err := mainTemplate.Execute(w, conditionsMap); err != nil {
                log.Println(err)
        }
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
        //read from session
        session, _ := loggedUserSession.Get(r, "authenticated-user-session")

        // remove the username
        session.Values["username"] = ""
        err := session.Save(r, w)

        if err != nil {
                log.Println(err)
        }

        w.Write([]byte("Logged out!"))
}

func (s *Server) Run(drv *driver.Driver) {
	fileschema, err := ioutil.ReadFile("service/schema")

	if err != nil {
		log.Fatalf("failed read schema")
	}

    db := s.DB
    schema := graphql.MustParseSchema(string(fileschema), &Resolver{db: db, drv: drv}, graphql.UseStringDescriptions())

	mux := http.NewServeMux()

	mux.Handle("/", DashBoardPageHandler())
	mux.Handle("/login", LoginPageHandler())
	mux.Handle("/dashboard", http.HandlerFunc(MainHandler))
	mux.Handle("/logout", http.HandlerFunc(LogoutHandler))
	mux.Handle("/query", &relay.Handler{Schema: schema})

	log.WithFields(log.Fields{"time": time.Now()}).Info("starting server")

	log.Fatal(http.ListenAndServe("localhost:8080", logged(mux)))
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


