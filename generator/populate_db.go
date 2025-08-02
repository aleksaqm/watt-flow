package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

// --- CONFIGURATION ---
const (
	dbConnectionString = "postgres://postgres:postgres@localhost:5432/watt-flow?sslmode=disable"
	totalUsers         = 1_000_000
	totalHouseholds    = 1_000_000
	insertBatchSize    = 10_000
)

var cities = []string{
	"Beograd", "Novi Sad", "Niš", "Kragujevac", "Subotica", "Zrenjanin", "Pančevo",
	"Čačak", "Kraljevo", "Novi Pazar", "Smederevo", "Leskovac", "Užice", "Vranje", "Šabac",
}

var streets = []string{
	"Kneza Miloša", "Bulevar kralja Aleksandra", "Nemanjina", "Knez Mihailova", "Glavna", "Cara Dušana",
	"Njegoševa", "Svetog Save", "Maksima Gorkog", "Kralja Petra", "Vojvode Stepe", "Takovska",
}

var serbianFirstNames = []string{
	"Marko", "Nikola", "Luka", "Stefan", "Lazar", "Uroš", "Filip", "Jovan", "Dušan", "Petar",
	"Milica", "Ana", "Teodora", "Jovana", "Sara", "Anđela", "Katarina", "Nevena", "Marija", "Jelena",
	"Aleksandar", "Nemanja", "Miloš", "Ivan", "Milan", "Dragan", "Zoran", "Goran", "Slobodan", "Vladimir",
	"Sofija", "Dunja", "Iva", "Mila", "Tara", "Anja", "Nađa", "Maša", "Lena", "Kristina",
}

var serbianSurnames = []string{
	"Jovanović", "Petrović", "Nikolić", "Marković", "Đorđević", "Stojanović", "Ilić", "Stanković",
	"Pavlović", "Milošević", "Todorović", "Vasić", "Lukić", "Popović", "Ristić", "Kostić",
	"Živković", "Krstić", "Stevanović", "Tomić", "Simić", "Janković", "Filipović", "Mladenović",
}

var nameReplacer = strings.NewReplacer(
	"š", "s", "đ", "dj", "č", "c", "ć", "c", "ž", "z",
	"Š", "s", "Đ", "dj", "Č", "c", "Ć", "c", "Ž", "z",
)

func cleanAndFormatName(name string) string {
	asciiName := nameReplacer.Replace(name)
	return strings.ToLower(asciiName)
}

func generateSerbianUsername(firstName string, lastName string) string {
	cleanFirstName := cleanAndFormatName(firstName)
	cleanLastName := cleanAndFormatName(lastName)

	randomNumber := rand.Intn(999)

	const (
		formatFullNameNum           = "%s.%s%d"
		formatFullNameUnderscoreNum = "%s_%s%d"
		formatLastNameNum           = "%s%d"
		formatFullNameConcatNum     = "%s%s%d"
		formatInitialLastNum        = "initial"
	)

	formats := []string{
		formatFullNameNum,
		formatFullNameUnderscoreNum,
		formatLastNameNum,
		formatFullNameConcatNum,
		formatInitialLastNum,
	}
	chosenFormat := formats[rand.Intn(len(formats))]

	var username string
	switch chosenFormat {
	case formatFullNameNum, formatFullNameUnderscoreNum, formatFullNameConcatNum:
		username = fmt.Sprintf(chosenFormat, cleanFirstName, cleanLastName, randomNumber)

	case formatLastNameNum:
		username = fmt.Sprintf(chosenFormat, cleanLastName, randomNumber)

	case formatInitialLastNum:
		if len(cleanFirstName) > 0 {
			initial := string(cleanFirstName[0])
			username = fmt.Sprintf("%s.%s%d", initial, cleanLastName, randomNumber)
		} else {
			username = fmt.Sprintf("%s%d", cleanLastName, randomNumber)
		}
	}

	return username
}

