package main

import (
	"github.com/sumory/runbot/common"
	"fmt"
	"runtime"
	"time"
)


func do(har *common.StatusAPI)func(){
	return func(){
		e,har,statusCode, r:=common.RunHar(har)
		fmt.Println(common.Current(), e, har.HarId, statusCode, r)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	tw := common.NewTimeWheel(50*time.Millisecond, 20, 500)

	db := common.NewMongo()
	hars := common.GetAllStatusApiOfUser(db, "db46161c-917a-40d8-b1fd-242e7cc8f4b3")
	fmt.Println(hars.Len())


	for e := hars.Front(); e != nil; e = e.Next() {
		element := e.Value
		h := element.(*common.StatusAPI)
		fmt.Printf("%+v\n", h)

		taskId, _ := tw.Loop(2*time.Second, do(h))
		fmt.Println("开启任务", taskId, h.HarId )
	}

	select {

	}
}
