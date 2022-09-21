package router

import (
	"chatDemo/api"
	"chatDemo/cache"
	"chatDemo/conf"
	"chatDemo/service"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	cache.Init()
	conf.Init()

	r := gin.Default()

	v1 := r.Group("/")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "SUCCESS")
		})
		v1.POST("user/register", api.UserRegister)
		v1.GET("ws", service.WsHandler)
	}
	return r
}
