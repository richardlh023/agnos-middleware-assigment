package api

import (
	"agnos-middleware/internal/models"
	"agnos-middleware/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PatientController struct {
	patientService *services.PatientService
}

func NewPatientController(patientService *services.PatientService) *PatientController {
	return &PatientController{
		patientService: patientService,
	}
}

// @Summary      Search for patients
// @Description  Search for patients by optional criteria. Requires JWT authentication. Staff can only access patients from their own hospital.
// @Tags         Patient
// @Accept       json
// @Produce      json
// @Param        id query string false "Patient ID (must be national_id or passport_id). Examples: Hospital A - 1234567890123, 9876543210987, AB1234567; Hospital B - 1111222233334, 4455667788990" default(1234567890123)
// @Param        patient_hn query string false "Hospital Number"
// @Param        national_id query string false "National ID"
// @Param        passport_id query string false "Passport ID"
// @Param        first_name query string false "First name (partial match)"
// @Param        middle_name query string false "Middle name (partial match)"
// @Param        last_name query string false "Last name (partial match)"
// @Param        date_of_birth query string false "Date of birth (YYYY-MM-DD)"
// @Param        phone_number query string false "Phone number"
// @Param        email query string false "Email"
// @Param        gender query string false "Gender (M/F)"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "Patients found"
// @Failure      400  {object}  utils.PatientSearchErrorResponse  "Bad request - at least one search criteria must be provided"
// @Failure      401  {object}  utils.AuthErrorResponse  "Unauthorized - authorization header required or invalid token"
// @Failure      403  {object}  utils.AccessDeniedErrorResponse  "Access denied - patient does not belong to your hospital"
// @Failure      404  {object}  utils.NotFoundErrorResponse  "Patient not found"
// @Router       /patient/search [get]
func (ctrl *PatientController) SearchPatient(ctx *gin.Context) {
	var req models.PatientSearchRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ID == nil && req.PatientHN == nil && req.NationalID == nil && req.PassportID == nil &&
		req.FirstName == nil && req.MiddleName == nil && req.LastName == nil &&
		req.DateOfBirth == nil && req.PhoneNumber == nil && req.Email == nil && req.Gender == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "at least one search criteria must be provided"})
		return
	}

	staffHospital, exists := ctx.Get("staff_hospital")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "staff information not found"})
		return
	}

	patients, err := ctrl.patientService.SearchPatient(&req, staffHospital.(string))
	if err != nil {
		if err.Error() == "access denied: patient does not belong to your hospital" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(patients) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"patients": patients,
		"count":    len(patients),
	})
}
