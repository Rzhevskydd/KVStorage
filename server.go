package main

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"sync"
	"time"
)

var sessions SafeMap = SafeMap{
	mu: sync.RWMutex{},
	v:  make(map[string]interface{}),
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	var val string
	cookie, err := r.Cookie("session_id")

	writeNoLogin := func(w http.ResponseWriter) {
		fmt.Fprintln(w, `<a href="/login">login</a>`)
		fmt.Fprintln(w, "You need to login")
	}

	if err != http.ErrNoCookie {
		var ok bool
		sid, _ := sessions.Get(cookie.Value)
		val, ok = sid.(string)
		if !ok {
			writeNoLogin(w)
			return
		}
	}


	if val != "" {
		fmt.Fprintln(w, `<a href="/logout">logout</a>`)
		fmt.Fprintln(w, "Welcome, "+val)
	} else {
		writeNoLogin(w)
	}
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	name := "Guest"

	sid := uuid.New().String()
	_ = sessions.Put(sid, name)

	expiration := time.Now().Add(30 * time.Minute)
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sid,
		Expires: expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutPage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	_ = sessions.Delete(cookie.Value)
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/logout", logoutPage)
	http.HandleFunc("/", mainPage)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

