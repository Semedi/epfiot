package service

import (
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/semedi/epfiot/core"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var encryptionKey = "something-very-secret"
var loggedUserSession = sessions.NewCookieStore([]byte(encryptionKey))

func init() {
	loggedUserSession.Options = &sessions.Options{
		// change domain to match your machine. Can be localhost
		// IF the Domain name doesn't match, your session will be EMPTY!
		Path:     "/",
		MaxAge:   3600 * 3, // 3 hours
		HttpOnly: true,
	}
}

func Authenticate(username string, res *core.Resolver, w http.ResponseWriter, r *http.Request) bool {
	session, _ := loggedUserSession.New(r, "authenticated-user-session")

	u, err := res.Db.Find_user(username)

	if err != nil || u == nil {
		return false
	}

	session.Values["username"] = username
	session.Values["userid"] = u.ID

	session.Save(r, w)

	return true
}

func Close(w http.ResponseWriter, r *http.Request) {
	//read from session
	session, _ := loggedUserSession.Get(r, "authenticated-user-session")

	// remove the username
	session.Values["username"] = ""
	err := session.Save(r, w)

	if err != nil {
		log.Println(err)
	}
}

func Current(r *http.Request) (bool, string) {
	session, err := loggedUserSession.Get(r, "authenticated-user-session")
	logged_user := ""
	success := false

	if err != nil {
		log.Println("Unable to retrieve session data!", err)
	}

	if session != nil {
		logged_user = fmt.Sprintf("%v", session.Values["username"])
		success = logged_user != "<nil>" && logged_user != ""
	}

	return success, logged_user
}

func Retrieve_session(r *http.Request) (bool, *http.Request) {
	session, err := loggedUserSession.Get(r, "authenticated-user-session")
	logged_user := ""

	if err != nil {
		log.Println("Unable to retrieve session data!", err)
		return false, nil
	}

	if session != nil {
		logged_user = fmt.Sprintf("%v", session.Values["username"])
		log.Println("", logged_user)

		//ctx := context.WithValue(r.Context(), "user", logged_user)
		ctx := context.WithValue(r.Context(), "userid", session.Values["userid"])

		return true, r.WithContext(ctx)
	}

	return false, nil
}
