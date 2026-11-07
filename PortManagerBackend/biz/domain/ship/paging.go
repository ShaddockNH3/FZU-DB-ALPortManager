package ship

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// Paging 是分页参数，Offset 由规范化的页码算出来，防止 page=0 算出负数偏移。
type Paging struct {
	Page     int
	PageSize int
	Offset   int
}

// NormalizePaging 把前端传来的任意页码夹到合法区间。
func NormalizePaging(page, pageSize int) Paging {
	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return Paging{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}
}
