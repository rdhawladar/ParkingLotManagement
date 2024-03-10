package model

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
