# postgres-stats-puller
Executes specified statistics queries (like database size, number of transactions etc.) on Postgres database
and saves the results to InfluxDB as a time series.
Additionally, provides WebSocket and HTTP endpoints for respectively live updates and historical data access.

## Endpoints
### `GET /data/<QueryName>?after=<ISO date time>&before=<ISO date time>`
Returns historical data for given query name between specified two dates and times. To avoid too much data at once, all requests always return data averaged to 30 data points only. Periods without any data are filled with zeroes. See example below.

#### Example
Request: `GET /data/DatabaseSize?after=2020-05-01T00:00:00&before=2020-05-30T23:59:59`  
Response:
```json5
{
  "status": "success",
  "code": 200,
  "data": {
    "2020-05-01T00:00:00": 18.0,
    "2020-05-02T00:00:00": 0,
    /*...*/
    "2020-05-30T00:00:00": 53.1415
  }
}
```

###  WebSocket live data push stream
Every client connected to web socket server automatically subscribes to that stream.  
Data:
```json5
{
  "status": "datapoint",
  "code": 200,
  "data": {
    "DatabaseSize": 12345
  }
}
```


## Requirements
In order to build, following modules have to be imported:
```
go get github.com/joho/godotenv
go get github.com/lib/pq
go get github.com/influxdata/influxdb1-client
go get github.com/influxdata/influxdb1-client/v2
go get github.com/gorilla/mux
go get github.com/gorilla/websocket
```
Additionally, Postgres and Influx databases have to be accessible.

## Building
1. Install dependencies (see [Requirements](#requirements) )
2. `make all`
3. Executable `server` will be created in `bin/` directory

## Runtime env variables
All of the following environmental variables are required for the server to run. They can be specified either in
`.env` file in the same directory as executable file or passed during execution, like:  
`POSTGRES_HOST=postgres POSTGRES_PORT=1234 server`

| Name | Possible values | Description |
| ---- | --------------- | ----------- |
| POSTGRES_HOST | string (hostname) | Postgres host. Can be either FQDN or IP address |
| POSTGRES_PORT | int (port number) | Postgres port number |
| POSTGRES_USER | string | Username to Postges |
| POSTGRES_PASSWORD | string | Password to Postgres |
| INFLUX_HOST | string (hostname) | InfluxDB host. Can be either FQDN or IP address |
| INFLUX_PORT | int (port number) | InfluxDB port number |
| INFLUX_USER= | string | Username to InfluxDB |
| INFLUX_PASSWORD | string | Password to InfluxDB |
| INFLUX_DB_NAME | string | Database name in InfluxDB |
| REFRESH_RATE_SEC | int | How often statistics should be pulled from Postgres database, in seconds. |
| ACTIVE_QUERIES | Any comma-separated combination of the following:  <ul><li>DatabaseSize</li><li>NumberOfConnections</li><li>NumberOfRunningQueries</li><li>NumberOfRunningQueriesOver15sec</li><li>TotalNumberOfTransactions</li> | Defines which queries should be executed |
| DEBUG | true/false | Should debug information be printed to stdout |
