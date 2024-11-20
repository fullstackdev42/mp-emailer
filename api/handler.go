package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/mp-emailer/campaign"
	"github.com/jonesrussell/mp-emailer/shared"
	"github.com/jonesrussell/mp-emailer/user"
	"github.com/labstack/echo/v4"
)

// Handler is the API handler
type Handler struct {
	campaignService campaign.ServiceInterface
	userService     user.ServiceInterface
	logger          loggo.LoggerInterface
	errorHandler    shared.ErrorHandlerInterface
	jwtExpiry       int
}

func (h *Handler) GetCampaigns(c echo.Context) error {
	campaigns, err := h.campaignService.GetCampaigns()
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error fetching campaigns", http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, campaigns)
}

func (h *Handler) GetCampaign(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}

	cmpn, err := h.campaignService.GetCampaignByID(campaign.GetCampaignParams{ID: id})
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error fetching campaign", http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, cmpn)
}

func (h *Handler) CreateCampaign(c echo.Context) error {
	dto := new(campaign.CreateCampaignDTO)
	if err := c.Bind(dto); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}

	createdCampaign, err := h.campaignService.CreateCampaign(dto)
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error creating campaign", http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, createdCampaign)
}

func (h *Handler) UpdateCampaign(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
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
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}

	if err := h.campaignService.DeleteCampaign(campaign.DeleteCampaignDTO{ID: id}); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error deleting campaign", http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

// RegisterUser User-related handlers
func (h *Handler) RegisterUser(c echo.Context) error {
	dto := new(user.RegisterDTO)
	if err := c.Bind(dto); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}

	createdUser, err := h.userService.RegisterUser(dto)
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error registering user", http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, createdUser)
}

func (h *Handler) LoginUser(c echo.Context) error {
	dto := new(user.LoginDTO)
	if err := c.Bind(dto); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}

	token, err := h.userService.LoginUser(dto)
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

func (h *Handler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		h.logger.Debug("Auth header: " + authHeader)

		if authHeader == "" {
			h.logger.Warn("Missing Authorization header")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing Authorization header"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		h.logger.Debug("Token: " + tokenString)

		claims, err := shared.ValidateToken(tokenString, os.Getenv("JWT_SECRET"))
		if err != nil {
			h.logger.Error("Token validation error: ", err)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		if shared.IsTokenExpired(claims) {
			h.logger.Warn("Token is expired")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired"})
		}

		h.logger.Info("Token validated successfully for user: " + claims.Username)
		c.Set("username", claims.Username)
		return next(c)
	}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	// Public routes
	e.POST("/api/user/login", h.LoginUser)

	// Protected routes
	api := e.Group("/api")
	api.Use(h.authMiddleware) // Apply the middleware to all routes in this group

	api.GET("/campaigns", h.GetCampaigns)
	api.GET("/campaign/:id", h.GetCampaign)
	api.POST("/campaign", h.CreateCampaign)
	api.PUT("/campaign/:id", h.UpdateCampaign)
	api.DELETE("/campaign/:id", h.DeleteCampaign)
	api.GET("/user/:username", h.GetUser)
}
