package ship

import (
	"FZU-DB-ALPortManager/biz/model/api"
)

const (
	// EquipmentSlotCount 是普通装备槽的数量（0..4）。
	EquipmentSlotCount = 5
	// AugmentSlotIndex 是兵装槽的下标，满星后才开放。
	AugmentSlotIndex = 5
)

// equipmentTypeNames 只保留一份枚举 → 数据库字符串的映射；
// 反方向直接用 thrift 生成的 EquipmentTypeFromString，避免两处维护。
var equipmentTypeNames = map[api.EquipmentType]string{
	api.EquipmentType_SMALL_CALIBER_MAIN_GUN:  "SMALL_CALIBER_MAIN_GUN",
	api.EquipmentType_MEDIUM_CALIBER_MAIN_GUN: "MEDIUM_CALIBER_MAIN_GUN",
	api.EquipmentType_LARGE_CALIBER_MAIN_GUN:  "LARGE_CALIBER_MAIN_GUN",
	api.EquipmentType_BATTLESHIP_MAIN_GUN:     "BATTLESHIP_MAIN_GUN",
	api.EquipmentType_TORPEDO:                 "TORPEDO",
	api.EquipmentType_SUBMARINE_TORPEDO:       "SUBMARINE_TORPEDO",
	api.EquipmentType_ANTI_AIR_GUN:            "ANTI_AIR_GUN",
	api.EquipmentType_AUXILIARY:               "AUXILIARY",
	api.EquipmentType_DIVE_BOMBER:             "DIVE_BOMBER",
	api.EquipmentType_TORPEDO_BOMBER:          "TORPEDO_BOMBER",
	api.EquipmentType_FIGHTER:                 "FIGHTER",
	api.EquipmentType_AUGMENT:                 "AUGMENT",
}

// EquipmentTypeName 把枚举换成数据库里的字符串。
func EquipmentTypeName(t api.EquipmentType) (string, error) {
	name, ok := equipmentTypeNames[t]
	if !ok {
		return "", ErrEquipmentTypeUnknown
	}
	return name, nil
}

// ParseEquipmentType 把数据库里的字符串换成枚举。
func ParseEquipmentType(name string) (api.EquipmentType, error) {
	t, err := api.EquipmentTypeFromString(name)
	if err != nil {
		return api.EquipmentType(0), ErrEquipmentTypeUnknown
	}
	return t, nil
}

// IsValidSlot 判断槽位下标是否落在普通槽或兵装槽上。
func IsValidSlot(slot int) bool {
	return slot >= 0 && slot <= AugmentSlotIndex
}

