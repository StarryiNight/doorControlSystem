package middlewares

import (
	"door/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const CtxUserName = "username"
const CtxCardId="card_id"

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {

		/*authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请先登陆",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		// 使用定义的解析JWT的函数来解析parts[1],获取到的tokenString*/
		token:=c.Query("token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请先登陆",
			})
			c.Abort()
			return
		}
		mc, err := pkg.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上 方便获取用户名和卡号
		c.Set(CtxUserName, mc.UserName)
		c.Set(CtxCardId,mc.CardId)
		fmt.Println(mc.UserName)
		c.Next()
	}
}
