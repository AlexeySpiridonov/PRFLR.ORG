package timer

import(
    "prflr.org/stringHelper"
    "prflr.org/config"
    "prflr.org/db"
    "labix.org/v2/mgo/bson"
)

/**
 * UDP Package struct
 */
type Timer struct {
    Thrd   string
    Src    string
    Timer  string
    Time   float32
    Info   string
    Apikey string
}

/**
 * Web panel Struct
 */
type Stat struct {
    Src   string
    Timer string
    Count int
    Total float32
    Min   float32
    Avg   float32
    Max   float32
}

func GetList(apiKey string, criteria map[string]interface{}) (*[]Timer, error) {
    // @TODO: add validation and Error Handling
    session, err := db.GetConnection()
    if err != nil {
        return nil, err
    }
    defer session.Close()

    collectionName := stringHelper.GetCappedCollectionNameForApiKey(apiKey)

    db  := session.DB(config.DBName)
    dbc := db.C(collectionName)

    // Query All
    var results []Timer

    //TODO add criteria builder
    criteria["apikey"] = apiKey
    err = dbc.Find(criteria).Sort("-_id").Limit(100).All(&results)

    return &results, err
}

func Aggregate(apiKey string, criteria map[string]interface{}, groupBy map[string]interface{}, sortBy string) (*[]Stat, error) {
    // @TODO: add validation and Error Handling
    session, err := db.GetConnection()
    if err != nil {
        return nil, err
    }
    defer session.Close()

    collectionName := stringHelper.GetCappedCollectionNameForApiKey(apiKey)

    db  := session.DB(config.DBName)
    dbc := db.C(collectionName)

    var results []Stat

    grouplist  := make(map[string]interface{})
    groupparam := make(map[string]interface{})

    grouplist["count"] = bson.M{"$sum": 1}
    grouplist["total"] = bson.M{"$sum": "$time"}
    grouplist["min"] = bson.M{"$min": "$time"}
    grouplist["avg"] = bson.M{"$avg": "$time"}
    grouplist["max"] = bson.M{"$max": "$time"}

    // group params
    for i := range groupBy {
        groupparam[i] = "$" + i
        grouplist[i] = bson.M{"$first": "$" + i}
    }

    criteria["apikey"] = apiKey

    grouplist["_id"] = groupparam
    group := bson.M{"$group": grouplist}
    sort  := bson.M{"$sort": bson.M{sortBy: -1}}
    match := bson.M{"$match": criteria}
    aggregate := []bson.M{match, group, sort}

    err = dbc.Pipe(aggregate).All(&results)

    return &results, err
}

func SetApiKey(oldApiKey, newApiKey string) error {
    // @TODO: add validation and Error Handling
    session, err := db.GetConnection()
    if err != nil {
        return err
    }
    defer session.Close()

    db2 := session.DB(config.DBName)
    dbc := db2.C(config.DBTimers)

    selector := make(map[string]interface{})
    selector["apikey"] = oldApiKey

    modifier := &bson.M{"$set": bson.M{"apikey": newApiKey}}

    _, err = dbc.UpdateAll(selector, modifier)

    return err
}

func (timer *Timer) Save() error {
    // @TODO: add validation and Error Handling
    session, err := db.GetConnection()
    if err != nil {
        return err
    }
    defer session.Close()

    // define User's Collection Name to Save Timers to
    collectionName := stringHelper.GetCappedCollectionNameForApiKey(timer.Apikey)

    db2 := session.DB(config.DBName)
    dbc := db2.C(collectionName)

    err = dbc.Insert(timer)

    return err
}