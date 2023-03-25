package controller

import (
	"github.com/gin-gonic/gin"
	"httpsql/service"
	"strconv"
)

var enginee = gin.Default()

func MainController() {
	sqlgroup := enginee.Group("/sql/")
	sqlgroup.GET("get", service.Get)
	sqlgroup.POST("set", service.Set)
}
func ServerStart(port int) {
	MainController()
	enginee.Run(":" + strconv.Itoa(port))
}
