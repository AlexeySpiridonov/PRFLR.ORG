package db

import(
    "prflr.org/config"
    "prflr.org/PRFLRLogger"
    "labix.org/v2/mgo"
)

func GetConnection() (*mgo.Session, error) {
    db, err := mgo.Dial(config.DBHosts)
    if err != nil {
        PRFLRLogger.Error(err)
        return nil, err
    }
    defer db.Close()

    // Optional. Switch the session to a monotonic behavior.
    db.SetMode(mgo.Monotonic, true)

    return db.Clone(), nil
}