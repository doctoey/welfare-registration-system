package controllers

import (
	"net/http"
	"strconv"
	"welfare-registration-backend/internal/entities"

	"github.com/gin-gonic/gin"
)

type ApplicationController struct {
	usecase entities.ApplicationUsecase
}

func NewApplicationController(usecase entities.ApplicationUsecase) *ApplicationController {
	return &ApplicationController{usecase: usecase}
}

func (ctrl *ApplicationController) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.POST("/applications", ctrl.CreateApplication)
		api.GET("/applications/status/:citizenId", ctrl.GetStatus)
		api.GET("/officer/applications", ctrl.GetOfficerApplications)
		api.PATCH("/officer/applications/:id/status", ctrl.UpdateStatus)
	}
}

func (ctrl *ApplicationController) CreateApplication(c *gin.Context) {
	var dto entities.CreateApplicationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	refNumber, err := ctrl.usecase.Register(dto)
	if err != nil {
		switch err.Error() {
		case "CITIZEN_ALREADY_REGISTERED":
			c.JSON(http.StatusConflict, gin.H{"error": "CITIZEN_ALREADY_REGISTERED"})
		case "INVALID_CITIZEN_ID_CHECKSUM":
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid citizen ID checksum"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"referenceNumber": refNumber,
		"citizenId":       dto.CitizenID,
		"birthDate":       dto.BirthDate,
	})
}

func (ctrl *ApplicationController) GetStatus(c *gin.Context) {
	citizenID := c.Param("citizenId")
	birthDate := c.Query("birthDate")
	refNumber := c.Query("ref")

	res, err := ctrl.usecase.GetStatus(citizenID, birthDate, refNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "APPLICATION_NOT_FOUND"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (ctrl *ApplicationController) GetOfficerApplications(c *gin.Context) {
	statusFilter := c.Query("status")
	apps, err := ctrl.usecase.GetOfficerApplications(statusFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, apps)
}

func (ctrl *ApplicationController) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	var dto entities.UpdateStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	app, err := ctrl.usecase.ReviewApplication(id, dto)
	if err != nil {
		if err.Error() == "REASON_REQUIRED_FOR_REJECT" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "REASON_REQUIRED_FOR_REJECT"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, app)
}
