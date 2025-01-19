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

// CreateSupplier handles the creation of a new supplier
func CreateSupplier(w http.ResponseWriter, r *http.Request) {
	var supplier models.Supplier
	if err := json.NewDecoder(r.Body).Decode(&supplier); err != nil {
		http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
		return
	}

	// Insert the supplier into the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := config.SupplierCollection.InsertOne(ctx, supplier)
	if err != nil {
		http.Error(w, `{"error": "Failed to create supplier"}`, http.StatusInternalServerError)
		return
	}

	// Assign generated ID to supplier
	supplier.ID = result.InsertedID.(primitive.ObjectID).Hex()

	response := map[string]interface{}{
		"message": "Supplier created successfully",
		"supplier": supplier,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAllSuppliers retrieves all suppliers from the database
func GetAllSuppliers(w http.ResponseWriter, r *http.Request) {
	cursor, err := config.SupplierCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch suppliers", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var suppliers []models.Supplier
	for cursor.Next(context.TODO()) {
		var supplier models.Supplier
		if err := cursor.Decode(&supplier); err != nil {
			http.Error(w, "Failed to read supplier data", http.StatusInternalServerError)
			return
		}
		suppliers = append(suppliers, supplier)
	}

	if len(suppliers) == 0 {
		http.Error(w, "No suppliers found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suppliers)
}

// GetSupplierByID retrieves a supplier by their ID
func GetSupplierByID(w http.ResponseWriter, r *http.Request) {
	supplierID := r.URL.Query().Get("id")
	if supplierID == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(supplierID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var supplier models.Supplier
	err = config.SupplierCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&supplier)
	if err != nil {
		http.Error(w, "Supplier not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(supplier)
}

// UpdateSupplier updates a supplier's data by their ID
func UpdateSupplier(w http.ResponseWriter, r *http.Request) {
	supplierID := r.URL.Query().Get("id")
	if supplierID == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(supplierID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedSupplier models.Supplier
	if err := json.NewDecoder(r.Body).Decode(&updatedSupplier); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"supplier_name":   updatedSupplier.SupplierName,
			"phone_number":    updatedSupplier.PhoneNumber,
			"address":         updatedSupplier.Address,
			"email":           updatedSupplier.Email,
			"supplied_products": updatedSupplier.SuppliedProducts,
		},
	}

	result, err := config.SupplierCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, "Failed to update supplier", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Supplier not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Supplier updated successfully"})
}

// DeleteSupplier deletes a supplier by their ID
func DeleteSupplier(w http.ResponseWriter, r *http.Request) {
	supplierID := r.URL.Query().Get("id")
	if supplierID == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(supplierID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := config.SupplierCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Failed to delete supplier", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Supplier not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Supplier deleted successfully"})
}
