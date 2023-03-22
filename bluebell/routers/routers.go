package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pro911/gin-demo/bluebell/middlewares"
	"github.com/pro911/gin-demo/bluebell/settings"
	"net/http"
	"time"
)

func Setup() *gin.Engine {
	r := gin.New()

	r.Use(middlewares.GinLogger(), middlewares.GinRecovery(true))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprintf("ok:%v", time.Now().Unix()))
	})
	r.POST("/version", func(c *gin.Context) {
		c.String(http.StatusOK, settings.Conf.AppConfig.Version)
	})

	return r
}
