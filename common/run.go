package common


import (
	"fmt"
	"time"
	"strconv"
	"gopkg.in/mgo.v2"
)

var TW *TimeWheel
var DB *mgo.Database
var RunningMap map[string]int64
var MyConfig *Config

func InitContext(c *Config){
	MyConfig = c
	DB = NewMongo(c)
	RunningMap = make(map[string]int64)
	TW = NewTimeWheel(50*time.Millisecond, 20, 500)
}

func do(db *mgo.Database, api *StatusAPI) func() {
	return func() {
		t1:=time.Now().UnixNano()
		e, har, statusCode, response := RunStatusAPI(api)
		if MyConfig.LOG_Debug{
			fmt.Println(Current(), e, har.Id, statusCode,response)
		}else{
			fmt.Println(Current(), e, har.Id, statusCode)
		}
		t2:=time.Now().UnixNano()

		log :=&StatusAPILog{
			StatusAPIId: api.Id,
			StatusCode :statusCode,
			Spent :t2-t1,
			UserId:api.UserId,
			Response :response,
			Date:time.Now(),
		}
		SaveStatusAPILog(db, log)
	}
}

func StartRunAll() {
	apis := GetAllStatusApi(DB)
	fmt.Printf("目前总共有%d个要运行的Status API \n", apis.Len())

	for e := apis.Front(); e != nil; e = e.Next() {
		element := e.Value
		api := element.(*StatusAPI)
		fmt.Printf("api: %+v\n", api)

		cron, err := strconv.Atoi(api.Cron)
		if err!=nil {
			fmt.Printf("error to parse `cron` of status api: %s\n", api.Id)
			continue
		}

		timeDuration := time.Duration(int64(cron*1000) * int64(time.Millisecond))
		taskId, _ := TW.Loop(timeDuration, do(DB, api))
		RunningMap[api.Id] = taskId
		fmt.Printf("初始化开启任务, apiId: %s taskId: %d \n", api.Id, taskId)
	}
}

func StartRun(statusAPIId string){
	StopRun(statusAPIId)//先停止


	api := GetStatusApi(DB, statusAPIId)
	fmt.Printf("要开启的任务: %+v\n", api)

	if api!=nil && api.Id!=""{
		cron, err := strconv.Atoi(api.Cron)
		if err!=nil {
			fmt.Printf("开启任务失败, apiId: %s \n", api.Id)
		}else{
			timeDuration := time.Duration(int64(cron*1000) * int64(time.Millisecond))
			taskId, _ := TW.Loop(timeDuration, do(DB, api))
			RunningMap[api.Id] = taskId
			fmt.Printf("开启任务成功, apiId: %s taskId: %d \n", api.Id, taskId)

		}
	}else{
		fmt.Printf("找不到要开启的任务, apiId: %s \n", statusAPIId)
	}
}

func StopRun(statusAPIId string){
	if taskId,ok := RunningMap[statusAPIId];ok{
		TW.Remove(taskId)
		fmt.Printf("停止任务, apiId: %s taskId: %d \n",statusAPIId, taskId)
	}else{
		fmt.Printf("找不到要停止的任务, apiId: %s \n",statusAPIId)
	}
}
