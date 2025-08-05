package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Worker struct {
	influxClient influxdb2.Client
	conn         *amqp.Connection
	channel      *amqp.Channel
	emailSender  *EmailSender
	pgDB         *sql.DB
	amqpURI      string
	influxOrg    string
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

func (worker *Worker) GetTotalConsumptionForMonth(deviceID string, year int, month int) (float64, error) {
	queryAPI := worker.influxClient.QueryAPI(worker.influxOrg)
	startTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
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

func NewWorker(amqpURI, influxURI, influxToken, pgConnStr, emailSecret string) (*Worker, error) {
	// Connect to InfluxDB
	influxClient := influxdb2.NewClient(influxURI, influxToken)

	// Connect to PostgreSQL
	pgDB, err := sql.Open("postgres", pgConnStr)
	if err != nil {
		influxClient.Close()
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	err = pgDB.Ping()
	if err != nil {
		influxClient.Close()
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}
	fmt.Println("Successfully connected to postgres!")

	return &Worker{
		influxClient: influxClient,
		pgDB:         pgDB,
		amqpURI:      amqpURI,
		emailSender:  NewEmailSender(emailSecret),
		influxOrg:    "watt-flow",
	}, nil
}

func (c *Worker) ConnectToBroker() error {
	conn, err := amqp.Dial(c.amqpURI)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %v", err)
	}

	c.conn = conn
	c.channel = channel

	_, err = c.channel.QueueDeclare(
		"billing_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare billing queue: %v", err)
	}
	fmt.Println("Successfully initialized rabbitmq connection and exchange!")
	return nil
}

func (c *Worker) Start(ctx context.Context) error {
	bill_msgs, err := c.channel.Consume(
		"billing_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start bill consumer: %v", err)
	}
	go c.processBills(ctx, bill_msgs)
	return nil
}

func (c *Worker) processBills(ctx context.Context, msgs <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgs:
			var billTask BillTaskDto
			if len(msg.Body) == 0 {
				continue
			}
			err := json.Unmarshal(msg.Body, &billTask)
			if err != nil {
				fmt.Printf("Error decoding message: %v\n", err)
				continue
			}

			var year, month int
			fmt.Sscanf(billTask.BillingDate, "%d-%02d", &year, &month)

			// Query InfluxDB for consumption
			spentPower, err := c.GetTotalConsumptionForMonth(billTask.PowerMeterID, year, month)
			if err != nil {
				fmt.Printf("Failed querying consumption for %s - %s, %s\n", billTask.BillingDate, billTask.PowerMeterID, err.Error())
				continue
			}
			fmt.Println(spentPower)
			fmt.Println(string(msg.Body))
			fmt.Println(billTask)

			// Calculate Price
			calculatedPrice := calculatePrice(spentPower, billTask.Pricelist)

			// Save Bill to Database
			bill := Bill{
				BillingDate:      billTask.BillingDate,
				IssueDate:        billTask.IssueDate,
				PricelistID:      billTask.Pricelist.ID,
				OwnerID:          billTask.OwnerID,
				SpentPower:       spentPower,
				Price:            calculatedPrice,
				Status:           "Delivered",
				HouseholdID:      billTask.HouseHoldID,
				PaymentReference: uuid.NewString(),
			}

			// Send Email
			err = c.emailSender.SendMonthlyBillPDF(billTask.OwnerEmail, bill, billTask.Pricelist, billTask.OwnerUsername, billTask.HouseholdCN)
			if err != nil {
				fmt.Printf("Failed sending email for %s - %s, %s\n", billTask.PowerMeterID, billTask.OwnerUsername, err.Error())
				bill.Status = "Not Delivered"
				continue
			}
			err = c.InsertBill(ctx, &bill)
			if err != nil {
				fmt.Printf("Failed saving bill to database for %s - %s, %s\n", billTask.PowerMeterID, billTask.OwnerUsername, err.Error())
				continue
			}
			fmt.Printf("Successfully processed bill for %s\n", billTask.OwnerEmail)

			if billTask.Last {
				fmt.Printf("Received last email task, updating status")
				err = c.UpdateStatus(ctx, billTask.MonthlyBillID, "Completed")
				if err != nil {
					fmt.Printf("Failed updating monthly bill status for %d: %s\n", billTask.MonthlyBillID, err.Error())
				} else {
					fmt.Printf("Successfully updated monthly bill status for %d\n", billTask.MonthlyBillID)
				}
			}
		}
	}
}

func (c *Worker) InsertBill(ctx context.Context, bill *Bill) error {
	query := `INSERT INTO bills (issue_date, billing_date, pricelist_id, spent_power, price, owner_id, status, household_id, payment_reference)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := c.pgDB.ExecContext(ctx, query,
		bill.IssueDate,
		bill.BillingDate,
		bill.PricelistID,
		bill.SpentPower,
		bill.Price,
		bill.OwnerID,
		bill.Status,
		bill.HouseholdID,
		bill.PaymentReference,
	)
	if err != nil {
		log.Printf("Failed to insert bill: %v", err)
		return err
	}

	log.Println("Bill inserted successfully")
	return nil
}

func (c *Worker) UpdateStatus(ctx context.Context, billID uint64, status string) error {
	query := `UPDATE monthly_bills SET status = $1 WHERE id = $2`

	_, err := c.pgDB.ExecContext(ctx, query,
		status,
		billID,
	)
	if err != nil {
		log.Printf("Failed to update monthly bill: %v", err)
		return err
	}

	log.Println("Bill status updated successfully")
	return nil
}

func (c *Worker) Shutdown(ctx context.Context) error {
	log.Println("Shutting down worker...")

	// Close connections
	if err := c.channel.Close(); err != nil {
		log.Printf("Error closing RabbitMQ channel: %v", err)
	}
	if err := c.conn.Close(); err != nil {
		log.Printf("Error closing RabbitMQ connection: %v", err)
	}
	if err := c.pgDB.Close(); err != nil {
		log.Printf("Error closing PostgreSQL connection: %v", err)
	}
	c.influxClient.Close()

	log.Println("Worker shutdown completed.")
	return nil
}
