package routes

import (
	"decoration_project/controllers"
	admincontroller "decoration_project/controllers/admin_controller"
	restorantcontrollers "decoration_project/controllers/restorant_controllers"
	staffcontrollers "decoration_project/controllers/staff_controllers"
	usercontroller "decoration_project/controllers/user_controller"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func InitializeRoutes() http.Handler {
	router := mux.NewRouter()

	// Admin Apis.
	router.HandleFunc("/admin/get_category", controllers.GetCategories).Methods("GET")
	router.HandleFunc("/admin/add_category", controllers.CreateCategory).Methods("POST")
	router.HandleFunc("/admin/create-service", admincontroller.AddService).Methods("POST")
	router.HandleFunc("/admin/service", admincontroller.GetServiceDetails).Methods("GET")
	router.HandleFunc("/admin/service-list", admincontroller.GetAllServiceCategoryWise).Methods("GET")
	router.HandleFunc("/admin/get-restorant", admincontroller.GetRestorants).Methods("GET")
	router.HandleFunc("/admin/get-active-bookings", admincontroller.GetAllBookings).Methods("POST")
	router.HandleFunc("/admin/booking-details", admincontroller.GetBookingsDetails).Methods("GET")
	router.HandleFunc("/admin/update-status",admincontroller.UpdateBookingStatus).Methods("POST");

	//User apis.
	router.HandleFunc("/users/create_account", usercontroller.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/users/login", usercontroller.LoginUserHandler).Methods("POST")
	router.HandleFunc("/users/get-category", usercontroller.GetCategoryForUser).Methods("GET")
	router.HandleFunc("/users/get-services", usercontroller.GetServicesBasedOnCategory).Methods("POST")
	router.HandleFunc("/users/book-service", usercontroller.CreateBooking).Methods("POST")
	router.HandleFunc("/users/get-bookings", usercontroller.GetUsersBookings).Methods("GET")
	router.HandleFunc("/users/get-service-details", usercontroller.GetServiceDetails).Methods("POST")
	router.HandleFunc("/users/get-bookings/{id}", usercontroller.GetBookingDetails).Methods("GET")
	router.HandleFunc("/users/booking/{id}/location", usercontroller.GetPartnerLiveLocation).Methods("GET")

	
	// Restorant APIS.

	router.HandleFunc("/restorant/create_account", restorantcontrollers.RegisterRestaurant).Methods("POST")
	router.HandleFunc("/restorant/login", restorantcontrollers.LoginRestaurant).Methods("POST")
	// router.HandleFunc("/restorant/get-category", restorantcontrollers.GetCategoryRestorant).Methods("GET")
	// router.HandleFunc("/restorant/create-services", restorantcontrollers.AddService).Methods("POST")
	// router.HandleFunc("/restorant/get-service", restorantcontrollers.GetServiceDetails).Methods("GET")
	// router.HandleFunc("/restorant/get-all-services", restorantcontrollers.GetAllServicesForRestaurant).Methods("GET")
	router.HandleFunc("/restorant/add-staff", restorantcontrollers.AddStaff).Methods("POST")
	router.HandleFunc("/restorant/get-all-staff", restorantcontrollers.GetAllStaff).Methods("GET")
	router.HandleFunc("/restorant/get-all-bookings", restorantcontrollers.GetAllBookedServices).Methods("GET")
	router.HandleFunc("/restorant/update-booking", restorantcontrollers.HandleBookingAction).Methods("POST")
	router.HandleFunc("/restorant/assign-staff-booking", restorantcontrollers.AssignStaffToBooking).Methods("POST")
	router.HandleFunc("/restorant/booking-details/{id}", restorantcontrollers.GetBookingDetails).Methods("GET")

	// Staff Apis.
	router.HandleFunc("/staff/login", staffcontrollers.LoginStaffHandler).Methods("POST")
	router.HandleFunc("/staff/get-bookings", staffcontrollers.GetAllAssignedBookings).Methods("GET")
	router.HandleFunc("/staff/get-bookings-details/{id}", staffcontrollers.GetAssignedServicesDetails).Methods("GET")
	router.HandleFunc("/staff/start-service", staffcontrollers.StartService).Methods("POST")
	router.HandleFunc("/staff/update-location", staffcontrollers.UpdateStaffLocation).Methods("POST")
	router.HandleFunc("/staff/complete-service", staffcontrollers.VerifyCompletionOTP).Methods("POST")

	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Or restrict to frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	return c.Handler(router)
}
