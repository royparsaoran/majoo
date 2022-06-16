package model

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model `swaggerignore:"true"`
	UserName   string
	Password   string
	// CreatedBy  int
	// UpdatedBy  int
	// Merchants  []Merchants `swaggerignore:"true" gorm:"references:UserID"`
}

type UserMerchant struct {
	UserId     uint64
	MerchantId uint64
}

func (Users) TableName() string {
	return "users"
}
