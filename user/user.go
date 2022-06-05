package user

import (
	"errors"
	"github.com/op/go-logging"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"prflr.org/config"
	"prflr.org/db"
	"prflr.org/helpers"
	"time"
)

var log = logging.MustGetLogger("user")

/**
 * User struct
 */
type User struct {
	Email    string
	Password string
	ApiKey   string
	Token    string
	Name     string
	Info     string
	Status   string
}

func (user *User) GetUserByApiKey(apiKey string) error {
	session, err := db.GetConnection()
	if err != nil {
		return err
	}
	defer session.Close()

	db := session.DB(config.DBName)
	dbc := db.C(config.DBUsers)

	c := make(map[string]interface{})
	c["apikey"] = apiKey

	err = dbc.Find(c).One(&user)

	if err != nil {
		return errors.New("No such a User with given ApiKey")
	}

	return nil
}

func GetUsers() (users []User, err error) {
	session, _ := db.GetConnection()
	defer session.Close()
	db := session.DB(config.DBName)
	dbc := db.C(config.DBUsers)
	users = make([]User, 0)
	err = dbc.Find(nil).All(&users)
	return
}

func (user *User) GetCurrentUser(r *http.Request) error {
	// check Cookies first
	cookie, err := r.Cookie(config.UserCookieName)
	if err != nil {
		return errors.New("No ApiKey passed in Cookies or empty Cookies")
	}

	if len(cookie.Value) <= 0 {
		return errors.New("Empty ApiKey passed in Cookies")
	}

	// getting User from DB via ApiKey
	err = user.GetUserByApiKey(cookie.Value)
	if err != nil {
		return errors.New("No such User with given ApiKey")
	}

	return nil
}

func (user *User) Register() (*User, error) {
	// check if a user with given Email already exists
	if _, err := getUserByEmail(user.Email); err == nil {
		return nil, errors.New("User with such Email already exists! Please provide another one")
	}

	// @TODO: validate user's attributes

	user.Token = user.GenerateToken()
	user.ApiKey = user.GenerateApiKey()
	user.Status = "enabled"

	err := user.Save(true)

	return user, err
}
func Auth(email string, password string, w http.ResponseWriter) (*User, error) {
	// getting User from DB via Email & Password
	user, err := getUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, errors.New("No such user with given Email and Password")
	}

	if user.Password != password {
		return nil, errors.New("No such user with given Email and Password")
	}

	if len(user.Status) > 0 && user.Status != "enabled" {
		return nil, errors.New("Your account is not active at the moment. Please try again later")
	}

	user.saveUserToCookie(w)

	return user, nil
}
func (user *User) Recover() (*User, error) {
	// check if a user with given Email exists
	return getUserByEmail(user.Email)
}

/* NOT EXPORTED!!! */
func getUserByEmail(email string) (*User, error) {
	session, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	db := session.DB(config.DBName)
	dbc := db.C(config.DBUsers)

	c := make(map[string]interface{})
	c["email"] = email

	var user *User
	err = dbc.Find(c).One(&user)

	if err != nil {
		return nil, errors.New("No such a User with given Email and Password")
	}

	return user, nil
}

/* NOT EXPORTED!!! */
func (user *User) saveUserToCookie(w http.ResponseWriter) {
	// setting Cookie to Response
	expire := time.Now().AddDate(0, 0, 1)
	raw := "ApiKey=" + user.ApiKey
	cookie := http.Cookie{
		config.UserCookieName,
		user.ApiKey,
		"/",
		config.DomainName,
		expire,
		expire.Format(time.UnixDate),
		86400,
		false,
		false,
		http.SameSiteDefaultMode,
		raw,
		[]string{raw},
	}
	http.SetCookie(w, &cookie)
}

func (user *User) Logout(w http.ResponseWriter) {
	// setting Cookie to Response
	expire := time.Now().AddDate(0, 0, 1)
	raw := "ApiKey=" + user.ApiKey
	cookie := http.Cookie{
		config.UserCookieName,
		"",
		"/",
		config.DomainName,
		expire,
		expire.Format(time.UnixDate),
		86400,
		false,
		false,
		http.SameSiteDefaultMode,
		raw,
		[]string{raw},
	}
	http.SetCookie(w, &cookie)
}

func (user *User) Save(insertIfNotExists bool) error {
	if len(user.ApiKey) <= 0 {
		return errors.New("Cannot save user: given ApiKey is empty!")
	}

	var err error

	// @TODO: add validation and Error Handling
	session, err := db.GetConnection()
	if err != nil {
		return err
	}
	defer session.Close()

	db := session.DB(config.DBName)
	dbc := db.C(config.DBUsers)

	selector := make(map[string]interface{})
	selector["apikey"] = user.ApiKey

	if insertIfNotExists {
		_, err = dbc.Upsert(selector, user)
	} else {
		err = dbc.Update(selector, user)
	}

	return err
}

func (user *User) SetApiKey(apiKey string, w http.ResponseWriter) error {
	if len(user.ApiKey) <= 0 {
		return errors.New("Cannot save user: given ApiKey is empty!")
	}
	if len(apiKey) <= 0 {
		return errors.New("Cannot save empty ApiKey!")
	}

	// @TODO: add validation and Error Handling
	session, err := db.GetConnection()
	if err != nil {
		return err
	}
	defer session.Close()

	db := session.DB(config.DBName)
	dbc := db.C(config.DBUsers)

	selector := make(map[string]interface{})
	selector["apikey"] = user.ApiKey

	modifier := &bson.M{"$set": bson.M{"apikey": apiKey}}

	err = dbc.Update(selector, modifier)

	// saving new User's ApiKey to Cookies
	user.ApiKey = apiKey
	user.saveUserToCookie(w)

	return err
}

func (user *User) GenerateToken() string {
	// @TODO: use md5(microtime())
	return helpers.RandomString(32)
}
func (user *User) GenerateApiKey() string {
	// @TODO: check if Token is given
	// @TODO: use User.Token + md5(microtime())
	return helpers.RandomString(32)
}

func (user *User) CreatePrivateStorage() {
	collectionName := helpers.GetCappedCollectionNameForApiKey(user.ApiKey)

	session, err := db.GetConnection()
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer session.Close()

	// creating capped collection
	err = db.CreateCappedCollection(session.DB(config.DBName).C(collectionName), config.CappedCollectionMaxByte, config.CappedCollectionMaxDocs)
	if err != nil {
		log.Error(err.Error())
	}
}

func (user *User) RemovePrivateStorage() {
	collectionName := helpers.GetCappedCollectionNameForApiKey(user.ApiKey)
	RemoveStorage(collectionName)
}

func RemoveStorage(c string) {
	session, err := db.GetConnection()
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer session.Close()

	//  drop capped collection
	err = session.DB(config.DBName).C(c).DropCollection()
	if err != nil {
		log.Error(err.Error())
	}
}

func (user *User) RemovePrivateStorageData() {
	// remove current
	user.RemovePrivateStorage()

	// create a new one
	user.CreatePrivateStorage()
}
