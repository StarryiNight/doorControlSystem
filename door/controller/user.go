package controller

import (
	"door/models"
	"door/mysql"
	"door/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const UserName = "username"
const CardId = "card_id"

func LoginHandler(c *gin.Context) {
	//1.获取请求参数
	p := new(models.ParamUser)
	if err := c.ShouldBind(p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "请求参数错误",
			"err": err.Error(),
		})
		return
	}
	//2.业务逻辑处理
	user, err := mysql.Login(p)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "登录失败",
			"err": err.Error(),
		})
		return
	}

	token, err := pkg.GenToken(user.CardId, user.UserName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "登录失败",
			"err": "服务繁忙",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "登录成功",
		"token":token,
	})


}


func ScanUsersHandler(c *gin.Context){
	users,err:=mysql.ScanUser()
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"msg":"查询失败",
			"err":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,users)
}

func ScanRecordsHandler(c *gin.Context){
	records,err:=mysql.ScanRecord()
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"msg":"查询失败",
			"err":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,records)
}

