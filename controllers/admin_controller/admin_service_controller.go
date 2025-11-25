package admincontroller

import (
	adminmodel "decoration_project/models/admin_model"
	restorantmodels "decoration_project/models/restorant_models"
	adminserices "decoration_project/services/admin_serices"
	"decoration_project/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)


// ======================================================
//                     ADD SERVICE
// ======================================================

// AddService handles adding a new restaurant service.
// Expects multipart form data:
//   - "service"  → JSON string containing service details
//   - "images[]" → One or more image files
//
// Validations:
//   - ServiceName, Description, CategoryID, Price must be provided
//
// On success:
//   - Returns newly created service_id and the uploaded image paths
func AddService(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20) // Max upload size: 10MB
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Failed to parse form data")
		return
	}

	// Extract JSON ("service") from form-data
	var service restorantmodels.RestaurantService
	serviceData := r.FormValue("service")
	if serviceData == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Service JSON is required")
		return
	}

	// Decode service JSON
	err = json.Unmarshal([]byte(serviceData), &service)
	if err != nil {
		fmt.Println(err.Error())
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Invalid service JSON")
		return
	}

	// Required fields validation
	if service.ServiceName == "" || service.ServiceDescription == "" ||
		service.ServicePrice == 0 || service.CategoryId == "" {

		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Service name, description, Category ID, and price are required")
		return
	}

	// Upload images
	var imagePaths []string
	files := r.MultipartForm.File["images"]
	for _, f := range files {
		path, err := utils.SaveFile(f, "uploads/services")
		if err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to save image")
			return
		}
		imagePaths = append(imagePaths, path)
	}

	// Save service in DB
	serviceID, err := adminserices.CreateService(service, imagePaths)
	if err != nil {
		fmt.Println(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, err.Error())
		return
	}

	// Response
	response := map[string]interface{}{
		"service_id": serviceID,
		"images":     imagePaths,
	}

	utils.SendResponse(w, http.StatusCreated, true, response, "Service created successfully")
}



// ======================================================
//                 GET SERVICE DETAILS
// ======================================================

// GetServiceDetails returns full details of a specific service.
// Requires:
//   - Query param: service_id
//
// Returns:
//   - Basic details + images + timestamps
func GetServiceDetails(w http.ResponseWriter, r *http.Request) {

	serviceID := r.URL.Query().Get("service_id")
	if serviceID == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "service_id query parameter is required")
		return
	}

	service, err := adminserices.GetServicesDetails(serviceID)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, err.Error())
		return
	}

	response := map[string]interface{}{
		"service_id":           service.ServiceID,
		"category_id":          service.CategoryId,
		"service_name":         service.ServiceName,
		"service_description":  service.ServiceDescription,
		"service_price":        service.ServicePrice,
		"images":               service.Images,
		"created_at":           service.CreatedAt,
		"updated_at":           service.UpdatedAt,
		"proposed_restorant_id": service.ProposedRestorantId,
	}

	utils.SendResponse(w, http.StatusOK, true, response, "Service details fetched successfully")
}



// ======================================================
//           GET ALL SERVICES BY CATEGORY
// ======================================================

// GetAllServiceCategoryWise fetches all services under a category.
// Requires:
//   - Query param: category-id
//
// Returns:
//   - A list of services with images
func GetAllServiceCategoryWise(w http.ResponseWriter, r *http.Request) {

	categoryId := r.URL.Query().Get("category-id")
	if categoryId == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "category_id query parameter is required")
		return
	}

	services, err := adminserices.GetAllServiceCategoryWise(categoryId)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, err.Error())
		return
	}

	// Prepare output list
	response := []map[string]interface{}{}
	for _, s := range services {
		response = append(response, map[string]interface{}{
			"service_id":           s.ServiceID,
			"category_id":          s.CategoryId,
			"service_name":         s.ServiceName,
			"service_description":  s.ServiceDescription,
			"service_price":        s.ServicePrice,
			"images":               s.Images,
			"created_at":           s.CreatedAt,
			"updated_at":           s.UpdatedAt,
			"proposed_restorant_id": s.ProposedRestorantId,
		})
	}

	utils.SendResponse(w, http.StatusOK, true, response, "All services fetched successfully")
}



