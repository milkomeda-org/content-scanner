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
	if v, ok := (interface{}(&Request{condition.URL, condition.Key, condition.Token, condition.Cat, r.template, content})).(statement.Request); ok {
		return &v
	}
	return nil
}

// 发送文档创建请求
func (req *Request) Post() {
	urlValues := url.Values{}
	urlValues.Add("api_key", req.Key)
	urlValues.Add("api_token", req.Token)
	if req._content.Class != "" {
		urlValues.Add("cat_name", req.Cat+req._content.Class)
	} else {
		urlValues.Add("cat_name", req.Cat)
	}
	urlValues.Add("page_title", req._content.Title)
	urlValues.Add("page_content", req.buildTemplate())
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
		fmt.Println(result)
		return
	}
	fmt.Println(req._content.Class+req._content.Title, "-", rs)
}

type Request struct {
	URL       string
	Key       string
	Token     string
	Cat       string
	_template *template.Template
	_content  *model.Content
}

func (req *Request) buildTemplate() string {
	content := &content{
		req._content.Catalog,
		req._content.Class,
		req._content.Title,
		req._content.Description,
		req._content.Method,
		req._content.URL,
		"",
		"",
		"",
		req._content.Return,
		"",
		req._content.Remark,
		req._content.Number,
	}
	for _, qs := range req._content.Header {
		q := `| ` + qs[0] + ` | ` + qs[1] + ` | ` + qs[2] + ` | ` + qs[3] + ` |` + "\n"
		content.Header = q
	}
	for _, qs := range req._content.Query {
		q := `| ` + qs[0] + ` | ` + qs[1] + ` | ` + qs[2] + ` | ` + qs[3] + ` |` + "\n"
		content.Query += q
	}
	for _, qs := range req._content.Body {
		q := `| ` + qs[0] + ` | ` + qs[1] + ` | ` + qs[2] + ` | ` + qs[3] + ` |` + "\n"
		content.Body += q
	}
	for _, qs := range req._content.ReturnParam {
		q := `| ` + qs[0] + ` | ` + qs[1] + ` | ` + qs[2] + ` |` + "\n"
		content.ReturnParam += q
	}
	var buf bytes.Buffer
	err := req._template.Execute(&buf, content)
	if nil != err {
		fmt.Println(err)
		return ""
	}
	return buf.String()
}

type content struct {
	Catalog     string // catalog xx/xx/
	Class       string // class string
	Title       string // title title
	Description string // description description
	Method      string // method method
	URL         string // url url
	Header      string // @header 名称 必选 类型 释义
	Query       string // @param 名称 必选 类型 释义
	Body        string // @body 名称 必选 类型 释义
	Return      string // @return str
	ReturnParam string // @return_param 名称 类型 释义
	Remark      string // @remark str
	Number      string // number int
}
