package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"apkclaundry/config"
	"apkclaundry/models"
	"apkclaundry/utils"

	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
)

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
		return
	}

	// Check if username already exists
	var existingUser models.User
	err := config.UserCollection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		http.Error(w, `{"error": "Username already exists"}`, http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error": "Failed to hash password"}`, http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Insert the user into the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := config.UserCollection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, `{"error": "Failed to create user"}`, http.StatusInternalServerError)
		return
	}

	// Retrieve the inserted user including the new ID
	err = config.UserCollection.FindOne(context.TODO(), bson.M{"_id": result.InsertedID}).Decode(&user)
	if err != nil {
		http.Error(w, `{"error": "Failed to retrieve user"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User registered successfully",
		"user": map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"role":       user.Role,
			"phone":      user.Phone,
			"address":    user.Address,
			"hired_date": user.HiredDate.Format(time.RFC3339), // ISO 8601 date format
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


// Login handles user authentication and JWT token generation
func Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Find the user in the database
	var user models.User
	err := config.UserCollection.FindOne(context.TODO(), bson.M{"username": creds.Username}).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := map[string]interface{}{
		"message": "Login successful",
		"token":   token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
			"phone":    user.Phone,
			"address":  user.Address,
			"hired_date": user.HiredDate.Format(time.RFC3339), // format tanggal ke dalam format string yang dapat dibaca manusia
		},
	}

	w.Header().Set("Content-Type", "application/json") // Set content type to JSON
	// Return the response as JSON
	json.NewEncoder(w).Encode(response)
}

