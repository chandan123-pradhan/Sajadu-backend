package restorantcontrollers

import (
	"decoration_project/models"
	restorantmodels "decoration_project/models/restorant_models"
	"decoration_project/repository"
	restorantservices "decoration_project/services/restorant_services"
	"decoration_project/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// ======================= ADD SERVICE =======================
func AddService(w http.ResponseWriter, r *http.Request) {
	// Validate restaurant token
	restaurantID, err := utils.ValidateRestaurantToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, err.Error())
		return
	}

	// Parse multipart form data (for images + JSON fields)
	err = r.ParseMultipartForm(10 << 20) // 10MB max upload
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Failed to parse form data")
		return
	}

	// Extract service details from "service" field (JSON string)
	var service restorantmodels.RestaurantService
	serviceData := r.FormValue("service")
	if serviceData == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Service JSON is required")
		return
	}

	err = json.Unmarshal([]byte(serviceData), &service)
	if err != nil {
		fmt.Println(err.Error())
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Invalid service JSON")
		return
	}

	// Assign restaurant ID from token
	service.RestaurantID = restaurantID

	// Validate required fields
	if service.ServiceName == "" || service.ServiceDescription == "" || service.ServicePrice == 0 || service.CategoryId == ""{
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Service name, description, Category ID, and price are required")
		return
	}

	// Validate items JSON
	var items map[string]interface{}
	if err := json.Unmarshal([]byte(service.Items), &items); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Items must be a valid JSON string")
		return
	}

	// Handle multiple image files
	var imagePaths []string
	files := r.MultipartForm.File["images"]
	for _, fileHeader := range files {
		path, err := utils.SaveFile(fileHeader, "uploads/services")
		if err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to save image")
			return
		}
		imagePaths = append(imagePaths, path)
	}

	// Save service
	serviceID, err := restorantservices.CreateService(service, imagePaths)
	if err != nil {
		fmt.Println(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, err.Error())
		return
	}

	// Prepare response
	responseData := map[string]interface{}{
		"service_id": serviceID,
		"images":     imagePaths,
	}

	utils.SendResponse(w, http.StatusCreated, true, responseData, "Service created successfully")
}

// ======================= GET SERVICE DETAILS =======================
func GetServiceDetails(w http.ResponseWriter, r *http.Request) {
	// Validate restaurant token
	_, err := utils.ValidateRestaurantToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, err.Error())
		return
	}

	// Get serviceID from query parameter
	serviceID := r.URL.Query().Get("service_id")
	if serviceID == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "service_id query parameter is required")
		return
	}

	// Fetch service details from service layer
	service, err := restorantservices.GetService(serviceID)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, err.Error())
		return
	}

	// Prepare response
	responseData := map[string]interface{}{
		"service_id":          service.ServiceID,
		"restaurant_id":       service.RestaurantID,
		"category_id":			service.CategoryId,
		"service_name":        service.ServiceName,
		"service_description": service.ServiceDescription,
		"service_price":       service.ServicePrice,
		"items":               service.Items,  // JSON string
		"images":              service.Images, // list of image URLs
		"created_at":          service.CreatedAt,
		"updated_at":          service.UpdatedAt,
	}

	utils.SendResponse(w, http.StatusOK, true, responseData, "Service details fetched successfully")
}



func GetAllServicesForRestaurant(w http.ResponseWriter, r *http.Request) {
    // Validate restaurant token and get restaurant ID
    restaurantID, err := utils.ValidateRestaurantToken(r)
    if err != nil {
        utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, "Unauthorized: "+err.Error())
        return
    }

    // Fetch all services for this restaurant
    services, err := restorantservices.GetAllServicesForRestaurant(restaurantID)
    if err != nil {
        utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, err.Error())
        return
    }

    // Prepare response
    responseData := make([]map[string]interface{}, 0) 
    for _, service := range services {
        responseData = append(responseData, map[string]interface{}{
            "service_id":          service.ServiceID,
            "restaurant_id":       service.RestaurantID,
			"category_id":			service.CategoryId,
            "service_name":        service.ServiceName,
            "service_description": service.ServiceDescription,
            "service_price":       service.ServicePrice,
            "items":               service.Items,  // JSON string
            "images":              service.Images, // list of image URLs
            "created_at":          service.CreatedAt,
            "updated_at":          service.UpdatedAt,
        })
    }

    utils.SendResponse(w, http.StatusOK, true, responseData, "All services fetched successfully")
}


func GetCategoryRestorant(w http.ResponseWriter, r *http.Request) {
    // Validate JWT token
    _, err := utils.ValidateRestaurantToken(r)
    if err != nil {
        utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, "Unauthorized: "+err.Error())
        return
    }

    // Fetch categories from DB
    categories, err := repository.GetAllCategories()
    if err != nil {
        utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{
            "category": []interface{}{},
        }, "Failed to fetch categories")
        return
    }

    if categories == nil {
        categories = []models.ProductCategory{}
    }

    data := map[string]interface{}{
        "category": categories,
    }

    utils.SendResponse(w, http.StatusOK, true, data, "Categories fetched successfully")
}
