package impl

import (
	"context"

	"github.com/nhatquangsin/game-service/app/api"
	"github.com/nhatquangsin/game-service/cache"
	"github.com/nhatquangsin/game-service/domain/repo"
	"github.com/nhatquangsin/game-service/infra/utils"
)

// ItemService implements all use cases of Item.
type ItemService struct {
	itemRepo    repo.ItemRepo
	cachedItems *cache.CachedItems
}

// NewItemService creates and returns new instance of ItemService.
func NewItemService(
	itemRepo repo.ItemRepo,
	cachedItems *cache.CachedItems,
) api.ItemService {
	svc := &ItemService{
		itemRepo:    itemRepo,
		cachedItems: cachedItems,
	}

	return svc
}

// ListItems lists all items.
func (s *ItemService) ListItems(ctx context.Context, req *api.ListItemsRequest) (*api.ListItemsResponse, error) {
	var result []*api.Item
	var pageMetadata utils.PageMetadata

	if len(req.ItemIDs) == 0 {
		filteredResult := utils.PaginationOnMem(s.cachedItems.Items, req.Offset, req.Limit)
		for _, i := range filteredResult.Items {
			result = append(result, &api.Item{
				ID:          i.ID,
				Name:        i.Name,
				Description: i.Description,
				Category:    i.Category,
			})
		}

		pageMetadata = filteredResult.Metadata
	} else {
		data, err := s.itemRepo.FindByItemIDs(ctx, req.ItemIDs)
		if err != nil {
			return nil, err
		}

		filteredResult := utils.PaginationOnMem(data, req.Offset, req.Limit)
		for _, i := range filteredResult.Items {
			result = append(result, &api.Item{
				ID:          i.ID,
				Name:        i.Name,
				Description: i.Description,
				Category:    i.Category,
			})
		}

		pageMetadata = filteredResult.Metadata
	}

	return &api.ListItemsResponse{
		Items:    result,
		Metadata: pageMetadata,
	}, nil
}
