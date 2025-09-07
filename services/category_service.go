package services

import (
	"decoration_project/models"
	"decoration_project/repository"
)

func GetCategories() ([]models.ProductCategory, error) {
	return repository.GetAllCategories()
}

func AddCategory(category models.ProductCategory) (string, error) {
    return repository.AddCategory(category)
}