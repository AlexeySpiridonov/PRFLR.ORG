package db

import(
    "prflr.org/config"
    "prflr.org/PRFLRLogger"
    "labix.org/v2/mgo"
//    "log"
)

type dbSingleton struct {
    Session *mgo.Session
}

var dbSingletonVar *dbSingleton = nil

func Init() {
    GetConnection()    
}

func GetConnection() (*mgo.Session, error) {
    err := connect()
    if err != nil {
        PRFLRLogger.Error(err)
        return nil, err
    }

    return dbSingletonVar.Session.Copy(), nil
}

func CreateCappedCollection(dbc *mgo.Collection, cappedCollectionMaxByte, cappedCollectionMaxDocs int) error {
    // creating capped collection
    return dbc.Create(&mgo.CollectionInfo{Capped: true, MaxBytes: cappedCollectionMaxByte, MaxDocs: cappedCollectionMaxDocs})
}

/* Not Exported */
func connect() error {
    if dbSingletonVar == nil {
        var err error

        dbSingletonVar = new(dbSingleton)
        dbSingletonVar.Session, err = mgo.Dial(config.DBHosts)

        if err != nil {
            PRFLRLogger.Error(err)
            return err
        }

        // Optional. Switch the session to a monotonic behavior.
        dbSingletonVar.Session.SetMode(mgo.Monotonic, true)
        dbSingletonVar.Session.SetSafe(nil)
        dbSingletonVar.Session.Fsync(false)
    }

    return nil
}