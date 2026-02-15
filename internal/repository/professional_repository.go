package repository

import (
	"context"
	"fmt"

	"github.com/GuilhermeTofino/mao-na-roda-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type professionalRepo struct {
	db *pgxpool.Pool
}

func NewProfessionalRepository(db *pgxpool.Pool) domain.ProfessionalRepository {
	return &professionalRepo{db: db}
}

func (r *professionalRepo) Create(ctx context.Context, pro *domain.Professional) error {
	// Tabela profissionais
	// Colunas: id, categoria_id, bio, valor_hora, verificado, localizacao, nota_media
	query := `
		INSERT INTO profissionais (id, categoria_id, bio, valor_hora, verificado, localizacao, nota_media)
		VALUES ($1, $2, $3, $4, $5, ST_SetSRID(ST_MakePoint($6, $7), 4326)::geography, $8)
	`
	// Nota: id é FK de perfis(id). Deve ser passado na criação.

	_, err := r.db.Exec(ctx, query,
		pro.ID, pro.CategoryID, pro.Bio, pro.PriceHour,
		pro.Verified, pro.Longitude, pro.Latitude, pro.Rating,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar profissional: %w", err)
	}
	return nil
}

func (r *professionalRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Professional, error) {
	// Join com categorias e perfis para trazer dados completos
	query := `
		SELECT 
			p.id, p.categoria_id, p.bio, p.valor_hora, p.verificado, p.nota_media, ST_X(p.localizacao::geometry), ST_Y(p.localizacao::geometry),
			c.nome as categoria_nome,
			u.nome_completo, u.avatar_url
		FROM profissionais p
		JOIN categorias c ON p.categoria_id = c.id
		JOIN perfis u ON p.id = u.id
		WHERE p.id = $1
	`
	var pro domain.Professional
	pro.Category = &domain.Category{}
	pro.User = &domain.User{}

	var avatar *string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&pro.ID, &pro.CategoryID, &pro.Bio, &pro.PriceHour, &pro.Verified, &pro.Rating, &pro.Longitude, &pro.Latitude,
		&pro.Category.Name,
		&pro.User.Name, &avatar,
	)
	if err != nil {
		return nil, fmt.Errorf("profissional não encontrado: %w", err)
	}

	pro.Category.ID = pro.CategoryID
	pro.User.ID = pro.ID // ID é o mesmo
	if avatar != nil {
		pro.User.AvatarURL = *avatar
	}

	return &pro, nil
}

func (r *professionalRepo) Search(ctx context.Context, filters domain.SearchFilters) ([]*domain.Professional, error) {
	query := `
		SELECT 
			p.id, p.categoria_id, p.bio, p.valor_hora, p.verificado, p.nota_media, ST_X(p.localizacao::geometry), ST_Y(p.localizacao::geometry),
			c.nome, c.icone,
			u.nome_completo, u.avatar_url
		FROM profissionais p
		JOIN categorias c ON p.categoria_id = c.id
		JOIN perfis u ON p.id = u.id
		WHERE 1=1
	`
	args := []interface{}{}
	argID := 1

	if filters.CategoryID != uuid.Nil {
		query += fmt.Sprintf(" AND p.categoria_id = $%d", argID)
		args = append(args, filters.CategoryID)
		argID++
	}

	if filters.MinRating > 0 {
		query += fmt.Sprintf(" AND p.nota_media >= $%d", argID)
		args = append(args, filters.MinRating)
		argID++
	}

	if filters.RadiusKM > 0 {
		query += fmt.Sprintf(" AND ST_DWithin(p.localizacao, ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography, $%d)", argID, argID+1, argID+2)
		args = append(args, filters.Longitude, filters.Latitude, filters.RadiusKM*1000)
		argID += 3
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro na busca: %w", err)
	}
	defer rows.Close()

	var pros []*domain.Professional
	for rows.Next() {
		var p domain.Professional
		p.Category = &domain.Category{}
		p.User = &domain.User{}
		var avatar *string

		err := rows.Scan(
			&p.ID, &p.CategoryID, &p.Bio, &p.PriceHour, &p.Verified, &p.Rating, &p.Longitude, &p.Latitude,
			&p.Category.Name, &p.Category.Icon,
			&p.User.Name, &avatar,
		)
		if err != nil {
			continue
		}
		if avatar != nil {
			p.User.AvatarURL = *avatar
		}
		p.Category.ID = p.CategoryID
		p.User.ID = p.ID // ID do perfil é o mesmo do profissional
		pros = append(pros, &p)
	}
	return pros, nil
}

func (r *professionalRepo) ListCategories(ctx context.Context) ([]*domain.Category, error) {
	query := `SELECT id, nome, icone FROM categorias`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []*domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Icon); err != nil {
			continue
		}
		cats = append(cats, &c)
	}
	return cats, nil
}
