package main


import (
	"github.com/jrallison/go-workers"
	"fmt"
	"github.com/sumory/go_relation/common"
)

func Get(message *workers.Msg) {
	fmt.Printf("%+v\n", message)
}


func Post(message *workers.Msg) {
	fmt.Printf("%+v\n", message)
}

type myMiddleware struct{}

func (r *myMiddleware) Call(queue string, message *workers.Msg, next func() bool) (acknowledge bool) {
	fmt.Printf("message %s start %s \n", message.Jid(), common.Now())
	acknowledge = next()
	fmt.Printf("message %s finish %s \n", message.Jid(), common.Now())
	return
}

func main() {
	workers.Configure(map[string]string{
		"server":  "192.168.100.185:6389",
		"database":  "0",
		"pool":    "30",
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		"process": "1",
	})

	//workers.Middleware.Append(&myMiddleware{})

	workers.Process("get_queue", Get, 10)
	workers.Process("post_queue", Post, 10)

	go workers.StatsServer(8080)

	workers.Run()
}