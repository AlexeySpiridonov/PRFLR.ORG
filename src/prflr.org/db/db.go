package db

import(
    "prflr.org/config"
    "labix.org/v2/mgo"
    //"time"
    "log"
)

func GetConnection() (*mgo.Session) {
    //maxWait := time.Duration(5 * time.Second)
    //db, err := mgo.DialWithTimeout(config.DBHosts, maxWait)
    db, err := mgo.Dial(config.DBHosts)
    if err != nil {
        //log.Print("! Db.go::GetConnection !")
        log.Fatal(err)
    }
    defer db.Close()

    // Optional. Switch the session to a monotonic behavior.
    db.SetMode(mgo.Monotonic, true)

    return db.Clone()
}