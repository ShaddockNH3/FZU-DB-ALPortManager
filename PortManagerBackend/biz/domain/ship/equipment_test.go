package ship

import (
	"testing"

	"FZU-DB-ALPortManager/biz/model/api"
)

func TestEquipmentTypeNameRoundTrip(t *testing.T) {
	types := []api.EquipmentType{
		api.EquipmentType_SMALL_CALIBER_MAIN_GUN,
		api.EquipmentType_MEDIUM_CALIBER_MAIN_GUN,
		api.EquipmentType_LARGE_CALIBER_MAIN_GUN,
		api.EquipmentType_BATTLESHIP_MAIN_GUN,
		api.EquipmentType_TORPEDO,
		api.EquipmentType_SUBMARINE_TORPEDO,
		api.EquipmentType_ANTI_AIR_GUN,
		api.EquipmentType_AUXILIARY,
		api.EquipmentType_DIVE_BOMBER,
		api.EquipmentType_TORPEDO_BOMBER,
		api.EquipmentType_FIGHTER,
		api.EquipmentType_AUGMENT,
	}

	for _, equipmentType := range types {
		name, err := EquipmentTypeName(equipmentType)
		if err != nil {
			t.Fatalf("EquipmentTypeName(%d) 失败: %v", equipmentType, err)
		}

		parsed, err := ParseEquipmentType(name)
		if err != nil {
			t.Fatalf("ParseEquipmentType(%q) 失败: %v", name, err)
		}
		if parsed != equipmentType {
			t.Errorf("往返不一致: %d → %q → %d", equipmentType, name, parsed)
		}
	}
}

func TestEquipmentTypeUnknown(t *testing.T) {
	if _, err := EquipmentTypeName(api.EquipmentType(999)); err == nil {
		t.Error("未知枚举应该报错")
	}
	if _, err := ParseEquipmentType("NOT_A_REAL_TYPE"); err == nil {
		t.Error("未知字符串应该报错")
	}
}

func TestAcceptableTypesAlwaysHasFiveSlots(t *testing.T) {
	for _, shipType := range []string{"驱逐", "轻巡", "重巡", "战列", "战巡", "航母", "轻航", "潜艇", "维修", "不存在的舰种"} {
		slots := AcceptableTypes(shipType)
		for index := 0; index < EquipmentSlotCount; index++ {
			if _, ok := slots[index]; !ok {
				t.Errorf("%s 的 %d 号槽缺失", shipType, index)
			}
		}
	}

	if len(AcceptableTypes("不存在的舰种")[0]) != 0 {
		t.Error("未知舰种不该给出可装类型")
	}
}

func TestSlotAccepts(t *testing.T) {
	cases := []struct {
		name      string
		shipType  string
		slot      int
		equipment api.EquipmentType
		want      bool
	}{
		{"驱逐 0 号槽收小口径主炮", "驱逐", 0, api.EquipmentType_SMALL_CALIBER_MAIN_GUN, true},
		{"驱逐 0 号槽不收战列主炮", "驱逐", 0, api.EquipmentType_BATTLESHIP_MAIN_GUN, false},
		{"航母 0 号槽只收战斗机", "航母", 0, api.EquipmentType_FIGHTER, true},
		{"航母 0 号槽不收轰炸机", "航母", 0, api.EquipmentType_DIVE_BOMBER, false},
		{"维修 3 号槽收设备", "维修", 3, api.EquipmentType_AUXILIARY, true},
		{"战列 1 号槽收大口径主炮当副炮", "战列", 1, api.EquipmentType_LARGE_CALIBER_MAIN_GUN, true},
		{"兵装槽只收兵装", "驱逐", AugmentSlotIndex, api.EquipmentType_AUGMENT, true},
		{"兵装槽不收防空炮", "驱逐", AugmentSlotIndex, api.EquipmentType_ANTI_AIR_GUN, false},
		{"越界槽位一律拒绝", "驱逐", 6, api.EquipmentType_AUGMENT, false},
		{"负数槽位一律拒绝", "驱逐", -1, api.EquipmentType_AUGMENT, false},
	}

	for _, testCase := range cases {
		got := SlotAccepts(testCase.shipType, testCase.slot, testCase.equipment)
		if got != testCase.want {
			t.Errorf("%s: 期望 %v，实际 %v", testCase.name, testCase.want, got)
		}
	}
}

func TestIsValidSlot(t *testing.T) {
	for slot := 0; slot <= AugmentSlotIndex; slot++ {
		if !IsValidSlot(slot) {
			t.Errorf("%d 应该是合法槽位", slot)
		}
	}
	if IsValidSlot(AugmentSlotIndex+1) || IsValidSlot(-1) {
		t.Error("越界槽位不算合法")
	}
}

// ShipTypes 和 AcceptableTypes 的 switch 是两份东西，容易改一处漏一处。
// 这条用例保证列表里的每个舰种都真能算出槽位。
func TestShipTypesMatchAcceptableTypes(t *testing.T) {
	for _, shipType := range ShipTypes {
		slots := AcceptableTypes(shipType)
		if len(slots[0]) == 0 {
			t.Errorf("%s 在 ShipTypes 里，但 AcceptableTypes 没给它槽位", shipType)
		}
		if !ValidShipType(shipType) {
			t.Errorf("%s 应该是合法舰种", shipType)
		}
	}

	if ValidShipType("不存在的舰种") {
		t.Error("未知舰种不该通过校验")
	}
	if ValidShipType("") {
		t.Error("空舰种不该通过校验")
	}
}
