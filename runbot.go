package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"os"
	"./common"
)

type User struct {
	Uid int
}

type Relation struct {
	Uid    int
	Target int
	Time   float64
}

var config *common.Config

var buildstamp string = ""
var githash string = ""

func main() {
	args := os.Args
	if len(args)==2 && (args[1]=="--version" || args[1] =="-v") {
		fmt.Printf("Git Commit Hash: %s\n", githash)
		fmt.Printf("UTC Build Time: %s\n", buildstamp)
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("os args:", os.Args)
	if args == nil || len(args) < 2 {
		panic("缺少配置文件")
	}

	config = common.InitConfig(args[1])
	common.InitContext(config)
	common.StartRunAll()//初始化时运行所有需要monitor的status api

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/start", handleStartStatusAPI)
	router.GET("/stop", handleStopStatusAPI)
	router.Run(fmt.Sprintf(":%d", config.Port))
}


func handleStartStatusAPI(c *gin.Context) {
	statusAPIId:= c.Query("statusAPIId")
	if statusAPIId == "" {
		c.JSON(http.StatusOK, gin.H{
			"success":false,
			"msg":"params error",
		})
		return
	}

	common.StartRun(statusAPIId)

	c.JSON(http.StatusOK, gin.H{
		"success":true,
		"msg":"ok",
	})
}
func handleStopStatusAPI(c *gin.Context) {
	statusAPIId:= c.Query("statusAPIId")
	if statusAPIId == "" {
		c.JSON(http.StatusOK, gin.H{
			"success":false,
			"msg":"params error",
		})
		return
	}

	common.StopRun(statusAPIId)

	c.JSON(http.StatusOK, gin.H{
		"success":true,
		"msg":"ok",
	})
}