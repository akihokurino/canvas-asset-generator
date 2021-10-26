package datastore

import "math"

const defaultLimit = 20

type Pager struct {
	page  int
	limit int
}

func NewPager(page int, limit int) *Pager {
	return &Pager{
		page:  page,
		limit: limit,
	}
}

func (p *Pager) Page() int {
	return p.page
}

func (p *Pager) Offset() int {
	page := int(math.Max(float64(p.page), 1))
	offset := p.Limit()
	return page*offset - offset
}

func (p *Pager) Limit() int {
	if p.limit > 0 {
		return p.limit
	} else {
		return defaultLimit
	}
}
