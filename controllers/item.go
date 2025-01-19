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

// CreateItem handles the creation of a new item
func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := config.ItemCollection.InsertOne(ctx, item)
	if err != nil {
		http.Error(w, `{"error": "Failed to create item"}`, http.StatusInternalServerError)
		return
	}

	item.ID = result.InsertedID.(primitive.ObjectID).Hex()

	response := map[string]interface{}{
		"message": "Item created successfully",
		"item":    item,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAllItems retrieves all items from the database
func GetAllItems(w http.ResponseWriter, r *http.Request) {
	cursor, err := config.ItemCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var items []models.Item
	for cursor.Next(context.TODO()) {
		var item models.Item
		if err := cursor.Decode(&item); err != nil {
			http.Error(w, "Failed to read item data", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if len(items) == 0 {
		http.Error(w, "No items found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// GetItemByID retrieves an item by its ID
func GetItemByID(w http.ResponseWriter, r *http.Request) {
	itemID := r.URL.Query().Get("id")
	if itemID == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var item models.Item
	err = config.ItemCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&item)
	if err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// UpdateItem updates an item's data by its ID
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	itemID := r.URL.Query().Get("id")
	if itemID == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedItem models.Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"item_name": updatedItem.ItemName,
			"quantity":  updatedItem.Quantity,
			"price":     updatedItem.Price,
		},
	}

	result, err := config.ItemCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Item updated successfully"})
}

// DeleteItem deletes an item by its ID
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	itemID := r.URL.Query().Get("id")
	if itemID == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := config.ItemCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Item deleted successfully"})
}
