package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserType define se o usuário é Cliente ou Profissional
type UserType string

const (
	ClientType       UserType = "cliente"      // Alterado para lowercase para bater com check constraint
	ProfessionalType UserType = "profissional" // Alterado para lowercase
)

// User agora mapeia para a tabela 'perfis'
// O ID é o mesmo do auth.users do Supabase
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"nome_completo"` // nome_completo
	Type      UserType  `json:"tipo_usuario"`  // tipo_usuario
	AvatarURL string    `json:"avatar_url"`
	Phone     string    `json:"telefone"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRepository define os métodos de acesso a dados para perfis
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	// GetByEmail removido pois o email fica no auth.users, não no perfis (geralmente)
	// Se precisarmos, teríamos que fazer join com auth.users ou confiar no ID do token.
}

// UserService adaptado
type UserService interface {
	// Register agora cria um perfil vinculado a um usuário Auth existente (ou criado via client)
	CreateProfile(ctx context.Context, user *User) error
	GetProfile(ctx context.Context, id uuid.UUID) (*User, error)
}
