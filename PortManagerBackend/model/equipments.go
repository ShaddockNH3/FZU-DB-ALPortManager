package model

import(
	"time"
)

type Equipment struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt *time.Time `gorm:"index"`

	Name string `gorm:"size:255;uniqueIndex"` 
	
	Type string `gorm:"size:64;index"` 
}

func (Equipment) TableName() string {
	return "equipments"
}
