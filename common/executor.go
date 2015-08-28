package common
import (
	"net/http"
	"io/ioutil"
	"fmt"
	"errors"
	"net/url"
	"strings"
	"io"
)


func RunStatusAPI(api *StatusAPI) (err error, h *StatusAPI, statusCode int, responseBody string) {
	if api.Type=="GET" {
		e, s, r := runGetHar(api)
		return e, api, s, strings.TrimRight(r,"\n")
	}else if api.Type=="POST" {
		e, s, r := runPostHar(api)
		return e, api, s,  strings.TrimRight(r,"\n")
	}else {
		return errors.New("unsupports http method"), api, 0, ""
	}
}

func runGetHar(api *StatusAPI) (err error, statusCode int, result string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("##RunGetHar error: %v\n", r)
			if _, ok := r.(error); ok {
				err = r.(error)
			}else {
				err = errors.New("runGetHar undefined error")
			}
		}
	}()

	client := &http.Client{}
	getContent := api.GetContent

	req, errRequest := http.NewRequest("GET", getContent.Url, nil)
	if errRequest!=nil {
		return errRequest, 0, ""
	}

	if getContent.Headers!=nil && len(getContent.Headers)>0 {
		var headers []KV = getContent.Headers
		for _, header := range headers {
			req.Header.Set(header.Name, header.Value)
		}
	}

	if getContent.Cookies!=nil && len(getContent.Cookies)>0 {
		var cookies []KV = getContent.Cookies
		var cookieStr = ""
		for _, cookie := range cookies {
			cookieStr+=cookie.Name+"="+cookie.Value+";"
		}
		req.Header.Set("cookie", cookieStr)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, errReadAll := ioutil.ReadAll(resp.Body)

	if errReadAll!=nil {
		return errReadAll, 0, ""
	}


	return nil, resp.StatusCode, string(body)
}



func runPostHar(api *StatusAPI) (err error, statusCode int, result string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("##RunPostHar error: %v \n", r)
			if _, ok := r.(error); ok {
				err = r.(error)
			}else {
				err = errors.New("runPostHar undefined error")
			}
		}
	}()

	client := &http.Client{}
	postContent := api.PostContent

	ps := url.Values{}
	var body io.ReadCloser

	var postData = postContent.PostData
	if postData.MimeType == "application/x-www-form-urlencoded" ||  postData.MimeType == "multipart/form-data" {
		if postData.Params!=nil && len(postData.Params)>0 {
			var params []KV = postData.Params
			for _, p := range params {
				ps.Set(p.Name, p.Value)
			}
		}
		body = ioutil.NopCloser(strings.NewReader(ps.Encode())) //把form数据编下码
	}else if postData.MimeType=="application/json" {
		body = ioutil.NopCloser(strings.NewReader(postData.Text))
	}


	req, errRequest := http.NewRequest("POST", postContent.Url, body)
	if errRequest!=nil {
		return errRequest, 0, ""
	}


	if postContent.Headers!=nil && len(postContent.Headers)>0 {
		var headers []KV = postContent.Headers
		for _, header := range headers {
			req.Header.Set(header.Name, header.Value)
		}
	}

	if postContent.Cookies!=nil && len(postContent.Cookies)>0 {
		var cookies []KV = postContent.Cookies
		var cookieStr = ""
		for _, cookie := range cookies {
			cookieStr+=cookie.Name+"="+cookie.Value+";"
		}
		req.Header.Set("cookie", cookieStr)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	r, errReadAll := ioutil.ReadAll(resp.Body)

	if errReadAll!=nil {
		return errReadAll, 0, ""
	}


	return nil, resp.StatusCode, string(r)
}