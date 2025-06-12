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
					&api.Item{
						ID:          "item_1",
						Description: "desc 1",
					},
					&api.Item{
						ID:          "item_2",
						Description: "desc 2",
					},
					&api.Item{
						ID:          "item_3",
						Description: "desc 3",
					},
					&api.Item{
						ID:          "item_4",
						Description: "desc 4",
					},
					&api.Item{
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
					&api.Item{
						ID:          "item_1",
						Description: "desc 1",
					},
					&api.Item{
						ID:          "item_2",
						Description: "desc 2",
					},
					&api.Item{
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
					&api.Item{
						ID:          "item_2",
						Description: "desc 2",
					},
					&api.Item{
						ID:          "item_3",
						Description: "desc 3",
					},
					&api.Item{
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
					&api.Item{
						ID:          "item_4",
						Description: "desc 4",
					},
					&api.Item{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			itemRepo := repo.NewMockItemRepo(t)

			svc := &ItemService{
				itemRepo: itemRepo,
				cachedItems: &cache.CachedItems{
					Items: []*entity.Item{
						&entity.Item{
							ID:          "item_1",
							Description: "desc 1",
						},
						&entity.Item{
							ID:          "item_2",
							Description: "desc 2",
						},
						&entity.Item{
							ID:          "item_3",
							Description: "desc 3",
						},
						&entity.Item{
							ID:          "item_4",
							Description: "desc 4",
						},
						&entity.Item{
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
