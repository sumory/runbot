package test


import (
	"github.com/jrallison/go-workers"
	"time"
	"encoding/json"
	"github.com/sumory/runbot/common"
)

var getHar = `
{
	"method": "GET",
	"url": "http:192.168.100.122:8010/intersect?type=1&uid=2&targets=3",
	"httpVersion": "HTTP/1.1",
	"queryString": [
		{
			"name": "type",
			"value": "1"
		},
		{
			"name": "uid",
			"value": "2"
		},
		{
			"name": "targets",
			"value": "3"
		}
	],
	"headers": [
	{
		"name": "Accept",
		"value": "*/*"
	}
	],
	"cookies": []
}
`

var postHar = `
{
    "method": "POST",
    "url": "http:192.168.100.122:8001/user/save",
    "httpVersion": "HTTP/1.1",
    "queryString": [],
    "headers": [
        {
            "name": "Content-type",
            "value": "application/x-www-form-urlencoded"
        },
        {
            "name": "Accept",
            "value": "*/*"
        }
    ],
    "cookies": [],
    "postData": {
        "mimeType": "application/x-www-form-urlencoded",
        "params": [
            {
                "name": "name",
                "value": "ss"
            },
            {
                "name": "age",
                "value": "45"
            }
        ]
    }
}
`


func test() {
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
			var getContent = new (common.GetContent)
			json.Unmarshal([]byte(getHar), getContent )
			workers.Enqueue("get_queue", "Get", map[string]interface{}{
				"harId":"abc",
				"userId":"user123",
				"collectionId":"collection1",
				"name":"/user/save",
				"content": getContent,
			})
		}
	}()

	//	go func() {
	//		for {
	//			time.Sleep(3*time.Second)
	//			workers.Enqueue("post_queue", "Post", []int{1, 2})
	//		}
	//	}()


	go workers.StatsServer(8081)

	workers.Run()

	select {

	}
}