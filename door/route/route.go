package route

import (
	"door/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)
import "door/controller"

func Setup() *gin.Engine {

	r := gin.Default()

	r.Static("/image","./static/image")

	r.LoadHTMLGlob("static/templates/*")


	r.Use(middlewares.Cors())
	//渲染页面
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK,"login.html",nil)
	})
	r.GET("/home",func(c *gin.Context) {
		c.HTML(http.StatusOK,"home.html",nil)
	})
	r.GET("register", func(c *gin.Context) {
		c.HTML(http.StatusOK,"register.html",nil)
	})

	//建立tcp连接
	r.GET("/connect", controller.TcpInitHandler)
	//登陆
	r.POST("/login", controller.LoginHandler)
	//注册
	r.POST("/register",controller.RegisterHandler)
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
