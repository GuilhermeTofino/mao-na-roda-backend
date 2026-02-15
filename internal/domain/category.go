package domain

import "github.com/google/uuid"

// Category mapeia a tabela 'categorias'
type Category struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"nome"`
	Icon string    `json:"icone"`
}
