package api

import (
	"agnos-middleware/internal/models"
	"agnos-middleware/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StaffController struct {
	authService *services.AuthService
}

func NewStaffController(authService *services.AuthService) *StaffController {
	return &StaffController{
		authService: authService,
	}
}

// @Summary      Create a new staff member
// @Description  Create a new hospital staff account with employee details. All fields will be pre-filled with example values in Swagger UI.
// @Tags         Staff
// @Accept       json
// @Produce      json
// @Param        request body models.CreateStaffRequest true "Staff creation request"
// @Success      201  {object}  map[string]interface{}  "Staff created successfully"
// @Failure      400  {object}  utils.CreateStaffErrorResponse  "Bad request - validation error"
// @Failure      409  {object}  utils.CreateStaffErrorResponse  "Conflict - username/email/employee_id already exists"
// @Failure      500  {object}  utils.CreateStaffErrorResponse  "Internal server error"
// @Router       /staff/create [post]
func (ctrl *StaffController) CreateStaff(ctx *gin.Context) {
	var req models.CreateStaffRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	staff, err := ctrl.authService.CreateStaff(&req)
	if err != nil {
		if err.Error() == "username already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create staff"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":          staff.ID,
		"employee_id": staff.EmployeeID,
		"username":    staff.Username,
		"first_name":  staff.FirstName,
		"last_name":   staff.LastName,
		"email":       staff.Email,
		"role":        staff.Role,
		"department":  staff.Department,
		"hospital":    staff.Hospital,
	})
}

// @Summary      Staff login
// @Description  Authenticate staff member and receive JWT token
// @Tags         Staff
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "Login credentials"
// @Success      200  {object}  models.LoginResponse  "Login successful"
// @Failure      400  {object}  utils.LoginErrorResponse  "Bad request - validation error"
// @Failure      401  {object}  utils.LoginErrorResponse  "Unauthorized - invalid credentials"
// @Router       /staff/login [post]
func (ctrl *StaffController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.authService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
