package domain

import (
	"context"

	"github.com/google/uuid"
)

// Professional agora mapeia para a tabela 'profissionais'
// ID é Foreign Key para perfis(id)
type Professional struct {
	ID         uuid.UUID `json:"id"` // PK, refere-se a perfis.id
	CategoryID uuid.UUID `json:"categoria_id"`
	Bio        string    `json:"bio"`
	PriceHour  float64   `json:"valor_hora"`
	Verified   bool      `json:"verificado"`
	Rating     float64   `json:"nota_media"`
	// ReviewCount removido do schema do usuário (teremos que calcular count ou adicionar se ele permitir)
	// Vamos assumir que "nota_media" é persistido.

	Latitude  float64 `json:"latitude"`  // Mapeado de 'localizacao'
	Longitude float64 `json:"longitude"` // Mapeado de 'localizacao'

	// Joins para retorno na API
	User     *User     `json:"perfil,omitempty"`
	Category *Category `json:"categoria,omitempty"`
}

// SearchFilters filtros adaptados
type SearchFilters struct {
	CategoryID uuid.UUID // Agora busca por ID da categoria
	MinRating  float64
	Latitude   float64
	Longitude  float64
	RadiusKM   float64
}

type ProfessionalRepository interface {
	Create(ctx context.Context, pro *Professional) error
	GetByID(ctx context.Context, id uuid.UUID) (*Professional, error)
	Search(ctx context.Context, filters SearchFilters) ([]*Professional, error)
	ListCategories(ctx context.Context) ([]*Category, error) // Novo método
}

type ProfessionalService interface {
	RegisterConfigs(ctx context.Context, pro *Professional) error
	GetProfile(ctx context.Context, id uuid.UUID) (*Professional, error)
	Search(ctx context.Context, categoryID string, rating float64, lat, long, radius float64) ([]*Professional, error)
	ListCategories(ctx context.Context) ([]*Category, error)
}
