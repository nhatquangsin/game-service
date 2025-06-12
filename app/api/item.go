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

// Item rest resource.
//
// +smkit:rest:resource=true
type Item struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
}

// ListItemsRequest represents a request for list item.
type ListItemsRequest struct {
	ItemIDs []string `json:"-" query:"itemIDs" validate:"max=100"`
	Offset  int      `json:"-" query:"offset" field:"offset" validate:"gte=0"`
	Limit   int      `json:"-" query:"limit" field:"limit" validate:"gte=1,lte=100"`
}

// ListItemsResponse represents a response for list item.
type ListItemsResponse struct {
	Items    []*entity.Item     `field:"_items" json:"_items,omitempty"`
	Metadata utils.PageMetadata `field:"_metadata" json:"_metadata,omitempty"`
}
