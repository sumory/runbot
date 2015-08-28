package common

import (
	"time"
)

type KV struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type GetContent struct {
	Method      string `json:"method" bson:"method"`
	Url         string `json:"url" bson:"url"`
	HttpVersion string `json:"httpVersion" bson:"httpVersion"`
	QueryString []KV `json:"queryString" bson:"queryString"`
	Headers     []KV `json:"headers" bson:"headers"`
	Cookies     []KV `json:"cookies" bson:"cookies"`
}

type PostData struct {
	Params   []KV
	MimeType string
	Text     string
}

type PostContent struct {
	Method      string `json:"method" bson:"method"`
	Url         string `json:"url" bson:"url"`
	HttpVersion string `json:"httpVersion" bson:"httpVersion"`
	QueryString []KV `json:"queryString" bson:"queryString"`
	Headers     []KV `json:"headers" bson:"headers"`
	Cookies     []KV `json:"cookies" bson:"cookies"`
	PostData    PostData `json:"postData" bson:"postData"`
}

type StatusAPI struct {
	Id          string `json:"id" bson:"id"`
	Monitor     bool `json:"monitor" bson:"monitor"`
	Cron        string `json:"cron" bson:"cron"`
	UserId      string `json:"userId" bson:"userId"`
	Name        string `json:"name" bson:"name"`
	Type        string `json:"type" bson:"type"`
	Content     interface{} `json:"content" bson:"content"` //GetContent or PostContent
	Date        time.Time `json:"date" bson:"date"`
	GetContent  GetContent
	PostContent PostContent
}

//status api请求响应日志
type StatusAPILog struct {
	StatusAPIId string  `json:"statusAPIId" bson:"statusAPIId"`
	StatusCode int `json:"statusCode" bson:"statusCode"`
	Spent int64  `json:"spent" bson:"spent"`
	UserId string `json:"userId" bson:"userId"`
	Response string `json:"response" bson:"response"`
	Date time.Time `json:"date" bson:"date"`
}

