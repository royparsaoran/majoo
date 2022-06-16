package controller

import (
	"net/http"

	"majoo/biz/auth"
	lib "majoo/lib"
	"majoo/model"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @tags Auth
// @Accept  json
// @Produce  json
// @Param Login body model.Users true "Login"
// @Success 200 {object} lib.OutputFormat{Data=auth.Token}
// @Router /login [post]
func Login(c *gin.Context) {
	var param model.Users

	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, lib.ResponseFormat(false, err.Error(), nil, nil))
		return
	}

	u, code, err := auth.Login(param.UserName, param.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, lib.ResponseFormatCode(false, err.Error(), nil, nil, code))
		return
	}
	c.JSON(http.StatusOK, lib.ResponseFormat(true, "ok", u, nil))
}

// Register godoc
// @tags Auth
// @Accept  json
// @Produce  json
// @Param Register body model.Users true "Register"
// @Success 200 {object} lib.OutputFormat{Data=auth.Token}
// @Router /register [post]
func Register(c *gin.Context) {
	var param model.Users
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, lib.ResponseFormat(false, err.Error(), nil, nil))
		return
	}

	u, code, err := auth.Register(param)

	if err != nil {
		c.JSON(http.StatusBadRequest, lib.ResponseFormatCode(false, err.Error(), nil, nil, code))
		return
	}
	c.JSON(http.StatusOK, lib.ResponseFormat(true, "ok", u, nil))
}