type Property struct {
	ID        int64
	Floors    int
	Status    int64
	OwnerID   sql.NullInt64
	City      string
	Street    string
	Number    string
	Latitude  float64
	Longitude float64
}

type Household struct {
	ID              int64
	Floor           int64
	Suite           string
	Status          int64
	SqFootage       float64
	OwnerID         sql.NullInt64
	PropertyID      int64
	CadastralNumber string
	DeviceStatusID  string
}

type DeviceStatus struct {
	DeviceID    string
	IsActive    bool
	Address     string
	HouseholdID sql.NullInt64
}

func writeSimulatorsToCSV(statuses []DeviceStatus, count int) error {
	if count > len(statuses) {
		count = len(statuses)
		log.Printf("Warning: Requested %d simulators, but only %d are available. Exporting all.", count, len(statuses))
	}
	if count == 0 {
		return nil
	}

	rand.Shuffle(len(statuses), func(i, j int) {
		statuses[i], statuses[j] = statuses[j], statuses[i]
	})

	selectedStatuses := statuses[:count]

	fileName := "simulators.csv"
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"device_id", "city"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	for _, status := range selectedStatuses {
		parts := strings.Split(status.Address, ",")
		var city string
		if len(parts) > 0 {
			city = strings.TrimSpace(parts[1])
		} else {
			city = "Unknown"
		}

		row := []string{status.DeviceID, city}
		if err := writer.Write(row); err != nil {
			log.Printf("Warning: failed to write row for device %s: %v", status.DeviceID, err)
		}
	}

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}
	fmt.Println("Successfully connected to PostgreSQL.")

	truncateTables(db)
	userIDs := populateUsers(db)

	fmt.Println("\nPhase 1: Generating all data in memory...")
	startTime := time.Now()
	properties, households, deviceStatuses := generateAllData(userIDs)
	fmt.Printf("Data generation finished in %v.\n", time.Since(startTime))

	fmt.Println("\nPhase 2: Inserting data into the database...")
	insertDeviceStatuses(db, deviceStatuses)
	insertProperties(db, properties)
	insertHouseholds(db, households)
	updateDeviceStatusHouseholdIDs(db)

	fmt.Println("\nPhase 3: Exporting 1000 random simulators to CSV file...")
	if err := writeSimulatorsToCSV(deviceStatuses, 1000); err != nil {
		log.Fatalf("Fatal error while writing simulators to CSV: %v", err)
	}
	fmt.Println("Successfully wrote simulator data to simulators.csv")

	fmt.Println("\nDatabase population completed successfully!")
}

