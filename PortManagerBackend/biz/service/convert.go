package service

import (
	"FZU-DB-ALPortManager/biz/model/api"
	"FZU-DB-ALPortManager/model"
	"time"
)

func ToDBModel(req *api.Ship) *model.ShipInfo {
	if req == nil {
		return nil
	}
	return &model.ShipInfo{
		ID:       uint(req.GetID()),
		ShipName: req.GetShipName(),
		Rarity:   req.GetRarity(),
		ShipType: req.GetShipType(),
		Faction:  req.GetFaction(),
		Level:    int(req.GetLevel()),
	}
}

func ToAPIModel(db *model.ShipInfo) *api.Ship {
	if db == nil {
		return nil
	}

	createdAtStr := db.CreatedAt.Format(time.DateTime)
	updatedAtStr := db.UpdatedAt.Format(time.DateTime)
	id := int64(db.ID)
	level := int32(db.Level)

	return &api.Ship{
		ID:        &id,
		ShipName:  db.ShipName,
		Rarity:    db.Rarity,
		ShipType:  db.ShipType,
		Faction:   db.Faction,
		Level:     &level,
		CreatedAt: &createdAtStr, // 赋值指针
		UpdatedAt: &updatedAtStr,
	}
}

func ToAPIModelList(dbList []*model.ShipInfo) []*api.Ship {
	result := make([]*api.Ship, 0, len(dbList))
	for _, dbItem := range dbList {
		result = append(result, ToAPIModel(dbItem))
	}
	return result
}
