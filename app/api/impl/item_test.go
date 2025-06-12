package impl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhatquangsin/game-service/app/api"
	"github.com/nhatquangsin/game-service/domain/repo"
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
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			itemRepo := repo.NewMockItemRepo(t)

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
