package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/maxilovera/go-crud-example/clients/responses"
)

func SetUserInContext(c *gin.Context, user *responses.UserInfo) {
	c.Set("UserInfo", user)
}

func GetUserInfoFromContext(c *gin.Context) *responses.UserInfo {
	userInfo, _ := c.Get("UserInfo")

	user, _ := userInfo.(*responses.UserInfo)

	return user
}
