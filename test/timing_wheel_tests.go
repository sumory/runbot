package main

import (
	"time"
	"github.com/sumory/runbot/common"
	"fmt"
	"runtime"
)

func get(){
	fmt.Println("jj")
}

func loop(i int) func(){
	return func(){
		fmt.Println("loop", i , common.Current())
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	tw := common.NewTimeWheel(50*time.Millisecond, 20, 500)

//	taskId, err := tw.After(3*time.Second, get)
//	fmt.Println(taskId, err)

//	var taskId int64
	for  i :=0;i<20;i++ {
		//time.Sleep(2*time.Second)

		_, _ = tw.Loop(1*time.Second, loop(i))
		//fmt.Println("taskId", taskId, "start")
	}


//	go func(){//打印监控
//		ticker := time.NewTicker(1 * time.Second)
//		for {
//			select {
//			case <-ticker.C:
//				fmt.Println(tw.Monitor())
//			}
//		}
//	}()

//	select {
//	case <- time.After(10*time.Second):
//		fmt.Println("remove",taskId)
//		tw.Remove(taskId)
//	case <- time.After(15*time.Second):
//		tw.Loop(4*time.Second, loop)
//	}

	select {

	}

}