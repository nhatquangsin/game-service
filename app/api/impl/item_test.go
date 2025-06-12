package impl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhatquangsin/game-service/app/api"
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
			name: "TC01 - Should return 2 items when total 2 items, limit 10 offset 0",
			args: args{
				req: &api.ListItemsRequest{
					Limit:  10,
					Offset: 0,
				},
			},
			wantErr: false,
			wantRes: &api.ListItemsResponse{
				Items: []*entity.Item{
					&entity.Item{
						ID:   "id_1",
						Name: "name_1",
					},
					&entity.Item{
						ID:   "id_2",
						Name: "name_2",
					},
				},
				Metadata: utils.PageMetadata{
					Total:   utils.Of(2),
					HasNext: utils.Of(false),
					Limit:   utils.Of(10),
					Offset:  utils.Of(0),
				},
			},
		},
		{
			name: "TC02 - Should return no item when total 2 items, limit 10 offset 1",
			args: args{
				req: &api.ListItemsRequest{
					Limit:  10,
					Offset: 1,
				},
			},
			wantErr: false,
			wantRes: &api.ListItemsResponse{
				Items: []*entity.Item{},
				Metadata: utils.PageMetadata{
					Total:   utils.Of(2),
					HasNext: utils.Of(false),
					Limit:   utils.Of(10),
					Offset:  utils.Of(1),
				},
			},
		},
		{
			name: "TC03 - Should return 1 item when total 3 items, limit 1 offset 0",
			args: args{
				req: &api.ListItemsRequest{
					Limit:  1,
					Offset: 1,
				},
			},
			wantErr: false,
			wantRes: &api.ListItemsResponse{
				Items: []*entity.Item{
					&entity.Item{
						ID:   "id_1",
						Name: "name_1",
					},
					&entity.Item{
						ID:   "id_2",
						Name: "name_2",
					},
				},
				Metadata: utils.PageMetadata{
					Total:   utils.Of(2),
					HasNext: utils.Of(false),
					Limit:   utils.Of(10),
					Offset:  utils.Of(1),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			itemRepo := repo.NewMockItemRepo(t)

			itemRepo.EXPECT().FindByItemIDs(ctx, tt.args.req.ItemIDs, tt.args.req.Limit, tt.args.req.Offset).
				Return(&repo.ListItemResult{
					Items:    tt.wantRes.Items,
					Metadata: tt.wantRes.Metadata,
				}, nil).Once()

			svc := &ItemService{
				itemRepo: itemRepo,
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
