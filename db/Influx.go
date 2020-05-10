package db

import (
	"errors"
	"fmt"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"os"
	"time"
)

var connection client.Client

func CreateInfluxConnection() {
	host := fmt.Sprintf(
		"http://%s:%s",
		os.Getenv("INFLUX_HOST"),
		os.Getenv("INFLUX_PORT"))
	var err error
	connection, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     host,
		Username: os.Getenv("INFLUX_USER"),
		Password: os.Getenv("INFLUX_PASSWORD"),
	})

	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
}

func StoreDatapoint(name string, value int) {
	// Create a new point batch
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database: os.Getenv("INFLUX_DB_NAME"),
	})

	fields := map[string]interface{}{
		"val": value,
	}

	pt, err := client.NewPoint(name, nil, fields, time.Now())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	bp.AddPoint(pt)

	err = connection.Write(bp)
	if err != nil {
		log.Fatal("Error: ", err.Error())
	}
}

func QueryData(table string, after time.Time, before time.Time) (client.Result, error) {

	groupBy := before.Sub(after) / 30

	q := client.NewQuery("SELECT mean(\"val\") AS \"value\" FROM \""+table+"\" WHERE time > '"+after.Format(time.RFC3339)+"' AND time < '"+before.Format(time.RFC3339)+"' GROUP BY time("+groupBy.String()+")", os.Getenv("INFLUX_DB_NAME"), "")
	response, err := connection.Query(q)
	if err != nil || response.Error() != nil {
		log.Print(err, response.Error())
		return client.Result{}, errors.New("Can't fetch data from db.")
	}

	return response.Results[0], nil
}

func GetInfluxConnection() client.Client {
	return connection
}

func CloseInfluxConnection() {
	connection.Close()
}
