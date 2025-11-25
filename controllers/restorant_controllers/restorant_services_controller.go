package restorantcontrollers

import (
	"decoration_project/models"
	restorantmodels "decoration_project/models/restorant_models"
	"decoration_project/repository"
	adminserices "decoration_project/services/admin_serices"
	restorantservices "decoration_project/services/restorant_services"
	"decoration_project/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// ======================= ADD SERVICE =======================
func AddServiceByRestorant(w http.ResponseWriter, r *http.Request) {
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

	// Validate required fields
	if service.ServiceName == "" || service.ServiceDescription == "" || service.ServicePrice == 0 || service.CategoryId == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Service name, description, Category ID,  Proposed Restorant ID and price are required")
		return
	}

	service.ProposedRestorantId = restaurantID

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
	serviceID, err := adminserices.CreateService(service, imagePaths)
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

func GetRestorantProposedServices(w http.ResponseWriter, r *http.Request) {
	// Validate restaurant token
	restaurantID, err := utils.ValidateRestaurantToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, err.Error())
		return
	}

	// Fetch all services for this restaurant
	services, err := restorantservices.GetAllProposedServices(restaurantID)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, err.Error())
		return
	}

	// Prepare response
	responseData := make([]map[string]interface{}, 0)
	for _, service := range services {
		responseData = append(responseData, map[string]interface{}{
			"service_id":            service.ServiceID,
			"category_id":           service.CategoryId,
			"service_name":          service.ServiceName,
			"service_description":   service.ServiceDescription,
			"service_price":         service.ServicePrice,
			"images":                service.Images, // list of image URLs
			"created_at":            service.CreatedAt,
			"updated_at":            service.UpdatedAt,
			"proposed_restorant_id": service.ProposedRestorantId,
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



func DeleteProposedService(w http.ResponseWriter, r *http.Request) {
	// Validate restaurant token
	restaurantID, err := utils.ValidateRestaurantToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, err.Error())
		return
	}

	// Read service ID from query params
	serviceID := r.URL.Query().Get("service_id")
	if serviceID == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "service_id is required")
		return
	}

	// Call service layer
	err = restorantservices.DeleteProposedService(serviceID, restaurantID)
	if err != nil {
		utils.SendResponse(w, http.StatusForbidden, false, nil, err.Error())
		return
	}

	utils.SendResponse(w, http.StatusOK, true, nil, "Service deleted successfully")
}