// ======================================================
//                    ADD FESTIVAL
// ======================================================

// AddFestival creates a new festival.
// Required fields:
//   - festival_name
//   - start_date
//   - end_date
//
// Handles:
//   - Duplicate festival_name (returns 409 Conflict)
func AddFestival(w http.ResponseWriter, r *http.Request) {

	var f adminmodel.Festival

	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid JSON")
		return
	}

	if f.FestivalName == "" || f.StartDate == "" || f.EndDate == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "All fields are required")
		return
	}

	err = adminserices.CreateFestival(f)
	if err != nil {

		if err.Error() == "Festival name already exists" {
			utils.SendResponse(w, http.StatusConflict, false, nil, "Festival already exists. Please choose a different name.")
			return
		}

		utils.SendResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	utils.SendResponse(w, http.StatusCreated, true, []interface{}{}, "Festival created successfully")
}



// ======================================================
//                  GET ALL FESTIVALS
// ======================================================

// GetAllFestivals returns all festivals.
// If no festival exists → returns empty list ([]) instead of null.
func GetAllFestivals(w http.ResponseWriter, r *http.Request) {

	festivals, err := adminserices.GetAllFestivals()
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, []interface{}{}, err.Error())
		return
	}

	if len(festivals) == 0 {
		utils.SendResponse(w, http.StatusOK, true, []interface{}{}, "Festivals fetched successfully")
		return
	}

	utils.SendResponse(w, http.StatusOK, true, festivals, "Festivals fetched successfully")
}



// ======================================================
//             ADD SERVICE TO FESTIVAL
// ======================================================

// AddServiceToFestival links a service to a festival.
// Required fields:
//   - festival_id
//   - service_id
// Optional:
//   - discount_percent (default 0)
//
// Handles:
//   - Duplicate festival-service mapping (returns 409)
func AddServiceToFestival(w http.ResponseWriter, r *http.Request) {

	var req adminmodel.AddFestivalServiceRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, []interface{}{}, "Invalid JSON")
		return
	}

	if req.FestivalID == "" || req.ServiceID == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, []interface{}{}, "festival_id and service_id are required")
		return
	}

	if req.DiscountPercent == 0 {
		req.DiscountPercent = 0
	}

	err = adminserices.AddServiceToFestival(req.FestivalID, req.ServiceID, req.DiscountPercent)
	if err != nil {

		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			utils.SendResponse(w, http.StatusConflict, false, []interface{}{}, "This service is already added to this festival")
			return
		}

		utils.SendResponse(w, http.StatusInternalServerError, false, []interface{}{}, err.Error())
		return
	}

	utils.SendResponse(w, http.StatusCreated, true, []interface{}{}, "Service added to festival successfully")
}



// ======================================================
//              GET FESTIVAL SERVICE LIST
// ======================================================

// GetAllFestivalServices fetches all services added to a festival.
// Requires:
//   - Query param: festival_id
//
// Returns:
//   - List of services linked with the festival
func GetAllFestivalServices(w http.ResponseWriter, r *http.Request) {

	festivalId := r.URL.Query().Get("festival_id")
	if festivalId == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "festival_id query parameter is required")
		return
	}

	services, err := adminserices.GetFestivalServices(festivalId)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	response := []map[string]interface{}{}
	for _, s := range services {
		response = append(response, map[string]interface{}{
			"service_id":           s.ServiceID,
			"category_id":          s.CategoryId,
			"service_name":         s.ServiceName,
			"service_description":  s.ServiceDescription,
			"service_price":        s.ServicePrice,
			"images":               s.Images,
			"created_at":           s.CreatedAt,
			"updated_at":           s.UpdatedAt,
			"proposed_restorant_id": s.ProposedRestorantId,
			"Discount": s.DiscountPercent,
		})
	}

	utils.SendResponse(w, http.StatusOK, true, response, "Festival services fetched successfully")
}

