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

type Har struct {
	Monitor bool `json:"monitor" bson:"monitor"`
	HarId        string `json:"harId" bson:"harId"`
	UserId       string `json:"userId" bson:"userId"`
	CollectionId string `json:"collectionId" bson:"collectionId"`
	Name         string `json:"name" bson:"name"`
	Type         string `json:"type" bson:"type"`
	Content      interface{} `json:"content" bson:"content"` //GetContent or PostContent
	Date         time.Time `json:"date" bson:"date"`
	GetContent   GetContent
	PostContent  PostContent
}

