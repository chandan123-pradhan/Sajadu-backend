package repository

import (
	"decoration_project/config"
	"decoration_project/models"
	"github.com/google/uuid"
)

func GetAllCategories() ([]models.ProductCategory, error) {
	query := "SELECT category_id, category_name, image_url FROM Product_Category"

	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.ProductCategory

	for rows.Next() {
		var category models.ProductCategory
		if err := rows.Scan(&category.CategoryID, &category.CategoryName, &category.ImageURL); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}


func AddCategory(category models.ProductCategory) (string, error) {
	// Generate UUID
	newID := uuid.New().String()

	query := "INSERT INTO Product_Category (category_id, category_name, image_url) VALUES (?, ?, ?)"
	_, err := config.DB.Exec(query, newID, category.CategoryName, category.ImageURL)
	if err != nil {
		return "", err
	}
	return newID, nil
}
