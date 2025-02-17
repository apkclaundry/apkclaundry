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
			controllers.GetAllUsers(w, r) // Mengambil semua data users
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	securedRouter.Handle("/employeename", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetAllEmployeesIDName(w, r) // Mengambil semua data users
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

	securedRouter.Handle("/customers-name", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetAllCustomersIDName(w, r) // Mengambil data user berdasarkan ID
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	securedRouter.Handle("/name-id", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetCustomerNameByID(w, r) // Mengambil nama customer berdasarkan ID
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	securedRouter.Handle("/supplier/transaction", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controllers.AddSupplierTransaction(w, r) // Mengambil nama customer berdasarkan ID
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

	// // Rute untuk Stock berdasarkan ID
	securedRouter.Handle("/item-name", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetItemTransactions(w, r) // Mengambil data user berdasarkan ID
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	// // Rute untuk transaksi
	securedRouter.Handle("/transaction", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetAllTransactions(w, r) // Mengambil semua data customer
		case http.MethodPost:
			controllers.CreateTransaction(w, r) // Membuat data karyawan baru
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	// // Rute untuk treansaksi berdasarkan ID
	securedRouter.Handle("/transaction-id", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetTransactionByID(w, r) // Mengambil data user berdasarkan ID
		case http.MethodPut:
			controllers.UpdateTransaction(w, r) // Mengupdate data user berdasarkan ID
		case http.MethodDelete:
			controllers.DeleteTransaction(w, r) // Menghapus data user berdasarkan ID
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Rute untuk transaksi item
	securedRouter.Handle("/item-transaction", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetAllItemTransactions(w, r) // Mengambil semua transaksi item
		case http.MethodPost:
			controllers.CreateItemTransaction(w, r) // Membuat transaksi item baru
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Rute untuk transaksi item berdasarkan ID
	securedRouter.Handle("/item-transaction-id", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetItemTransactionByID(w, r) // Mengambil transaksi item berdasarkan ID
		case http.MethodPut:
			controllers.UpdateItemTransaction(w, r) // Mengupdate transaksi item berdasarkan ID
		case http.MethodDelete:
			controllers.DeleteItemTransaction(w, r) // Menghapus transaksi item berdasarkan ID
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
