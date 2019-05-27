package repository

import (
	"context"
	"../models"
)


type LinkRepo interface {
	Create(ctx context.Context, l *models.Link) (int64, error)
	Fetch(ctx context.Context) ([]*models.Link, error)
	Get(ctx context.Context, shortened string) ([]*models.Link, error)
	RegisterClick(ctx context.Context, c *models.Click) error
	GetClicks(ctx context.Context, owner string) (int, error)
}
