package adminserices

import adminrepo "decoration_project/repository/admin_repo"

func SaveOrUpdateAdminFCMToken(token string) error {
	return adminrepo.SaveOrUpdateAdminFCM(token)
}

// GetAdminFCMToken fetches the current admin FCM token
func GetAdminFCMToken() (string, error) {
	return adminrepo.GetAdminFCMToken()
}


