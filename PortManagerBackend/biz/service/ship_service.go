package service

import (
	"context"
	"errors"
	"sort"

	"FZU-DB-ALPortManager/biz/convert"
	"FZU-DB-ALPortManager/biz/domain/ship"
	"FZU-DB-ALPortManager/biz/model/api"
	"FZU-DB-ALPortManager/model"
	"FZU-DB-ALPortManager/query"

	"gorm.io/gorm"
)

// ShipService 是用例层：只做"取数据 → 过一遍领域规则 → 存回去"，
// 装备槽能不能装、星级几颗这些判断都在 biz/domain/ship。
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

	// 只创建舰船本体，装备槽记录在装备时按需创建
	if err := query.ShipInfo.WithContext(s.ctx).Create(dbModel); err != nil {
		return nil, err
	}

	return s.GetShipByID(int64(dbModel.ID))
}

func (s *ShipService) GetShipByID(id int64) (*api.Ship, error) {
	q := query.ShipInfo
	dbModel, err := q.WithContext(s.ctx).
		Preload(q.Equipments.Equipment).
		Where(q.ID.Eq(uint(id))).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ship.ErrShipNotFound
		}
		return nil, err
	}
	return convert.DBToAPIModel(dbModel), nil
}

func (s *ShipService) GetShipList(req *api.GetShipListReq) (*api.GetShipListResp, error) {
	q := query.ShipInfo.WithContext(s.ctx)

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

	paging := ship.NormalizePaging(int(req.GetPage()), int(req.GetPageSize()))

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	dbList, err := q.Offset(paging.Offset).Limit(paging.PageSize).Order(query.ShipInfo.ID.Desc()).Find()
	if err != nil {
		return nil, err
	}

	return &api.GetShipListResp{
		Total:    total,
		Ships:    convert.DBToAPIModelList(dbList),
		Page:     int32(paging.Page),
		PageSize: int32(paging.PageSize),
	}, nil
}

func (s *ShipService) UpdateShip(id int64, req *api.Ship) (*api.Ship, error) {
	if _, err := s.GetShipByID(id); err != nil {
		return nil, err
	}

	updates := convert.APIToDBUpdateMap(req)
	if len(updates) == 0 {
		return s.GetShipByID(id)
	}

	if _, err := query.ShipInfo.WithContext(s.ctx).
		Where(query.ShipInfo.ID.Eq(uint(id))).
		Updates(updates); err != nil {
		return nil, err
	}

	return s.GetShipByID(id)
}

func (s *ShipService) DeleteShip(id int64) error {
	shipID := uint(id)

	return query.Q.Transaction(func(tx *query.Query) error {
		if _, err := tx.ShipEquipment.WithContext(s.ctx).
			Where(query.ShipEquipment.ShipID.Eq(shipID)).
			Delete(); err != nil {
			return err
		}

		result, err := tx.ShipInfo.WithContext(s.ctx).
			Where(query.ShipInfo.ID.Eq(shipID)).
			Delete()
		if err != nil {
			return err
		}
		if result.RowsAffected == 0 {
			return ship.ErrShipNotFound
		}
		return nil
	})
}

func (s *ShipService) GetStatistics() (*api.GetStatisticsResp, error) {
	ships, err := query.ShipInfo.WithContext(s.ctx).Find()
	if err != nil {
		return nil, err
	}

	return &api.GetStatisticsResp{
		TotalShips:    int64(len(ships)),
		FactionStats:  countBy(ships, func(m *model.ShipInfo) string { return m.Faction }),
		RarityStats:   countBy(ships, func(m *model.ShipInfo) string { return m.Rarity }),
		ShipTypeStats: countBy(ships, func(m *model.ShipInfo) string { return m.ShipType }),
	}, nil
}

