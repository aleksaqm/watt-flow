package util

import (
	"context"
	"fmt"
	"log"
	"watt-flow/config"
	"watt-flow/dto"

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

func (inf *InfluxQueryHelper) SendStatusQuery(queryParams dto.FluxQueryStatusDto) (*dto.StatusQueryResult, error) {
	client := influxdb2.NewClient(inf.influxURL, inf.influxToken)
	defer client.Close()

	// Create a query API client
	queryAPI := client.QueryAPI(inf.organization)

	// Define the Flux query
	fluxQuery := fmt.Sprintf(`
  import "array"
  import "experimental"
  now = array.from(rows: [
    { _time: experimental.subDuration(from: now(), d: 1m), _value: false, _field: "value", _measurement: "online_status", device_id: "%s", _stop: now(), _start: now()}
  ])

  data = from(bucket: "device_status")
    |> range(start: -%s)
    |> filter(fn: (r) => r._measurement == "online_status" and r._field == "value" and r.device_id=="%s")

  last_period = 
        from(bucket: "device_status")
          |> range(start: -30d, stop: -%s)
          |> filter(fn: (r) => r._measurement == "online_status" and r._field == "value" and r.device_id=="%s")
          |> last()

  all_data = union(tables: [data, now, last_period])
    all_data 
    |> group()
    |> sort(columns: ["_time"])
    |> range(start: -%s)
    |> elapsed(unit: 1%s)
    |> filter(fn: (r) => r._value == false )
    |> aggregateWindow(every: %s, fn: sum, column: "elapsed")
    |> map(fn: (r) => ({
      _time: r._time,
      _value: float(v: r.elapsed)
    }))
`, queryParams.DeviceId, queryParams.TimePeriod, queryParams.DeviceId, queryParams.TimePeriod, queryParams.DeviceId, queryParams.TimePeriod, queryParams.Precision, queryParams.GroupPeriod)

	result, err := queryAPI.Query(context.Background(), fluxQuery)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	results := dto.StatusQueryResult{
		Rows: make([]dto.StatusQueryResultRow, 52),
	}

	for result.Next() {
		results.Rows = append(results.Rows, dto.StatusQueryResultRow{
			TimeField: result.Record().Time(),
			Value:     result.Record().Value().(float32),
		})
	}

	// Check for any errors after the query execution
	if result.Err() != nil {
		log.Fatalf("Query execution error: %v", result.Err())
		return nil, err
	}
	return &results, nil
}
