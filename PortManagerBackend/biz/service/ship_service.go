package service

import (
	"FZU-DB-ALPortManager/biz/convert"
	"FZU-DB-ALPortManager/biz/model/api"
	"FZU-DB-ALPortManager/model"
	"FZU-DB-ALPortManager/query"
	"context"
)

type ShipService struct {
	ctx context.Context
}

func NewShipService(ctx context.Context) *ShipService {
	return &ShipService{ctx: ctx}
}

func (s *ShipService) CreateShip(req *api.Ship) (*api.Ship, error) {
	dbModel, err := convert.APIToDBCreateModel(req)
	if err != nil {
		return nil, err
	}

	// 只创建舰船本体，不创建空的装备栏记录
	// 装备栏记录会在装备时按需创建
	if err := query.ShipInfo.WithContext(s.ctx).Create(dbModel); err != nil {
		return nil, err
	}

	return s.GetShipByID(int64(dbModel.ID))
}

func (s *ShipService) GetShipByID(id int64) (*api.Ship, error) {
	q := query.ShipInfo
	dbModel, err := q.WithContext(s.ctx).
		Preload(q.Equipments.Equipment). // 这个 Preload 魔法是对的！
		Where(q.ID.Eq(uint(id))).
		First()

	if err != nil {
		return nil, err
	}
	return convert.DBToAPIModel(dbModel), nil
}

