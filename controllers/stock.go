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

// GetStockTransactions retrieves all stock transactions
func GetStockTransactions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := config.StockTransactionCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch transactions"}`, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var transactions []models.StockTransaction
	for cursor.Next(ctx) {
		var transaction models.StockTransaction
		if err := cursor.Decode(&transaction); err != nil {
			http.Error(w, `{"error": "Failed to decode transaction data"}`, http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, transaction)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// GetStockTransactionByID retrieves a stock transaction by ID
func GetStockTransactionByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error": "ID is required"}`, http.StatusBadRequest)
		return
	}

	var transaction models.StockTransaction
	err := config.StockTransactionCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&transaction)
	if err != nil {
		http.Error(w, `{"error": "Transaction not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

// CreateStockTransaction creates a new stock transaction
func CreateStockTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.StockTransaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	transaction.IDTransaction = primitive.NewObjectID().Hex()
	transaction.Date = time.Now()

	_, err := config.StockTransactionCollection.InsertOne(context.TODO(), transaction)
	if err != nil {
		http.Error(w, `{"error": "Failed to create transaction"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

// UpdateStockTransaction updates a stock transaction by ID
func UpdateStockTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error": "ID is required"}`, http.StatusBadRequest)
		return
	}

	var updatedTransaction models.StockTransaction
	if err := json.NewDecoder(r.Body).Decode(&updatedTransaction); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": updatedTransaction}

	_, err := config.StockTransactionCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, `{"error": "Failed to update transaction"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTransaction)
}

// DeleteStockTransaction deletes a stock transaction by ID
func DeleteStockTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error": "ID is required"}`, http.StatusBadRequest)
		return
	}

	_, err := config.StockTransactionCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, `{"error": "Failed to delete transaction"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
