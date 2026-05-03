package main

import (
	"fmt"
	"log"

	"documentManage/config"
	"documentManage/database"
	"documentManage/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()

	if config.AppConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	database.InitDB()

	r := routers.SetupRouter()

	port := config.AppConfig.Server.Port
	log.Printf("服务器启动在端口 %d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
