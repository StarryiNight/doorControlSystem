package route

import (
	"door/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)
import "door/controller"

func Setup() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("static/*")
	r.StaticFS("/image",http.Dir("/image"))
	r.StaticFS("/static",http.Dir("/static"))
	r.Use(middlewares.Cors())
	//渲染页面
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK,"/templates/login.html",nil)
	})
	r.GET("/home",middlewares.JWTAuthMiddleware(),func(c *gin.Context) {
		c.HTML(http.StatusOK,"/templates/home.html",nil)
	})

	//建立tcp连接
	r.GET("/connect", controller.TcpInitHandler)
	//登陆
	r.POST("/login", controller.LoginHandler)

	//发送消息
	r.POST("/send",middlewares.JWTAuthMiddleware(), controller.SendMsgHandler)

	//网络开门
	r.POST("/open",middlewares.JWTAuthMiddleware(),controller.OpenDoorHandler)

	//更新卡信息
	r.GET("/update",controller.UpdateCardHandler)

	r.GET("/users",controller.ScanUsersHandler)
	r.GET("/records",controller.ScanRecordsHandler)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
