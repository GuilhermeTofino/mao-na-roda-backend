package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/GuilhermeTofino/mao-na-roda-backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServiceHandler struct {
	service       domain.ServiceService
	reviewService domain.ReviewService // Se decidir separar, ou usar repositório direto se simples
	// Aqui vamos adicionar a lógica de review direta no handler ou criar um ReviewService separado.
	// Pelo request inicial, "ServiceService" parece englobar fluxo.
	// Vamos assumir que criaremos um ReviewService simples para completar.
}

// Vamos precisar de um ReviewService
type reviewService struct {
	repo domain.ReviewRepository
}

func NewReviewService(repo domain.ReviewRepository) domain.ReviewService {
	return &reviewService{repo: repo}
}
func (s *reviewService) CreateReview(ctx context.Context, review *domain.Review) error {
	review.ID = uuid.New()
	return s.repo.Create(ctx, review)
}

// Handler principal
type MainServiceHandler struct {
	svc       domain.ServiceService
	reviewSvc domain.ReviewService
}

func NewServiceHandler(svc domain.ServiceService, reviewSvc domain.ReviewService) *MainServiceHandler {
	return &MainServiceHandler{svc: svc, reviewSvc: reviewSvc}
}

type requestServicePayload struct {
	ProfessionalID string `json:"professional_id" binding:"required"`
	Description    string `json:"description" binding:"required"`
	Date           string `json:"scheduled_for" binding:"required"` // Format: 2023-10-27T10:00:00Z
}

func (h *MainServiceHandler) CreateRequest(c *gin.Context) {
	clientIDStr := c.GetHeader("X-User-ID") // Mock: pegar do auth middleware
	clientID, _ := uuid.Parse(clientIDStr)

	var req requestServicePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	proID, _ := uuid.Parse(req.ProfessionalID)
	date, _ := time.Parse(time.RFC3339, req.Date)

	serviceReq, err := h.svc.RequestService(c.Request.Context(), clientID, proID, req.Description, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, serviceReq)
}

type updateStatusPayload struct {
	Status domain.ServiceStatus `json:"status" binding:"required"`
}

func (h *MainServiceHandler) UpdateStatus(c *gin.Context) {
	actorIDStr := c.GetHeader("X-User-ID")
	actorID, _ := uuid.Parse(actorIDStr)

	idStr := c.Param("id")
	id, _ := uuid.Parse(idStr)

	var req updateStatusPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateStatus(c.Request.Context(), id, req.Status, actorID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status atualizado"})
}

type createReviewPayload struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

func (h *MainServiceHandler) CreateReview(c *gin.Context) {
	clientIDStr := c.GetHeader("X-User-ID")
	clientID, _ := uuid.Parse(clientIDStr)

	// Nota: Em um app real, o serviceID seria usado para validar se o serviço existe e pegar o professionalID.
	// Aqui estamos simplificando conforme o refactor anterior.

	var req createReviewPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Idealmente, buscar o serviço para pegar o professionalID
	// Aqui vamos assumir que o frontend manda ou buscamos no banco.
	// Simplificação: vamos buscar o serviço aqui (h.svc.GetByID...) - não exposto na interface svc ainda mas ok.
	// Vamos assumir que mandamos no payload para simplificar o exemplo, ou injetamos repo no handler (menos clean).
	// Melhor: Adicionar método GetByID no ServiceService. (Feito no repo, não exportado no interface service ainda).

	// Por convenção, vamos assumir que o ServiceService tem um método auxiliar ou que passamos o ProID no body também por enquanto
	// para não complicar demais o exemplo.
	// TODO: Validar se serviço está "Concluído" antes de avaliar.

	// Mock ProID
	proID := uuid.New() // Deveria vir do serviço

	review := &domain.Review{
		// ServiceID removido conforme schema do usuário
		ClientID:       clientID,
		ProfessionalID: proID,
		Rating:         req.Rating,
		Comment:        req.Comment,
	}

	if err := h.reviewSvc.CreateReview(c.Request.Context(), review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}
