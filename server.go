package main

import (
	db "ask/db"
	"ask/queries"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	loadEnvs()
	db.CreatePostgresConnection()
	db.CreateInfluxConnection()
	defer db.CloseInfluxConnection()
	defer db.ClosePostgresConnection()

	numOfSeconds, _ := strconv.Atoi(os.Getenv("REFRESH_RATE_SEC"))
	pullDataEvery(time.Duration(numOfSeconds)*time.Second, getAllData)
}

func pullDataEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func getAllData(t time.Time) {
	postgresConnection := db.GetPostgresConnection()
	databaseSize := queries.DatabaseSize(postgresConnection)
	numberOfConnections := queries.NumberOfConnections(postgresConnection)
	numberOfRunningQueries := queries.NumberOfRunningQueries(postgresConnection)
	numberOfRunningQueriesOver15Sec := queries.NumberOfRunningQueriesOver15sec(postgresConnection)
	totalNumberOfTransactions := queries.TotalNumberOfTransactions(postgresConnection)

	db.StoreDatapoint("database_size", databaseSize);
	db.StoreDatapoint("no_of_connections", numberOfConnections);
	db.StoreDatapoint("no_of_running_queries", numberOfRunningQueries);
	db.StoreDatapoint("no_of_running_queries_over_15s", numberOfRunningQueriesOver15Sec);
	db.StoreDatapoint("total_no_of_transactions", totalNumberOfTransactions);

	log.Print("Saved datapoint: ", databaseSize, numberOfConnections, numberOfRunningQueries, numberOfRunningQueriesOver15Sec, totalNumberOfTransactions)
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
