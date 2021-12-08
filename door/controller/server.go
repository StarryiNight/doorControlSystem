package controller

import (
	"bytes"
	"door/models"
	"door/mysql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"time"
)

var Conn net.Conn

const address = "192.168.43.166:8266"

func TcpInitHandler(c *gin.Context) {
	client, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("listener create error", err)
	}
	fmt.Println("waiting for client...")
	for {
		Conn, err = client.Accept()
		if err != nil {
			fmt.Println("accept error", err)

		}
		go ConnectionHandle(Conn, c)
	}
}

func ConnectionHandle(conn net.Conn, c *gin.Context) {
	defer conn.Close()

	//fmt.Println("client address: ", conn.RemoteAddr())

	buffer := make([]byte, 0)
	var (
		strBuffer string
		err       error

	)

	edFlag := false
	sumLen := 0
	flag := false
	//读取json{}数据 防止tcp粘包而读取数据错误
	for true {
		recBuffer := make([]byte, 1023)
		_, err = conn.Read(recBuffer)
		if err != nil {
			fmt.Println("Read error", err)
		}
		aviLen := 0
		first := 0
		firstFlag:=false
		for index, value := range recBuffer {
			//fmt.Println("index,value",index,string(value))
			if value == '{' {
				flag = true
			}
			if !flag {
				continue
			}

			if flag&&value == '}' {
				edFlag = true
				break;
			}

			if value != 0 {
				if first == 0 &&!firstFlag{
					first = index
					firstFlag=true
				}
				aviLen++
				sumLen++
			}
		}
		if flag {
			buffer = append(buffer, recBuffer[first:first+aviLen+1]...)
		}

		ed:=bytes.IndexByte(buffer, 0)
		if ed!=-1 {
			buffer=buffer[:ed]
		}
		if edFlag {
			break
		}

	}

	strBuffer = string(buffer)

	var user models.User
	err = json.Unmarshal(buffer, &user)
	if err != nil {
		fmt.Println("json.Unmarshal() failed ", err)
		return
	}
	fmt.Println(user)
	err = mysql.AddRecord(user)
	if err != nil {
		fmt.Println("mysql.AddRecord() failed ", err)
		return
	}

	fmt.Println("Message: ", strBuffer)
	time.Sleep(time.Second * 1)
	/*	_, err = conn.Write([]byte("server received"))
		if err != nil {
			fmt.Println("send message error",err)
		}
		fmt.Println("send message success")*/
}

func OpenDoorHandler(c *gin.Context) {
	msg := "open"
	_, err := Conn.Write([]byte(msg))
	_, err = Conn.Write([]byte("*"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "send failed",
			"err": err,
		})
		return
	}

	user := models.User{
		CardId:   c.GetString(CardId),
		UserName: c.GetString(UserName),
	}
	err = mysql.AddRecord(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "send failed",
			"err": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "send success",
	})
}

func SendMsgHandler(c *gin.Context) {
	msg := c.PostForm("msg")
	_, err := Conn.Write([]byte(msg))
	_, err = Conn.Write([]byte("*"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "send failed",
			"err": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "send success",
	})
}

func UpdateCardHandler(c *gin.Context) {
	//查询数据
	users,err:=mysql.ScanUser()
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"msg":"查询失败",
			"err":err.Error(),
		})
		return
	}

	//发送数据
	jsonData, err := json.Marshal(users)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"msg":"json处理失败",
			"err":err.Error(),
		})
		return
	}

	var str string
	for i,r:=range jsonData{
		str+=string(r)
		if (i+1)%63==0 {
			_, err = Conn.Write([]byte(str))
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"msg": "send failed",
					"err": err,
				})
				return
			}

			str=""
		}
	}
	_, err = Conn.Write([]byte(str))

/*	_, err = Conn.Write(jsonData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "send failed",
			"err": err,
		})
		return
	}*/

	_, err = Conn.Write([]byte("*"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "send failed",
			"err": err,
		})
		return
	}
	/*for _, value := range users{
		//处理成json数据
		jsonData, err := json.Marshal(value)
		if err != nil {
			c.JSON(http.StatusOK,gin.H{
				"msg":"json处理失败",
				"err":err.Error(),
			})
			return
		}
		_, err = Conn.Write(jsonData)
		_, err = Conn.Write([]byte("\n"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg": "send failed",
				"err": err,
			})
			return
		}
		time.Sleep(time.Millisecond*200)
	}*/


	c.JSON(http.StatusOK, gin.H{
		"msg": "send success",
	})
}