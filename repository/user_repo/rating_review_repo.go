package userrepo

import (
	"decoration_project/config"
	usermodels "decoration_project/models/user_models"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

func AddServiceReview(review usermodels.ServiceReview) error {
    query := `
        INSERT INTO Service_Reviews 
            (review_id, service_id, user_id, user_name, rating, review_text)
        VALUES 
            (UUID(), ?, ?, ?, ?, ?)
    `

    _, err := config.DB.Exec(
        query,
        review.ServiceID,
        review.UserID,
        review.UserName,
        review.Rating,
        review.ReviewText,
    )

    if err != nil {
        // Check if duplicate entry error (MySQL 1062)
        if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
            return fmt.Errorf("you have already added a review for this service")
        }
        return fmt.Errorf("failed to insert service review: %w", err)
    }

    // Update average rating (optional)
    avgQuery := `
        UPDATE Our_Services 
        SET average_rating = (
            SELECT ROUND(AVG(rating),2) FROM Service_Reviews WHERE service_id = ?
        )
        WHERE service_id = ?
    `
    _, err = config.DB.Exec(avgQuery, review.ServiceID, review.ServiceID)
    if err != nil {
        return fmt.Errorf("failed to update average rating: %w", err)
    }

    return nil
}
