package dal

import (
	"context"
	"errors"

	"FZU-DB-ALPortManager/model"
	"FZU-DB-ALPortManager/query"

	"gorm.io/gorm"
)

// Equipment 装备与槽位的数据访问层。
//
// 「这个装备能不能装到这个槽」是业务规则，不在这里；这里只负责读写。
type Equipment struct{}

func NewEquipment() *Equipment { return &Equipment{} }

// GetByID 取一件装备，找不到返回 ErrNotFound。
func (Equipment) GetByID(ctx context.Context, id uint) (*model.Equipment, error) {
	row, err := query.Equipment.WithContext(ctx).
		Where(query.Equipment.ID.Eq(id)).
		First()
	if err != nil {
		return nil, translate(err)
	}
	return row, nil
}

// List 取装备列表。typeName 为空表示不过滤类型。
func (Equipment) List(ctx context.Context, typeName string) ([]*model.Equipment, error) {
	q := query.Equipment.WithContext(ctx)
	if typeName != "" {
		q = q.Where(query.Equipment.Type.Eq(typeName))
	}

	rows, err := q.Order(query.Equipment.ID).Find()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// UpsertSlot 把某件装备写进舰船的某个槽位：有记录就改，没有就插。
func (Equipment) UpsertSlot(ctx context.Context, shipID uint, slot int, equipmentID uint) error {
	_, err := query.ShipEquipment.WithContext(ctx).
		Where(query.ShipEquipment.ShipID.Eq(shipID)).
		Where(query.ShipEquipment.SlotIndex.Eq(slot)).
		First()

	switch {
	case err == nil:
		_, err = query.ShipEquipment.WithContext(ctx).
			Where(query.ShipEquipment.ShipID.Eq(shipID)).
			Where(query.ShipEquipment.SlotIndex.Eq(slot)).
			Update(query.ShipEquipment.EquipmentID, equipmentID)
		return err

	case errors.Is(err, gorm.ErrRecordNotFound):
		return query.ShipEquipment.WithContext(ctx).Create(&model.ShipEquipment{
			ShipID:      shipID,
			SlotIndex:   slot,
			EquipmentID: equipmentID,
		})

	default:
		return err
	}
}

// ClearSlot 卸下某格装备。本来就没装也算成功。
func (Equipment) ClearSlot(ctx context.Context, shipID uint, slot int) error {
	_, err := query.ShipEquipment.WithContext(ctx).
		Where(query.ShipEquipment.ShipID.Eq(shipID)).
		Where(query.ShipEquipment.SlotIndex.Eq(slot)).
		Delete()
	return err
}
