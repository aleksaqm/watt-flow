package util

import (
	"context"
	"fmt"
	"log"
	"watt-flow/config"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxQueryHelper struct {
	influxURL         string
	influxToken       string
	organization      string
	statusBucket      string
	measurementBucket string
}

func NewInfluxQueryHelper(env *config.Environment) *InfluxQueryHelper {
	return &InfluxQueryHelper{
		influxURL:         env.InfluxURL,
		influxToken:       env.InfluxToken,
		organization:      env.InlfuxOrg,
		statusBucket:      env.InfluxStatusBucket,
		measurementBucket: env.InfluxMeasurementBucket,
	}
}

func (inf *InfluxQueryHelper) SendStatusQuery() error {
	client := influxdb2.NewClient(inf.influxURL, inf.influxToken)
	defer client.Close()

	// Create a query API client
	queryAPI := client.QueryAPI(inf.organization)

	// Define the Flux query
	fluxQuery := `
  import "array"
  import "experimental"
now = array.from(rows: [
  { _time: experimental.subDuration(from: now(), d: 1m), _value: false, _field: "value", _measurement: "online_status", device_id: "ff781b42-c3b0-475b-bdc5-cb467d0f4222", _stop: now(), _start: now()}
])

data = from(bucket: "device_status")
  |> range(start: -32h)
  |> filter(fn: (r) => r._measurement == "online_status" and r._field == "value" and r.device_id=="ff781b42-c3b0-475b-bdc5-cb467d0f4222")

last_period = 
      from(bucket: "device_status")
        |> range(start: -30d, stop: -32h)
        |> filter(fn: (r) => r._measurement == "online_status" and r._field == "value" and r.device_id=="ff781b42-c3b0-475b-bdc5-cb467d0f4222")
        |> last()

all_data = union(tables: [data, now, last_period])
  all_data 
  |> group()
  |> sort(columns: ["_time"])
  |> range(start: -32h)
  |> elapsed(unit: 1s)
  |> filter(fn: (r) => r._value == false )
  |> aggregateWindow(every: 1h, fn: sum, column: "elapsed")
  |> map(fn: (r) => ({
    _time: r._time,
    _value: float(v: r.elapsed) / 60.0
  }))
`

	result, err := queryAPI.Query(context.Background(), fluxQuery)
	if err != nil {
		return err
	}
	defer result.Close()

	for result.Next() {
		// Access the results, here we print each record
		fmt.Printf("Time: %v, Duration (minutes): %v\n", result.Record().Time(), result.Record().Value())
	}

	// Check for any errors after the query execution
	if result.Err() != nil {
		log.Fatalf("Query execution error: %v", result.Err())
	}
	return nil
}
