package utils

// PaginationOnMem paginate a list of items on memory by offset and limit
func PaginationOnMem[T any](items []*T, offset, limit int) *PaginationResult[T] {
	start := offset
	end := offset + limit

	if start > len(items) {
		start = len(items)
	}

	if end > len(items) {
		end = len(items)
	}

	return &PaginationResult[T]{
		Items: items[start:end],
		Metadata: PageMetadata{
			Limit:   Of(limit),
			Offset:  Of(offset),
			Total:   Of(len(items)),
			HasNext: Of(len(items) > offset+limit),
		},
	}
}