func generateAllData(userIDs []int64) ([]Property, []Household, []DeviceStatus) {
	var propertyIDCounter, householdIDCounter int64

	allProperties := make([]Property, 0, totalHouseholds/5)
	allHouseholds := make([]Household, 0, totalHouseholds)
	allDeviceStatuses := make([]DeviceStatus, 0, totalHouseholds)

	var generatedHouseholds int64
	for generatedHouseholds < totalHouseholds {
		propID := propertyIDCounter + 1
		propertyIDCounter++

		propAddress := fmt.Sprintf("%s %s, %s", streets[rand.Intn(len(streets))], fmt.Sprintf("%d", rand.Intn(200)+1), cities[rand.Intn(len(cities))])

		prop := Property{
			ID:        propID,
			Floors:    rand.Intn(25) + 1,
			Status:    2,
			OwnerID:   sql.NullInt64{Int64: userIDs[rand.Intn(len(userIDs))], Valid: true},
			City:      cities[rand.Intn(len(cities))],
			Street:    streets[rand.Intn(len(streets))],
			Number:    fmt.Sprintf("%d", rand.Intn(200)+1),
			Latitude:  44.7866 + (rand.Float64()-0.5)*0.5,
			Longitude: 20.4489 + (rand.Float64()-0.5)*0.5,
		}
		allProperties = append(allProperties, prop)

		householdsInProperty := rand.Intn(25) + 1
		for i := 0; i < householdsInProperty && generatedHouseholds < totalHouseholds; i++ {
			householdID := householdIDCounter + 1
			householdIDCounter++

			var householdOwnerID sql.NullInt64
			var householdStatus int64
			if rand.Float64() < 0.70 {
				householdStatus = 1
				householdOwnerID = sql.NullInt64{Int64: userIDs[rand.Intn(len(userIDs))], Valid: true}
			} else {
				householdStatus = 2
				householdOwnerID = sql.NullInt64{Valid: false}
			}

			deviceID := fmt.Sprintf("household_%d", householdID)

			household := Household{
				ID:              householdID,
				Floor:           int64(rand.Intn(prop.Floors)),
				Suite:           fmt.Sprintf("Stan %d", rand.Intn(100)+1),
				Status:          householdStatus,
				SqFootage:       30.0 + rand.Float64()*120.0,
				OwnerID:         householdOwnerID,
				PropertyID:      prop.ID,
				CadastralNumber: fmt.Sprintf("CAD-%d-%d", prop.ID, householdID),
				DeviceStatusID:  deviceID,
			}
			allHouseholds = append(allHouseholds, household)

			deviceStatus := DeviceStatus{
				DeviceID:    deviceID,
				IsActive:    rand.Float32() < 0.95,
				Address:     propAddress,
				HouseholdID: sql.NullInt64{Int64: householdID, Valid: true},
			}
			allDeviceStatuses = append(allDeviceStatuses, deviceStatus)

			generatedHouseholds++
		}
	}
	return allProperties, allHouseholds, allDeviceStatuses
}

func truncateTables(db *sql.DB) {
	fmt.Println("Truncating tables...")
	query := `TRUNCATE TABLE public.device_status, public.households, public.properties RESTART IDENTITY CASCADE;`
	if _, err := db.Exec(query); err != nil {
		log.Fatalf("Error truncating tables: %v", err)
	}
}

func populateUsers(db *sql.DB) []int64 {
	fmt.Printf("Populating %d users with realistic names...\n", totalUsers)
	userIDs := make([]int64, totalUsers)
	var userIDCounter int64 = 10
	usedUsernames := make(map[string]bool)

	for i := 0; i < totalUsers; i += insertBatchSize {
		end := i + insertBatchSize
		if end > totalUsers {
			end = totalUsers
		}

		txn, _ := db.Begin()
		stmt, _ := txn.Prepare(pq.CopyIn("users", "id", "username", "password", "email", "role", "first_name", "last_name", "status"))

		for j := i; j < end; j++ {
			userID := userIDCounter + 1
			userIDCounter++
			userIDs[j] = userID

			firstName := serbianFirstNames[rand.Intn(len(serbianFirstNames))]
			lastName := serbianSurnames[rand.Intn(len(serbianSurnames))]
			var username string
			for {
				username = generateSerbianUsername(firstName, lastName)
				if !usedUsernames[username] {
					usedUsernames[username] = true
					break
				}
			}
			email := fmt.Sprintf("%s@example.com", username)
			_, _ = stmt.Exec(userID, username, "hashed_password", email, 0, firstName, lastName, 0)
		}

		_, _ = stmt.Exec()
		_ = stmt.Close()
		if err := txn.Commit(); err != nil {
			log.Fatalf("User commit error: %v", err)
		}
		fmt.Printf("\r  -> Inserted users up to %d", end)
	}
	fmt.Println("\nUsers populated.")
	return userIDs
}

