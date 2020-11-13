// Copyright The ZHIYUN Co. All rights reserved.
// Created by vinson on 2020/11/13.

package parsor

import (
	"annotation-parse/model"
	"annotation-parse/statement"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type YApi struct {
	Queue *[]*statement.Request
}

// 获取实例

func (r *YApi) New(condition *model.Condition) *statement.Context {
	r.Queue = &[]*statement.Request{}
	if v, ok := (interface{}(r)).(statement.Context); ok {
		return &v
	}
	return nil
}

// 获取请求队列
func (r *YApi) RequestQueue() *[]*statement.Request {
	if v, ok := (interface{}(*r.Queue)).([]*statement.Request); ok {
		return &v
	}
	return nil
}

// 获取Request
func (r *YApi) NewRequest(condition *model.Condition, content *model.Content) *statement.Request {
	if v, ok := (interface{}(&YApiRequest{condition, content})).(statement.Request); ok {
		return &v
	}
	return nil
}

// 发送文档创建请求
func (req *YApiRequest) Post() {
	br := strings.NewReader(req.build())
	request, err := http.NewRequest("POST", req.condition.URL, br)
	if nil != err {
		fmt.Println(err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(request)
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
	if result["errcode"].(float64) != 0 {
		rs = "Failed"
	}
	fmt.Println(req.content.Class+req.content.Title, "-", rs)
}

type YApiRequest struct {
	condition *model.Condition
	content   *model.Content
}

func (req YApiRequest) build() string {
	raw := make(map[string]interface{}, 0)
	raw["token"] = req.condition.Token
	raw["title"] = req.content.Title
	raw["catid"] = req.condition.Cat
	raw["path"] = req.content.URL
	raw["status"] = "done"        // TODO
	raw["res_body_type"] = "json" //TODO
	raw["res_body"] = req.content.Return
	raw["req_query"] = []interface{}{}
	raw["req_headers"] = []interface{}{}
	raw["req_body_form"] = []interface{}{}
	raw["desc"] = req.content.Description
	raw["method"] = req.content.Method
	raw["req_params"] = []interface{}{}
	bs, err := json.Marshal(raw)
	if err == nil {
		return string(bs)
	}
	fmt.Println(err)
	return ""
}
