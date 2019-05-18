package repository

import (
	"context"
	"../models"
)


type LinkRepo interface {
	Create(ctx context.Context, l *models.Link) (int64, error)
	Fetch(ctx context.Context) ([]*models.Link, error)
}
