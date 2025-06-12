package repo

import (
	"context"

	"github.com/nhatquangsin/game-service/domain/entity"
	"github.com/nhatquangsin/game-service/infra/utils"
)

// ItemRepo exposed all function interact with items data.
type ItemRepo interface {
	FindAll(ctx context.Context, limit, offset int) (*ListItemResult, error)
	FindByItemIDs(ctx context.Context, itemIDs []string, limit, offset int) (*ListItemResult, error)
}

type ListItemResult struct {
	Items    []*entity.Item
	Metadata utils.PageMetadata
}
