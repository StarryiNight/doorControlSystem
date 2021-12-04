package main

import (
	"door/mysql"
	"door/route"
	"fmt"
)

func main() {

	if err := mysql.Init(); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}

	r:=route.Setup()

	r.Run(":8080")


}