// EquipShip 装/卸装备。
//
// 装之前先过三道规则：槽位是否合法、装备类型是否匹配舰种、兵装槽是否已满星。
// 原来这些判断一个都没有，前端传什么就存什么。
func (s *ShipService) EquipShip(req *api.EquipShipReq) (*api.Ship, error) {
	shipID := uint(req.ShipId)
	slotIndex := int(req.SlotIndex)
	equipmentID := uint(req.GetEquipmentId())

	if !ship.IsValidSlot(slotIndex) {
		return nil, ship.ErrSlotOutOfRange
	}

	dbShip, err := query.ShipInfo.WithContext(s.ctx).
		Where(query.ShipInfo.ID.Eq(shipID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ship.ErrShipNotFound
		}
		return nil, err
	}

	// equipmentID 为 0 表示卸下这一格
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

	if err := s.checkSlotRules(dbShip, slotIndex, equipmentID); err != nil {
		return nil, err
	}

	if err := s.upsertEquipment(shipID, slotIndex, equipmentID); err != nil {
		return nil, err
	}

	return s.GetShipByID(req.ShipId)
}

func (s *ShipService) checkSlotRules(dbShip *model.ShipInfo, slotIndex int, equipmentID uint) error {
	equipment, err := query.Equipment.WithContext(s.ctx).
		Where(query.Equipment.ID.Eq(equipmentID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ship.ErrEquipmentNotFound
		}
		return err
	}

	equipmentType, err := ship.ParseEquipmentType(equipment.Type)
	if err != nil {
		return err
	}

	if !ship.SlotAccepts(dbShip.ShipType, slotIndex, equipmentType) {
		return ship.ErrEquipmentTypeNotAllowed
	}

	if slotIndex == ship.AugmentSlotIndex && !ship.CanEquipAugment(dbShip.Rarity, dbShip.Stars) {
		return ship.ErrAugmentLocked
	}

	return nil
}

func (s *ShipService) upsertEquipment(shipID uint, slotIndex int, equipmentID uint) error {
	_, err := query.ShipEquipment.WithContext(s.ctx).
		Where(query.ShipEquipment.ShipID.Eq(shipID)).
		Where(query.ShipEquipment.SlotIndex.Eq(slotIndex)).
		First()

	switch {
	case err == nil:
		_, err = query.ShipEquipment.WithContext(s.ctx).
			Where(query.ShipEquipment.ShipID.Eq(shipID)).
			Where(query.ShipEquipment.SlotIndex.Eq(slotIndex)).
			Update(query.ShipEquipment.EquipmentID, equipmentID)
		return err
	case errors.Is(err, gorm.ErrRecordNotFound):
		return query.ShipEquipment.WithContext(s.ctx).Create(&model.ShipEquipment{
			ShipID:      shipID,
			SlotIndex:   slotIndex,
			EquipmentID: equipmentID,
		})
	default:
		return err
	}
}

func (s *ShipService) GetEquipmentList(req *api.GetEquipmentListReq) (*api.GetEquipmentListResp, error) {
	q := query.Equipment.WithContext(s.ctx)

	if req.Type != nil {
		if typeName, err := ship.EquipmentTypeName(*req.Type); err == nil {
			q = q.Where(query.Equipment.Type.Eq(typeName))
		}
	}

	dbList, err := q.Order(query.Equipment.ID).Find()
	if err != nil {
		return nil, err
	}

	equipments := make([]*api.Equipment, 0, len(dbList))
	for _, item := range dbList {
		equipmentType, err := ship.ParseEquipmentType(item.Type)
		if err != nil {
			continue
		}
		equipments = append(equipments, &api.Equipment{
			ID:   int64(item.ID),
			Name: item.Name,
			Type: equipmentType,
		})
	}

	return &api.GetEquipmentListResp{Equipments: equipments}, nil
}

// countBy 按某个字段分组计数，按数量倒序、同数量按名称升序，界面不会每次刷新换顺序。
func countBy(ships []*model.ShipInfo, pick func(*model.ShipInfo) string) []*api.StatItem {
	counts := make(map[string]int64, len(ships))
	for _, item := range ships {
		counts[pick(item)]++
	}

	stats := make([]*api.StatItem, 0, len(counts))
	for name, count := range counts {
		stats = append(stats, &api.StatItem{Name: name, Count: count})
	}
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Count != stats[j].Count {
			return stats[i].Count > stats[j].Count
		}
		return stats[i].Name < stats[j].Name
	})
	return stats
}
