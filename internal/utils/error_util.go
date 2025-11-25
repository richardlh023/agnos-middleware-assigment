package utils

type ErrorResponse struct {
	Error string `json:"error"`
}

type LoginErrorResponse struct {
	Error string `json:"error" example:"invalid credentials"`
}

type CreateStaffErrorResponse struct {
	Error string `json:"error" example:"username already exists"`
}

type PatientSearchErrorResponse struct {
	Error string `json:"error" example:"at least one search criteria must be provided"`
}

type AuthErrorResponse struct {
	Error string `json:"error" example:"authorization header required"`
}

type AccessDeniedErrorResponse struct {
	Error string `json:"error" example:"access denied: patient does not belong to your hospital"`
}

type NotFoundErrorResponse struct {
	Error string `json:"error" example:"patient not found"`
}
