package db

import (
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

func CloseInfluxConnection() {
	connection.Close()
}
