package routes

import (
	"apkclaundry/controllers"
	"apkclaundry/middleware"
	"net/http"
)

func InitRoutes() *http.ServeMux {
	router := http.NewServeMux()

	// // Rute Auth
	// router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodPost:
	// 		controllers.Register(w, r) // Untuk registrasi staff
	// 	default:
	// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 	}
	// })

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

	// Rute Register
	securedRouter.Handle("/Register", middleware.AuthMiddleware(middleware.RoleMiddleware("admin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controllers.Register(w, r) // Membuat data karyawan baru
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))))

	// // Rute untuk employee
	securedRouter.Handle("/employee", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetAllUsers(w, r) // Mengambil semua data customer
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	// // Rute untuk Users berdasarkan ID
	securedRouter.Handle("/employee-id", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetUserByID(w, r) // Mengambil data user berdasarkan ID
		case http.MethodPut:
			controllers.UpdateUser(w, r) // Mengupdate data user berdasarkan ID
		case http.MethodDelete:
			controllers.DeleteUser(w, r) // Menghapus data user berdasarkan ID
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	// // Rute untuk customer
	securedRouter.Handle("/customer", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetAllCustomers(w, r) // Mengambil semua data customer
		case http.MethodPost:
			controllers.CreateCustomer(w, r) // Membuat data karyawan baru
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	// // Rute untuk Users berdasarkan ID
	securedRouter.Handle("/customer-id", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetCustomerByID(w, r) // Mengambil data user berdasarkan ID
		case http.MethodPut:
			controllers.UpdateCustomer(w, r) // Mengupdate data user berdasarkan ID
		case http.MethodDelete:
			controllers.DeleteCustomer(w, r) // Menghapus data user berdasarkan ID
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	// // Rute untuk supplier
	securedRouter.Handle("/supplier", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetAllSuppliers(w, r) // Mengambil semua data customer
		case http.MethodPost:
			controllers.CreateSupplier(w, r) // Membuat data karyawan baru
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	// // Rute untuk supplier berdasarkan ID
	securedRouter.Handle("/supplier-id", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetSupplierByID(w, r) // Mengambil data user berdasarkan ID
		case http.MethodPut:
			controllers.UpdateSupplier(w, r) // Mengupdate data user berdasarkan ID
		case http.MethodDelete:
			controllers.DeleteSupplier(w, r) // Menghapus data user berdasarkan ID
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))


		// // Rute untuk Stock
		securedRouter.Handle("/stock", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				controllers.GetAllItems(w, r) // Mengambil semua data customer
			case http.MethodPost:
				controllers.CreateItem(w, r) // Membuat data karyawan baru
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})))
	
		// // Rute untuk Stock berdasarkan ID
		securedRouter.Handle("/stock-id", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				controllers.GetItemByID(w, r) // Mengambil data user berdasarkan ID
			case http.MethodPut:
				controllers.UpdateItem(w, r) // Mengupdate data user berdasarkan ID
			case http.MethodDelete:
				controllers.DeleteItem(w, r) // Menghapus data user berdasarkan ID
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})))

	// // Rute untuk pembayaran
	// securedRouter.Handle("/create-payment", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodPost:
	// 		controllers.CreatePayment(w, r) // Membuat pembayaran
	// 	default:
	// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 	}
	// })))

	// // // Rute untuk karyawan dengan RoleMiddleware khusus admin
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
