package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pro911/gin-demo/bluebell/logic"
	"github.com/pro911/gin-demo/bluebell/models/request"
	"go.uber.org/zap"
	"net/http"
)

func SignUpHandler(c *gin.Context) {
	//接收参数
	var p request.UserSignUpReq
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误,直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": "请求参数有误",
		})
		return
	}
	//处理业务
	logic.SignUp()
	//返回参数
	c.JSON(http.StatusOK, gin.H{
		"x": 1,
	})
	return
}
