package controller

import (
	"majoo/biz/outlet"
	lib "majoo/lib"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetOutletOmzet godoc
// @tags Outlet
// @Accept  json
// @Produce  json
// @Param merchant_id path int true "{1}"
// @Param outlet_id path int true "{1}"
// @Param bulan query string false "2022-03"
// @Param limit query int false "20"
// @Param page query int false "1"
// @Success 200 {object} lib.OutputFormat{Data=[]outlet.OutletOmzet}
// @Security BearerAuth
// @Router /merchant/{merchant_id}/outlet/{outlet_id}/omzet [get]
func GetOutletOmzet(c *gin.Context) {
	merchantID, _ := strconv.Atoi(c.Param("merchant_id"))
	outletID, _ := strconv.Atoi(c.Param("outlet_id"))
	bulan := c.DefaultQuery("bulan", time.Now().Format("2006-01"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	bearer := c.Request.Header.Get("Authorization")

	repayment, err := outlet.GetOutletOmzet(merchantID, outletID, bulan, bearer, limit, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.ResponseFormat(false, err.Error(), nil, nil))
		return
	}

	c.JSON(http.StatusOK, lib.ResponseFormat(true, "OK", repayment, nil))
}
