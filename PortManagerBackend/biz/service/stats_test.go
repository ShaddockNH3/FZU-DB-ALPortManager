package service

import (
	"testing"

	"FZU-DB-ALPortManager/biz/dal"
	"FZU-DB-ALPortManager/biz/model/api"
)

// 统计的口径：数量倒序，数量相同按名称升序。界面靠这个顺序稳定，不能每次刷新都跳。
func TestToStatItemsSortsByCountThenName(t *testing.T) {
	items := toStatItems([]dal.GroupCount{
		{Name: "皇家", Count: 2},
		{Name: "白鹰", Count: 2},
		{Name: "重樱", Count: 5},
		{Name: "铁血", Count: 1},
	})

	want := []string{"重樱", "白鹰", "皇家", "铁血"}
	for i, name := range want {
		if items[i].Name != name {
			t.Fatalf("第 %d 项应该是 %s，实际 %s（顺序 %v）", i, name, items[i].Name, namesOf(items))
		}
	}
}

func TestToStatItemsEmpty(t *testing.T) {
	if got := toStatItems(nil); len(got) != 0 {
		t.Errorf("空输入应该返回空切片，实际 %d 项", len(got))
	}
}

func namesOf(items []*api.StatItem) []string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		names = append(names, item.Name)
	}
	return names
}
