package web

import (
    "prflr.org/config"
    "prflr.org/db"
    "prflr.org/structures"
    "labix.org/v2/mgo/bson"
	"encoding/json"
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
    http.HandleFunc("/auth/", authHandler)
    http.HandleFunc("/logout/", logoutHandler)
    http.HandleFunc("/", mainHandler)

    go http.ListenAndServe(config.HTTPPort, nil)
}

/* HTTP Handlers */
func mainHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(config.BaseDir + "web/assets/main.html")
	if err != nil {
	    log.Fatal(err)
	}
	t.Execute(w, nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(config.BaseDir + "web/assets/register.html")
	t.Execute(w, nil)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(config.BaseDir + "web/assets/auth.html")
	t.Execute(w, nil)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/", 301)
}

func lastHandler(w http.ResponseWriter, r *http.Request) {
	/*

	db, err := mgo.Dial(dbHosts)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	db.SetMode(mgo.Monotonic, true)
	dbc := db.DB(dbName).C(dbTimers)

	*/


	session := db.GetConnection()
    db      := session.DB(config.DBName)
    dbc     := db.C(config.DBTimers)

	// Query All
	var results []structures.Timer

	//TODO add criteria builder
	err := dbc.Find(makeCriteria(r.FormValue("filter"))).Sort("-_id").Limit(100).All(&results)
	if err != nil {
		log.Panic(err)
	}

	j, err := json.Marshal(&results)
	if err != nil {
		log.Panic(err)
	}
	fmt.Fprintf(w, "%s", j)

    // @TODO
	session.Close()
}

func aggregateHandler(w http.ResponseWriter, r *http.Request) {
	//TODO
	/*

	db, err := mgo.Dial(dbHosts)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	db.SetMode(mgo.Monotonic, true)
	dbc := db.DB(dbName).C(dbTimers)

	*/


	session := db.GetConnection()
	db      := session.DB(config.DBName)
    dbc     := db.C(config.DBTimers)

	// Query All
	var results []structures.Stat

	grouplist  := make(map[string]interface{})
	groupparam := make(map[string]interface{})

	grouplist["count"] = bson.M{"$sum": 1}
	grouplist["total"] = bson.M{"$sum": "$time"}
	grouplist["min"] = bson.M{"$min": "$time"}
	grouplist["avg"] = bson.M{"$avg": "$time"}
	grouplist["max"] = bson.M{"$max": "$time"}

	q := strings.Split(r.FormValue("groupby"), ",")

	if len(q) >= 1 && q[0] == "src" {
		groupparam["src"] = "$src"
		grouplist["src"] = bson.M{"$first": "$src"}
	}
	if len(q) >= 2 && q[1] == "timer" {
		grouplist["timer"] = bson.M{"$first": "$timer"}
		groupparam["timer"] = "$timer"
	}
	grouplist["_id"] = groupparam
	group := bson.M{"$group": grouplist}
	sort  := bson.M{"$sort": bson.M{r.FormValue("sortby"): -1}}
	match := bson.M{"$match": makeCriteria(r.FormValue("filter"))}
	aggregate := []bson.M{match, group, sort}

	err := dbc.Pipe(aggregate).All(&results)

	if err != nil {
		log.Panic(err)
	}

	j, err := json.Marshal(results)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", j)

	// @TODO
    session.Close()
}

func makeCriteria(filter string) interface{} {
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