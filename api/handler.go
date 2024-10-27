package api

import (
	"net/http"
	"strconv"

	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	campaignService campaign.ServiceInterface
	userService     user.ServiceInterface
	logger          loggo.LoggerInterface
	errorHandler    *shared.ErrorHandler
}

func (h *Handler) GetCampaigns(c echo.Context) error {
	campaigns, err := h.campaignService.GetAllCampaigns()
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error fetching campaigns", http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, campaigns)
}

func (h *Handler) GetCampaign(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}

	campaign, err := h.campaignService.GetCampaignByID(campaign.GetCampaignParams{ID: id})
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error fetching campaign", http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, campaign)
}

func (h *Handler) CreateCampaign(c echo.Context) error {
	dto := new(campaign.CreateCampaignDTO)
	if err := c.Bind(dto); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}

	if err := h.campaignService.CreateCampaign(dto); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error creating campaign", http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, dto)
}

func (h *Handler) UpdateCampaign(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}

	dto := new(campaign.UpdateCampaignDTO)
	if err := c.Bind(dto); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}
	dto.ID = id

	if err := h.campaignService.UpdateCampaign(dto); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error updating campaign", http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, dto)
}

func (h *Handler) DeleteCampaign(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}

	if err := h.campaignService.DeleteCampaign(campaign.DeleteCampaignParams{ID: id}); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error deleting campaign", http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

// User-related handlers
func (h *Handler) RegisterUser(c echo.Context) error {
	dto := new(user.CreateDTO)
	if err := c.Bind(dto); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}

	if err := h.userService.RegisterUser(dto); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error registering user", http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, dto)
}

func (h *Handler) LoginUser(c echo.Context) error {
	dto := new(user.LoginDTO)
	if err := c.Bind(dto); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}

	token, err := h.userService.VerifyUser(dto)
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid credentials", http.StatusUnauthorized)
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) GetUser(c echo.Context) error {
	username := c.Param("username")
	dto := &user.GetDTO{Username: username}

	userDetails, err := h.userService.GetUser(dto)
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error fetching user", http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, userDetails)
}
