package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"apkclaundry/config"
	"apkclaundry/models"
	"apkclaundry/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Register handles user registration
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

	// Set hired_date and salary_date
	user.HiredDate = time.Now()
	user.SalaryDate = nil // SalaryDate bisa null

	// Insert the user into the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert user into the UserCollection
	result, err := config.UserCollection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, `{"error": "Failed to create user"}`, http.StatusInternalServerError)
		return
	}

	// Assign generated ID to user
	user.ID = result.InsertedID.(primitive.ObjectID).Hex()

	// Insert user into the EmployeeCollection
	_, err = config.EmployeeCollection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, `{"error": "Failed to create employee"}`, http.StatusInternalServerError)
		return
	}

	// Handle SalaryDate safely
	var salaryDate string
	if user.SalaryDate != nil {
		salaryDate = user.SalaryDate.Format("02/01/2006")
	} else {
		salaryDate = ""
	}

	response := map[string]interface{}{
		"message": "User registered successfully",
		"user": map[string]interface{}{
			"id":          user.ID,
			"username":    user.Username,
			"role":        user.Role,
			"phone":       user.Phone,
			"address":     user.Address,
			"hired_date":  user.HiredDate.Format("02/01/2006"),
			"salary_date": salaryDate,
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
			"id":         user.ID,
			"username":   user.Username,
			"role":       user.Role,
			"phone":      user.Phone,
			"address":    user.Address,
			"salary":     user.Salary,
			"hired_date": user.HiredDate.Format("02/01/2006"), // Format tanggal ke dd/mm/yyyy
		},
	}

	w.Header().Set("Content-Type", "application/json") // Set content type to JSON
	// Return the response as JSON
	json.NewEncoder(w).Encode(response)
}

// GetAllUsers mengambil semua data user
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Query semua user dari MongoDB
	cursor, err := config.UserCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Gagal mengambil data user", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	// Slice untuk menyimpan hasil query
	var users []models.User
	for cursor.Next(context.TODO()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			http.Error(w, "Gagal membaca data user", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Periksa apakah ada user yang ditemukan
	if len(users) == 0 {
		http.Error(w, "Tidak ada user yang ditemukan", http.StatusNotFound)
		return
	}

	// Format hired_date dan salary_date
	type UserResponse struct {
		ID         string  `json:"id"`
		Username   string  `json:"username"`
		Role       string  `json:"role"`
		Phone      string  `json:"phone"`
		Address    string  `json:"address"`
		Salary     float64 `json:"salary"`
		HiredDate  string  `json:"hired_date"`
		SalaryDate string  `json:"salary_date"`
	}

	var formattedUsers []UserResponse
	for _, user := range users {
		var salaryDate string
		if user.SalaryDate != nil {
			salaryDate = user.SalaryDate.Format("02/01/2006")
		} else {
			salaryDate = ""
		}

		formattedUsers = append(formattedUsers, UserResponse{
			ID:         user.ID,
			Username:   user.Username,
			Role:       user.Role,
			Phone:      user.Phone,
			Address:    user.Address,
			Salary:     user.Salary,
			HiredDate:  user.HiredDate.Format("02/01/2006"),
			SalaryDate: salaryDate,
		})
	}

	// Kirim response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(formattedUsers)
}

// GetUserByID mengambil data user berdasarkan ID
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL parameter
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "ID tidak disediakan", http.StatusBadRequest)
		return
	}

	// Convert ID menjadi ObjectID MongoDB
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	// Query MongoDB untuk mencari user berdasarkan ID
	var user models.User
	err = config.UserCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		http.Error(w, "User tidak ditemukan", http.StatusNotFound)
		return
	}

	// Handle SalaryDate safely
	var salaryDate string
	if user.SalaryDate != nil {
		salaryDate = user.SalaryDate.Format("02/01/2006")
	} else {
		salaryDate = ""
	}

	// Format response
	type UserResponse struct {
		ID         string  `json:"id"`
		Username   string  `json:"username"`
		Role       string  `json:"role"`
		Phone      string  `json:"phone"`
		Address    string  `json:"address"`
		Salary     float64 `json:"salary"`
		SalaryDate string  `json:"salary_date"`
		HiredDate  string  `json:"hired_date"`
	}

	formattedUser := UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Role:       user.Role,
		Phone:      user.Phone,
		Address:    user.Address,
		Salary:     user.Salary,
		HiredDate:  user.HiredDate.Format("02/01/2006"),
		SalaryDate: salaryDate,
	}

	// Kirim response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(formattedUser)
}


// UpdateUser memperbarui data user berdasarkan ID
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL parameter
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "ID tidak disediakan", http.StatusBadRequest)
		return
	}

	// Convert ID menjadi ObjectID MongoDB
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	// Decode data JSON dari body request
	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Input tidak valid", http.StatusBadRequest)
		return
	}

	// Query untuk update user berdasarkan ID
	update := bson.M{
		"$set": bson.M{
			"username":    updatedUser.Username,
			"role":        updatedUser.Role,
			"phone":       updatedUser.Phone,
			"address":     updatedUser.Address,
			"salary":      updatedUser.Salary,
			"salary_date": updatedUser.SalaryDate,
		},
	}

	// Update user di database
	result, err := config.UserCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, "Gagal memperbarui user", http.StatusInternalServerError)
		return
	}

	// Periksa apakah ada perubahan
	if result.MatchedCount == 0 {
		http.Error(w, "User tidak ditemukan", http.StatusNotFound)
		return
	}

	// Kirim response sukses
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User berhasil diperbarui"})
}

// DeleteUser menghapus data user berdasarkan ID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL parameter
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "ID tidak disediakan", http.StatusBadRequest)
		return
	}

	// Convert ID menjadi ObjectID MongoDB
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	// Hapus user dari koleksi UserCollection
	result, err := config.UserCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Gagal menghapus user", http.StatusInternalServerError)
		return
	}

	// Periksa apakah ada data yang dihapus
	if result.DeletedCount == 0 {
		http.Error(w, "User tidak ditemukan", http.StatusNotFound)
		return
	}

	// Kirim response sukses
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User berhasil dihapus"})
}

// GetAllEmployeesIDName mengambil semua data karyawan dan hanya mengembalikan id dan nama
func GetAllEmployeesIDName(w http.ResponseWriter, r *http.Request) {
	// Query semua karyawan dari MongoDB
	cursor, err := config.EmployeeCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Gagal mengambil data karyawan", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	// Slice untuk menyimpan hasil query
	type EmployeeResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	var employees []EmployeeResponse
	for cursor.Next(context.TODO()) {
		var employee models.User // Menggunakan models.User karena data karyawan disimpan dalam User
		if err := cursor.Decode(&employee); err != nil {
			http.Error(w, "Gagal membaca data karyawan", http.StatusInternalServerError)
			return
		}
		employees = append(employees, EmployeeResponse{
			ID:   employee.ID,
			Name: employee.Username,
		})
	}

	// Kirim response dengan daftar karyawan yang sudah diformat
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}
