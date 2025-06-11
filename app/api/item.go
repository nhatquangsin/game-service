package api

import (
	"context"

	"github.com/nhatquangsin/game-service/domain/entity"
	"github.com/nhatquangsin/game-service/infra/utils"
)

// ItemService exposes all available use cases of item.
type ItemService interface {
	ListItems(ctx context.Context, req *ListItemsRequest) (*ListItemsResponse, error)
}

// ListItemsRequest represents a request for list item.
type ListItemsRequest struct {
	Offset int `json:"-" query:"offset" field:"offset" validate:"gte=0"`
	Limit  int `json:"-" query:"limit" field:"limit" validate:"gte=1,lte=100"`
}

// ListItemsResponse represents a response for list item.
type ListItemsResponse struct {
	Items    []*entity.Item     `field:"_items" json:"_items,omitempty"`
	Metadata utils.PageMetadata `field:"_metadata" json:"_metadata,omitempty"`
}
