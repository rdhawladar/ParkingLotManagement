package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var db *sqlx.DB

func initDB() {
	// Construct the connection string
	// Example: "host=localhost user=postgres password=myPassword dbname=parking_lot sslmode=disable"
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	var err error
	db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database successfully!")
}

// ParkingLotRequest defines the expected request body for creating a parking lot
type ParkingLotRequest struct {
	Name        string `json:"name"`
	TotalSpaces int    `json:"total_spaces"`
	ManagerID   int    `json:"manager_id"`
}

// ParkingLot represents the ParkingLots table structure
type ParkingLot struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	TotalSpaces int    `db:"total_spaces"`
	ManagerID   int    `db:"manager_id"`
}

// VehicleParkingRequest defines the expected request body for parking a vehicle
type VehicleParkingRequest struct {
	LotID        int    `json:"lot_id"`
	LicensePlate string `json:"license_plate"`
	OwnerID      int    `json:"owner_id"`
}

// ParkingSlot represents the ParkingSlots table structure
type ParkingSlot struct {
	ID                 int  `db:"id"`
	ParkingLotID       int  `db:"parking_lot_id"`
	SlotNumber         int  `db:"slot_number"`
	IsAvailable        bool `db:"is_available"`
	IsUnderMaintenance bool `db:"is_under_maintenance"`
}

// VehicleUnparkingRequest defines the expected request body for unparking a vehicle
type VehicleUnparkingRequest struct {
	LicensePlate string `json:"license_plate"`
}

// ParkingSlotStatus represents the status of a parking slot.
type ParkingSlotStatus struct {
	SlotNumber         int    `db:"slot_number" json:"slot_number"`
	IsAvailable        bool   `db:"is_available" json:"is_available"`
	IsUnderMaintenance bool   `db:"is_under_maintenance" json:"is_under_maintenance"`
	VehicleLicense     string `db:"license_plate" json:"vehicle_license,omitempty"` // Empty if no vehicle is parked
}

// MaintenanceModeRequest represents the status of a parking slot.
type MaintenanceModeRequest struct {
	IsUnderMaintenance bool `json:"is_under_maintenance"`
}

func main() {
	// Initialize the database connection
	initDB()

	// Set up a Gin router
	router := gin.Default()

	// Define routes for the API
	router.POST("/parkinglots", createParkingLot)               // Create a parking lot
	router.POST("/park", parkVehicle)                           // Park a vehicle
	router.POST("/unpark", unparkVehicle)                       // Unpark a vehicle
	router.GET("/parkinglots/:id", viewParkingLotStatus)        // View parking lot status
	router.PUT("/slots/:id/maintenance", toggleMaintenanceMode) // Toggle maintenance mode for a slot
	router.GET("/reports/daily", dailyReport)                   // Get daily parking report

	// Start the Gin server on port 8080
	router.Run(":8080")
}

func createParkingLot(c *gin.Context) {
	var request ParkingLotRequest

	// Bind the JSON body to the ParkingLotRequest struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Using a named query with RETURNING clause to get the inserted ID
	var lotID int
	err := db.QueryRowx(
		"INSERT INTO ParkingLots (name, total_spaces, manager_id) VALUES ($1, $2, $3) RETURNING id",
		request.Name, request.TotalSpaces, request.ManagerID).Scan(&lotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert parking lot"})
		return
	}

	// Now lotID contains the ID of the newly created parking lot
	// Continue to create parking slots for the new parking lot
	tx, err := db.Beginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	for i := 1; i <= request.TotalSpaces; i++ {
		_, err = tx.Exec("INSERT INTO ParkingSlots (parking_lot_id, slot_number, is_available, is_under_maintenance) VALUES ($1, $2, TRUE, FALSE)", lotID, i)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert parking slots"})
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Parking lot created successfully", "lot_id": lotID})
}

func parkVehicle(c *gin.Context) {
	var request VehicleParkingRequest

	// Bind the JSON body to the VehicleParkingRequest struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.Beginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Check if the vehicle is already parked (has an ongoing parking session)
	var ongoingSessionCount int
	err = tx.Get(&ongoingSessionCount, "SELECT COUNT(*) FROM ParkingSessions WHERE vehicle_id = (SELECT id FROM Vehicles WHERE license_plate = $1) AND unparked_at IS NULL", request.LicensePlate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for existing parking session"})
		return
	}
	if ongoingSessionCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vehicle is already parked"})
		return
	}

	// Find the first available parking slot in the requested lot
	var slot ParkingSlot
	err = tx.QueryRowx(`
		SELECT id, parking_lot_id, slot_number, is_available, is_under_maintenance 
		FROM ParkingSlots 
		WHERE parking_lot_id = $1 AND is_available = TRUE AND is_under_maintenance = FALSE 
		ORDER BY slot_number 
		LIMIT 1`,
		request.LotID).StructScan(&slot)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No available parking slots"})
		return
	}

	// Mark the parking slot as occupied
	_, err = tx.Exec("UPDATE ParkingSlots SET is_available = FALSE WHERE id = $1", slot.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update parking slot status"})
		return
	}

	// Insert a new parking session record
	_, err = tx.Exec(`
		INSERT INTO ParkingSessions (vehicle_id, parking_lot_id, parking_slot_id, parked_at) 
		VALUES ((SELECT id FROM Vehicles WHERE license_plate = $1), $2, $3, CURRENT_TIMESTAMP)`,
		request.LicensePlate, request.LotID, slot.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create parking session. Please check if vehicle is exist."})
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vehicle parked successfully", "slot_id": slot.ID, "lot_id": request.LotID})
}

