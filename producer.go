package main


import (
	"github.com/jrallison/go-workers"
	"time"
)


func main() {
	workers.Configure(map[string]string{
		"server":  "192.168.100.185:6389",
		"database":  "0",
		"pool":    "30",
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		"process": "1",
	})


	go func() {
	for {
		time.Sleep(2*time.Second)
		workers.Enqueue("get_queue", "Get", map[string]interface{}{
			"harId":"abc",
			"harCotent":map[string]interface{}{
				"url":"http://baidu.com",
				"headers":"abc",
			},
		})
	}
	}()

	go func() {
		for {
			time.Sleep(3*time.Second)
			workers.Enqueue("post_queue", "Post", []int{1, 2})
		}
	}()


	go workers.StatsServer(8081)

	workers.Run()

	select {

	}
}