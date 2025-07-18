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
	client            influxdb2.Client
	organization      string
	statusBucket      string
	measurementBucket string
}

func NewInfluxQueryHelper(env *config.Environment) *InfluxQueryHelper {
	client := influxdb2.NewClient(env.InfluxURL, env.InfluxToken)
	return &InfluxQueryHelper{
		client:            client,
		organization:      env.InlfuxOrg,
		statusBucket:      env.InfluxStatusBucket,
		measurementBucket: env.InfluxMeasurementBucket,
	}
}

func (inf *InfluxQueryHelper) GetTotalConsumptionForMonth(deviceID string, year int, month int) (float64, error) {
	queryAPI := inf.client.QueryAPI(inf.organization)
	startTime := time.Date(year-1, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.AddDate(0, 1, 0) // First day of the next month

	fluxQuery := generateMonthConsumptionQueryString(deviceID, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	result, err := queryAPI.Query(context.Background(), fluxQuery)
	if err != nil {
		return -1.0, err
	}
	defer result.Close()
	var totalConsuption float64

	for result.Next() {
		if value, ok := result.Record().Value().(float64); ok {
			totalConsuption = value
		}
	}

	if result.Err() != nil {
		return 0, result.Err()
	}

	return totalConsuption, nil
}

func (inf *InfluxQueryHelper) SendStatusQuery(queryParams dto.FluxQueryStatusDto) (*dto.StatusQueryResult, error) {
	queryAPI := inf.client.QueryAPI(inf.organization)
	fluxQuery := ""

	if queryParams.Realtime {
		fluxQuery = generateRealtimeQuery(queryParams)
	} else {
		if queryParams.TimePeriod == "custom" {
			fluxQuery = generateRangeQueryString(queryParams)
		} else {
			fluxQuery = generateQueryString(queryParams)
		}
	}

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
		case bool:
			if v {
				floatVal = 1
			} else {
				floatVal = 0
			}
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

func (inf *InfluxQueryHelper) SendConsumptionQuery(queryParams dto.FluxQueryConsumptionDto) (*dto.StatusQueryResult, error) {
	queryAPI := inf.client.QueryAPI(inf.organization)
	fluxQuery := ""
	log.Print(inf)
	if queryParams.Realtime {
		fluxQuery = generateRealtimePowerConsumptionQuery(queryParams)
	} else {
		if queryParams.TimePeriod == "custom" {
			fluxQuery = generatePowerConsumptionRangeQueryString(queryParams)
		} else {
			fluxQuery = generatePowerConsumptionQuery(queryParams)
		}
	}

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
		case bool:
			if v {
				floatVal = 1
			} else {
				floatVal = 0
			}
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

func generateMonthConsumptionQueryString(deviceID string, startMonth string, endMonth string) string {
	fluxQuery := fmt.Sprintf(`
  from(bucket: "power_measurements")
    |> range(start: %s, stop: %s)
    |> filter(fn: (r) => r["_measurement"] == "power_consumption" and r.device_id == "%s")
    |> sum(column: "_value")
    `, startMonth, endMonth, deviceID)

	return fluxQuery
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

func generateRealtimeQuery(params dto.FluxQueryStatusDto) string {
	fluxQuery := fmt.Sprintf(`

  import "array"
  import "experimental"

  data = from(bucket: "device_status")
    |> range(start: -3h)
    |> filter(fn: (r) => r._measurement == "online_status" and r._field == "value" and r.device_id=="%s")
    |> sort(columns: ["_time"])

  first_from_data = data |> first()   |> findRecord(fn: (key) => true, idx: 0)

  first = array.from(rows: [
    { _time: experimental.subDuration(from: now(), d: 3h), _value: if first_from_data._value == true then false else true, _field: "value", _measurement: "online_status", device_id: "%s", _stop: now(), _start: now()},
  ])

  all_data = union(tables: [data, first])


  all_data
    |> aggregateWindow(createEmpty: true, every: 1m, fn: last)
    |>fill(usePrevious: true)
    `, params.DeviceId, params.DeviceId)
	return fluxQuery
}

func generatePowerConsumptionQuery(params dto.FluxQueryConsumptionDto) string {
	fluxQuery := fmt.Sprintf(`
  import "array"
  import "experimental"

  startTime = experimental.subDuration(from: now(), d: %s)

  data = from(bucket: "power_measurements")
    |> range(start: startTime)
    |> filter(fn: (r) => r._measurement == "power_consumption" and r._field == "value" and r.city == "%s")

  bounds = array.from(rows: [
    { 
      _time: startTime, 
      _value: 0.0, // Vrednost 0 da ne utiče na sumu
      _field: "value", 
      _measurement: "power_consumption", 
      city: "%s", 
      _start: startTime, // Dodajemo _start i _stop da se poklopi shema
      _stop: now()
    }
  ])

  union(tables: [data, bounds])
    |> group(columns: ["city"])
    |> aggregateWindow(every: %s, fn: sum) // createEmpty: true više nije neophodno
    |> map(fn: (r) => ({
      _time: r._time,
      _value: float(v: r._value),
      city: r.city
    }))
    |> yield(name: "power_consumption_summary")
`, params.TimePeriod, params.City, params.City, params.GroupPeriod)
	return fluxQuery
}

func generatePowerConsumptionRangeQueryString(params dto.FluxQueryConsumptionDto) string {
	startDate := params.StartDate.Format(time.RFC3339)
	endDate := params.EndDate.Format(time.RFC3339)
	fluxQuery := fmt.Sprintf(`
  import "array"

  data = from(bucket: "power_measurements")
    |> range(start: %s, stop: %s)
    |> filter(fn: (r) => r._measurement == "power_consumption" and r._field == "value" and r.city == "%s")
    
  bounds = array.from(rows: [
    { 
      _time: %s, 
      _value: 0.0, 
      _field: "value", 
      _measurement: "power_consumption", 
      city: "%s", 
      _start: %s, 
      _stop: %s 
    }
  ])

  union(tables: [data, bounds])
    |> group(columns: ["city"])
    |> sort(columns: ["_time"]) 
    |> aggregateWindow(every: %s, fn: sum)
    |> map(fn: (r) => ({
      _time: r._time,
      _value: float(v: r._value),
      city: r.city
    }))
    |> yield(name: "power_consumption_summary")
`, startDate, endDate, params.City, startDate, params.City, startDate, endDate, params.GroupPeriod)
	return fluxQuery
}

func generateRealtimePowerConsumptionQuery(params dto.FluxQueryConsumptionDto) string {
	fluxQuery := fmt.Sprintf(`
  import "array"
  import "experimental"

  startTime = experimental.subDuration(from: now(), d: 1h)

  data = from(bucket: "power_measurements")
      |> range(start: startTime)
      |> filter(fn: (r) => r._measurement == "power_consumption" and r._field == "value" and r.city == "%s")
  
  bounds = array.from(rows: [
    { 
      _time: startTime, 
      _value: 0.0, 
      _field: "value", 
      _measurement: "power_consumption", 
      city: "%s", 
      _start: startTime,
      _stop: now()
    }
  ])

  union(tables: [data, bounds])
    |> group(columns: ["city"])
    |> aggregateWindow(every: 5m, fn: sum)
    |> yield(name: "realtime_power_consumption")
  `, params.City, params.City)
	return fluxQuery
}
