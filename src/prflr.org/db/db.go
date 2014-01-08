package db

import(
    "prflr.org/config"
    "prflr.org/timerstruct"
    "labix.org/v2/mgo"
    "log"
)

//func GetConnection(dbname string) (*mgo.Database) {
func GetConnection() (*mgo.Session) {
    db, err := mgo.Dial(config.DBHosts)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Optional. Switch the session to a monotonic behavior.
    db.SetMode(mgo.Monotonic, true)

    //return db.Clone().DB(dbname)
    return db.Clone()
}

func Init() {
    //db := GetConnection(config.DBName)
    session := GetConnection()
    db      := session.DB(config.DBName)

	err := db.DropDatabase()
	if err != nil {
		log.Fatal(err)
	}
	dbc := db.C(config.DBTimers)

	// creating capped collection
	dbc.Create(&mgo.CollectionInfo{Capped: true, MaxBytes: config.CappedCollectionMaxByte, MaxDocs: config.CappedCollectionMaxDocs})

	// Insert Test Datas
	err = dbc.Insert(&timerstruct.Timer{Thrd: "1234567890", Timer: "prflr.check69", Src: "test.src69", Time: 1, Info: "test data 69", Apikey: "PRFLRApiKey69"})
	if err != nil {
		log.Fatal(err)
	}

	session.Close()
}