package user

import(
    "labix.org/v2/mgo/bson"
    "prflr.org/config"
    "prflr.org/db"
    "net/http"
    "errors"
    "time"
    //"log"
    //"fmt"
)

/**
 * User struct
 */
type User struct {
    Email    string
    Password string
    ApiKey   string
    Token    float32
    Name     string
    Info     string
}

func (user *User) GetUser(email string, password string) error {
    session := db.GetConnection()
    db      := session.DB(config.DBName)
    dbc     := db.C(config.DBUsers)

    c := make(map[string]interface{})
    c["email"]    = email
    c["password"] = password // @TODO: add md5 ?

    err := dbc.Find(c).One(&user)

    session.Close()

    if err != nil {
        return errors.New("No such a User with given Email and Password")
    }

    return nil
}

func (user *User) GetUserByApiKey(apiKey string) error {
    session := db.GetConnection()
    db      := session.DB(config.DBName)
    dbc     := db.C(config.DBUsers)

    c := make(map[string]interface{})
    c["apikey"] = apiKey

    err := dbc.Find(c).One(&user)

    session.Close()

    if err != nil {
        return errors.New("No such a User with given ApiKey")
    }

    return nil
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

func (user *User) Auth(email string, password string, w http.ResponseWriter) error {
    // getting User from DB via Email & Password
    err := user.GetUser(email, password)
    if err != nil {
        return err
    }

    user.saveUserToCookie(w)

    return nil
}

/* NOT EXPORTED!!! */
func (user *User) saveUserToCookie(w http.ResponseWriter) {
    // setting Cookie to Response
    expire := time.Now().AddDate(0, 0, 1)
    raw    := "ApiKey="+user.ApiKey
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
        raw,
        []string{raw},
    }
    http.SetCookie(w, &cookie)
}

func (user *User) Logout(w http.ResponseWriter) {
    // setting Cookie to Response
    expire := time.Now().AddDate(0, 0, 1)
    raw    := "ApiKey="+user.ApiKey
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
        raw,
        []string{raw},
    }
    http.SetCookie(w, &cookie)
}

func (user *User) Save() error {
    if len(user.ApiKey) <= 0 {
        return errors.New("Cannot save user: given ApiKey is empty!")
    }

    // @TODO: add validation and Error Handling
    session := db.GetConnection()
    db      := session.DB(config.DBName)
    dbc     := db.C(config.DBUsers)

    selector := make(map[string]interface{})
    selector["apikey"] = user.ApiKey

    err := dbc.Update(selector, user)

    session.Close()

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
    session := db.GetConnection()
    db      := session.DB(config.DBName)
    dbc     := db.C(config.DBUsers)

    selector := make(map[string]interface{})
    selector["apikey"] = user.ApiKey

    modifier := &bson.M{"$set": bson.M{"apikey": apiKey}}

    err := dbc.Update(selector, modifier)

    // saving new User's ApiKey to Cookies
    user.ApiKey = apiKey
    user.saveUserToCookie(w)

    session.Close()

    return err
}