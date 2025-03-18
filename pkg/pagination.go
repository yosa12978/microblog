package pkg

type Page[T any] struct {
	Total    uint
	Current  uint
	Size     uint
	NextPage uint
	PrevPage uint
	HasNext  bool
	HasPrev  bool
	Content  []T
}

func NewPage[T any](content []T, total, current, size uint) Page[T] {
	nextPage := current
	if current < total {
		nextPage = current + 1
	}
	prevPage := current
	if current > 1 {
		prevPage = current - 1
	}
	return Page[T]{
		Content:  content,
		Total:    total,
		Current:  current,
		HasNext:  current < total,
		HasPrev:  current > 1,
		Size:     size,
		NextPage: nextPage,
		PrevPage: prevPage,
	}
}
