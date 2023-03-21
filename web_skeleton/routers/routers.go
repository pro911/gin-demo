package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/pro911/gin-demo/web_skeleton/middlewares"
	"net/http"
)

func Setup() *gin.Engine {
	r := gin.New()

	r.Use(middlewares.GinLogger(), middlewares.GinRecovery(true))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	return r
}
