package convert

import (
	"FZU-DB-ALPortManager/biz/domain/ship"
	"FZU-DB-ALPortManager/biz/model/api"
	"FZU-DB-ALPortManager/model"
	"strings"
)

// 装备类型字符串 ↔ 枚举、槽位规则、星级规则都搬到了 biz/domain/ship，
// 这里只做"数据库行 ↔ API 结构"的形状转换。

// DBToAPIModel 把一行舰船（含预加载的装备）转成接口结构。
func DBToAPIModel(m *model.ShipInfo) *api.Ship {
	if m == nil {
		return nil
	}

	id := int64(m.ID)
	level := int32(m.Level)
	stars := int32(m.Stars)
	createdAt := m.CreatedAt.Format("2006-01-02 15:04:05")
	updatedAt := m.UpdatedAt.Format("2006-01-02 15:04:05")

	slots, augmentSlot := buildSlots(m)

	return &api.Ship{
		ID:             &id,
		ShipName:       m.ShipName,
		Rarity:         m.Rarity,
		ShipType:       m.ShipType,
		Faction:        m.Faction,
		Level:          &level,
		Stars:          &stars,
		CreatedAt:      &createdAt,
		UpdatedAt:      &updatedAt,
		EquipmentSlots: slots,
		AugmentSlot:    augmentSlot,
	}
}

// buildSlots 把已装装备填进 5 个普通槽和 1 个兵装槽。
// 槽位列表总是齐的，未知装备类型直接跳过，不往接口里塞空对象。
func buildSlots(m *model.ShipInfo) ([]*api.EquipmentSlot, *api.EquipmentSlot) {
	acceptable := ship.AcceptableTypes(m.ShipType)

	slots := make([]*api.EquipmentSlot, ship.EquipmentSlotCount)
	for index := range slots {
		slots[index] = &api.EquipmentSlot{AcceptableTypes: acceptable[index]}
	}
	augmentSlot := &api.EquipmentSlot{AcceptableTypes: []api.EquipmentType{api.EquipmentType_AUGMENT}}

	for _, equipped := range m.Equipments {
		if equipped.EquipmentID == 0 || equipped.Equipment.ID == 0 {
			continue
		}
		equipmentType, err := ship.ParseEquipmentType(equipped.Equipment.Type)
		if err != nil {
			continue
		}

		item := &api.Equipment{
			ID:   int64(equipped.Equipment.ID),
			Name: equipped.Equipment.Name,
			Type: equipmentType,
		}

		switch {
		case equipped.SlotIndex >= 0 && equipped.SlotIndex < ship.EquipmentSlotCount:
			slots[equipped.SlotIndex].EquippedItem = item
		case equipped.SlotIndex == ship.AugmentSlotIndex:
			augmentSlot.EquippedItem = item
		}
	}

	return slots, augmentSlot
}

// DBToAPIModelList 批量转换。
func DBToAPIModelList(dbList []*model.ShipInfo) []*api.Ship {
	apiList := make([]*api.Ship, len(dbList))
	for i, m := range dbList {
		apiList[i] = DBToAPIModel(m)
	}
	return apiList
}

// APIToDBCreateModel 新建舰船：等级固定 1 级，星级按稀有度给，不预先建装备槽。
func APIToDBCreateModel(req *api.Ship) (*model.ShipInfo, error) {
	if req == nil {
		return nil, ship.ErrEmptyShipName
	}
	name := strings.TrimSpace(req.ShipName)
	if name == "" {
		return nil, ship.ErrEmptyShipName
	}

	return &model.ShipInfo{
		ShipName: name,
		Rarity:   req.Rarity,
		ShipType: req.ShipType,
		Faction:  req.Faction,
		Level:    1,
		Stars:    ship.InitialStars(req.Rarity),
	}, nil
}

// APIToDBUpdateMap 只更新明确传了值的字段，空字符串视为"不改"。
func APIToDBUpdateMap(req *api.Ship) map[string]interface{} {
	updates := make(map[string]interface{})
	if name := strings.TrimSpace(req.ShipName); name != "" {
		updates["ship_name"] = name
	}
	if req.Rarity != "" {
		updates["rarity"] = req.Rarity
	}
	// 舰种此前漏了：前端能改，接口返回 200，库里没动。一并补上。
	if req.ShipType != "" {
		updates["ship_type"] = req.ShipType
	}
	if req.Faction != "" {
		updates["faction"] = req.Faction
	}
	if req.Level != nil && *req.Level > 0 {
		updates["level"] = *req.Level
	}
	if req.Stars != nil && *req.Stars > 0 {
		updates["stars"] = *req.Stars
	}
	return updates
}
