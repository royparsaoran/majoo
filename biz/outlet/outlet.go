package outlet

import (
	"errors"
	"majoo/biz/merchants"
	"majoo/conn"
	lib "majoo/lib"
	"majoo/model"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OutletOmzet struct {
	Date         string
	MerchantName string
	OutletName   string
	Omzet        decimal.Decimal
}

type Nresult struct {
	N decimal.Decimal
}

func Get(ID int) (model.Outlets, error) {
	var b model.Outlets

	if err := conn.DBConnection().Find(&b, "id = ?", ID).Error; err != nil {
		return b, err
	}
	return b, nil
}

func GetOutletOmzet(merchantID int, outletID int, flag string, bearer string, limit int, page int) (ret []OutletOmzet, err error) {
	var (
		models []OutletOmzet
		t, _   = time.Parse("2006-01", flag)
		end    = t.AddDate(0, 0, page*limit-1)
		start  = end.AddDate(0, 0, -limit+1)
	)

	err = lib.MerchantMiddleware(merchantID, bearer)
	if err != nil {
		return ret, nil
	}

	merchant, err := merchants.Get(merchantID)
	if err != nil {
		return ret, nil
	}

	outlet, err := Get(outletID)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ret, errors.New("outlet not found")
		}

		return ret, nil
	}

	for start.Before(end) || start.Equal(end) {
		omzet := GetOmzet(start.Format("2006-01-02"), merchantID)

		m := OutletOmzet{
			Date:         start.Format("2006-01-02"),
			MerchantName: merchant.MerchantName,
			OutletName:   outlet.OutletName,
			Omzet:        omzet,
		}

		models = append(models, m)
		start = start.AddDate(0, 0, 1)
	}

	return models, nil
}

func GetOmzet(date string, outletID int) (count decimal.Decimal) {
	var n Nresult

	if err := conn.DBConnection().
		Table("transactions").
		Select("sum(bill_total) as n").
		Where("outlet_id = ? AND created_at >= ? AND created_at <= ?", outletID, date+" 00:00:00", date+" 23:59:59").
		Scan(&n).
		Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return count
		}

		return count
	}

	return n.N
}
