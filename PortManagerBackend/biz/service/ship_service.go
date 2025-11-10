package service

import (
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
	dbModel := ToDBModel(req)

	err := query.ShipInfo.WithContext(s.ctx).Create(dbModel)
	if err != nil {
		return nil, err
	}

	return ToAPIModel(dbModel), nil
}

func (s *ShipService) GetShipByID(id int64) (*api.Ship, error) {
	dbModel, err := query.ShipInfo.WithContext(s.ctx).
		Where(query.ShipInfo.ID.Eq(uint(id))).
		First()

	if err != nil {
		return nil, err
	}

	return ToAPIModel(dbModel), nil
}

func (s *ShipService) GetShipList(req *api.GetShipListReq) (*api.GetShipListResp, error) {
	q := query.ShipInfo.WithContext(s.ctx)

	if req.ShipName != nil && *req.ShipName != "" {
		q = q.Where(query.ShipInfo.ShipName.Like("%" + *req.ShipName + "%"))
	}
	if req.Faction != nil && *req.Faction != "" {
		q = q.Where(query.ShipInfo.Faction.Eq(*req.Faction))
	}
	if req.Rarity != nil && *req.Rarity != "" {
		q = q.Where(query.ShipInfo.Rarity.Eq(*req.Rarity))
	}
	if req.ShipType != nil && *req.ShipType != "" {
		q = q.Where(query.ShipInfo.ShipType.Eq(*req.ShipType))
	}

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	offset := (page - 1) * pageSize

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	dbList, err := q.Offset(offset).Limit(pageSize).Order(query.ShipInfo.ID.Desc()).Find()
	if err != nil {
		return nil, err
	}

	return &api.GetShipListResp{
		Total:    total,
		Ships:    ToAPIModelList(dbList),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

func (s *ShipService) UpdateShip(id int64, req *api.Ship) (*api.Ship, error) {
	_, err := query.ShipInfo.WithContext(s.ctx).Where(query.ShipInfo.ID.Eq(uint(id))).First()
	if err != nil {
		return nil, err // 不存在就返回错误
	}

	dbModel := ToDBModel(req)
	dbModel.ID = uint(id) // 确保 ID 是对的

	err = query.ShipInfo.WithContext(s.ctx).Where(query.ShipInfo.ID.Eq(uint(id))).Save(dbModel)
	if err != nil {
		return nil, err
	}

	return s.GetShipByID(id)
}

func (s *ShipService) DeleteShip(id int64) error {
	_, err := query.ShipInfo.WithContext(s.ctx).Where(query.ShipInfo.ID.Eq(uint(id))).Delete()
	return err
}

func (s *ShipService) GetStatistics() (*api.GetStatisticsResp, error) {
	resp := &api.GetStatisticsResp{}

	total, _ := query.ShipInfo.WithContext(s.ctx).Count()
	resp.TotalShips = total

	type Result struct {
		Name  string
		Count int64
	}
	db := query.ShipInfo.WithContext(s.ctx).UnderlyingDB() // 拿到底层的 GORM DB 对象

	var factionResults []Result
	db.Model(&model.ShipInfo{}).Select("faction as name, count(*) as count").Group("faction").Scan(&factionResults)
	for _, r := range factionResults {
		resp.FactionStats = append(resp.FactionStats, &api.StatItem{Name: r.Name, Count: r.Count})
	}

	var rarityResults []Result
	db.Model(&model.ShipInfo{}).Select("rarity as name, count(*) as count").Group("rarity").Scan(&rarityResults)
	for _, r := range rarityResults {
		resp.RarityStats = append(resp.RarityStats, &api.StatItem{Name: r.Name, Count: r.Count})
	}

	var typeResults []Result
	db.Model(&model.ShipInfo{}).Select("ship_type as name, count(*) as count").Group("ship_type").Scan(&typeResults)
	for _, r := range typeResults {
		resp.ShipTypeStats = append(resp.ShipTypeStats, &api.StatItem{Name: r.Name, Count: r.Count})
	}

	return resp, nil
}
