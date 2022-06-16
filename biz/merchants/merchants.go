package merchants

import (
	"errors"
	"majoo/conn"
	lib "majoo/lib"
	"majoo/model"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type MerchantOmzet struct {
	Date         string
	MerchantName string
	Omzet        decimal.Decimal
}

type Nresult struct {
	N decimal.Decimal
}

func Get(ID int) (model.Merchants, error) {
	var b model.Merchants

	if err := conn.DBConnection().Find(&b, "id = ?", ID).Error; err != nil {
		return b, err
	}
	return b, nil
}

func GetMerchantOmzet(merchantID int, flag string, bearer string, limit int, page int) (ret []MerchantOmzet, err error) {
	var (
		models []MerchantOmzet
		t, _   = time.Parse("2006-01", flag)
		end    = t.AddDate(0, 0, page*limit-1)
		start  = end.AddDate(0, 0, -limit+1)
	)

	err = lib.MerchantMiddleware(merchantID, bearer)
	if err != nil {
		return ret, nil
	}

	merchant, err := Get(merchantID)
	if err != nil {
		return ret, nil
	}

	for start.Before(end) || start.Equal(end) {
		omzet := GetOmzet(start.Format("2006-01-02"), merchantID)

		m := MerchantOmzet{
			Date:         start.Format("2006-01-02"),
			MerchantName: merchant.MerchantName,
			Omzet:        omzet,
		}

		models = append(models, m)
		start = start.AddDate(0, 0, 1)
	}

	return models, nil
}

func GetOmzet(date string, merchantID int) (count decimal.Decimal) {
	var n Nresult

	if err := conn.DBConnection().
		Table("transactions").
		Select("sum(bill_total) as n").
		Where("merchant_id = ? AND created_at >= ? AND created_at <= ?", merchantID, date+" 00:00:00", date+" 23:59:59").
		Scan(&n).
		Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return count
		}

		return count
	}

	return n.N
}
