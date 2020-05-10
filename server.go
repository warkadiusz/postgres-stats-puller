package main

import (
	db "ask/db"
	"ask/http_service"
	"ask/queries"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var debug bool
var httpService *http_service.HttpService
var queryFactory queries.QueryFactory

func main() {
	loadEnvs()
	setupDebug()
	httpService = http_service.CreateHttpService()
	db.CreateInfluxConnection()
	db.CreatePostgresConnection()
	defer db.CloseInfluxConnection()
	defer db.ClosePostgresConnection()
	queryFactory = queries.CreateQueryFactory(db.GetPostgresConnection())

	pullRate := getPullRate()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go pullDataEvery(pullRate)
	wg.Wait()
}

func pullDataEvery(pullRate time.Duration) {
	for range time.Tick(pullRate) {
		getAllData()
	}
}

func getAllData() {
	definedQueries := getQueriesToUse()

	for _, query := range definedQueries {
		qValue := query.GetValue()
		db.StoreDatapoint(query.GetName(), qValue)
		httpService.BroadcastDatapoint(query.GetName(), qValue)

		debugLog(" Saved datapoint: Name=" + query.GetName() + " Value=" + strconv.Itoa(qValue))
	}

}

func getQueriesToUse() []queries.Query {
	definedQueriesToUse := strings.Split(os.Getenv("ACTIVE_QUERIES"), ",")
	var queriesToUse []queries.Query
	for _, queryName := range definedQueriesToUse {
		createdQuery, err := queryFactory.Create(queryName)
		if err != nil {
			log.Fatal(err)
		}

		queriesToUse = append(queriesToUse, createdQuery)
	}

	return queriesToUse
}

func getPullRate() time.Duration {
	numOfSeconds, err := strconv.Atoi(os.Getenv("REFRESH_RATE_SEC"))
	if err != nil {
		log.Fatal("Refresh (pull) rate is not valid integer")
	}

	return time.Duration(numOfSeconds) * time.Second
}

func loadEnvs() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file not found. Ensure it exists and it's readable. Exit...")
	}
}

func setupDebug() {
	debug = os.Getenv("DEBUG") == "true"
}

func debugLog(msg string) {
	if debug {
		log.Print(msg)
	}
}
