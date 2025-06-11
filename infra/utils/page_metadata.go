package utils

// BatchMetadata represents batch actions metadata of result for batch actions
// in service.
type BatchMetadata struct {
	NUpdated int `json:"nUpdated" proto:"1"`
	NDeleted int `json:"nDeleted" proto:"2"`
	NCreated int `json:"nCreated" proto:"3"`
	NFailed  int `json:"nFailed" proto:"4"`
}

// RangeMetadata represents a range metadata of result for range query in
// service.
type RangeMetadata struct {
	HasNext bool
	Limit   uint
	Offset  uint
}

// PageMetadata represents a page metadata of result for page query in service.
type PageMetadata struct {
	Limit   *int  `json:"limit,omitempty" proto:"1"`
	Offset  *int  `json:"offset,omitempty" proto:"2"`
	Total   *int  `json:"total,omitempty" proto:"3"`
	HasNext *bool `json:"hasNext,omitempty" proto:"4"`
}

// PaginationResult represents a pagination result for page query in service.
type PaginationResult[T any] struct {
	Items    []*T         `json:"_items,omitempty" proto:"1"`
	Metadata PageMetadata `json:"_metadata,omitempty" proto:"2"`
}
