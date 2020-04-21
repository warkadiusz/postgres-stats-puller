package main

import (
	db "ask/db"
	"ask/queries"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var debug bool

func main() {
	loadEnvs()
	db.CreatePostgresConnection()
	db.CreateInfluxConnection()
	defer db.CloseInfluxConnection()
	defer db.ClosePostgresConnection()

	debug = os.Getenv("DEBUG") == "true"

	postgresConnection := db.GetPostgresConnection()
	queryFactory := queries.CreateQueryFactory(postgresConnection)

	queriesToUseDefinition := strings.Split(os.Getenv("ACTIVE_QUERIES"), ",")
	var queriesToUse []queries.Query
	for _, queryName := range queriesToUseDefinition {
		createdQuery, err := queryFactory.Create(queryName)
		if err != nil {
			log.Fatal(err)
		}

		queriesToUse = append(queriesToUse, createdQuery)
	}

	numOfSeconds, _ := strconv.Atoi(os.Getenv("REFRESH_RATE_SEC"))
	pullDataEvery(time.Duration(numOfSeconds)*time.Second, getAllData, queriesToUse)
}

func pullDataEvery(d time.Duration, f func(time.Time, []queries.Query), args []queries.Query) {
	for x := range time.Tick(d) {
		f(x, args)
	}
}

func getAllData(t time.Time, queries []queries.Query) {
	for _, query := range queries {
		qValue := query.GetValue()
		db.StoreDatapoint(query.GetName(), qValue)
		if debug {
			log.Print(" Saved datapoint: Name=" + query.GetName() + " Value=" + strconv.Itoa(qValue))
		}
	}

}

func loadEnvs() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file not found. Ensure it exists and it's readable. Exit...")
	}
}

/**
func transactionsHandler(writer http.ResponseWriter, request *http.Request) {
	rows, err := postgresConnection.Query("SELECT * FROM users")

	if err != nil {
		log.Fatal(err.Error())
	}

	users := []User{}

	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	jsonOutput, _ := json.Marshal(users)

	fmt.Fprintln(writer, string(jsonOutput))
}
*/