func (s *ShipService) GetShipList(req *api.GetShipListReq) (*api.GetShipListResp, error) {
	q := query.ShipInfo.WithContext(s.ctx)

	// 筛选条�?
	if req.ShipName != nil && *req.ShipName != "" {
		q = q.Where(query.ShipInfo.ShipName.Like("%" + *req.ShipName + "%"))
	}
	if req.Rarity != nil && *req.Rarity != "" {
		q = q.Where(query.ShipInfo.Rarity.Eq(*req.Rarity))
	}
	if req.ShipType != nil && *req.ShipType != "" {
		q = q.Where(query.ShipInfo.ShipType.Eq(*req.ShipType))
	}
	if req.Faction != nil && *req.Faction != "" {
		q = q.Where(query.ShipInfo.Faction.Eq(*req.Faction))
	}

	q = q.Preload(query.ShipInfo.Equipments.Equipment)

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	offset := (page - 1) * pageSize
	total, _ := q.Count()
	dbList, _ := q.Offset(offset).Limit(pageSize).Order(query.ShipInfo.ID.Desc()).Find()

	return &api.GetShipListResp{
		Total:    total,
		Ships:    convert.DBToAPIModelList(dbList),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

func (s *ShipService) UpdateShip(id int64, req *api.Ship) (*api.Ship, error) {
	if _, err := query.ShipInfo.WithContext(s.ctx).Where(query.ShipInfo.ID.Eq(uint(id))).First(); err != nil {
		return nil, err
	}

	updatesMap := convert.APIToDBUpdateMap(req)

	if _, err := query.ShipInfo.WithContext(s.ctx).Where(query.ShipInfo.ID.Eq(uint(id))).Updates(updatesMap); err != nil {
		return nil, err
	}

	return s.GetShipByID(id)
}

func (s *ShipService) DeleteShip(id int64) error {
	shipID := uint(id)
	return query.Q.Transaction(func(tx *query.Query) error {
		if _, err := tx.ShipEquipment.WithContext(s.ctx).Where(query.ShipEquipment.ShipID.Eq(shipID)).Delete(); err != nil {
			return err
		}
		if _, err := tx.ShipInfo.WithContext(s.ctx).Where(query.ShipInfo.ID.Eq(shipID)).Delete(); err != nil {
			return err
		}
		return nil
	})
}

func (s *ShipService) GetStatistics() (*api.GetStatisticsResp, error) {
	q := query.ShipInfo.WithContext(s.ctx)

	// 统计总数
	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	// 获取所有舰�?
	ships, err := q.Find()
	if err != nil {
		return nil, err
	}

	// 统计阵营
	factionMap := make(map[string]int64)
	rarityMap := make(map[string]int64)
	shipTypeMap := make(map[string]int64)

	for _, ship := range ships {
		factionMap[ship.Faction]++
		rarityMap[ship.Rarity]++
		shipTypeMap[ship.ShipType]++
	}

	// 转换为StatItem数组
	factionStats := make([]*api.StatItem, 0, len(factionMap))
	for name, count := range factionMap {
		factionStats = append(factionStats, &api.StatItem{
			Name:  name,
			Count: count,
		})
	}

	rarityStats := make([]*api.StatItem, 0, len(rarityMap))
	for name, count := range rarityMap {
		rarityStats = append(rarityStats, &api.StatItem{
			Name:  name,
			Count: count,
		})
	}

	shipTypeStats := make([]*api.StatItem, 0, len(shipTypeMap))
	for name, count := range shipTypeMap {
		shipTypeStats = append(shipTypeStats, &api.StatItem{
			Name:  name,
			Count: count,
		})
	}

	return &api.GetStatisticsResp{
		TotalShips:    total,
		FactionStats:  factionStats,
		RarityStats:   rarityStats,
		ShipTypeStats: shipTypeStats,
	}, nil
}

func (s *ShipService) EquipShip(req *api.EquipShipReq) (*api.Ship, error) {
	shipID := uint(req.ShipId)
	slotIndex := int(req.SlotIndex)
	equipmentID := uint(req.GetEquipmentId())

	// 如果equipmentID�?，表示卸下装备，删除记录
	if equipmentID == 0 {
		_, err := query.ShipEquipment.WithContext(s.ctx).
			Where(query.ShipEquipment.ShipID.Eq(shipID)).
			Where(query.ShipEquipment.SlotIndex.Eq(slotIndex)).
			Delete()
		if err != nil {
			return nil, err
		}
		return s.GetShipByID(req.ShipId)
	}

	// 检查该装备栏是否存�?
	_, err := query.ShipEquipment.WithContext(s.ctx).
		Where(query.ShipEquipment.ShipID.Eq(shipID)).
		Where(query.ShipEquipment.SlotIndex.Eq(slotIndex)).
		First()

	if err == nil {
		// 记录存在，更新装备ID
		_, err = query.ShipEquipment.WithContext(s.ctx).
			Where(query.ShipEquipment.ShipID.Eq(shipID)).
			Where(query.ShipEquipment.SlotIndex.Eq(slotIndex)).
			Update(query.ShipEquipment.EquipmentID, equipmentID)
		if err != nil {
			return nil, err
		}
	} else {
		// 记录不存在，创建新记�?
		err = query.ShipEquipment.WithContext(s.ctx).Create(&model.ShipEquipment{
			ShipID:      shipID,
			SlotIndex:   slotIndex,
			EquipmentID: equipmentID,
		})
		if err != nil {
			return nil, err
		}
	}

	return s.GetShipByID(req.ShipId)
}

func (s *ShipService) GetEquipmentList(req *api.GetEquipmentListReq) (*api.GetEquipmentListResp, error) {
	q := query.Equipment.WithContext(s.ctx)

	if req.Type != nil {
		typeStr := equipmentTypeToString(*req.Type)
		if typeStr != "" {
			q = q.Where(query.Equipment.Type.Eq(typeStr))
		}
	}

	dbList, err := q.Order(query.Equipment.ID).Find()
	if err != nil {
		return nil, err
	}

	equipments := make([]*api.Equipment, len(dbList))
	for i, e := range dbList {
		equipments[i] = &api.Equipment{
			ID:   int64(e.ID),
			Name: e.Name,
			Type: stringToEquipmentTypeAPI(e.Type),
		}
	}

	return &api.GetEquipmentListResp{
		Equipments: equipments,
	}, nil
}

// 辅助函数：将API的EquipmentType转换为数据库的字符串
func equipmentTypeToString(t api.EquipmentType) string {
	switch t {
	case api.EquipmentType_SMALL_CALIBER_MAIN_GUN:
		return "SMALL_CALIBER_MAIN_GUN"
	case api.EquipmentType_MEDIUM_CALIBER_MAIN_GUN:
		return "MEDIUM_CALIBER_MAIN_GUN"
	case api.EquipmentType_LARGE_CALIBER_MAIN_GUN:
		return "LARGE_CALIBER_MAIN_GUN"
	case api.EquipmentType_BATTLESHIP_MAIN_GUN:
		return "BATTLESHIP_MAIN_GUN"
	case api.EquipmentType_TORPEDO:
		return "TORPEDO"
	case api.EquipmentType_SUBMARINE_TORPEDO:
		return "SUBMARINE_TORPEDO"
	case api.EquipmentType_ANTI_AIR_GUN:
		return "ANTI_AIR_GUN"
	case api.EquipmentType_AUXILIARY:
		return "AUXILIARY"
	case api.EquipmentType_DIVE_BOMBER:
		return "DIVE_BOMBER"
	case api.EquipmentType_TORPEDO_BOMBER:
		return "TORPEDO_BOMBER"
	case api.EquipmentType_FIGHTER:
		return "FIGHTER"
	case api.EquipmentType_AUGMENT:
		return "AUGMENT"
	default:
		return ""
	}
}

// 辅助函数：将数据库的字符串转换为API的EquipmentType
func stringToEquipmentTypeAPI(typeStr string) api.EquipmentType {
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
	case "DIVE_BOMBER":
		return api.EquipmentType_DIVE_BOMBER
	case "TORPEDO_BOMBER":
		return api.EquipmentType_TORPEDO_BOMBER
	case "FIGHTER":
		return api.EquipmentType_FIGHTER
	case "AUGMENT":
		return api.EquipmentType_AUGMENT
	default:
		return api.EquipmentType(0)
	}
}
