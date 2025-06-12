package cache

import (
	"context"

	"github.com/nhatquangsin/game-service/domain/entity"
	"github.com/nhatquangsin/game-service/domain/repo"
)

// CachedItems represent for all items load from db.
type CachedItems struct {
	ItemsMap map[string]*entity.Item
	Items    []*entity.Item
}

// LoadAllItems load all items from db.
func LoadAllItems(itemRepo repo.ItemRepo) (*CachedItems, error) {
	innerMap := make(map[string]*entity.Item)
	var inner []*entity.Item
	ctx := context.Background()
	items, err := itemRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		innerMap[item.ID] = item
		inner = append(inner, item)
	}
	return &CachedItems{
		ItemsMap: innerMap,
		Items:    inner,
	}, nil
}
