package util

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
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
	fluxQuery := ""
	log.Println(queryParams)
	if queryParams.TimePeriod == "custom" {
		fluxQuery = generateRangeQueryString(queryParams)
	} else {
		fluxQuery = generateQueryString(queryParams)
	}
	log.Println(fluxQuery)

	result, err := queryAPI.Query(context.Background(), fluxQuery)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	results := dto.StatusQueryResult{
		Rows: []dto.StatusQueryResultRow{},
	}

	for result.Next() {
		value := result.Record().Value()
		var floatVal float64

		switch v := value.(type) {
		case float64:
			floatVal = v
		case int64:
			floatVal = float64(v)
		case int:
			floatVal = float64(v)
		case string:
			// Try to parse string as float if needed
			parsed, err := strconv.ParseFloat(v, 64)
			if err != nil {
				log.Printf("Error converting string to float: %v", err)
				continue
			}
			floatVal = parsed
		case nil:
			floatVal = 0
		default:
			log.Printf("Unexpected value type: %T", value)
			continue
		}

		results.Rows = append(results.Rows, dto.StatusQueryResultRow{
			TimeField: result.Record().Time(),
			Value:     floatVal,
		})
	}
	if result.Err() != nil {
		log.Fatalf("Query execution error: %v", result.Err())
		return nil, err
	}
	return &results, nil
}

func generateQueryString(params dto.FluxQueryStatusDto) string {
	fluxQuery := fmt.Sprintf(`
  import "array"
  import "experimental"


  data = from(bucket: "device_status")
    |> range(start: -%s)
    |> filter(fn: (r) => r._measurement == "online_status" and r._field == "value" and r.device_id=="%s")

  last_from_data = data |> last()   |> findRecord(fn: (key) => true, idx: 0)

  bounds = array.from(rows: [
    { _time: experimental.subDuration(from: now(), d: 1m),  _value: if last_from_data._value == true then false else true, _field: "value", _measurement: "online_status", device_id: "%s", _stop: now(), _start: now()},
  
  { _time: experimental.subDuration(from: now(), d: %s), _value: false, _field: "value", _measurement: "online_status", device_id: "%s", _stop: now(), _start: now()}

    ])


  all_data = union(tables: [data, bounds])
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
`, params.TimePeriod, params.DeviceId, params.DeviceId, params.TimePeriod, params.DeviceId, params.TimePeriod, params.Precision, params.GroupPeriod)
	return fluxQuery
}

func generateRangeQueryString(params dto.FluxQueryStatusDto) string {
	startDate := params.StartDate.Format(time.RFC3339)
	endDate := params.EndDate.Format(time.RFC3339)
	fluxQuery := fmt.Sprintf(`
  import "array"
  import "experimental"



  data = from(bucket: "device_status")
    |> range(start: %s, stop: %s)
    |> filter(fn: (r) => r._measurement == "online_status" and r._field == "value" and r.device_id=="%s")
    

  now = array.from(rows: [
    { _time: %s, _value: false, _field: "value", _measurement: "online_status", device_id: "%s", _stop: now(), _start: now()}
  ])


  all_data = union(tables: [data, now])

    all_data 
    |> group()
    |> sort(columns: ["_time"])
    |> range(start: %s, stop: %s)
    |> elapsed(unit: 1%s)
    |> filter(fn: (r) => r._value == false )
    |> aggregateWindow(every: %s, fn: sum, column: "elapsed")
    |> map(fn: (r) => ({
      _time: r._time,
      _value: float(v: r.elapsed)
    }))
`, startDate, endDate, params.DeviceId, startDate, params.DeviceId, startDate, endDate, params.Precision, params.GroupPeriod)
	return fluxQuery
}
