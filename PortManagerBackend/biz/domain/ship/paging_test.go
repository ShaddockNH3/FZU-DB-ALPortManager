package ship

import "testing"

func TestNormalizePaging(t *testing.T) {
	cases := []struct {
		name                           string
		page, pageSize                 int
		wantPage, wantSize, wantOffset int
	}{
		{"正常值原样保留", 2, 10, 2, 10, 10},
		{"页码 0 归到第一页", 0, 10, 1, 10, 0},
		{"负数页码归到第一页", -5, 10, 1, 10, 0},
		{"页大小 0 用默认值", 3, 0, 3, DefaultPageSize, 2 * DefaultPageSize},
		{"页大小超上限被夹住", 1, 5000, 1, MaxPageSize, 0},
		{"负数页大小用默认值", 1, -20, 1, DefaultPageSize, 0},
	}

	for _, testCase := range cases {
		got := NormalizePaging(testCase.page, testCase.pageSize)
		if got.Page != testCase.wantPage || got.PageSize != testCase.wantSize || got.Offset != testCase.wantOffset {
			t.Errorf("%s: 期望 {%d %d %d}，实际 {%d %d %d}",
				testCase.name,
				testCase.wantPage, testCase.wantSize, testCase.wantOffset,
				got.Page, got.PageSize, got.Offset)
		}
	}
}

func TestNormalizePagingOffsetNeverNegative(t *testing.T) {
	for page := -10; page <= 3; page++ {
		if got := NormalizePaging(page, 20); got.Offset < 0 {
			t.Fatalf("page=%d 的 offset 变成了负数: %d", page, got.Offset)
		}
	}
}
