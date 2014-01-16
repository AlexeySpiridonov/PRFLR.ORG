package user

import(
    "labix.org/v2/mgo/bson"
    "prflr.org/stringHelper"
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
    Token    string
    Name     string
    Info     string
    Status   string
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

func (user *User) Register() (*User, error) {
    // check if a user with given Email already exists
    if _, err := getUserByEmail(user.Email); err == nil {
        return nil, errors.New("User with such Email already exists! Please provide another one")
    }

    // @TODO: validate user's attributes

    user.Token  = user.GenerateToken()
    user.ApiKey = user.GenerateApiKey()
    user.Status = "disabled"

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

/* NOT EXPORTED!!! */
func getUserByEmail(email string) (*User, error) {
    session := db.GetConnection()
    db      := session.DB(config.DBName)
    dbc     := db.C(config.DBUsers)

    c := make(map[string]interface{})
    c["email"] = email

    var user *User
    err := dbc.Find(c).One(&user)

    session.Close()

    if err != nil {
        return nil, errors.New("No such a User with given Email and Password")
    }

    return user, nil
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

func (user *User) Save(insertIfNotExists bool) error {
    if len(user.ApiKey) <= 0 {
        return errors.New("Cannot save user: given ApiKey is empty!")
    }

    // @TODO: add validation and Error Handling
    session := db.GetConnection()
    db      := session.DB(config.DBName)
    dbc     := db.C(config.DBUsers)

    selector := make(map[string]interface{})
    selector["apikey"] = user.ApiKey

    var err error
    if insertIfNotExists {
        _, err = dbc.Upsert(selector, user)
    } else {
        err = dbc.Update(selector, user)
    }

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

func (user *User) GenerateToken() string {
    // @TODO: use md5(microtime())
    return stringHelper.RandomString(32)
}

func (user *User) GenerateApiKey() string {
    // @TODO: check if Token is given
    // @TODO: use User.Token + md5(microtime())
    return stringHelper.RandomString(32)
}