package service

import (
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"

	"html/template"

	"github.com/semedi/epfiot/core"
	"github.com/semedi/epfiot/driver"
	"golang.org/x/crypto/bcrypt"
)

// templates:
var frontend *template.Template

type Server struct {
	db *core.DB
	c  driver.Provider
}

func New(drv driver.Provider) *Server {
	var r error
	frontend, r = template.ParseGlob("service/front/*.html")
	if r != nil {
		panic(r)
	}

	s := new(Server)

	database, err := core.NewDB("./db.sqlite")
	if err != nil {
		panic(err)
	}

	err = database.CreateHostdevs(driver.Usb_info())
	if err != nil {
		panic(err)
	}

	s.db = database
	s.c = drv

	return s
}

func render(cond string, file string, w http.ResponseWriter, r *http.Request) {
	conditionsMap := map[string]interface{}{}
	islogged, user := Current(r)
	if islogged == true {
		conditionsMap["Username"] = user
		conditionsMap[cond] = true

		if err := frontend.ExecuteTemplate(w, file, conditionsMap); err != nil {
			log.Println(err)
		}

	} else {

		log.WithFields(log.Fields{"time": time.Now()}).Warn("Trying to access with an expired session!")
		http.Redirect(w, r, "/login", http.StatusForbidden)
	}

}

func ConsoleHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render("Console", "graphql.html", w, r)
	})
}

func ServerPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render("Server", "server.html", w, r)
	})
}

func DocsPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render("Docs", "docs.html", w, r)
	})
}

func DashBoardPageHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render("Dashboard", "dashboard.html", w, r)
	})
}

func LoginPageHandler(res *core.Resolver) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conditionsMap := map[string]interface{}{}

		// check if session is active
		islogged, user := Current(r)

		if islogged == true {

			conditionsMap["Username"] = user
			http.Redirect(w, r, "/dashboard", http.StatusFound)

		} else {

			// verify username and password
			if r.FormValue("Username") != "" && r.FormValue("Password") != "" {
				username := r.FormValue("Username")
				password := r.FormValue("Password")

				hashedPasswordFromDatabase := []byte("$2a$10$4Yhs5bfGgp4vz7j6ScujKuhpRTA4l4OWg7oSukRbyRN7dc.C1pamu")
				if err := bcrypt.CompareHashAndPassword(hashedPasswordFromDatabase, []byte(password)); err != nil {
					log.Println("Either username or password is wrong")
					conditionsMap["LoginError"] = true
				} else {
					if Authenticate(username, res, w, r) {
						http.Redirect(w, r, "/dashboard", http.StatusFound)

						conditionsMap["LoginError"] = false
						conditionsMap["Username"] = username
					} else {
						conditionsMap["LoginError"] = true
					}
				}
			}

			if err := frontend.ExecuteTemplate(w, "login.html", conditionsMap); err != nil {
				log.Println(err)
			}
		}
	})

}

func BootstrapHandler(res *core.Resolver) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()
		thing := r.Form.Get("name")
		log.WithFields(log.Fields{"time": time.Now()}).Info("Thing pairing completed: ", thing)

		err := res.ThingBootstrapped(thing)
		if err != nil {
			log.WithFields(log.Fields{"time": time.Now()}).Error("Error pairing ", thing)
		}
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	Close(w, r)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (s *Server) Run() {

	fileschema, err := ioutil.ReadFile("schema")

	fs := http.FileServer(http.Dir("service/front/static"))

	if err != nil {
		log.Fatalf("failed read schema")
	}

	database := s.db
	r := core.NewResolver(database, s.c)

	schema := graphql.MustParseSchema(string(fileschema), r, graphql.UseStringDescriptions())

	mux := http.NewServeMux()

	mux.Handle("/", DashBoardPageHandler())
	mux.Handle("/console", ConsoleHandler())
	mux.Handle("/login", LoginPageHandler(r))
	mux.Handle("/server", ServerPage())
	mux.Handle("/docs", DocsPage())
	mux.Handle("/logout", http.HandlerFunc(LogoutHandler))
	mux.Handle("/query", authenticated(&relay.Handler{Schema: schema}))
	mux.Handle("/bootstrap", BootstrapHandler(r))

	mux.Handle("/static/", http.StripPrefix("/static", fs))

	log.WithFields(log.Fields{"time": time.Now()}).Info("starting server")

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", logged(mux)))
}

func authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//islogged, user := Current(r)
		islogged, request := Retrieve_session(r)

		if islogged == true {
			next.ServeHTTP(w, request)
		} else {
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
