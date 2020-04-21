# postgres-stats-puller
Executes specified statistics queries (like database size, number of transactions etc.) on Postgres database
and saves the results to InfluxDB as a time series.

### Requirements
In order to build, following modules have to be imported:
```
go get github.com/joho/godotenv
go get github.com/lib/pq
go get github.com/influxdata/influxdb1-client
go get github.com/influxdata/influxdb1-client/v2
```
Additionally, Postgres and Influx databases have to be accessible.

### Building
1. Install dependencies (see [Requirements](#requirements) )
2. `make all`
3. Executable `server` will be created in `bin/` directory

### Runtime env variables
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
| REFRESH_RATE_SEC | int | How often statistics should be pulled from Postgres database |
| ACTIVE_QUERIES | Any comma-separated combination of the following:  <ul><li>DatabaseSize</li><li>NumberOfConnections</li><li>NumberOfRunningQueries</li><li>NumberOfRunningQueriesOver15sec</li><li>TotalNumberOfTransactions</li> | Defines which queries should be executed |
| DEBUG | true/false | Should debug information be printed to stdout |
