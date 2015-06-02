package timer

import(
    "prflr.org/stringHelper"
    "prflr.org/config"
    "prflr.org/db"
    "labix.org/v2/mgo/bson"
    "strconv"
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
    Timestamp int64
}

/**
 * Web panel Struct
 */
type Stat struct {
    Src   string
    Timer string
    Timestamp int64
    Count int
    Total float32
    Min   float32
    Avg   float32
    Max   float32
}

/**
 * UI Graph Struct 
 */
type Graph struct {
    Min    int64
    Max    int64
    Median map[string]int
    Avg    map[string]int
    RPS    map[string]int
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

func FormatGraph(apiKey string, criteria map[string]interface{}) (*Graph, error) {
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

    // group params
    grouplist  := map[string]interface{} {
        "count": bson.M{"$sum": 1},
        "total": bson.M{"$sum": "$time"},
        "avg"  : bson.M{"$avg": "$time"},
        "timestamp": bson.M{"$first": "$timestamp"},
        "_id": map[string]string {
            "timestamp": "$timestamp",
        },
    }

    // criteria
    criteria["apikey"] = apiKey
    /*criteria := map[string]string {
        "apikey": apiKey,
        "src": "node1.mag.ndmsystems.com",
    }*/

    group := bson.M{"$group": grouplist}
    sort  := bson.M{"$sort" : bson.M{"timestamp": 1}}
    match := bson.M{"$match": criteria}
    aggregate := []bson.M{match, group, sort}

    err = dbc.Pipe(aggregate).All(&results)
    if err != nil {
        return nil, err
    }

    graph := &Graph{
        Avg: make(map[string]int), 
        RPS: make(map[string]int),
        Median: make(map[string]int),
        Min: 0,
        Max: 0,
    }
    for _, stat := range results {
        key := "key_" + strconv.FormatInt(stat.Timestamp, 10)
        graph.Avg[key] = int(stat.Avg)
        graph.RPS[key] = int(stat.Count)
    }

    if len(results) > 0 {
        graph.Min = results[0].Timestamp
        graph.Max = results[len(results)-1].Timestamp
    }

    return graph, nil
}

/**
 * @TODO: use ApiKey for CollectionName
 */
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
