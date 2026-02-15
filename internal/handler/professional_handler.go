package handler

import (
	"net/http"
	"strconv"

	"github.com/GuilhermeTofino/mao-na-roda-backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfessionalHandler struct {
	service domain.ProfessionalService
}

func NewProfessionalHandler(service domain.ProfessionalService) *ProfessionalHandler {
	return &ProfessionalHandler{service: service}
}

type createProfessionalRequest struct {
	CategoryID string  `json:"categoria_id" binding:"required"`
	Bio        string  `json:"bio" binding:"required"`
	PriceHour  float64 `json:"valor_hora" binding:"required"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

func (h *ProfessionalHandler) Create(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var req createProfessionalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	catID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de categoria inválido"})
		return
	}

	pro := &domain.Professional{
		ID:         userID, // ID do profissional = ID do perfil (1:1)
		CategoryID: catID,
		Bio:        req.Bio,
		PriceHour:  req.PriceHour,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
	}

	if err := h.service.RegisterConfigs(c.Request.Context(), pro); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pro)
}

func (h *ProfessionalHandler) Search(c *gin.Context) {
	categoryID := c.Query("category_id")
	ratingStr := c.Query("rating")
	latStr := c.Query("lat")
	longStr := c.Query("long")
	radiusStr := c.Query("radius")

	var rating, lat, long, radius float64
	var err error

	if ratingStr != "" {
		rating, _ = strconv.ParseFloat(ratingStr, 64)
	}
	if latStr != "" {
		lat, err = strconv.ParseFloat(latStr, 64)
	}
	if longStr != "" {
		long, err = strconv.ParseFloat(longStr, 64)
	}
	if radiusStr != "" {
		radius, err = strconv.ParseFloat(radiusStr, 64)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros numéricos inválidos"})
		return
	}

	pros, err := h.service.Search(c.Request.Context(), categoryID, rating, lat, long, radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pros)
}

func (h *ProfessionalHandler) ListCategories(c *gin.Context) {
	cats, err := h.service.ListCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cats)
}
