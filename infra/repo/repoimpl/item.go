package repoimpl

import (
	"context"

	"entgo.io/ent/dialect/sql"

	"github.com/nhatquangsin/game-service/domain/entity"
	"github.com/nhatquangsin/game-service/domain/repo"
	"github.com/nhatquangsin/game-service/infra/config"
	"github.com/nhatquangsin/game-service/infra/repo/database"
	"github.com/nhatquangsin/game-service/infra/repo/entc/item"
	"github.com/nhatquangsin/game-service/infra/utils"
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
func (r *ItemRepo) FindAll(ctx context.Context, limit, offset int) (*repo.ListItemResult, error) {
	builder := r.client.Slave(ctx).Item.Query()
	total, err := builder.Modify(func(s *sql.Selector) {
		s.Select("COUNT(1)")
	}).Int(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := builder.Modify(func(s *sql.Selector) {
		s.Select(item.Columns...)
	}).Offset(offset).Limit(limit).All(ctx)
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

	hasNext := total > limit+offset

	return &repo.ListItemResult{
		Items: res,
		Metadata: utils.PageMetadata{
			Total:   &total,
			HasNext: &hasNext,
			Offset:  &offset,
			Limit:   &limit,
		},
	}, nil
}

// FindByItemIDs to find item by item id
func (r *ItemRepo) FindByItemIDs(ctx context.Context, itemIDs []string, limit, offset int) (*repo.ListItemResult, error) {
	builder := r.client.Slave(ctx).Item.Query().Where(item.IDIn(itemIDs...))
	total, err := builder.Modify(func(s *sql.Selector) {
		s.Select("COUNT(1)")
	}).Int(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := builder.Modify(func(s *sql.Selector) {
		s.Select(item.Columns...)
	}).Offset(offset).Limit(limit).All(ctx)
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

	hasNext := total > limit+offset

	return &repo.ListItemResult{
		Items: res,
		Metadata: utils.PageMetadata{
			Total:   &total,
			HasNext: &hasNext,
			Offset:  &offset,
			Limit:   &limit,
		},
	}, nil
}
