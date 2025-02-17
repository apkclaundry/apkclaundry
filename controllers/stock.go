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

// CreateItemTransaction handles the creation of a new item transaction
func CreateItemTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.ItemTransaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
		return
	}

	transaction.Date = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := config.ItemTransactionCollection.InsertOne(ctx, transaction)
	if err != nil {
		http.Error(w, `{"error": "Failed to create transaction"}`, http.StatusInternalServerError)
		return
	}

	transaction.ID = result.InsertedID.(primitive.ObjectID).Hex()

	response := map[string]interface{}{
		"message": "Transaction created successfully",
		"transaction": transaction,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAllItemTransactions retrieves all item transactions from the database
func GetAllItemTransactions(w http.ResponseWriter, r *http.Request) {
	cursor, err := config.ItemTransactionCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch transactions"}`, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var transactions []models.ItemTransaction
	for cursor.Next(context.TODO()) {
		var transaction models.ItemTransaction
		if err := cursor.Decode(&transaction); err != nil {
			http.Error(w, `{"error": "Failed to read transaction data"}`, http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, transaction)
	}

	if len(transactions) == 0 {
		http.Error(w, `{"error": "No transactions found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// GetItemTransactionByID retrieves a transaction by its ID
func GetItemTransactionByID(w http.ResponseWriter, r *http.Request) {
	transactionID := r.URL.Query().Get("id")
	if transactionID == "" {
		http.Error(w, `{"error": "ID not provided"}`, http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	var transaction models.ItemTransaction
	err = config.ItemTransactionCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&transaction)
	if err != nil {
		http.Error(w, `{"error": "Transaction not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

// UpdateItemTransaction updates a transaction's data by its ID
func UpdateItemTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID := r.URL.Query().Get("id")
	if transactionID == "" {
		http.Error(w, `{"error": "ID not provided"}`, http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	var updatedTransaction models.ItemTransaction
	if err := json.NewDecoder(r.Body).Decode(&updatedTransaction); err != nil {
		http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"item_id":         updatedTransaction.ItemID,
			"item_name":       updatedTransaction.ItemName,
			"date":            updatedTransaction.Date,
			"transaction_type": updatedTransaction.TransactionType,
			"quantity":        updatedTransaction.Quantity,
			"stock_after":     updatedTransaction.StockAfter,
		},
	}

	result, err := config.ItemTransactionCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, `{"error": "Failed to update transaction"}`, http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, `{"error": "Transaction not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction updated successfully"})
}

// DeleteItemTransaction deletes a transaction by its ID
func DeleteItemTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID := r.URL.Query().Get("id")
	if transactionID == "" {
		http.Error(w, `{"error": "ID not provided"}`, http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	result, err := config.ItemTransactionCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, `{"error": "Failed to delete transaction"}`, http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, `{"error": "Transaction not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction deleted successfully"})
}

// GetItemTransactions retrieves item transactions with only selected fields (id, itemid, itemname)
func GetItemTransactions(w http.ResponseWriter, r *http.Request) {
	// Query to find all item transactions
	cursor, err := config.ItemTransactionCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch item transactions", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	// Define a struct for the response
	var transactions []struct {
		ID       primitive.ObjectID `bson:"_id"`
		ItemID   primitive.ObjectID `bson:"item_id"`
		ItemName string             `bson:"item_name"`
	}

	// Iterate over the cursor and decode only the selected fields
	for cursor.Next(context.TODO()) {
		var transaction struct {
			ID       primitive.ObjectID `bson:"_id"`
			ItemID   primitive.ObjectID `bson:"item_id"`
			ItemName string             `bson:"item_name"`
		}
		if err := cursor.Decode(&transaction); err != nil {
			http.Error(w, "Failed to read transaction data", http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, transaction)
	}

	// Check if no transactions were found
	if len(transactions) == 0 {
		http.Error(w, "No transactions found", http.StatusNotFound)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
