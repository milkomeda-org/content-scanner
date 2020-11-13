// Created by vinson on 2020/11/4.
// ShowDoc Scanner

package parsor

import (
	"annotation-parse/model"
	"annotation-parse/statement"
	"annotation-parse/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"text/template"
)

type ShowDoc struct {
	template *template.Template
	Queue    *[]*statement.Request
}

// 获取实例

func (r *ShowDoc) New(condition *model.Condition) *statement.Context {
	t, err := utils.GetTemplate(condition.TemplatePath)
	if nil != err {
		fmt.Println(err)
		panic("get template error")
	}
	r.template = t
	r.Queue = &[]*statement.Request{}
	if v, ok := (interface{}(r)).(statement.Context); ok {
		return &v
	}
	return nil
}

// 获取请求队列
func (r *ShowDoc) RequestQueue() *[]*statement.Request {
	if v, ok := (interface{}(*r.Queue)).([]*statement.Request); ok {
		return &v
	}
	return nil
}

// 获取Request
func (r *ShowDoc) NewRequest(condition *model.Condition, content *model.Content) *statement.Request {
	var buf bytes.Buffer
	err := r.template.Execute(&buf, content)
	if nil != err {
		fmt.Println(err)
		return nil
	}
	if v, ok := (interface{}(&Request{condition.URL, condition.Key, condition.Token, condition.Cat, content.Class, content.Title, buf.String()})).(statement.Request); ok {
		return &v
	}
	return nil
}

// 发送文档创建请求
func (req *Request) Post() {
	urlValues := url.Values{}
	urlValues.Add("api_key", req.Key)
	urlValues.Add("api_token", req.Token)
	if req.Class != "" {
		urlValues.Add("cat_name", req.Cat+req.Class)
	} else {
		urlValues.Add("cat_name", req.Cat)
	}
	urlValues.Add("page_title", req.Title)
	urlValues.Add("page_content", req.Content)
	var resp, err = http.PostForm(req.URL, urlValues)
	if nil != err {
		fmt.Println(err)
		return
	}
	defer func() {
		err := resp.Body.Close()
		if nil != err {
			fmt.Println(err)
			return
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		fmt.Println(err)
		return
	}
	var result map[string]interface{}
	_ = json.Unmarshal(body, &result)
	rs := "Ok"
	if result["error_code"].(float64) != 0 {
		rs = "Failed"
	}
	fmt.Println(req.Class+req.Title, "-", rs)
}

type Request struct {
	URL     string
	Key     string
	Token   string
	Cat     string
	Class   string
	Title   string
	Content string
}
