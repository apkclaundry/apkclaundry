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

// formatDate mengubah time.Time menjadi string dengan format dd/mm/yyyy
func formatDate(date time.Time) string {
	return date.Format("02/01/2006")
}

// CreateTransaction handles the creation of a new transaction
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
		return
	}

	// Set TransactionDate to the current time if not provided
	if transaction.TransactionDate.IsZero() {
		transaction.TransactionDate = time.Now()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := config.TransactionCollection.InsertOne(ctx, transaction)
	if err != nil {
		http.Error(w, `{"error": "Failed to create transaction"}`, http.StatusInternalServerError)
		return
	}

	transaction.ID = result.InsertedID.(primitive.ObjectID).Hex()

	response := map[string]interface{}{
		"message":     "Transaction created successfully",
		"transaction": transaction,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAllTransactions retrieves all transactions from the database
func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	cursor, err := config.TransactionCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var transactions []models.Transaction
	for cursor.Next(context.TODO()) {
		var transaction models.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			http.Error(w, "Failed to read transaction data", http.StatusInternalServerError)
			return
		}
		// Format TransactionDate ke dd/mm/yyyy
		transaction.TransactionDateFormatted = formatDate(transaction.TransactionDate)
		transactions = append(transactions, transaction)
	}

	if len(transactions) == 0 {
		http.Error(w, "No transactions found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// GetTransactionByID retrieves a transaction by its ID
func GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	transactionID := r.URL.Query().Get("id")
	if transactionID == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var transaction models.Transaction
	err = config.TransactionCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&transaction)
	if err != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	// Format TransactionDate ke dd/mm/yyyy
	transaction.TransactionDateFormatted = formatDate(transaction.TransactionDate)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

// UpdateTransaction updates a transaction's data by its ID
func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID := r.URL.Query().Get("id")
	if transactionID == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedTransaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&updatedTransaction); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"customer_name":   updatedTransaction.CustomerName,
			"phone_number":    updatedTransaction.PhoneNumber,
			"service_type":    updatedTransaction.ServiceType,
			"weight_per_kg":   updatedTransaction.WeightPerKg,
			"total_price":     updatedTransaction.TotalPrice,
			"payment_method":  updatedTransaction.PaymentMethod,
			"transaction_date": updatedTransaction.TransactionDate,
		},
	}

	result, err := config.TransactionCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction updated successfully"})
}

// DeleteTransaction deletes a transaction by its ID
func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID := r.URL.Query().Get("id")
	if transactionID == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := config.TransactionCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Failed to delete transaction", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction deleted successfully"})
}
