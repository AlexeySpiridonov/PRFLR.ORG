package web

import (
	"PRFLR.ORG/config"
	"PRFLR.ORG/helpers"
	"PRFLR.ORG/mailer"
	"PRFLR.ORG/timer"
	"PRFLR.ORG/user"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/op/go-logging"
	"html/template"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var log = logging.MustGetLogger("web")

// compile all templates and cache them
//var templates = template.Must(template.ParseGlob(config.BaseDir + "assets/landing/*.html"))

func Start() {
	/* Starting Web Server */
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(config.BaseDir+"assets"))))
	http.Handle("/assets/landing/", http.StripPrefix("/assets/landing/", http.FileServer(http.Dir(config.BaseDir+"assets/landing"))))
	http.Handle("/favicon.ico", http.FileServer(http.Dir(config.BaseDir+"assets"))) //cool code for favicon! :) it's very important!
	http.HandleFunc("/last/", lastHandler)
	http.HandleFunc("/aggregate/", aggregateHandler)
	http.HandleFunc("/graph/", graphHandler)
	http.HandleFunc("/signup/", registerHandler)
	http.HandleFunc("/signin/", loginHandler)
	http.HandleFunc("/forgotPassword/", forgotPasswordHandler)
	http.HandleFunc("/passwordRecovered/", passwordRecoveredHandler)
	http.HandleFunc("/thankyou/", thankyouHandler)
	http.HandleFunc("/resetApiKey/", resetApiKeyHandler)
	http.HandleFunc("/removeData/", removeDataHandler)
	http.HandleFunc("/logout/", logoutHandler)
	http.HandleFunc("/", mainHandler)

	go http.ListenAndServe(config.HTTPPort, nil)
}

/* HTTP Handlers */
func mainHandler(w http.ResponseWriter, r *http.Request) {
	tplVars := make(map[string]interface{})

	user := &user.User{}
	if err := user.GetCurrentUser(r); err != nil {
		//log.Notice(err.Error())
		// check for Auth Form Submit
		email := r.PostFormValue("email")
		pass := r.PostFormValue("password")

		// auth successful?..
		loginErr := auth(email, pass, w)
		if loginErr == nil {
			log.Error(loginErr.Error())
			http.Redirect(w, r, helpers.GenerateUrl("/"), http.StatusFound)
		}

		// ok, no user then show Auth Page
		t, err := template.ParseFiles(
			config.BaseDir+"assets/index.html",
			config.BaseDir+"assets/landing/header.html",
			config.BaseDir+"assets/landing/footer.html",
		)
		if err != nil {
			log.Error(err.Error())
			return
		}

		tplVars["loginErr"] = loginErr

		t.Execute(w, tplVars)
	} else {
		// we have user!
		// let's show Web Panel for this user
		t, err := template.ParseFiles(config.BaseDir + "assets/main.html")
		if err != nil {
			log.Error(err.Error())
		}

		tplVars["user"] = user
		tplVars["ApiKey"] = helpers.GetApiIDForApiKey(user.ApiKey)

		t.Execute(w, tplVars)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// check for Register Form Submit
	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	pass := r.PostFormValue("pass")
	confirmPass := r.PostFormValue("confirm_pass")
	info := r.PostFormValue("info")

	registerAttempt := r.PostFormValue("register")

	tplVars := make(map[string]interface{})
	tplVars["name"] = name
	tplVars["email"] = email
	tplVars["pass"] = pass
	tplVars["confirm_pass"] = confirmPass
	tplVars["info"] = info

	if len(registerAttempt) > 0 {
		user, registerErr := register(name, email, pass, confirmPass, info)
		if registerErr == nil {
			// Creating Capped Collection for this User
			go user.CreatePrivateStorage()

			// Sending Email Notifications
			//go sendRegistrationEmail(user)
			go sendRegistrationEmail(user)

			// Getting Hell out of here!!! Whheeeee!!!!111 =)
			http.Redirect(w, r, helpers.GenerateUrl("/thankyou"), http.StatusFound)
		}
		tplVars["registerErr"] = registerErr
	}

	// ok, no user then show Auth Page
	t, err := template.ParseFiles(
		config.BaseDir+"assets/register.html",
		config.BaseDir+"assets/landing/header.html",
		config.BaseDir+"assets/landing/footer.html",
	)
	if err != nil {
		log.Error(err.Error())
		return
	}

	t.Execute(w, tplVars)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// check for Auth Form Submit
	email := r.PostFormValue("email")
	pass := r.PostFormValue("password")

	loginAttempt := r.PostFormValue("login")

	tplVars := make(map[string]interface{})

	if len(loginAttempt) > 0 {
		// auth successful?..
		loginErr := auth(email, pass, w)
		if loginErr == nil {
			http.Redirect(w, r, helpers.GenerateUrl("/"), http.StatusFound)
		} else {
			log.Error(loginErr.Error())
		}

		tplVars["loginErr"] = loginErr
	}

	// ok, no user then show Auth Page
	t, err := template.ParseFiles(
		config.BaseDir+"assets/login.html",
		config.BaseDir+"assets/landing/header.html",
		config.BaseDir+"assets/landing/footer.html",
	)
	if err != nil {
		log.Error(err.Error())
		return
	}

	t.Execute(w, tplVars)
}

func forgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// check for Recovery Form Submit
	email := r.PostFormValue("email")
	recoveryAttempt := r.PostFormValue("recovery")

	tplVars := make(map[string]interface{})
	tplVars["email"] = email

	if len(recoveryAttempt) > 0 {
		user, recoveryErr := recoverPassword(email)
		if recoveryErr == nil {
			go sendRecoveryEmail(user)
			http.Redirect(w, r, helpers.GenerateUrl("/passwordRecovered"), http.StatusFound)
		}
		tplVars["recoveryErr"] = recoveryErr
	}

	t, err := template.ParseFiles(
		config.BaseDir+"assets/forgotPassword.html",
		config.BaseDir+"assets/landing/header.html",
		config.BaseDir+"assets/landing/footer.html",
	)
	if err != nil {
		log.Error(err.Error())
		return
	}

	t.Execute(w, tplVars)
}
func passwordRecoveredHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(
		config.BaseDir+"assets/passwordRecovered.html",
		config.BaseDir+"assets/landing/header.html",
		config.BaseDir+"assets/landing/footer.html",
	)
	t.Execute(w, nil)
}

func thankyouHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(
		config.BaseDir+"assets/thankyou.html",
		config.BaseDir+"assets/landing/header.html",
		config.BaseDir+"assets/landing/footer.html",
	)
	t.Execute(w, nil)
}

func resetApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	user := &user.User{}
	user.GetCurrentUser(r)

	if len(user.ApiKey) > 0 {
		//oldApiKey := user.ApiKey

		// Changing User's ApiKey and Cookies
		if err := user.SetApiKey(user.GenerateApiKey(), w); err != nil {
			log.Error(err.Error())
		}

		// Updating existing Timers with new ApiKey in order to not lose it!
		//timer.SetApiKey(oldApiKey, user.ApiKey)
	}

	// @TODO: make a helpers for generating URLs !!!
	http.Redirect(w, r, helpers.GenerateUrl("#settings"), http.StatusFound)
}

func removeDataHandler(w http.ResponseWriter, r *http.Request) {
	user := &user.User{}
	if err := user.GetCurrentUser(r); err == nil {
		//io.WriteString(w, "User: " + user.Email)

		user.RemovePrivateStorageData()
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	user := &user.User{}

	user.GetCurrentUser(r)
	user.Logout(w)

	http.Redirect(w, r, helpers.GenerateUrl("/"), http.StatusFound)
}

func lastHandler(w http.ResponseWriter, r *http.Request) {
	user := &user.User{}
	user.GetCurrentUser(r)

	criteria := makeCriteria(r.FormValue("filter"))

	results, err := timer.GetList(user.ApiKey, criteria)
	if err != nil {
		log.Error(err.Error())
		return
	}

	j, err := json.Marshal(&results)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Output JSON!
	fmt.Fprintf(w, "%s", j)
}

func aggregateHandler(w http.ResponseWriter, r *http.Request) {
	user := &user.User{}
	user.GetCurrentUser(r)

	// aggregate query parameters
	groupBy := make(map[string]interface{})
	sortBy := r.FormValue("sortby")

	// define criteria for current user
	criteria := makeCriteria(r.FormValue("filter"))

	// filling in GroupBy parameter
	q := strings.Split(r.FormValue("groupby"), ",")
	if len(q) >= 1 && q[0] == "src" {
		groupBy["src"] = 1
	}
	if len(q) >= 2 && q[1] == "timer" {
		groupBy["timer"] = 1
	}

	results, err := timer.Aggregate(user.ApiKey, criteria, groupBy, sortBy)
	if err != nil {
		log.Error(err.Error())
		return
	}

	j, err := json.Marshal(results)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Output JSON!
	fmt.Fprintf(w, "%s", j)
}

func graphHandler(w http.ResponseWriter, r *http.Request) {
	user := &user.User{}
	user.GetCurrentUser(r)

	// define criteria for current user
	criteria := makeCriteria(r.FormValue("filter"))

	// From / To Period Criteria
	start := r.FormValue("start")
	end := r.FormValue("end")
	var startTime, endTime time.Time
	var startTimeErr, endTimeErr error
	if len(start) > 0 {
		startTime, startTimeErr = time.Parse("02/01/2006 15:04:05", start)
	} else {
		startTimeErr = errors.New("Value is empty")
	}
	if len(end) > 0 {
		endTime, endTimeErr = time.Parse("02/01/2006 15:04:05", end)
	} else {
		endTimeErr = errors.New("Value is empty")
	}

	if startTimeErr == nil && endTimeErr == nil {
		criteria["timestamp"] = &bson.M{"$gte": startTime.Unix(), "$lte": endTime.Unix()}
		//fmt.Println(startTime, endTime)
	} else if startTimeErr == nil {
		criteria["timestamp"] = &bson.M{"$gte": startTime.Unix()}
	} else if endTimeErr == nil {
		criteria["timestamp"] = &bson.M{"$lte": endTime.Unix()}
	}

	graph, err := timer.FormatGraph(user.ApiKey, criteria)
	if err != nil {
		log.Error(err.Error())
		graph = &timer.Graph{}
	}

	j, err := json.Marshal(graph)
	if err != nil {
		log.Error(err.Error())
		return
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

func register(name, email, pass, confirmPass, info string) (*user.User, error) {
	if len(name) == 0 {
		return nil, errors.New("Please specify your Full Name")
	}
	if len(email) == 0 {
		return nil, errors.New("Please specify your Email address")
	}
	if len(pass) == 0 {
		return nil, errors.New("Please specify your Password")
	}
	if len(confirmPass) == 0 {
		return nil, errors.New("Please re-enter your Password")
	}
	if pass != confirmPass {
		return nil, errors.New("Password does not match to Confirmed Password")
	}

	user := &user.User{
		Name:     name,
		Email:    email,
		Password: pass,
		Info:     info,
	}

	return user.Register()
}
func recoverPassword(email string) (*user.User, error) {
	if len(email) == 0 {
		return nil, errors.New("Please specify your Email address")
	}

	user := &user.User{
		Email: email,
	}

	return user.Recover()
}

func auth(email, password string, w http.ResponseWriter) error {
	if len(email) == 0 || len(password) == 0 {
		return errors.New("Empty email or password")
	}

	_, err := user.Auth(email, password, w)

	return err
}

func sendRegistrationEmail(user *user.User) error {
	msg := "Greetings!\n\nThank you for your decision to make your projects performance even better!\n\n" +
		"Please find all the service information below:\n\n" +

		"Email: " + user.Email + "\n" +
		"Pass: " + user.Password + "\n" +
		"API Key: " + helpers.GetApiIDForApiKey(user.ApiKey) + "\n\n" +

		"Following links are for the SDK Integration into your application:\n\n" +

		"SDK: https://github.com/PRFLR/SDK\n" +
		"WebPanel: http://prflr.org\n" +
		"Tutorials: https://github.com/PRFLR/SDK/wiki\n\n" +

		"Good luck in neverending fight for milliseconds!\n" +
		"PRFLR Team © " + strconv.Itoa(time.Now().Year()) + ", info@prflr.org\n\n" +
		"Join our G+ Community: http://goo.gl/AqJV4V"

	// sending to the User
	mail := &mailer.Email{
		From:    config.RegisterEmailFrom,
		To:      user.Email,
		Subject: config.RegisterEmailSubject,
		Msg:     msg,
	}

	err := mail.Send()
	if err != nil {
		log.Error(err.Error())
	}

	// Sending to PRFLR Team!
	mail = &mailer.Email{
		From:    config.RegisterEmailFrom,
		To:      config.RegisterEmailTo,
		Subject: config.RegisterEmailSubject,
		Msg:     msg,
	}

	err = mail.Send()
	if err != nil {
		log.Error(err.Error())
	}

	return nil
}

func sendRecoveryEmail(user *user.User) error {
	msg := "Greetings!\n\n" +

		"Your Pass: " + user.Password + "\n\n" +

		"Please try to login at: http://prflr.org\n\n" +

		"Good luck in neverending fight for milliseconds!\n" +
		"PRFLR Team © " + strconv.Itoa(time.Now().Year()) + ", info@prflr.org\n\n" +
		"Join our G+ Community: http://goo.gl/AqJV4V"

	// sending to the User
	mail := &mailer.Email{
		From:    config.RecoveryEmailFrom,
		To:      user.Email,
		Subject: config.RecoveryEmailSubject,
		Msg:     msg,
	}

	err := mail.Send()
	if err != nil {
		log.Error(err.Error())
	}

	return nil
}
