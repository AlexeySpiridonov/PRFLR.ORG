package timer

import(
    "prflr.org/stringHelper"
    "prflr.org/config"
    "prflr.org/db"
    "labix.org/v2/mgo/bson"
    "sort"
    //"fmt"
    //"strconv"
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

type StatTimestampSorter []Stat
func (a StatTimestampSorter) Len() int           { return len(a) }
func (a StatTimestampSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a StatTimestampSorter) Less(i, j int) bool { return a[i].Timestamp < a[j].Timestamp }

/**
 * UI Graph Struct 
 */
type Graph struct {
    Median [][]int
    Avg    [][]interface{}
    TPS    [][]int // Timers per Second
    Min   float32
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

    graph := &Graph{
        Min: float32(999999999999),
        Max: float32(0),
    }
    var results []Stat

    // group params
    grouplist  := map[string]interface{} {
        "count": bson.M{"$sum": 1},
        "total": bson.M{"$sum": "$time"},
        "avg"  : bson.M{"$avg": "$time"},
        "min"  : bson.M{"$min": "$time"},
        "max"  : bson.M{"$max": "$time"},
        "timestamp": bson.M{"$first": "$timestamp"},
        "_id": map[string]interface{} {
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
    Sort  := bson.M{"$sort" : bson.M{"timestamp": 1}}
    match := bson.M{"$match": criteria}
    aggregate := []bson.M{match, group, Sort}

    err = dbc.Pipe(aggregate).All(&results)
    if err != nil {
        return nil, err
    }

    //k := int64(5)
    k := int64(len(results) / 500)
    if k <= 0 {
        k = 1
    }

    var normalizedResultsData = make(map[int64][]Stat)
    for _, stat := range results {
        if stat.Timestamp <= 0 {
            continue
        }
        var key = stat.Timestamp/k
        if _, exists := normalizedResultsData[key]; exists {
            normalizedResultsData[key] = append(normalizedResultsData[key], stat)
        } else {
            normalizedResultsData[key] = []Stat{stat}
            //fmt.Println(normalizedResultsData[key])
        }

        if stat.Min < graph.Min {
            graph.Min = stat.Min
        }
        if stat.Max > graph.Max {
            graph.Max = stat.Max
        }
    }
    var normalizedResults []Stat
    for key, statData := range normalizedResultsData {
        var count  = 0
        var median []float32
        for _, stat := range statData {
            count += stat.Count
            median = append(median, stat.Avg)
        }

        // calc median
        timerMedian := float32(0)
        medianLenth := len(median)
        if medianLenth % 2 == 0 {
            // even
            timerMedian = (median[medianLenth / 2] + median[(medianLenth / 2)-1]) / 2
        } else {
            // odd
            timerMedian = median[(medianLenth / 2)]
        }

        ts := key * k
        normalizedResults = append(normalizedResults, Stat{Timestamp: ts, Count: count, Avg: timerMedian})
    }

    sort.Sort(StatTimestampSorter(normalizedResults))

    for _, stat := range normalizedResults {
        graph.Avg = append(graph.Avg, []interface{}{stat.Timestamp, stat.Avg})
        graph.TPS = append(graph.TPS, []int{int(stat.Timestamp), stat.Count})
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
