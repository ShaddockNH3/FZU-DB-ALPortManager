package convert

import (
	"testing"

	"FZU-DB-ALPortManager/biz/domain/ship"
	"FZU-DB-ALPortManager/biz/model/api"
	"FZU-DB-ALPortManager/model"
)

func TestDBToAPIModelFillsSlots(t *testing.T) {
	dbShip := &model.ShipInfo{
		ID:       1,
		ShipName: "企业",
		Rarity:   "超稀有",
		ShipType: "航母",
		Faction:  "白鹰",
		Level:    100,
		Stars:    6,
		Equipments: []model.ShipEquipment{
			{
				SlotIndex:   0,
				EquipmentID: 10,
				Equipment:   model.Equipment{ID: 10, Name: "F6F 地狱猫", Type: "FIGHTER"},
			},
			{
				SlotIndex:   ship.AugmentSlotIndex,
				EquipmentID: 11,
				Equipment:   model.Equipment{ID: 11, Name: "航空兵装", Type: "AUGMENT"},
			},
		},
	}

	got := DBToAPIModel(dbShip)

	if len(got.EquipmentSlots) != ship.EquipmentSlotCount {
		t.Fatalf("槽位数量应该是 %d，实际 %d", ship.EquipmentSlotCount, len(got.EquipmentSlots))
	}
	if got.EquipmentSlots[0].EquippedItem == nil || got.EquipmentSlots[0].EquippedItem.Name != "F6F 地狱猫" {
		t.Error("0 号槽应该装上战斗机")
	}
	if got.AugmentSlot == nil || got.AugmentSlot.EquippedItem == nil {
		t.Error("兵装槽应该装上兵装")
	}
	if len(got.EquipmentSlots[1].AcceptableTypes) == 0 {
		t.Error("航母 1 号槽应该有可装类型")
	}
}

func TestDBToAPIModelSkipsUnknownEquipmentType(t *testing.T) {
	dbShip := &model.ShipInfo{
		ID:       2,
		ShipName: "未知装备船",
		ShipType: "驱逐",
		Equipments: []model.ShipEquipment{
			{SlotIndex: 0, EquipmentID: 20, Equipment: model.Equipment{ID: 20, Name: "坏数据", Type: ""}},
			{SlotIndex: 1, EquipmentID: 0, Equipment: model.Equipment{ID: 0}},
		},
	}

	got := DBToAPIModel(dbShip)

	for index, slot := range got.EquipmentSlots {
		if slot.EquippedItem != nil {
			t.Errorf("%d 号槽不该有装备", index)
		}
	}
	if got.AugmentSlot.EquippedItem != nil {
		t.Error("兵装槽不该有装备")
	}
}

func TestAPIToDBCreateModel(t *testing.T) {
	dbModel, err := APIToDBCreateModel(&api.Ship{
		ShipName: " 企业 ",
		Rarity:   "超稀有",
		ShipType: "航母",
		Faction:  "白鹰",
	})
	if err != nil {
		t.Fatalf("不该报错: %v", err)
	}
	if dbModel.ShipName != "企业" {
		t.Errorf("名字应该去掉首尾空格，实际 %q", dbModel.ShipName)
	}
	if dbModel.Level != 1 {
		t.Errorf("新船应该 1 级，实际 %d", dbModel.Level)
	}
	if dbModel.Stars != ship.InitialStars("超稀有") {
		t.Errorf("星级应该按稀有度给，实际 %d", dbModel.Stars)
	}

	if _, err := APIToDBCreateModel(&api.Ship{ShipName: "   "}); err == nil {
		t.Error("空名字应该报错")
	}
}

func TestAPIToDBUpdateMapSkipsEmptyFields(t *testing.T) {
	level := int32(50)
	updates := APIToDBUpdateMap(&api.Ship{ShipName: "", Level: &level})

	if _, ok := updates["ship_name"]; ok {
		t.Error("没传名字就不该更新名字")
	}
	if updates["level"] != int32(50) {
		t.Error("传了等级应该更新等级")
	}
}
