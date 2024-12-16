package routes

import (
	"apkclaundry/controllers"
	// "apkclaundry/middleware"
	"net/http"
)

func InitRoutes() *http.ServeMux {
	router := http.NewServeMux()

	// Rute Auth
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controllers.Register(w, r) // Untuk registrasi staff
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controllers.Login(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Rute dengan AuthMiddleware
	securedRouter := http.NewServeMux()

	// // Rute untuk customer
	// securedRouter.Handle("/customers", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodPost:
	// 		controllers.AddCustomer(w, r) // Membuat customer baru
	// 	case http.MethodGet:
	// 		controllers.GetAllCustomers(w, r) // Mengambil semua data customer
	// 	default:
	// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 	}
	// })))

	// // Rute untuk customer berdasarkan ID
	// securedRouter.Handle("/customer-id", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodGet:
	// 		controllers.GetCustomerByID(w, r) // Mengambil data customer berdasarkan ID
	// 	case http.MethodPut:
	// 		controllers.UpdateCustomer(w, r) // Mengupdate data customer berdasarkan ID
	// 	case http.MethodDelete:
	// 		controllers.DeleteCustomer(w, r) // Menghapus data customer berdasarkan ID
	// 	default:
	// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 	}
	// })))

	// // Rute untuk pembayaran
	// securedRouter.Handle("/create-payment", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodPost:
	// 		controllers.CreatePayment(w, r) // Membuat pembayaran
	// 	default:
	// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 	}
	// })))

	// // Rute untuk karyawan dengan RoleMiddleware khusus admin
	// securedRouter.Handle("/employees", middleware.RoleMiddleware("admin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodPost:
	// 		controllers.CreateEmployee(w, r) // Membuat data karyawan baru
	// 	case http.MethodGet:
	// 		controllers.GetAllEmployees(w, r) // Mengambil semua data karyawan
	// 	default:
	// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 	}
	// })))

	// // Rute untuk karyawan berdasarkan ID (Admin-only)
	// securedRouter.Handle("/employee-id", middleware.RoleMiddleware("admin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodGet:
	// 		controllers.GetEmployeeByID(w, r) // Mengambil data karyawan berdasarkan ID
	// 	case http.MethodPut:
	// 		controllers.UpdateEmployee(w, r) // Mengupdate data karyawan
	// 	case http.MethodDelete:
	// 		controllers.DeleteEmployee(w, r) // Menghapus data karyawan
	// 	default:
	// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 	}
	// })))

	// Gabungkan router utama dengan router yang dilindungi middleware
	router.Handle("/", securedRouter)

	return router
}
