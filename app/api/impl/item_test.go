package impl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhatquangsin/game-service/app/api"
	"github.com/nhatquangsin/game-service/cache"
	"github.com/nhatquangsin/game-service/domain/entity"
	"github.com/nhatquangsin/game-service/domain/repo"
	"github.com/nhatquangsin/game-service/infra/utils"
)

type itemRepoFindByItemIDsArgs struct {
	itemIDs []string
}

type itemRepoFindByItemIDsWant struct {
	items []*entity.Item
	err   error
}

func TestItemService_ListItems(t *testing.T) {
	type args struct {
		req *api.ListItemsRequest
	}

	tests := []struct {
		name    string
		args    args
		err     error
		wantErr bool
		wantRes *api.ListItemsResponse

		itemRepoFindByItemIDsArgs *itemRepoFindByItemIDsArgs
		itemRepoFindByItemIDsWant *itemRepoFindByItemIDsWant
	}{
		{
			name: "TC01 - request items empty, limit 10, offset 0 - should return all items",
			args: args{
				&api.ListItemsRequest{
					Limit:  10,
					Offset: 0,
				},
			},
			wantErr: false,
			wantRes: &api.ListItemsResponse{
				Items: []*api.Item{
					{
						ID:          "item_1",
						Description: "desc 1",
					},
					{
						ID:          "item_2",
						Description: "desc 2",
					},
					{
						ID:          "item_3",
						Description: "desc 3",
					},
					{
						ID:          "item_4",
						Description: "desc 4",
					},
					{
						ID:          "item_5",
						Description: "desc 5",
					},
				},
				Metadata: utils.PageMetadata{
					Total:   utils.Of(5),
					HasNext: utils.Of(false),
					Limit:   utils.Of(10),
					Offset:  utils.Of(0),
				},
			},
		},
		{
			name: "TC02 - request items empty, limit 3, offset 0 - should return 3 items",
			args: args{
				&api.ListItemsRequest{
					Limit:  3,
					Offset: 0,
				},
			},
			wantErr: false,
			wantRes: &api.ListItemsResponse{
				Items: []*api.Item{
					{
						ID:          "item_1",
						Description: "desc 1",
					},
					{
						ID:          "item_2",
						Description: "desc 2",
					},
					{
						ID:          "item_3",
						Description: "desc 3",
					},
				},
				Metadata: utils.PageMetadata{
					Total:   utils.Of(5),
					HasNext: utils.Of(true),
					Limit:   utils.Of(3),
					Offset:  utils.Of(0),
				},
			},
		},
		{
			name: "TC03 - request items empty, limit 2, offset 1 - should return 3 items",
			args: args{
				&api.ListItemsRequest{
					Limit:  3,
					Offset: 1,
				},
			},
			wantErr: false,
			wantRes: &api.ListItemsResponse{
				Items: []*api.Item{
					{
						ID:          "item_2",
						Description: "desc 2",
					},
					{
						ID:          "item_3",
						Description: "desc 3",
					},
					{
						ID:          "item_4",
						Description: "desc 4",
					},
				},
				Metadata: utils.PageMetadata{
					Total:   utils.Of(5),
					HasNext: utils.Of(true),
					Limit:   utils.Of(3),
					Offset:  utils.Of(1),
				},
			},
		},
		{
			name: "TC04 - request items empty, limit 2, offset 3 - should return 2 items",
			args: args{
				&api.ListItemsRequest{
					Limit:  3,
					Offset: 3,
				},
			},
			wantErr: false,
			wantRes: &api.ListItemsResponse{
				Items: []*api.Item{
					{
						ID:          "item_4",
						Description: "desc 4",
					},
					{
						ID:          "item_5",
						Description: "desc 5",
					},
				},
				Metadata: utils.PageMetadata{
					Total:   utils.Of(5),
					HasNext: utils.Of(false),
					Limit:   utils.Of(3),
					Offset:  utils.Of(3),
				},
			},
		},
		{
			name: "TC05 - request items 1, limit 2, offset 0 - should return 1 items",
			args: args{
				&api.ListItemsRequest{
					ItemIDs: []string{"item_3"},
					Limit:   2,
					Offset:  0,
				},
			},
			wantErr: false,
			wantRes: &api.ListItemsResponse{
				Items: []*api.Item{
					{
						ID:          "item_3",
						Description: "desc 3",
					},
				},
				Metadata: utils.PageMetadata{
					Total:   utils.Of(1),
					HasNext: utils.Of(false),
					Limit:   utils.Of(2),
					Offset:  utils.Of(0),
				},
			},
			itemRepoFindByItemIDsArgs: &itemRepoFindByItemIDsArgs{
				itemIDs: []string{"item_3"},
			},
			itemRepoFindByItemIDsWant: &itemRepoFindByItemIDsWant{
				items: []*entity.Item{
					{
						ID:          "item_3",
						Description: "desc 3",
					},
				},
			},
		},
		{
			name: "TC06 - request items 3, limit 2, offset 0 - should return 2 items",
			args: args{
				&api.ListItemsRequest{
					ItemIDs: []string{"item_3", "item_1", "item_5"},
					Limit:   2,
					Offset:  0,
				},
			},
			wantErr: false,
			wantRes: &api.ListItemsResponse{
				Items: []*api.Item{
					{
						ID:          "item_3",
						Description: "desc 3",
					},
					{
						ID:          "item_1",
						Description: "desc 1",
					},
				},
				Metadata: utils.PageMetadata{
					Total:   utils.Of(3),
					HasNext: utils.Of(true),
					Limit:   utils.Of(2),
					Offset:  utils.Of(0),
				},
			},
			itemRepoFindByItemIDsArgs: &itemRepoFindByItemIDsArgs{
				itemIDs: []string{"item_3", "item_1", "item_5"},
			},
			itemRepoFindByItemIDsWant: &itemRepoFindByItemIDsWant{
				items: []*entity.Item{
					{
						ID:          "item_3",
						Description: "desc 3",
					},
					{
						ID:          "item_1",
						Description: "desc 1",
					},
					{
						ID:          "item_5",
						Description: "desc 5",
					},
				},
			},
		},
		{
			name: "TC07 - request items 3, limit 2, offset 2 - should return 1 items",
			args: args{
				&api.ListItemsRequest{
					ItemIDs: []string{"item_3", "item_1", "item_5"},
					Limit:   2,
					Offset:  2,
				},
			},
			wantErr: false,
			wantRes: &api.ListItemsResponse{
				Items: []*api.Item{
					{
						ID:          "item_5",
						Description: "desc 5",
					},
				},
				Metadata: utils.PageMetadata{
					Total:   utils.Of(3),
					HasNext: utils.Of(false),
					Limit:   utils.Of(2),
					Offset:  utils.Of(2),
				},
			},
			itemRepoFindByItemIDsArgs: &itemRepoFindByItemIDsArgs{
				itemIDs: []string{"item_3", "item_1", "item_5"},
			},
			itemRepoFindByItemIDsWant: &itemRepoFindByItemIDsWant{
				items: []*entity.Item{
					{
						ID:          "item_3",
						Description: "desc 3",
					},
					{
						ID:          "item_1",
						Description: "desc 1",
					},
					{
						ID:          "item_5",
						Description: "desc 5",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			itemRepo := repo.NewMockItemRepo(t)

			if len(tt.args.req.ItemIDs) > 0 {
				itemRepo.EXPECT().FindByItemIDs(ctx, tt.itemRepoFindByItemIDsArgs.itemIDs).Return(
					tt.itemRepoFindByItemIDsWant.items,
					tt.itemRepoFindByItemIDsWant.err,
				).Once()
			}

			svc := &ItemService{
				itemRepo: itemRepo,
				cachedItems: &cache.CachedItems{
					Items: []*entity.Item{
						{
							ID:          "item_1",
							Description: "desc 1",
						},
						{
							ID:          "item_2",
							Description: "desc 2",
						},
						{
							ID:          "item_3",
							Description: "desc 3",
						},
						{
							ID:          "item_4",
							Description: "desc 4",
						},
						{
							ID:          "item_5",
							Description: "desc 5",
						},
					},
				},
			}

			res, err := svc.ListItems(ctx, tt.args.req)
			if tt.wantErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, res.Items, tt.wantRes.Items)
			assert.Equal(t, res.Metadata, tt.wantRes.Metadata)
		})
	}
}
