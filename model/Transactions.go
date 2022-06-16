package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transactions struct {
	gorm.Model `swaggerignore:"true"`
	Merchant   Merchants `gorm:"foreignKey:MerchantID"`
	MerchantID uint64    `gorm:"column:merchant_id"`
	Outlet     Outlets   `gorm:"foreignKey:OutletID"`
	OutletID   uint64    `gorm:"column:outlet_id"`
	BillTotal  decimal.Decimal
	// CreatedBy  int
	// UpdatedBy  int
}

func (Transactions) TableName() string {
	return "transactions"
}
