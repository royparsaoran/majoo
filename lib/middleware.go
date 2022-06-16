package lib

import (
	"errors"
	"strings"
)

func MerchantMiddleware(merchantId int, bearer string) (err error) {
	var allowed bool
	merchant, err := GetMerchantBearer(bearer)
	if err != nil {
		return err
	}

	for _, v := range merchant {
		if v == merchantId {
			allowed = true
		}
	}

	if !allowed {
		return errors.New("unauthorized")
	}

	return nil
}

func GetMerchantBearer(bearToken string) (merchants []int, err error) {
	strArr := strings.Split(bearToken, " ")
	str := ""
	if len(strArr) == 2 {
		str = strArr[1]
	} else {
		return merchants, errors.New("unauthorized - invalid token length")
	}

	claims, err := TokenValidation(str)
	if err != nil {
		return nil, errors.New("unauthorized - invalid token key " + err.Error())
	}

	mid := claims["data"].([]interface{})
	for _, v := range mid {
		merchants = append(merchants, int(v.(map[string]interface{})["MerchantId"].(float64)))
	}

	return merchants, nil
}
