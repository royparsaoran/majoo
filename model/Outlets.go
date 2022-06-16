package model

import (
	"gorm.io/gorm"
)

type Outlets struct {
	gorm.Model `swaggerignore:"true"`
	OutletName string
	Merchant   Merchants `gorm:"foreignKey:MerchantID"`
	MerchantID uint64    `gorm:"column:merchant_id"`
	// CreatedBy  int
	// UpdatedBy  int
}

func (Outlets) TableName() string {
	return "outlets"
}
