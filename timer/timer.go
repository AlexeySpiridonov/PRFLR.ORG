package timer

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"prflr.org/config"
	"prflr.org/db"
	"prflr.org/helpers"
	"sort"
	"time"
)

/**
 * UDP Package struct
 */
type Timer struct {
	Thrd      string
	Src       string
	Timer     string
	Time      float32
	Info      string
	Apikey    string
	Timestamp int64
}

/**
 * Web panel Struct
 */
type Stat struct {
	Src       string
	Timer     string
	Timestamp int64
	Count     int
	Total     float32
	Min       float32
	Avg       float32
	Max       float32
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
	Min    [][]interface{}
	Max    [][]interface{}
}

func GetList(apiKey string, criteria map[string]interface{}) (*[]Timer, error) {
	// @TODO: add validation and Error Handling
	session, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	collectionName := helpers.GetCappedCollectionNameForApiKey(apiKey)

	db := session.DB(config.DBName)
	dbc := db.C(collectionName)

	// Query All
	var results []Timer

	//TODO add criteria builder
	criteria["apikey"] = apiKey
	err = dbc.Find(criteria).Sort("-_id").Limit(300).All(&results)

	return &results, err
}

func Aggregate(apiKey string, criteria map[string]interface{}, groupBy map[string]interface{}, sortBy string) (*[]Stat, error) {
	// @TODO: add validation and Error Handling
	session, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	collectionName := helpers.GetCappedCollectionNameForApiKey(apiKey)

	db := session.DB(config.DBName)
	dbc := db.C(collectionName)

	var results []Stat

	grouplist := make(map[string]interface{})
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
	sort := bson.M{"$sort": bson.M{sortBy: -1}}
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

	collectionName := helpers.GetCappedCollectionNameForApiKey(apiKey)

	db := session.DB(config.DBName)
	dbc := db.C(collectionName)

	graph := &Graph{}
	var results []Stat

	// group params
	grouplist := map[string]interface{}{
		"count":     bson.M{"$sum": 1},
		"total":     bson.M{"$sum": "$time"},
		"avg":       bson.M{"$avg": "$time"},
		"min":       bson.M{"$min": "$time"},
		"max":       bson.M{"$max": "$time"},
		"timestamp": bson.M{"$first": "$timestamp"},
		"_id": map[string]interface{}{
			"timestamp": "$timestamp",
		},
	}

	// criteria
	criteria["apikey"] = apiKey
	if _, ok := criteria["timestamp"]; !ok {
		criteria["timestamp"] = &bson.M{"$gte": time.Now().Unix() - (86400 * 30)} // last month
	}

	group := bson.M{"$group": grouplist}
	Sort := bson.M{"$sort": bson.M{"timestamp": 1}}
	match := bson.M{"$match": criteria}
	//limit := bson.M{"$limit": 100000}
	aggregate := []bson.M{match, Sort /*limit,*/, group}

	mongoStart := time.Now().UnixNano()
	err = dbc.Pipe(aggregate).All(&results)
	fmt.Println("Mongo query time: ", time.Now().UnixNano()-mongoStart)
	if err != nil {
		return nil, err
	}

	// Reformat Result to String-Keys Map
	resultsMap := make(map[int64]Stat)
	reformatStart := time.Now().UnixNano()
	for _, stat := range results {
		if stat.Timestamp <= 0 {
			continue
		}
		resultsMap[stat.Timestamp] = stat
	}
	fmt.Println("Reformat time: ", time.Now().UnixNano()-reformatStart)

	// Add Zeros
	if len(results) > 0 {
		zerosStart := time.Now().UnixNano()
		// Sort by TS
		sort.Sort(StatTimestampSorter(results))
		minTS := results[0].Timestamp
		maxTS := results[len(results)-1].Timestamp
		fmt.Println("Adding zeros: ", maxTS-minTS)
		for i := minTS; i <= maxTS; i++ {
			if _, found := resultsMap[i]; !found {
				resultsMap[i] = Stat{Timestamp: i, Count: 0, Avg: 0, Min: 0, Max: 0}
			}
		}
		fmt.Println("Zeros time: ", time.Now().UnixNano()-zerosStart)
	}

	//k := int64(len(results) / 500)
	k := int64(len(resultsMap) / 500)
	if k <= 0 {
		k = 1
	}

	var normalizedResultsData = make(map[int64][]Stat)
	normalizeStart := time.Now().UnixNano()
	for _, stat := range resultsMap {
		if stat.Timestamp <= 0 {
			continue
		}
		var key = stat.Timestamp / k
		if _, exists := normalizedResultsData[key]; exists {
			normalizedResultsData[key] = append(normalizedResultsData[key], stat)
		} else {
			normalizedResultsData[key] = []Stat{stat}
			//fmt.Println(normalizedResultsData[key])
		}
	}
	fmt.Println("Normalize time: ", time.Now().UnixNano()-normalizeStart)
	var normalizedResults []Stat
	formatStart := time.Now().UnixNano()
	for key, statData := range normalizedResultsData {
		var count = 0
		var median []float32
		var min = float32(99999999)
		var max = float32(0)
		for _, stat := range statData {
			//count += stat.Count
			if count < stat.Count {
				count = stat.Count
			}
			median = append(median, stat.Avg)
			if min > stat.Min {
				min = stat.Min
			}
			if max < stat.Max {
				max = stat.Max
			}
		}

		// calc median
		timerMedian := float32(0)
		medianLenth := len(median)
		if medianLenth%2 == 0 {
			// even
			timerMedian = (median[medianLenth/2] + median[(medianLenth/2)-1]) / 2
		} else {
			// odd
			timerMedian = median[(medianLenth / 2)]
		}

		ts := key * k
		normalizedResults = append(normalizedResults, Stat{Timestamp: ts, Count: count /*count/len(statData)*/, Avg: timerMedian, Min: min, Max: max})
	}
	fmt.Println("Format time: ", time.Now().UnixNano()-formatStart)

	sort.Sort(StatTimestampSorter(normalizedResults))

	for _, stat := range normalizedResults {
		graph.Min = append(graph.Min, []interface{}{stat.Timestamp, stat.Min})
		graph.Avg = append(graph.Avg, []interface{}{stat.Timestamp, stat.Avg})
		graph.Max = append(graph.Max, []interface{}{stat.Timestamp, stat.Max})
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

func (timer Timer) Save() error {
	// @TODO: add validation and Error Handling
	session, err := db.GetConnection()
	if err != nil {
		return err
	}
	defer session.Close()

	// define User's Collection Name to Save Timers to
	collectionName := helpers.GetCappedCollectionNameForApiKey(timer.Apikey)

	db2 := session.DB(config.DBName)
	dbc := db2.C(collectionName)

	err = dbc.Insert(timer)

	return err
}
