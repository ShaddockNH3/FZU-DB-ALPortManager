package model

import(
	"time"
)

type ShipEquipment struct {
	ShipID    uint `gorm:"primaryKey;autoIncrement:false"`
	SlotIndex int  `gorm:"primaryKey;autoIncrement:false"` // 0-4 是常规栏, 5 是兵装栏

	EquipmentID uint `gorm:"index"`

	Equipment Equipment `gorm:"foreignKey:EquipmentID;references:ID"`

	CreatedAt time.Time `gorm:""`
}

func (ShipEquipment) TableName() string {
	return "ship_equipments"
}

