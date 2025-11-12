package convert

import (
	"FZU-DB-ALPortManager/biz/model/api"
	"FZU-DB-ALPortManager/model"
	"errors"
)

func DBToAPIModel(m *model.ShipInfo) *api.Ship {
	if m == nil {
		return nil
	}

	id := int64(m.ID)
	level := int32(m.Level)
	stars := int32(m.Stars)
	createdAt := m.CreatedAt.Format("2006-01-02 15:04:05")
	updatedAt := m.UpdatedAt.Format("2006-01-02 15:04:05")
	shipAPI := &api.Ship{
		ID:        &id,
		ShipName:  m.ShipName,
		Rarity:    m.Rarity,
		ShipType:  m.ShipType,
		Faction:   m.Faction,
		Level:     &level,
		Stars:     &stars,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}

	if m.Equipments != nil {
		acceptableTypesMap := getAcceptableTypesByShipType(m.ShipType)
		shipAPI.EquipmentSlots = make([]*api.EquipmentSlot, 5)
		for i := 0; i < 5; i++ {
			shipAPI.EquipmentSlots[i] = &api.EquipmentSlot{
				AcceptableTypes: acceptableTypesMap[i],
			}
		}
		shipAPI.AugmentSlot = &api.EquipmentSlot{
			AcceptableTypes: []api.EquipmentType{api.EquipmentType_AUGMENT},
		}

		for _, dbEquipSlot := range m.Equipments {
			if dbEquipSlot.EquipmentID == 0 || dbEquipSlot.Equipment.ID == 0 {
				continue
			}

			apiEquipType := stringToEquipmentType(dbEquipSlot.Equipment.Type)
			// 检查类型字符串是否为空（表示未知类型），而不是检查枚举值
			if dbEquipSlot.Equipment.Type == "" {
				continue // 如果类型为空，就跳过这个装备
			}

			apiEquip := &api.Equipment{
				ID:   int64(dbEquipSlot.Equipment.ID),
				Name: dbEquipSlot.Equipment.Name,
				Type: api.EquipmentType(apiEquipType),
			}

			if dbEquipSlot.SlotIndex < 5 {
				shipAPI.EquipmentSlots[dbEquipSlot.SlotIndex].EquippedItem = apiEquip
			} else if dbEquipSlot.SlotIndex == 5 {
				shipAPI.AugmentSlot.EquippedItem = apiEquip
			}
		}
	}

	return shipAPI
}

func DBToAPIModelList(dbList []*model.ShipInfo) []*api.Ship {
	apiList := make([]*api.Ship, len(dbList))
	for i, m := range dbList {
		apiList[i] = DBToAPIModel(m)
	}
	return apiList
}

func APIToDBCreateModel(req *api.Ship) (*model.ShipInfo, error) {
	if req == nil {
		return nil, errors.New("请求不能为空")
	}

	stars := calculateInitialStarsByRarity(req.Rarity)
	dbModel := &model.ShipInfo{
		ShipName: req.ShipName,
		Rarity:   req.Rarity,
		ShipType: req.ShipType,
		Faction:  req.Faction,
		Level:    1, // 新船都是1级
		Stars:    stars,
	}

	// 不初始化装备栏，避免 GORM 自动保存 equipment_id=0 的记录
	// 装备栏会在装备时按需创建

	return dbModel, nil
}

func APIToDBUpdateMap(req *api.Ship) map[string]interface{} {
	updates := make(map[string]interface{})
	if req.ShipName != "" {
		updates["ship_name"] = req.ShipName
	}
	if req.Rarity != "" {
		updates["rarity"] = req.Rarity
		// 如果没有明确指定星级，则根据稀有度自动计算
		if req.Stars == nil {
			stars := calculateStarsByRarity(req.Rarity)
			updates["stars"] = stars
		}
	}
	if req.Faction != "" {
		updates["faction"] = req.Faction
	}
	if req.Level != nil && *req.Level > 0 {
		updates["level"] = *req.Level
	}
	// 支持独立更新星级
	if req.Stars != nil && *req.Stars > 0 {
		updates["stars"] = *req.Stars
	}
	return updates
}

// calculateInitialStarsByRarity 根据稀有度计算初始星级
func calculateInitialStarsByRarity(rarity string) int {
	switch rarity {
	case "普通":
		return 1
	case "稀有", "精锐":
		return 2
	case "超稀有", "海上传奇", "最高方案", "决战方案":
		return 3
	default:
		return 1
	}
}

// calculateStarsByRarity 根据稀有度计算最大星级
func calculateStarsByRarity(rarity string) int {
	switch rarity {
	case "普通":
		return 4
	case "稀有", "精锐":
		return 5
	case "超稀有", "海上传奇", "最高方案", "决战方案":
		return 6
	default:
		return 3
	}
}

