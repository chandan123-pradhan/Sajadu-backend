package staffservices

import (
	staffrepo "decoration_project/repository/staff_repo"
	"decoration_project/utils"
	"errors"
)


func LoginStaff(email, password string) (map[string]interface{}, error) {
    // Fetch staff by email
    staff, err := staffrepo.GetStaffByEmail(email)
    if err != nil || staff.StaffID == "" {
        return nil, errors.New("invalid email or password") // staff not found
    }

    // Compare plain text password
    if staff.Password != password {
        return nil, errors.New("invalid email or password")
    }

    // Generate auth token
    token, err := utils.GenerateToken(staff.StaffID)
    if err != nil {
        return nil, err
    }

    // Prepare response (exclude password)
    resp := map[string]interface{}{
        "staff_id":     staff.StaffID,
        "restaurant_id": staff.RestaurantID,
        "name":         staff.Name,
        "email":        staff.Email,
        "whatsapp_no":  staff.WhatsappNo,
        "designation":  staff.Designation,
        "description":  staff.Description,
        "image_url":    staff.ImageURL,
        "auth_token":   token,
    }

    return resp, nil
}
