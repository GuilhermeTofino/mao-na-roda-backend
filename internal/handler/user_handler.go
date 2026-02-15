package handler

import (
	"net/http"

	"github.com/GuilhermeTofino/mao-na-roda-backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service domain.UserService
}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{service: service}
}

type createProfileRequest struct {
	Name      string          `json:"nome_completo" binding:"required"`
	Type      domain.UserType `json:"tipo_usuario" binding:"required,oneof=cliente profissional"`
	AvatarURL string          `json:"avatar_url"`
	Phone     string          `json:"telefone"`
}

func (h *UserHandler) CreateProfile(c *gin.Context) {
	// ID deve vir do Token Auth (Supabase Auth)
	// Como MOCK, vamos pegar do Header ou gerar novo se for teste local sem auth real
	userIDStr := c.GetHeader("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID de usuário ausente ou inválido"})
		return
	}

	var req createProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		ID:        userID,
		Name:      req.Name,
		Type:      req.Type,
		AvatarURL: req.AvatarURL,
		Phone:     req.Phone,
	}

	if err := h.service.CreateProfile(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	user, err := h.service.GetProfile(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Removemos Login/Register pois isso é feito no Client via Supabase SDK
// O backend só lida com o perfil adicional na tabela 'perfis'
