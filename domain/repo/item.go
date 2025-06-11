package repo

import (
	"context"

	"github.com/nhatquangsin/game-service/domain/entity"
)

// ItemRepo exposed all function interact with items data.
type ItemRepo interface {
	FindAll(ctx context.Context) ([]*entity.Item, error)
	FindByItemIDs(ctx context.Context, itemIDs []string) ([]*entity.Item, error)
}