func stringToEquipmentType(typeStr string) api.EquipmentType {
	switch typeStr {
	case "SMALL_CALIBER_MAIN_GUN":
		return api.EquipmentType_SMALL_CALIBER_MAIN_GUN
	case "MEDIUM_CALIBER_MAIN_GUN":
		return api.EquipmentType_MEDIUM_CALIBER_MAIN_GUN
	case "LARGE_CALIBER_MAIN_GUN":
		return api.EquipmentType_LARGE_CALIBER_MAIN_GUN
	case "BATTLESHIP_MAIN_GUN":
		return api.EquipmentType_BATTLESHIP_MAIN_GUN
	case "TORPEDO":
		return api.EquipmentType_TORPEDO
	case "SUBMARINE_TORPEDO":
		return api.EquipmentType_SUBMARINE_TORPEDO
	case "ANTI_AIR_GUN":
		return api.EquipmentType_ANTI_AIR_GUN
	case "AUXILIARY":
		return api.EquipmentType_AUXILIARY
	case "FIGHTER":
		return api.EquipmentType_FIGHTER
	case "DIVE_BOMBER":
		return api.EquipmentType_DIVE_BOMBER
	case "TORPEDO_BOMBER":
		return api.EquipmentType_TORPEDO_BOMBER
	case "AUGMENT":
		return api.EquipmentType_AUGMENT
	default:
		return api.EquipmentType(0) // Return zero value for unknown types
	}
}

func getAcceptableTypesByShipType(shipType string) map[int][]api.EquipmentType {
	var (
		auxiliaryX2 = []api.EquipmentType{api.EquipmentType_AUXILIARY}
		antiAirGun  = []api.EquipmentType{api.EquipmentType_ANTI_AIR_GUN}
	)

	// 舰种：驱逐/轻巡/重巡/战列/战巡/航母/轻航/潜艇/维修
	switch shipType {
	case "驱逐":
		return map[int][]api.EquipmentType{
			0: {api.EquipmentType_SMALL_CALIBER_MAIN_GUN}, // 驱逐主炮
			1: {api.EquipmentType_TORPEDO},                // 鱼雷
			2: antiAirGun,                                 // 防空炮
			3: auxiliaryX2,                                // 设备1
			4: auxiliaryX2,                                // 设备2
		}
	case "轻巡":
		return map[int][]api.EquipmentType{
			0: {api.EquipmentType_SMALL_CALIBER_MAIN_GUN, api.EquipmentType_MEDIUM_CALIBER_MAIN_GUN}, // 轻巡主炮
			1: {api.EquipmentType_TORPEDO, api.EquipmentType_SMALL_CALIBER_MAIN_GUN},                 // 副炮/鱼雷
			2: antiAirGun,
			3: auxiliaryX2,
			4: auxiliaryX2,
		}
	case "重巡":
		return map[int][]api.EquipmentType{
			0: {api.EquipmentType_MEDIUM_CALIBER_MAIN_GUN, api.EquipmentType_LARGE_CALIBER_MAIN_GUN},                            // 重巡主炮
			1: {api.EquipmentType_TORPEDO, api.EquipmentType_SMALL_CALIBER_MAIN_GUN, api.EquipmentType_MEDIUM_CALIBER_MAIN_GUN}, // 副炮/鱼雷
			2: antiAirGun,
			3: auxiliaryX2,
			4: auxiliaryX2,
		}
	case "战列", "战巡":
		return map[int][]api.EquipmentType{
			0: {api.EquipmentType_BATTLESHIP_MAIN_GUN},                                                                                         // 战列主炮
			1: {api.EquipmentType_SMALL_CALIBER_MAIN_GUN, api.EquipmentType_MEDIUM_CALIBER_MAIN_GUN, api.EquipmentType_LARGE_CALIBER_MAIN_GUN}, // 副炮
			2: antiAirGun,
			3: auxiliaryX2,
			4: auxiliaryX2,
		}
	case "航母":
		return map[int][]api.EquipmentType{
			0: {api.EquipmentType_FIGHTER},        // 战斗机
			1: {api.EquipmentType_DIVE_BOMBER},    // 轰炸机
			2: {api.EquipmentType_TORPEDO_BOMBER}, // 鱼雷机
			3: auxiliaryX2,
			4: auxiliaryX2,
		}
	case "轻航":
		return map[int][]api.EquipmentType{
			0: {api.EquipmentType_FIGHTER, api.EquipmentType_DIVE_BOMBER}, // 飞机1
			1: {api.EquipmentType_TORPEDO_BOMBER},                         // 飞机2 (通常是鱼雷机)
			2: antiAirGun,
			3: auxiliaryX2,
			4: auxiliaryX2,
		}
	case "潜艇":
		return map[int][]api.EquipmentType{
			0: {api.EquipmentType_SUBMARINE_TORPEDO},      // 潜艇鱼雷1
			1: {api.EquipmentType_SUBMARINE_TORPEDO},      // 潜艇鱼雷2
			2: {api.EquipmentType_SMALL_CALIBER_MAIN_GUN}, // 潜艇炮
			3: auxiliaryX2,
			4: auxiliaryX2,
		}
	case "维修":
		return map[int][]api.EquipmentType{
			0: antiAirGun,
			1: antiAirGun,
			2: antiAirGun,
			3: auxiliaryX2,
			4: auxiliaryX2,
		}
	}
	return map[int][]api.EquipmentType{} // 如果是不认识的舰种，就返回空的规则
}
