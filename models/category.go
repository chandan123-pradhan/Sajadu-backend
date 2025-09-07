// models/product_category.go
package models

type ProductCategory struct {
    CategoryID   string `json:"category_id"`   // use string, not int64
    CategoryName string `json:"category_name"`
    ImageURL     string `json:"image_url"`
}
