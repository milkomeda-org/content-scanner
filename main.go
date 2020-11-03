// Copyright The ZHIYUN Co. All rights reserved.
// Created by vinson on 2020/10/30.

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"text/template"
)

// 带有此标识的匹配注释才会进行文档解析
const CiFlag = "#doc"

// 支持的文件格式
var ExtSupported = []string{".go", ".js"}

var URL string
var key string
var token string
var cat string
var speed int
var searchPath string
var templatePath string

func main() {
	flag.StringVar(&URL, "url", "http://172.16.2.101:4999/server/index.php?s=/api/item/updateByApi", "show doc api")
	flag.StringVar(&key, "key", "e9f0bdd396a768399c63ef86d70ccc322044412143", "show doc api_key")
	flag.StringVar(&token, "token", "834e06eb69e21565d997cf15a1159da21794468976", "show doc api_token")
	flag.StringVar(&cat, "cat", "", "doc cat, End with second /")
	flag.StringVar(&searchPath, "searchPath", "./controller", "search path")
	flag.StringVar(&templatePath, "templatePath", "./template.txt", "customer template file path")
	flag.IntVar(&speed, "speed", 1, "for Concurrent requests")
	flag.Parse()
	fmt.Println("start")
	var fs []*Request
	var wg sync.WaitGroup
	t, err := getTemplate()
	if nil != err {
		fmt.Println(err)
		return
	}
	if ok, fileInfo := IsDir(searchPath); ok {
		files, err := ioutil.ReadDir(searchPath)
		if nil != err {
			fmt.Println(err)
			return
		}
		Scan(&files, &fs, t, &wg, searchPath)
	}else {
		ParseFile(fileInfo, &fs, t, &wg, false, searchPath)
	}
	if len(fs) > 0 {
		wg.Add(len(fs))
		for _, request := range fs {
			request := request
			go func() {
				Post(*request)
				wg.Done()
			}()
		}
		wg.Wait()
		fs = fs[:0]
	}
	fmt.Println("end")
}

// 递归扫描
func Scan(files *[]os.FileInfo, fs *[]*Request, t *template.Template, wg *sync.WaitGroup, parentPath string) {
	for _, f := range *files {
		if f.IsDir() {
			p := pathAssemble(parentPath, f.Name())
			files2, err := ioutil.ReadDir(p)
			if nil != err {
				fmt.Println(err)
				return
			}
			Scan(&files2, fs, t, wg, p)
		}else {
			ParseFile(&f, fs, t, wg, true, parentPath)
		}
	}
}

func ParseFile(f *os.FileInfo, fs *[]*Request, t *template.Template, wg *sync.WaitGroup, IsDir bool, parentPath string) {
	var fp string
	if IsDir {
		fp = pathAssemble(parentPath, (*f).Name())
	}else {
		fp = parentPath
	}
	if !IsContain(ExtSupported, path.Ext(fp)) {
		return
	}
	content, err := ReadAll(fp)
	if nil != err {
		fmt.Println(err)
		return
	}
	reg := regexp.MustCompile(`/\*\*[\w\W]*?\*/`)
	if reg != nil {
		s := reg.FindAllStringSubmatch(content, -1)
		if l := len(s); l > 0 {
			fmt.Println((*f).Name(), "*", l)
		}
		for _, part := range s {
			con,err := Parse(part[0])
			if nil != err {
				continue
			}
			var buf bytes.Buffer
			_ = t.Execute(&buf, con)
			func() {
				//开协程处理
				decoded, _ := base64.StdEncoding.DecodeString(cat)
				decodecat := string(decoded)
				v := &Request{URL, key, token, decodecat, con.Class, con.Title, buf.String()}
				*fs = append(*fs, v)
				if len(*fs) >= 5 {
					wg.Add(len(*fs))
					for _, request := range *fs {
						request := request
						go func() {
							Post(*request)
							wg.Done()
						}()
					}
					wg.Wait()
					*fs = (*fs)[:0]
				}
			}()
		}
	}
}

// 读取文件字符串
func ReadAll(filePth string) (string, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return "", err
	}
	readAll, err := ioutil.ReadAll(f)
	if nil == err {
		return string(readAll), err
	}
	return "", err
}

