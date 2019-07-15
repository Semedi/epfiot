package auth

import (
    "github.com/gorilla/sessions"
	"net/http"
	log "github.com/sirupsen/logrus"
    "fmt"
)

var encryptionKey     = "something-very-secret"
var loggedUserSession = sessions.NewCookieStore([]byte(encryptionKey))


func init() {
    loggedUserSession.Options = &sessions.Options{
         // change domain to match your machine. Can be localhost
         // IF the Domain name doesn't match, your session will be EMPTY!
         Domain:   "localhost",
         Path:     "/",
         MaxAge:   3600 * 3, // 3 hours
         HttpOnly: true,
    }
}

func Authenticate(username string, w http.ResponseWriter, r *http.Request) {
    session, _ := loggedUserSession.New(r, "authenticated-user-session")

    session.Values["username"] = username
    err := session.Save(r, w)

    if err != nil {
            log.Println(err)
    }
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
    user := ""
    success := false

    if err != nil {
            log.Println("Unable to retrieve session data!", err)
    }

    if session != nil {
        user = fmt.Sprintf("%v", session.Values["username"])
        success = true

    }

   return success, user
}
