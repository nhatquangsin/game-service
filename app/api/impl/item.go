package impl

import (
	"context"

	"github.com/nhatquangsin/game-service/app/api"
	"github.com/nhatquangsin/game-service/domain/repo"
)

// ItemService implements all use cases of Item.
type ItemService struct {
	itemRepo repo.ItemRepo
}

// NewItemService creates and returns new instance of ItemService.
func NewItemService(
	itemRepo repo.ItemRepo,
) api.ItemService {
	svc := &ItemService{
		itemRepo: itemRepo,
	}

	return svc
}

// ListItems lists all items.
func (s *ItemService) ListItems(ctx context.Context, req *api.ListItemsRequest) (*api.ListItemsResponse, error) {
	result, err := s.itemRepo.FindByItemIDs(ctx, req.ItemIDs, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	return &api.ListItemsResponse{
		Items:    result.Items,
		Metadata: result.Metadata,
	}, nil
}
