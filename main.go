package main

import (
	"fmt"
	"log"
	_ "net/http"
	"syscall"

	"github.com/fvbock/endless"

	"github.com/Songkun007/go-gin-blog/models"
	"github.com/Songkun007/go-gin-blog/pkg/gredis"
	"github.com/Songkun007/go-gin-blog/pkg/logging"
	"github.com/Songkun007/go-gin-blog/pkg/setting"
	"github.com/Songkun007/go-gin-blog/routers"
)

func main() {

	// 初始化配置
	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.Setup()

	// 方式一，常规启动
	//router := routers.InitRouter()
	//
	//s := &http.Server{
	//	Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
	//	Handler:        router,
	//	ReadTimeout:    setting.ReadTimeout,
	//	WriteTimeout:   setting.WriteTimeout,
	//	MaxHeaderBytes: 1 << 20,
	//}
	//
	//s.ListenAndServe()

	// 方式二，优雅重启，endless 热更新是采取创建子进程后，将原进程退出的方式
	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}