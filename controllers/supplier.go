package controllers

import (
	"context"
	"encoding/json"
	"log"
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

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := config.SupplierCollection.InsertOne(ctx, supplier)
    if err != nil {
        http.Error(w, `{"error": "Failed to create supplier"}`, http.StatusInternalServerError)
        return
    }

    // Konversi InsertedID ke ObjectID dengan pengecekan
    if objID, ok := result.InsertedID.(primitive.ObjectID); ok {
        supplier.ID = objID // Tetap sebagai ObjectID
    } else {
        http.Error(w, `{"error": "Failed to parse inserted ID"}`, http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "message":  "Supplier created successfully",
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
			"supplier_name":     updatedSupplier.SupplierName,
			"phone_number":      updatedSupplier.PhoneNumber,
			"address":           updatedSupplier.Address,
			"email":             updatedSupplier.Email,
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

func AddSupplierTransaction(w http.ResponseWriter, r *http.Request) {
	// Ambil ID Supplier dari query parameter
	supplierID := r.URL.Query().Get("supplier_id")
	if supplierID == "" {
		log.Println("Error: Supplier ID tidak disediakan")
		http.Error(w, `{"error": "Supplier ID tidak disediakan"}`, http.StatusBadRequest)
		return
	}
	log.Println("Received supplier_id:", supplierID)

	// Konversi supplierID ke ObjectID
	objID, err := primitive.ObjectIDFromHex(supplierID)
	if err != nil {
		log.Println("Error: Invalid Supplier ID:", err)
		http.Error(w, `{"error": "Supplier ID tidak valid"}`, http.StatusBadRequest)
		return
	}
	log.Println("Converted supplier_id to ObjectID:", objID)

	// Cek apakah supplier ada sebelum menambahkan transaksi
	var existingSupplier models.Supplier
	err = config.SupplierCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&existingSupplier)
	if err != nil {
		log.Println("Error: Supplier not found for supplier_id", supplierID)
		http.Error(w, `{"error": "Supplier tidak ditemukan"}`, http.StatusNotFound)
		return
	}
	log.Println("Supplier found:", existingSupplier.SupplierName)

	// Pastikan field transactions ada dan berupa array
	if existingSupplier.Transactions == nil {
		log.Println("Warning: transactions is null, initializing as empty array.")
		existingSupplier.Transactions = []models.SupplierTransaction{}
	}

	// Decode data transaksi dari request body
	var transaction models.SupplierTransaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		log.Println("Error: Invalid transaction input", err)
		http.Error(w, `{"error": "Input tidak valid"}`, http.StatusBadRequest)
		return
	}
	log.Println("Decoded transaction data:", transaction)

	// Set tanggal transaksi & buat ID unik
	transaction.Date = time.Now()
	transaction.TransactionID = primitive.NewObjectID().Hex()
	log.Println("Generated Transaction ID:", transaction.TransactionID)

	// Update supplier dengan menambahkan transaksi baru
	update := bson.M{"$push": bson.M{"transactions": transaction}}
	result, err := config.SupplierCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		log.Println("Error: Failed to add transaction to supplier:", err)
		http.Error(w, `{"error": "Gagal menambahkan transaksi ke database"}`, http.StatusInternalServerError)
		return
	}
	log.Println("MongoDB update result:", result)

	// Cek apakah ada dokumen yang berhasil diupdate
	if result.MatchedCount == 0 {
		log.Println("Error: Supplier not found for update")
		http.Error(w, `{"error": "Supplier tidak ditemukan"}`, http.StatusNotFound)
		return
	}

	// Kirim respons sukses
	log.Println("Transaction added successfully for supplier_id:", supplierID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaksi berhasil ditambahkan"})
}





