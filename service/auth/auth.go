package auth

import (
    "github.com/gorilla/sessions"
	"net/http"
)

var encryptionKey     = "something-very-secret"
var loggedUserSession = sessions.NewCookieStore([]byte(encryptionKey))


func p() {
}

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

func Authenticate(username string, r *http.Request) {
    session, _ := loggedUserSession.New(r, "authenticated-user-session")

    session.Values["username"] = username
    err := session.Save(r, w)

    if err != nil {
            log.Println(err)
    }
}

func Current() (bool, string) {
    session, err := loggedUserSession.Get(r, "authenticated-user-session")
    user := ""
    success := false

    if err != nil {
            log.Println("Unable to retrieve session data!", err)
    }

    if session != nil {
        user = session.Values["username"]
        success = true
    }

   return success, user
}
