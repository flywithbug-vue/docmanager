package handler

import (
	"doc-manager/web_server/common"
	"doc-manager/web_server/core/jwt"
	"doc-manager/web_server/model"
	"net/http"
	"strings"

	log "github.com/flywithbug/log4go"
	"github.com/gin-gonic/gin"
)

func loginHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.JSON(http.StatusOK, aRes)
	}()
	user := new(model.User)
	err := c.BindJSON(user)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	user, err = model.LoginUser(user.Account, user.Password)
	if err != nil {
		log.Error(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "account or password not right")
		return
	}
	claims := jwt.NewCustomClaims(user.UserId, user.Account)
	token, err := jwt.GenerateToken(claims)
	if err != nil {
		log.Error(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "token generate error"+err.Error())
		return
	}
	userAg := c.GetHeader(common.KeyUserAgent)
	_, err = model.UserLogin(user.UserId, userAg, token, c.ClientIP())
	if err != nil {
		log.Error(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "token generate error"+err.Error())
		return
	}
	aRes.SetResponseDataInfo("token", token)
}

func registerHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.JSON(aRes.Code, aRes)
	}()
	user := new(model.User)
	err := c.BindJSON(user)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	if user.Account == "" {
		aRes.SetErrorInfo(http.StatusBadRequest, "account can not be nil")
		return
	}
	if user.Password == "" {
		aRes.SetErrorInfo(http.StatusBadRequest, "Password can not be nil")
		return
	}
	if user.Email == "" {
		aRes.SetErrorInfo(http.StatusBadRequest, "email can not be nil")
		return
	}
	err = user.UserInsert()
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest, err.Error())
		return
	}
	aRes.AddResponseInfo("user", user)
}

func logoutHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.JSON(aRes.Code, aRes)
	}()
	token := common.UserToken(c)
	if token == "" {
		aRes.SetErrorInfo(http.StatusBadRequest, "token not found")
		return
	}
	err := model.UpdateLoginStatus(token, model.StatusLogout)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest, err.Error())
		return
	}
	aRes.SetSuccessInfo(http.StatusOK, "success")
}

func getUserInfoHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.JSON(aRes.Code, aRes)
	}()
	userId := common.UserId(c)
	if strings.EqualFold(userId, "") {
		aRes.SetErrorInfo(http.StatusUnauthorized, "user not found")
		return
	}
	user, err := model.FindByUserId(userId)
	if err != nil {
		aRes.SetErrorInfo(http.StatusUnauthorized, "user not found:"+err.Error())
		return
	}
	aRes.AddResponseInfo("user", user)
}

func getAllUserInfoHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.JSON(aRes.Code, aRes)
	}()
	userId := common.UserId(c)
	if strings.EqualFold(userId, "") {
		aRes.SetErrorInfo(http.StatusUnauthorized, "user not found")
		return
	}

	users, err := model.FindAllUsers()
	if err != nil {
		aRes.SetErrorInfo(http.StatusUnauthorized, "users find error"+err.Error())
		return
	}
	aRes.AddResponseInfo("users", users)
}

func searchUserHandler(c *gin.Context) {

}
