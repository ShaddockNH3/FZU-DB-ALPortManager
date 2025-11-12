package model

import (
	"time"
)

type ShipInfo struct {
	ID        uint       `gorm:"primarykey"`
	CreatedAt time.Time  `gorm:""`
	UpdatedAt time.Time  `gorm:""`
	DeletedAt *time.Time `gorm:"index"`

	ShipName string `gorm:"size:255;index"`
	Rarity   string `gorm:"size:64"`
	ShipType string `gorm:"size:64;index"`
	Faction  string `gorm:"size:64;index"`
	Level    int    `gorm:"default:1"`

	Stars int `gorm:"default:1"` // 当前星级，满星后可装备兵装

	Equipments []ShipEquipment `gorm:"foreignKey:ShipID;references:ID"`
}

func (ShipInfo) TableName() string {
	return "ship_info"
}