func unparkVehicle(c *gin.Context) {
	var request VehicleUnparkingRequest

	// Bind the JSON body to the VehicleUnparkingRequest struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.Beginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	var parkingSlotID int
	var parkedHours float64

	// Retrieve the parking slot ID and calculate parked hours for the given vehicle's active parking session
	err = tx.QueryRowx(`
		UPDATE ParkingSessions
		SET unparked_at = CURRENT_TIMESTAMP
		WHERE vehicle_id = (SELECT id FROM Vehicles WHERE license_plate = $1) AND unparked_at IS NULL
		RETURNING parking_slot_id, EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - parked_at))/3600 AS parked_hours`,
		request.LicensePlate).Scan(&parkingSlotID, &parkedHours)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve parking session, No unpark vehicle found."})
		return
	}

	// Calculate the parking fee, ensuring a minimum charge for 1 hour
	fee := calculateParkingFee(parkedHours)

	// Mark the parking slot as available
	_, err = tx.Exec("UPDATE ParkingSlots SET is_available = TRUE WHERE id = $1", parkingSlotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update parking slot status"})
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vehicle unparked successfully", "fee": fee})
}

func calculateParkingFee(parkedHours float64) float64 {
	// Ensure a minimum charge for 1 hour
	if parkedHours < 1 {
		parkedHours = 1
	}
	return 10 * math.Ceil(parkedHours)
}

func viewParkingLotStatus(c *gin.Context) {
	lotID := c.Param("id") // Get the parking lot ID from the route parameter

	// Prepare a slice to hold the results
	var slots []ParkingSlotStatus

	// Query the database for slot statuses in the specified lot
	err := db.Select(&slots, `
		SELECT 
			ps.slot_number, 
			ps.is_available, 
			ps.is_under_maintenance, 
			v.license_plate
		FROM ParkingSlots ps
		JOIN ParkingSessions psess ON ps.id = psess.parking_slot_id AND ps.is_available = FALSE AND psess.unparked_at IS NULL
		LEFT JOIN Vehicles v ON psess.vehicle_id = v.id
		WHERE ps.parking_lot_id = $1
		ORDER BY ps.slot_number`, lotID)

	if err != nil {
		log.Error().Err(err).Msg("viewParkingLotStatus: Failed to retrieve parking lot status")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve parking lot status"})
		return
	}

	// Return the slot statuses to the client
	c.JSON(http.StatusOK, gin.H{"parking_lot_id": lotID, "slots": slots})
}

func toggleMaintenanceMode(c *gin.Context) {
	var request MaintenanceModeRequest
	slotID := c.Param("id") // Assuming you're using the slot ID as a URL parameter

	// Bind the JSON body to the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the parking slot exists
	var exists bool
	err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM ParkingSlots WHERE id = $1)", slotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check slot existence"})
		return
	}
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slot ID"})
		return
	}

	// Perform the update on the specified parking slot
	result, err := db.Exec(`
        UPDATE ParkingSlots
        SET is_under_maintenance = $1
        WHERE id = $2
    `, request.IsUnderMaintenance, slotID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update maintenance mode"})
		return
	}

	// Check if any row was actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slot ID or no update required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Maintenance mode updated successfully", "slot_id": slotID, "is_under_maintenance": request.IsUnderMaintenance})
}

func dailyReport(c *gin.Context) {
	// Optionally, get the date for which the report is requested from query parameters
	reportDate := c.Query("date") // Assuming the date is passed as a query parameter in 'YYYY-MM-DD' format

	// If no date is provided, use the current date
	if reportDate == "" {
		reportDate = time.Now().Format("2006-01-02") // Use current date
	}

	// Prepare a struct to hold the report data
	var report struct {
		TotalVehicles      int     `db:"total_vehicles"`
		TotalParkingTime   float64 `db:"total_parking_time"` // In hours
		TotalFeesCollected float64 `db:"total_fees_collected"`
	}

	// Query to calculate the daily report
	err := db.Get(&report, `
        SELECT 
            COUNT(*) AS total_vehicles, 
            SUM(EXTRACT(EPOCH FROM (COALESCE(unparked_at, CURRENT_TIMESTAMP) - parked_at))/3600) AS total_parking_time,
            SUM(CEIL(EXTRACT(EPOCH FROM (COALESCE(unparked_at, CURRENT_TIMESTAMP) - parked_at))/3600) * 10) AS total_fees_collected
        FROM ParkingSessions
        WHERE DATE(parked_at) = $1 OR DATE(COALESCE(unparked_at, CURRENT_TIMESTAMP)) = $1
    `, reportDate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate daily report"})
		return
	}

	// Return the report data
	c.JSON(http.StatusOK, gin.H{
		"date":                 reportDate,
		"total_vehicles":       report.TotalVehicles,
		"total_parking_time":   report.TotalParkingTime,
		"total_fees_collected": report.TotalFeesCollected,
	})
}
