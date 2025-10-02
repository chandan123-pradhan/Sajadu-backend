package usercontroller

import (
	"decoration_project/models"
	restorantmodels "decoration_project/models/restorant_models"
	usermodels "decoration_project/models/user_models"
	"decoration_project/repository"
	userservices "decoration_project/services/user_services"
	"decoration_project/utils"
	"encoding/json"
	"net/http"
)

func GetCategoryForUser(w http.ResponseWriter, r *http.Request) {
	// Validate JWT token
	_, err := utils.ValidateToken(r)
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

func GetServicesBasedOnCategory(w http.ResponseWriter, r *http.Request) {
	// Validate JWT token
	_, err := utils.ValidateToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, "Unauthorized: "+err.Error())
		return
	}

	err = json.NewDecoder(r.Body).Decode(&usermodels.GetServicesBasedOnCategory)

	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Invalid request body")
		return
	}

	if usermodels.GetServicesBasedOnCategory.CategoryId == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Category ID Required")
		return
	}
	// Fetch services from service layer
	servicesList, err := userservices.GetServicesByCategory(usermodels.GetServicesBasedOnCategory.CategoryId)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to fetch services: "+err.Error())
		return
	}
	// Ensure empty slice instead of null
	if servicesList == nil {
		servicesList = []restorantmodels.RestaurantService{}
	}
	// Send successful response
	utils.SendResponse(w, http.StatusOK, true, map[string]interface{}{
		"services": servicesList,
	}, "Services fetched successfully")

}




func GetServiceDetails(w http.ResponseWriter, r *http.Request) {
	// Validate JWT token
	_, err := utils.ValidateToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, "Unauthorized: "+err.Error())
		return
	}

	// Decode request body
	err = json.NewDecoder(r.Body).Decode(&usermodels.GetServiceDetailsRequest)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Invalid request body")
		return
	}

	// Validate required fields
	if usermodels.GetServiceDetailsRequest.ServiceId == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Restaurant ID and Service ID are required")
		return
	}

	// Call service layer
	serviceWithRest, err := userservices.GetServiceDetails(
		usermodels.GetServiceDetailsRequest.ServiceId,
	)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to fetch details: "+err.Error())
		return
	}

	// If not found
	if serviceWithRest.Service.ServiceID == "" {
		utils.SendResponse(w, http.StatusNotFound, false, map[string]interface{}{}, "No service found for given Restaurant ID and Service ID")
		return
	}

	// Success response
	utils.SendResponse(w, http.StatusOK, true, map[string]interface{}{
		"service":    serviceWithRest.Service,
	}, "Service details fetched successfully")
}


func SearchServices(w http.ResponseWriter, r *http.Request) {
	// Validate JWT token
	_, err := utils.ValidateToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, "Unauthorized: "+err.Error())
		return
	}

	// Decode request body
	err = json.NewDecoder(r.Body).Decode(&usermodels.SearchServicesRequest)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Invalid request body")
		return
	}

	// Validate input
	if usermodels.SearchServicesRequest.Query == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Search query is required")
		return
	}

	// Call service layer
	services, err := userservices.SearchServices(usermodels.SearchServicesRequest.Query)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to search services: "+err.Error())
		return
	}

	// Ensure not nil
	if services == nil {
		services = []restorantmodels.RestaurantService{}
	}

	// Success response
	utils.SendResponse(w, http.StatusOK, true, map[string]interface{}{
		"services": services,
	}, "Search results fetched successfully")
}