package controllers

import (
	"decoration_project/models"
	"decoration_project/repository"
	"decoration_project/services"
	"decoration_project/utils"
	"fmt"
	"strings"

	"net/http"

	"github.com/go-sql-driver/mysql"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := repository.GetAllCategories()
	if err != nil {
		fmt.Print(err)
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{
			"category": []interface{}{}, // always array, even on error
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



func CreateCategory(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 5MB request size)
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Invalid form data.")
		return
	}

	// Validate category name
	categoryName := r.FormValue("category_name")
	if categoryName == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Category name must be present.")
		return
	}

	// Validate and save image
	file, handler, err := r.FormFile("image")
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Image is required.")
		return
	}
	defer file.Close()

	// Validate size (max 2MB)
	if err := utils.ValidateImageSize(file, 2<<20); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, err.Error())
		return
	}

	// Save image (helper returns URL or error)
	imageURL, err := utils.SaveUploadedFile(file, handler, "uploads")
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to save image.")
		return
	}

	// Prepare category model
	newCategory := models.ProductCategory{
		CategoryName: categoryName,
		ImageURL:     imageURL,
	}

	// Insert into DB (via service)
	categoryID, err := services.AddCategory(newCategory)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {

			parts := strings.Split(mysqlErr.Message, "'")
			duplicateValue := ""
			if len(parts) >= 2 {
				duplicateValue = parts[1]
			}

			utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{},
				"Category '"+duplicateValue+"' already exists.")
			return
		}

		// Other errors
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, err.Error())
		return
	}

	// Send response
	data := map[string]interface{}{
		"category_id": categoryID,
	}
	utils.SendResponse(w, http.StatusCreated, true, data, "Category added successfully.")
}