func getTemplate() (*template.Template, error) {
	files := []string{templatePath}
	return template.ParseFiles(files...)
}

func Post(req Request) {
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

type Content struct {
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

func Parse(annotation string) (*Content,error){
	// 解析@
	if !strings.Contains(annotation, CiFlag) {
		return nil, errors.New("no doc annotation")
	}
	tagReg := regexp.MustCompile(`@[\w\W]*?\n`)
	tagLine := tagReg.FindAllStringSubmatch(annotation, -1)
	var content = &Content{}
	for _, tv := range tagLine {
		tv[0] = strings.ReplaceAll(tv[0], "\r", "")
		tv[0] = strings.ReplaceAll(tv[0], "\n", "")
		if i := strings.LastIndex(tv[0], "@class "); i != -1 {
			// 名称
			content.Class = tv[0][i+7:]
			continue
		}
		if i := strings.LastIndex(tv[0], "@title "); i != -1 {
			// 名称
			content.Title = tv[0][i+7:]
			continue
		}
		if i := strings.LastIndex(tv[0], "@description "); i != -1 {
			// 描述
			content.Description = tv[0][i+13:]
			continue
		}
		if i := strings.LastIndex(tv[0], "@method "); i != -1 {
			// 方法
			content.Method = tv[0][i+8:]
			continue
		}
		if i := strings.LastIndex(tv[0], "@url "); i != -1 {
			// url
			content.URL = tv[0][i+5:]
			continue
		}
		if i := strings.LastIndex(tv[0], "@return "); i != -1 {
			// return str
			str := tv[0][i+8:]
			sb := []byte(str)
			if json.Valid(sb) {
				var sj bytes.Buffer
				_ = json.Indent(&sj, sb, "", "\t")
				content.Return = sj.String()
			} else {
				content.Return = str
			}
			continue
		}
		if i := strings.LastIndex(tv[0], "@remark "); i != -1 {
			// remark
			content.Remark = tv[0][i+8:]
			continue
		}
		if i := strings.LastIndex(tv[0], "@number "); i != -1 {
			// number
			content.Number = tv[0][i+8:]
			continue
		}
		if i := strings.LastIndex(tv[0], "@header "); i != -1 {
			// header
			func() {
				defer func() {
					if recover() != nil {
						return
					}
				}()
				qs := strings.Split(tv[0][i+8:], " ")
				q := `| ` + qs[0] + ` | ` + qs[1] + ` | ` + qs[2] + ` | ` + qs[3] + ` |` + "\n"
				content.Header += q
			}()
			continue
		}
		if i := strings.LastIndex(tv[0], "@query "); i != -1 {
			// param
			func() {
				defer func() {
					if recover() != nil {
						return
					}
				}()
				qs := strings.Split(tv[0][i+7:], " ")
				q := `| ` + qs[0] + ` | ` + qs[1] + ` | ` + qs[2] + ` | ` + qs[3] + ` |` + "\n"
				content.Query += q
			}()
			continue
		}
		if i := strings.LastIndex(tv[0], "@body "); i != -1 {
			// param
			func() {
				defer func() {
					if recover() != nil {
						return
					}
				}()
				qs := strings.Split(tv[0][i+6:], " ")
				q := `| ` + qs[0] + ` | ` + qs[1] + ` | ` + qs[2] + ` | ` + qs[3] + ` |` + "\n"
				content.Body += q
			}()
			continue
		}
		if i := strings.LastIndex(tv[0], "@return_param "); i != -1 {
			// return_param
			func() {
				defer func() {
					if recover() != nil {
						return
					}
				}()
				qs := strings.Split(tv[0][i+14:], " ")
				q := `| ` + qs[0] + ` | ` + qs[1] + ` | ` + qs[2] + ` |` + "\n"
				content.ReturnParam += q
			}()
			continue
		}
	}
	return content, nil
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func IsDir(path string) (bool, *os.FileInfo) {
	s, err := os.Stat(path)
	if err != nil {
		return false, nil
	}
	return s.IsDir(), &s
}

func pathAssemble(p, f string) string{
	if path.Base(p) == "/" || path.Base(searchPath) == "\\" {
		return p + f
	} else {
		return p + "/" + f
	}
}
