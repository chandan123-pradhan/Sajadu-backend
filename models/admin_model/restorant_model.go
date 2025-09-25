package adminmodel

type RestorantModel struct {
    RestorantId   string   `json:"restorant_id"`
    RestorantName string   `json:"restorant_name"`
    Images        []string `json:"images"`
}
