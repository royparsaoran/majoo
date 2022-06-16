package model

import (
	"gorm.io/gorm"
)

type Merchants struct {
	gorm.Model   `swaggerignore:"true"`
	MerchantName string
	Users        Users  `gorm:"foreignKey:UserID"`
	UserID       uint64 `gorm:"column:user_id"`
	// CreatedBy    int
	// UpdatedBy    int
}

func (Merchants) TableName() string {
	return "merchants"
}
