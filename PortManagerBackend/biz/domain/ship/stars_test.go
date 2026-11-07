package ship

import "testing"

func TestStarRules(t *testing.T) {
	cases := []struct {
		rarity       string
		initial, max int
	}{
		{"普通", 1, 4},
		{"稀有", 2, 5},
		{"精锐", 2, 5},
		{"超稀有", 3, 6},
		{"海上传奇", 3, 6},
		{"最高方案", 3, 6},
		{"决战方案", 3, 6},
		{"没见过的稀有度", 1, 3},
	}

	for _, testCase := range cases {
		if got := InitialStars(testCase.rarity); got != testCase.initial {
			t.Errorf("%s 初始星级：期望 %d，实际 %d", testCase.rarity, testCase.initial, got)
		}
		if got := MaxStars(testCase.rarity); got != testCase.max {
			t.Errorf("%s 满星星级：期望 %d，实际 %d", testCase.rarity, testCase.max, got)
		}
	}
}

func TestCanEquipAugment(t *testing.T) {
	if CanEquipAugment("普通", 3) {
		t.Error("普通 3 星还没满星，不该解锁兵装")
	}
	if !CanEquipAugment("普通", 4) {
		t.Error("普通满 4 星应该解锁兵装")
	}
	if !CanEquipAugment("决战方案", 6) {
		t.Error("决战方案满 6 星应该解锁兵装")
	}
	if CanEquipAugment("决战方案", 5) {
		t.Error("决战方案 5 星还没满星，不该解锁兵装")
	}
}
