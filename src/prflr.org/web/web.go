package web

import (
    "prflr.org/config"
    "prflr.org/user"
    "prflr.org/timer"
    "prflr.org/stringHelper"
    "prflr.org/urlHelper"
    "labix.org/v2/mgo/bson"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"log"
	"strings"
)

func Start() {
    /* Starting Web Server */
    http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(config.BaseDir + "web/assets"))))
    http.Handle("/favicon.ico", http.FileServer(http.Dir(config.BaseDir + "web/assets"))) //cool code for favicon! :) it's very important!
    http.HandleFunc("/last/", lastHandler)
    http.HandleFunc("/aggregate/", aggregateHandler)
    http.HandleFunc("/register/", registerHandler)
    http.HandleFunc("/forgotPassword/", forgotPasswordHandler)
    http.HandleFunc("/resetApiKey/", resetApiKeyHandler)
    http.HandleFunc("/logout/", logoutHandler)
    http.HandleFunc("/", mainHandler)

    go http.ListenAndServe(config.HTTPPort, nil)
}

/* HTTP Handlers */
func mainHandler(w http.ResponseWriter, r *http.Request) {
	user := &user.User{}
	if err := user.GetCurrentUser(r); err != nil {
	    // check for Auth Form Submit
	    email := r.FormValue("email")
	    pass  := r.FormValue("password")

        // auth successful?..
        loginErr := auth(email, pass, w)
	    if loginErr == nil {
	        //log.Print("No error! Redirect!")
            http.Redirect(w, r, urlHelper.GenerateUrl("/"), http.StatusFound)
	    }

	    // ok, no user then show Auth Page
	    t, err := template.ParseFiles(config.BaseDir + "web/assets/auth.html")
	    if err != nil {
            log.Fatal(err)
        }

        tplVars := make(map[string]interface{})
        tplVars["loginErr"] = loginErr

        t.Execute(w, tplVars)
	} else {
	    // we have user!
	    // let's show Web Panel for this user
	    t, err := template.ParseFiles(config.BaseDir + "web/assets/main.html")
    	if err != nil {
    	    log.Fatal(err)
    	}
    	t.Execute(w, user)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(config.BaseDir + "web/assets/register.html")
	t.Execute(w, nil)
}

func forgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(config.BaseDir + "web/assets/forgotPassword.html")
	t.Execute(w, nil)
}

func resetApiKeyHandler(w http.ResponseWriter, r *http.Request) {
    user := &user.User{}

    user.GetCurrentUser(r)

    if len(user.ApiKey) > 0 {
        oldApiKey := user.ApiKey

        // Changing User's ApiKey and Cookies
        if err := user.SetApiKey(stringHelper.RandomString(128), w); err != nil {
            log.Print(err)
        }

        // Updating existing Timers with new ApiKey in order to not lose it!
        timer.SetApiKey(oldApiKey, user.ApiKey)
    }

    // @TODO: make a urlHelper for generating URLs !!!
    http.Redirect(w, r, urlHelper.GenerateUrl("#settings"), http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
    user := &user.User{}

    user.GetCurrentUser(r)
    user.Logout(w)

    http.Redirect(w, r, urlHelper.GenerateUrl("/"), http.StatusFound)
}

func lastHandler(w http.ResponseWriter, r *http.Request) {
    user := &user.User{}
    user.GetCurrentUser(r)

	criteria := makeCriteria(r.FormValue("filter"))
	criteria["apikey"] = user.ApiKey

	results, err := timer.GetList(criteria);
	if err != nil {
		log.Panic(err)
	}

	j, err := json.Marshal(&results)
	if err != nil {
		log.Panic(err)
	}
	fmt.Fprintf(w, "%s", j)
}

func aggregateHandler(w http.ResponseWriter, r *http.Request) {
    user := &user.User{}
    user.GetCurrentUser(r)

    // aggregate query parameters
    groupBy  := make(map[string]interface{})
    sortBy   := r.FormValue("sortby")

    // define criteria for current user
    criteria := makeCriteria(r.FormValue("filter"))
    criteria["apikey"] = user.ApiKey

    // filling in GroupBy parameter
    q := strings.Split(r.FormValue("groupby"), ",")
    if len(q) >= 1 && q[0] == "src" {
        groupBy["src"] = 1
    }
    if len(q) >= 2 && q[1] == "timer" {
        groupBy["timer"] = 1
    }

    results, err := timer.Aggregate(criteria, groupBy, sortBy)
    if err != nil {
        log.Panic(err)
    }

    j, err := json.Marshal(results)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Fprintf(w, "%s", j)
}

func makeCriteria(filter string) map[string]interface{} {
	q := strings.Split(filter, "/")
	c := make(map[string]interface{})

	if len(q) >= 1 && q[0] != "" && q[0] != "*" {
		c["src"] = &bson.RegEx{q[0], "i"}
	}
	if len(q) >= 2 && q[1] != "" && q[1] != "*" {
		c["timer"] = &bson.RegEx{q[1], "i"}
	}
	if len(q) >= 3 && q[2] != "" && q[2] != "*" {
		c["info"] = &bson.RegEx{q[2], "i"}
	}
	if len(q) >= 4 && q[3] != "" && q[3] != "*" {
		c["thrd"] = q[3]
	}
	return c
}

func auth(email, password string, w http.ResponseWriter) error {
    if len(email) == 0 || len(password) == 0 {
        return errors.New("")
    }

    user := &user.User{}
    return user.Auth(email, password, w)
}