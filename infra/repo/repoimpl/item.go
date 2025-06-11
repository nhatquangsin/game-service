package repoimpl

import (
	"context"

	"github.com/nhatquangsin/game-service/domain/entity"
	"github.com/nhatquangsin/game-service/domain/repo"
	"github.com/nhatquangsin/game-service/infra/config"
	"github.com/nhatquangsin/game-service/infra/repo/database"
	"github.com/nhatquangsin/game-service/infra/repo/entc/item"
)

// ItemRepo implements interface ItemRepo.
type ItemRepo struct {
	client database.Client
	cfg    config.Config
}

// NewItemRepo creates and returns a new instance of repo.ItemRepo.
func NewItemRepo(
	client database.Client,
	cfg config.Config,
) repo.ItemRepo {
	return &ItemRepo{
		client: client,
		cfg:    cfg,
	}
}

// FindAll to list all items.
func (r *ItemRepo) FindAll(ctx context.Context) ([]*entity.Item, error) {
	rows, err := r.client.Slave(ctx).Item.Query().
		Unique(false).
		All(ctx)
	if err != nil {
		return nil, err
	}
	var res []*entity.Item
	for _, row := range rows {
		i := &entity.Item{
			ID:          row.ID,
			Name:        row.Name,
			Category:    row.Category,
			Description: row.Description,
		}

		res = append(res, i)
	}

	return res, nil
}

// FindByItemIDs to find item by item id
func (r *ItemRepo) FindByItemIDs(ctx context.Context, itemIDs []string) ([]*entity.Item, error) {
	rows, err := r.client.Slave(ctx).Item.Query().Where(item.IDIn(itemIDs...)).Unique(false).All(ctx)
	if err != nil {
		return nil, err
	}

	var res []*entity.Item
	for _, row := range rows {
		res = append(res, &entity.Item{
			ID:          row.ID,
			Name:        row.Name,
			Category:    row.Category,
			Description: row.Description,
		})
	}

	return res, nil
}
