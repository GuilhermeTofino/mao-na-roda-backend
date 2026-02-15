package main

import (
	"log"
	"os"

	"github.com/GuilhermeTofino/mao-na-roda-backend/internal/handler"
	"github.com/GuilhermeTofino/mao-na-roda-backend/internal/repository"
	"github.com/GuilhermeTofino/mao-na-roda-backend/internal/service"
	"github.com/GuilhermeTofino/mao-na-roda-backend/pkg/database"
	"github.com/GuilhermeTofino/mao-na-roda-backend/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente")
	}

	pool, err := database.InitDB()
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}
	defer database.CloseDB()

	userRepo := repository.NewUserRepository(pool)
	proRepo := repository.NewProfessionalRepository(pool)
	// serviceRepo e reviewRepo precisam ser refatorados para o schema novo se for usar
	// Por enquanto vamos focar em compilar User e Professional que foram alterados
	serviceRepo := repository.NewServiceRepository(pool)
	reviewRepo := repository.NewReviewRepository(pool)

	userService := service.NewUserService(userRepo)
	proService := service.NewProfessionalService(proRepo)
	serviceService := service.NewServiceService(serviceRepo, reviewRepo)
	reviewService := handler.NewReviewService(reviewRepo)

	userHandler := handler.NewUserHandler(userService)
	proHandler := handler.NewProfessionalHandler(proService)
	serviceHandler := handler.NewServiceHandler(serviceService, reviewService)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RecoveryMiddleware())

	// Rotas Públicas
	r.GET("/categorias", proHandler.ListCategories)
	r.GET("/profissionais/busca", proHandler.Search)

	// Rotas Protegidas
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Perfis
		protected.POST("/perfil", userHandler.CreateProfile) // Cria perfil na tabela 'perfis'
		protected.GET("/perfil/:id", userHandler.GetProfile)

		// Profissionais
		protected.POST("/profissionais", proHandler.Create)

		// Serviços (Schema ainda pendente de ajuste fino no Handler, mas lógica básica ok)
		protected.POST("/servicos", serviceHandler.CreateRequest)
		protected.PUT("/servicos/:id/status", serviceHandler.UpdateStatus)
		protected.POST("/servicos/:id/avaliar", serviceHandler.CreateReview)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
