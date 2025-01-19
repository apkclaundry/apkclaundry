package models

import "time"

// Customer represents a laundry customer
type Customer struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Phone     string    `json:"phone" bson:"phone"`
	Address   string    `json:"address" bson:"address"`
	Email     string	`json:"email" bson:"email"`
}

// User represents an employee or system user
type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Username  string    `json:"username" bson:"username"`
	Password  string    `json:"password" bson:"password"` // Should be hashed
	Role      string    `json:"role" bson:"role"`         // e.g., "admin" or "staff"
	Phone     string    `json:"phone" bson:"phone"`         // Contact number
	Address   string    `json:"address" bson:"address"`     // Home address
	HiredDate time.Time `json:"hired_date" bson:"hired_date"` // Date of hiring
}

// Employee represents an employee in the laundry business
type Employee struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`           // Full name of the employee
	Phone     string    `json:"phone" bson:"phone"`         // Contact number
	Address   string    `json:"address" bson:"address"`     // Home address
	Position  string    `json:"position" bson:"position"`   // Job position (e.g., "Admin", "Cashier", "Operator")
	Salary    float64   `json:"salary" bson:"salary"`       // Monthly salary
	HiredDate time.Time `json:"hired_date" bson:"hired_date"` // Date of hiring
}

// Supplier represents the supplier model
type Supplier struct {
	ID             string 			  `json:"id" bson:"_id,omitempty"`
	SupplierName   string             `json:"supplier_name" bson:"supplier_name"`
	PhoneNumber    string             `json:"phone_number" bson:"phone_number"`
	Address        string             `json:"address" bson:"address"`
	Email          string             `json:"email" bson:"email"`
	SuppliedProducts []string         `json:"supplied_products" bson:"supplied_products"`
}

// Inventory represents a stock item in the laundry
type Item struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	ItemName    string    `json:"item_name" bson:"item_name"`
	Quantity    int       `json:"quantity" bson:"quantity"`
	Price       float64   `json:"price" bson:"price"`
}



type Transaction struct {
	ID           string  `json:"id" bson:"_id,omitempty"`
	CustomerName string  `json:"customer_name" bson:"customer_name"`
	PhoneNumber  string  `json:"phone_number" bson:"phone_number"`
	ServiceType  string  `json:"service_type" bson:"service_type"`
	WeightPerKg  float64 `json:"weight_per_kg" bson:"weight_per_kg"`
	TotalPrice   float64 `json:"total_price" bson:"total_price"`
}


// Payment represents a payment transaction
type Payment struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	CustomerID  string    `json:"customer_id" bson:"customer_id"`
	Amount      float64   `json:"amount" bson:"amount"`
	PaymentType string    `json:"payment_type" bson:"payment_type"` // e.g., "cash", "card"
	Date        time.Time `json:"date" bson:"date"`
}
