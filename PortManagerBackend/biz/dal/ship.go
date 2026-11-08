package dal

import (
	"context"
	"errors"
	"fmt"

	"FZU-DB-ALPortManager/model"
	"FZU-DB-ALPortManager/query"

	"gorm.io/gorm"
)

// Ship 舰船的数据访问层。
//
// 这一层只跟数据库打交道：查、写、删。业务判断（槽位能不能装、星级够不够）
// 在 biz/domain/ship，编排在 biz/service。加索引、换 ORM、写原生 SQL 都只改这里。
type Ship struct{}

func NewShip() *Ship { return &Ship{} }

// ListFilter 列表筛选条件，零值表示这一维不过滤。
type ListFilter struct {
	ShipName string
	Rarity   string
	ShipType string
	Faction  string
}

// GetByID 取详情，带装备（含装备本体）。
func (Ship) GetByID(ctx context.Context, id uint) (*model.ShipInfo, error) {
	q := query.ShipInfo
	row, err := q.WithContext(ctx).
		Preload(q.Equipments.Equipment).
		Where(q.ID.Eq(id)).
		First()
	if err != nil {
		return nil, translate(err)
	}
	return row, nil
}

// List 按条件分页取列表，同时返回总数。列表也要装备数据，前端要显示星级和槽位。
func (Ship) List(ctx context.Context, filter ListFilter, offset, limit int) ([]*model.ShipInfo, int64, error) {
	q := query.ShipInfo.WithContext(ctx)

	if filter.ShipName != "" {
		q = q.Where(query.ShipInfo.ShipName.Like("%" + filter.ShipName + "%"))
	}
	if filter.Rarity != "" {
		q = q.Where(query.ShipInfo.Rarity.Eq(filter.Rarity))
	}
	if filter.ShipType != "" {
		q = q.Where(query.ShipInfo.ShipType.Eq(filter.ShipType))
	}
	if filter.Faction != "" {
		q = q.Where(query.ShipInfo.Faction.Eq(filter.Faction))
	}

	total, err := q.Count()
	if err != nil {
		return nil, 0, err
	}

	rows, err := q.Preload(query.ShipInfo.Equipments.Equipment).
		Offset(offset).
		Limit(limit).
		Order(query.ShipInfo.ID.Desc()).
		Find()
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

// GroupCount 分组计数的一行。
type GroupCount struct {
	Name  string
	Count int64
}

// StatField 允许拿来分组的列。只放开这三个，
// 别让调用方把任意列名当参数拼进 SQL。
type StatField string

const (
	StatByFaction  StatField = "faction"
	StatByRarity   StatField = "rarity"
	StatByShipType StatField = "ship_type"
)

var statFields = map[StatField]bool{
	StatByFaction:  true,
	StatByRarity:   true,
	StatByShipType: true,
}

// Count 舰船总数。
func (Ship) Count(ctx context.Context) (int64, error) {
	return query.ShipInfo.WithContext(ctx).Count()
}

// CountBy 按某列分组计数，交给数据库做 GROUP BY。
// 之前是把整张表 Find 到内存里再数，几条和几万条的开销差一个量级。
func (Ship) CountBy(ctx context.Context, column StatField) ([]GroupCount, error) {
	if !statFields[column] {
		return nil, fmt.Errorf("dal: 不支持的分组列 %q", column)
	}

	var rows []GroupCount
	db := query.ShipInfo.WithContext(ctx).UnderlyingDB()

	err := db.Model(&model.ShipInfo{}).
		Select(string(column)+" AS name", "COUNT(*) AS count").
		Group(string(column)).
		Order("count DESC, name ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// Create 新建舰船，只写本体，装备槽记录在装装备时按需创建。
func (Ship) Create(ctx context.Context, row *model.ShipInfo) error {
	return query.ShipInfo.WithContext(ctx).Create(row)
}

// Update 按字段更新。updates 为空表示没有要改的字段。
func (Ship) Update(ctx context.Context, id uint, updates map[string]any) error {
	if len(updates) == 0 {
		return nil
	}

	_, err := query.ShipInfo.WithContext(ctx).
		Where(query.ShipInfo.ID.Eq(id)).
		Updates(updates)
	return err
}

// Delete 删舰船，顺带清掉它的装备记录。两条语句放在一个事务里，
// 不然删到一半失败会留下指向不存在舰船的槽位记录。
func (Ship) Delete(ctx context.Context, id uint) error {
	return query.Q.Transaction(func(tx *query.Query) error {
		if _, err := tx.ShipEquipment.WithContext(ctx).
			Where(tx.ShipEquipment.ShipID.Eq(id)).
			Delete(); err != nil {
			return err
		}

		if _, err := tx.ShipInfo.WithContext(ctx).
			Where(tx.ShipInfo.ID.Eq(id)).
			Delete(); err != nil {
			return err
		}
		return nil
	})
}

// translate 把 gorm 的「没查到」换成 dal.ErrNotFound，其它错误原样往上抛。
func translate(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}
