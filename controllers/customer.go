package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"apkclaundry/config"
	"apkclaundry/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateCustomer handles the creation of a new customer
func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
		return
	}

	// Insert the customer into the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := config.CustomerCollection.InsertOne(ctx, customer)
	if err != nil {
		http.Error(w, `{"error": "Failed to create customer"}`, http.StatusInternalServerError)
		return
	}

	// Assign generated ID to customer
	customer.ID = result.InsertedID.(primitive.ObjectID).Hex()

	response := map[string]interface{}{
		"message":  "Customer created successfully",
		"customer": customer,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAllCustomers retrieves all customers from the database
func GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	cursor, err := config.CustomerCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch customers", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var customers []models.Customer
	for cursor.Next(context.TODO()) {
		var customer models.Customer
		if err := cursor.Decode(&customer); err != nil {
			http.Error(w, "Failed to read customer data", http.StatusInternalServerError)
			return
		}
		customers = append(customers, customer)
	}

	if len(customers) == 0 {
		http.Error(w, "No customers found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

// GetAllCustomersIDName retrieves all customers with only ID and name
func GetAllCustomersIDName(w http.ResponseWriter, r *http.Request) {
	cursor, err := config.CustomerCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch customers", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var customers []struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}
	for cursor.Next(context.TODO()) {
		var customer struct {
			ID    primitive.ObjectID `bson:"_id"`
			Name  string             `bson:"name"`
			Phone string             `json:"phone"`
		}
		if err := cursor.Decode(&customer); err != nil {
			http.Error(w, "Failed to read customer data", http.StatusInternalServerError)
			return
		}

		// Konversi ObjectID ke string
		customers = append(customers, struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Phone string `json:"phone"`
		}{
			ID:   customer.ID.Hex(),
			Name: customer.Name,
			Phone: customer.Phone,
		})
	}

	if len(customers) == 0 {
		http.Error(w, "No customers found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

// GetCustomerByID retrieves a customer by their ID
func GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("id")
	if customerID == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var customer models.Customer
	err = config.CustomerCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&customer)
	if err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

// GetCustomerNameByID retrieves only the name of a customer by their ID
func GetCustomerNameByID(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("id")
	if customerID == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var result struct {
		Name string `json:"name"`
		Phone string `json:"phone"`
	}

	err = config.CustomerCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result)
	if err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// UpdateCustomer updates a customer's data by their ID
func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("id")
	if customerID == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedCustomer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&updatedCustomer); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"name":    updatedCustomer.Name,
			"phone":   updatedCustomer.Phone,
			"address": updatedCustomer.Address,
			"email":   updatedCustomer.Email,
		},
	}

	result, err := config.CustomerCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, "Failed to update customer", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Customer updated successfully"})
}

// DeleteCustomer deletes a customer by their ID
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("id")
	if customerID == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := config.CustomerCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Failed to delete customer", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Customer deleted successfully"})
}