func insertDeviceStatuses(db *sql.DB, statuses []DeviceStatus) {
	fmt.Printf("Inserting %d device_status records...\n", len(statuses))

	for i := 0; i < len(statuses); i += insertBatchSize {
		end := i + insertBatchSize
		if end > len(statuses) {
			end = len(statuses)
		}
		batch := statuses[i:end]

		txn, _ := db.Begin()
		stmt, _ := txn.Prepare(pq.CopyIn("device_status", "device_id", "is_active", "address"))
		for _, ds := range batch {
			_, err := stmt.Exec(ds.DeviceID, ds.IsActive, ds.Address)
			if err != nil {
				log.Fatalf("DeviceStatus COPY exec error: %v", err)
			}
		}
		if _, err := stmt.Exec(); err != nil {
			log.Fatalf("DeviceStatus COPY finalize error: %v", err)
		}
		_ = stmt.Close()
		if err := txn.Commit(); err != nil {
			log.Fatalf("DeviceStatus commit error: %v", err)
		}
		fmt.Printf("\r  -> Inserted device_status up to %d", end)
	}
	fmt.Println("\nDevice Statuses inserted.")
}

func insertProperties(db *sql.DB, properties []Property) {
	fmt.Printf("Inserting %d properties...\n", len(properties))

	for i := 0; i < len(properties); i += insertBatchSize {
		end := i + insertBatchSize
		if end > len(properties) {
			end = len(properties)
		}
		batch := properties[i:end]

		txn, _ := db.Begin()
		stmt, _ := txn.Prepare(pq.CopyIn("properties", "id", "floors", "status", "owner_id", "city", "street", "number", "latitude", "longitude", "created_at"))
		for _, p := range batch {
			_, err := stmt.Exec(p.ID, p.Floors, p.Status, p.OwnerID, p.City, p.Street, p.Number, p.Latitude, p.Longitude, time.Now())
			if err != nil {
				log.Fatalf("Property COPY exec error: %v", err)
			}
		}
		if _, err := stmt.Exec(); err != nil {
			log.Fatalf("Property COPY finalize error: %v", err)
		}
		_ = stmt.Close()
		if err := txn.Commit(); err != nil {
			log.Fatalf("Property commit error: %v", err)
		}
		fmt.Printf("\r  -> Inserted properties up to %d", end)
	}
	fmt.Println("\nProperties inserted.")
}

func insertHouseholds(db *sql.DB, households []Household) {
	fmt.Printf("Inserting %d households...\n", len(households))

	for i := 0; i < len(households); i += insertBatchSize {
		end := i + insertBatchSize
		if end > len(households) {
			end = len(households)
		}
		batch := households[i:end]

		txn, _ := db.Begin()
		stmt, _ := txn.Prepare(pq.CopyIn("households", "id", "floor", "suite", "status", "sq_footage", "owner_id", "property_id", "cadastral_number", "device_status_id"))
		for _, h := range batch {
			sqFootageStr := strconv.FormatFloat(h.SqFootage, 'f', 2, 64)
			_, err := stmt.Exec(h.ID, h.Floor, h.Suite, h.Status, sqFootageStr, h.OwnerID, h.PropertyID, h.CadastralNumber, h.DeviceStatusID)
			if err != nil {
				log.Fatalf("Household COPY exec error: %v", err)
			}
		}
		if _, err := stmt.Exec(); err != nil {
			log.Fatalf("Household COPY finalize error: %v", err)
		}
		_ = stmt.Close()
		if err := txn.Commit(); err != nil {
			log.Fatalf("Household commit error: %v", err)
		}
		fmt.Printf("\r  -> Inserted households up to %d", end)
	}
	fmt.Println("\nHouseholds inserted.")
}

func updateDeviceStatusHouseholdIDs(db *sql.DB) {
	fmt.Println("Updating device_status records with household IDs...")
	startTime := time.Now()
	query := `
		UPDATE public.device_status ds
		SET household_id = h.id
		FROM public.households h
		WHERE ds.device_id = h.device_status_id;
	`
	result, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to update device_status with household IDs: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Update complete. %d device_status records updated in %v.\n", rowsAffected, time.Since(startTime))
}