// AcceptableTypes 返回该舰种每个槽位能装的装备类型。
//
// 一律返回 5 个槽（不认识的舰种给空列表），调用方不用再判 nil；
// 兵装槽单独成槽，任何舰种满星后都只收 AUGMENT。
func AcceptableTypes(shipType string) map[int][]api.EquipmentType {
	slots := make(map[int][]api.EquipmentType, EquipmentSlotCount)
	for i := 0; i < EquipmentSlotCount; i++ {
		slots[i] = nil
	}

	switch shipType {
	case "驱逐":
		slots[0] = []api.EquipmentType{api.EquipmentType_SMALL_CALIBER_MAIN_GUN}
		slots[1] = []api.EquipmentType{api.EquipmentType_TORPEDO}
		slots[2] = []api.EquipmentType{api.EquipmentType_ANTI_AIR_GUN}
		slots[3] = []api.EquipmentType{api.EquipmentType_AUXILIARY}
		slots[4] = []api.EquipmentType{api.EquipmentType_AUXILIARY}
	case "轻巡":
		slots[0] = []api.EquipmentType{
			api.EquipmentType_SMALL_CALIBER_MAIN_GUN,
			api.EquipmentType_MEDIUM_CALIBER_MAIN_GUN,
		}
		slots[1] = []api.EquipmentType{
			api.EquipmentType_TORPEDO,
			api.EquipmentType_SMALL_CALIBER_MAIN_GUN,
		}
		slots[2] = []api.EquipmentType{api.EquipmentType_ANTI_AIR_GUN}
		slots[3] = []api.EquipmentType{api.EquipmentType_AUXILIARY}
		slots[4] = []api.EquipmentType{api.EquipmentType_AUXILIARY}
	case "重巡":
		slots[0] = []api.EquipmentType{
			api.EquipmentType_MEDIUM_CALIBER_MAIN_GUN,
			api.EquipmentType_LARGE_CALIBER_MAIN_GUN,
		}
		slots[1] = []api.EquipmentType{
			api.EquipmentType_TORPEDO,
			api.EquipmentType_SMALL_CALIBER_MAIN_GUN,
			api.EquipmentType_MEDIUM_CALIBER_MAIN_GUN,
		}
		slots[2] = []api.EquipmentType{api.EquipmentType_ANTI_AIR_GUN}
		slots[3] = []api.EquipmentType{api.EquipmentType_AUXILIARY}
		slots[4] = []api.EquipmentType{api.EquipmentType_AUXILIARY}
	case "战列", "战巡":
		slots[0] = []api.EquipmentType{api.EquipmentType_BATTLESHIP_MAIN_GUN}
		slots[1] = []api.EquipmentType{
			api.EquipmentType_SMALL_CALIBER_MAIN_GUN,
			api.EquipmentType_MEDIUM_CALIBER_MAIN_GUN,
			api.EquipmentType_LARGE_CALIBER_MAIN_GUN,
		}
		slots[2] = []api.EquipmentType{api.EquipmentType_ANTI_AIR_GUN}
		slots[3] = []api.EquipmentType{api.EquipmentType_AUXILIARY}
		slots[4] = []api.EquipmentType{api.EquipmentType_AUXILIARY}
	case "航母":
		slots[0] = []api.EquipmentType{api.EquipmentType_FIGHTER}
		slots[1] = []api.EquipmentType{api.EquipmentType_DIVE_BOMBER}
		slots[2] = []api.EquipmentType{api.EquipmentType_TORPEDO_BOMBER}
		slots[3] = []api.EquipmentType{api.EquipmentType_AUXILIARY}
		slots[4] = []api.EquipmentType{api.EquipmentType_AUXILIARY}
	case "轻航":
		slots[0] = []api.EquipmentType{
			api.EquipmentType_FIGHTER,
			api.EquipmentType_DIVE_BOMBER,
		}
		slots[1] = []api.EquipmentType{api.EquipmentType_TORPEDO_BOMBER}
		slots[2] = []api.EquipmentType{api.EquipmentType_ANTI_AIR_GUN}
		slots[3] = []api.EquipmentType{api.EquipmentType_AUXILIARY}
		slots[4] = []api.EquipmentType{api.EquipmentType_AUXILIARY}
	case "潜艇":
		slots[0] = []api.EquipmentType{api.EquipmentType_SUBMARINE_TORPEDO}
		slots[1] = []api.EquipmentType{api.EquipmentType_SUBMARINE_TORPEDO}
		slots[2] = []api.EquipmentType{api.EquipmentType_SMALL_CALIBER_MAIN_GUN}
		slots[3] = []api.EquipmentType{api.EquipmentType_AUXILIARY}
		slots[4] = []api.EquipmentType{api.EquipmentType_AUXILIARY}
	case "维修":
		slots[0] = []api.EquipmentType{api.EquipmentType_ANTI_AIR_GUN}
		slots[1] = []api.EquipmentType{api.EquipmentType_ANTI_AIR_GUN}
		slots[2] = []api.EquipmentType{api.EquipmentType_ANTI_AIR_GUN}
		slots[3] = []api.EquipmentType{api.EquipmentType_AUXILIARY}
		slots[4] = []api.EquipmentType{api.EquipmentType_AUXILIARY}
	}

	return slots
}

// SlotAccepts 判断某舰种的某个槽位能否装这类装备。
func SlotAccepts(shipType string, slot int, equipmentType api.EquipmentType) bool {
	if !IsValidSlot(slot) {
		return false
	}
	if slot == AugmentSlotIndex {
		return equipmentType == api.EquipmentType_AUGMENT
	}
	for _, allowed := range AcceptableTypes(shipType)[slot] {
		if allowed == equipmentType {
			return true
		}
	}
	return false
}
