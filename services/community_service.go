package services

import (
	"context"
	"database/sql"
	"fmt"

	"reveil-api/models"

	"github.com/google/uuid"
)

type CommunityService struct {
	db *sql.DB
}

func NewCommunityService(db *sql.DB) *CommunityService {
	return &CommunityService{db: db}
}

func (s *CommunityService) GetCommunity(ctx context.Context, id uuid.UUID) (*models.Community, error) {
	var c models.Community
	err := s.db.QueryRowContext(ctx, "SELECT id, name, description, created_at FROM communities WHERE id = $1", id).Scan(
		&c.ID, &c.Name, &c.Description, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *CommunityService) ListCommunities(ctx context.Context) ([]models.Community, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, description, created_at FROM communities")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var communities []models.Community
	for rows.Next() {
		var c models.Community
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt); err != nil {
			return nil, err
		}
		communities = append(communities, c)
	}
	return communities, nil
}

func (s *CommunityService) CreateCommunity(ctx context.Context, name, description string) (*models.Community, error) {
	var c models.Community
	err := s.db.QueryRowContext(ctx,
		"INSERT INTO communities (name, description) VALUES ($1, $2) RETURNING id, name, description, created_at",
		name, description).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create community: %w", err)
	}
	return &c, nil
}
