package controller

import (
	"majoo/biz/merchants"
	lib "majoo/lib"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetMerchantOmzet godoc
// @tags Merchant
// @Accept  json
// @Produce  json
// @Param merchant_id path int true "{1}"
// @Param bulan query string false "2022-03"
// @Param limit query int false "20"
// @Param page query int false "1"
// @Success 200 {object} lib.OutputFormat{Data=[]merchants.MerchantOmzet}
// @Security BearerAuth
// @Router /merchant/{merchant_id}/omzet [get]
func GetMerchantOmzet(c *gin.Context) {
	ID, _ := strconv.Atoi(c.Param("merchant_id"))
	bulan := c.DefaultQuery("bulan", time.Now().Format("2006-01"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	bearer := c.Request.Header.Get("Authorization")

	repayment, err := merchants.GetMerchantOmzet(ID, bulan, bearer, limit, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.ResponseFormat(false, err.Error(), nil, nil))
		return
	}

	c.JSON(http.StatusOK, lib.ResponseFormat(true, "OK", repayment, nil))
}
