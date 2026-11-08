package service

import (
	"context"
	"errors"
	"sort"

	"FZU-DB-ALPortManager/biz/convert"
	"FZU-DB-ALPortManager/biz/dal"
	"FZU-DB-ALPortManager/biz/domain/ship"
	"FZU-DB-ALPortManager/biz/model/api"
	"FZU-DB-ALPortManager/model"
)

// ShipService 业务层：取数 → 过一遍领域规则 → 存回去。
//
// 数据访问全在 biz/dal，这一层不出现 gorm、query，也不用处理 ErrRecordNotFound。
// 装备槽能不能装、星级几颗这些判断在 biz/domain/ship。
type ShipService struct {
	ctx        context.Context
	ships      *dal.Ship
	equipments *dal.Equipment
}

func NewShipService(ctx context.Context) *ShipService {
	return &ShipService{
		ctx:        ctx,
		ships:      dal.NewShip(),
		equipments: dal.NewEquipment(),
	}
}

func (s *ShipService) CreateShip(req *api.Ship) (*api.Ship, error) {
	dbModel, err := convert.APIToDBCreateModel(req)
	if err != nil {
		return nil, err
	}

	if !ship.ValidShipType(req.ShipType) {
		return nil, ship.ErrUnknownShipType
	}

	if err := s.ships.Create(s.ctx, dbModel); err != nil {
		return nil, err
	}

	return s.GetShipByID(int64(dbModel.ID))
}

func (s *ShipService) GetShipByID(id int64) (*api.Ship, error) {
	dbModel, err := s.ships.GetByID(s.ctx, uint(id))
	if err != nil {
		if errors.Is(err, dal.ErrNotFound) {
			return nil, ship.ErrShipNotFound
		}
		return nil, err
	}
	return convert.DBToAPIModel(dbModel), nil
}

func (s *ShipService) GetShipList(req *api.GetShipListReq) (*api.GetShipListResp, error) {
	paging := ship.NormalizePaging(int(req.GetPage()), int(req.GetPageSize()))

	filter := dal.ListFilter{
		ShipName: derefString(req.ShipName),
		Rarity:   derefString(req.Rarity),
		ShipType: derefString(req.ShipType),
		Faction:  derefString(req.Faction),
	}

	dbList, total, err := s.ships.List(s.ctx, filter, paging.Offset, paging.PageSize)
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

	if req.ShipType != "" && !ship.ValidShipType(req.ShipType) {
		return nil, ship.ErrUnknownShipType
	}

	if err := s.ships.Update(s.ctx, uint(id), convert.APIToDBUpdateMap(req)); err != nil {
		return nil, err
	}

	return s.GetShipByID(id)
}

func (s *ShipService) DeleteShip(id int64) error {
	if _, err := s.GetShipByID(id); err != nil {
		return err
	}

	return s.ships.Delete(s.ctx, uint(id))
}

func (s *ShipService) GetStatistics() (*api.GetStatisticsResp, error) {
	total, err := s.ships.Count(s.ctx)
	if err != nil {
		return nil, err
	}

	factions, err := s.ships.CountBy(s.ctx, dal.StatByFaction)
	if err != nil {
		return nil, err
	}

	rarities, err := s.ships.CountBy(s.ctx, dal.StatByRarity)
	if err != nil {
		return nil, err
	}

	shipTypes, err := s.ships.CountBy(s.ctx, dal.StatByShipType)
	if err != nil {
		return nil, err
	}

	return &api.GetStatisticsResp{
		TotalShips:    total,
		FactionStats:  toStatItems(factions),
		RarityStats:   toStatItems(rarities),
		ShipTypeStats: toStatItems(shipTypes),
	}, nil
}

// EquipShip 装/卸装备。
//
// 装之前先过三道规则：槽位是否合法、装备类型是否匹配舰种、兵装槽是否已满星。
func (s *ShipService) EquipShip(req *api.EquipShipReq) (*api.Ship, error) {
	shipID := uint(req.ShipId)
	slotIndex := int(req.SlotIndex)
	equipmentID := uint(req.GetEquipmentId())

	if !ship.IsValidSlot(slotIndex) {
		return nil, ship.ErrSlotOutOfRange
	}

	dbShip, err := s.ships.GetByID(s.ctx, shipID)
	if err != nil {
		if errors.Is(err, dal.ErrNotFound) {
			return nil, ship.ErrShipNotFound
		}
		return nil, err
	}

	// equipmentID 为 0 表示卸下这一格
	if equipmentID == 0 {
		if err := s.equipments.ClearSlot(s.ctx, shipID, slotIndex); err != nil {
			return nil, err
		}
		return s.GetShipByID(req.ShipId)
	}

	if err := s.checkSlotRules(dbShip, slotIndex, equipmentID); err != nil {
		return nil, err
	}

	if err := s.equipments.UpsertSlot(s.ctx, shipID, slotIndex, equipmentID); err != nil {
		return nil, err
	}

	return s.GetShipByID(req.ShipId)
}

func (s *ShipService) checkSlotRules(dbShip *model.ShipInfo, slotIndex int, equipmentID uint) error {
	equipment, err := s.equipments.GetByID(s.ctx, equipmentID)
	if err != nil {
		if errors.Is(err, dal.ErrNotFound) {
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

func (s *ShipService) GetEquipmentList(req *api.GetEquipmentListReq) (*api.GetEquipmentListResp, error) {
	var typeName string
	if req.Type != nil {
		name, err := ship.EquipmentTypeName(*req.Type)
		if err != nil {
			return nil, err
		}
		typeName = name
	}

	dbList, err := s.equipments.List(s.ctx, typeName)
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

// toStatItems 转成接口结构，顺带按「数量倒序、同数量按名称升序」排一遍。
// 排序放在 Go 里而不是只信 SQL：不同数据库对中文的排序规则不一样，
// 同一份数据在 PostgreSQL 和 openGauss 上顺序可能不同，界面就会闪。
func toStatItems(rows []dal.GroupCount) []*api.StatItem {
	items := make([]*api.StatItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, &api.StatItem{Name: row.Name, Count: row.Count})
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Count != items[j].Count {
			return items[i].Count > items[j].Count
		}
		return items[i].Name < items[j].Name
	})

	return items
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
