// Copyright The ef Co. ltd All rights reserved.
// Created by vinson on 2020/10/30.

package main

import (
	"annotation-parse/parsor"
	"annotation-parse/statement"
	"annotation-parse/utils"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"
)

// 带有此标识的匹配注释才会进行文档解析
const CiFlag = "#doc"

// 支持的文件格式
var ExtSupported = []string{".go", ".js", ".java", ".php", ".py"}

var t string
var URL string
var key string
var token string
var cat string
var speed int
var searchPath string
var templatePath string

func main() {
	t1 := time.Now()
	defer func() {
		t2 := time.Now()
		fmt.Println(t2.Sub(t1))
	}()
	flag.StringVar(&t, "type", "showdoc", "Document system type, example for showdoc, yapi?")
	flag.StringVar(&URL, "url", "http://172.16.2.101:4999/server/index.php?s=/api/item/updateByApi", "show doc api")
	flag.StringVar(&key, "key", "e9f0bdd396a768399c63ef86d70ccc322044412143", "show doc api_key")
	flag.StringVar(&token, "token", "834e06eb69e21565d997cf15a1159da21794468976", "show doc api_token")
	flag.StringVar(&cat, "cat", "", "doc cat, End with second /")
	flag.StringVar(&searchPath, "searchPath", "./controller", "search path")
	flag.StringVar(&templatePath, "templatePath", "./template.txt", "customer template file path")
	flag.IntVar(&speed, "speed", 1, "for Concurrent requests")
	flag.Parse()
	fmt.Println(fmt.Sprintf("#\n"+
		"# start\n"+
		"# type: %s \n"+
		"# url: %s\n"+
		"# cat: %s\n"+
		"# search_path: %s\n"+
		"# template: %s\n"+
		"# multipost: %d\n"+
		"#",
		t, URL, cat, searchPath, templatePath, speed))
	var ps *statement.Context
	// get context
	switch t {
	case "showdoc":
		ps = parsor.NewShowDoc()
	default:
		panic("missing document type")
	}
	// parse args
	var rq = (*ps).RequestQueue()
	var wg sync.WaitGroup
	t, err := utils.GetTemplate(templatePath)
	if nil != err {
		fmt.Println(err)
		return
	}
	if ok, fileInfo := utils.IsDir(searchPath); ok {
		files, err := ioutil.ReadDir(searchPath)
		if nil != err {
			fmt.Println(err)
			return
		}
		Scan(ps, &files, rq, t, &wg, searchPath)
	} else {
		ParseFile(ps, fileInfo, rq, t, &wg, false, searchPath)
	}
	if len(*rq) > 0 {
		wg.Add(len(*rq))
		for _, request := range *rq {
			request := request
			go func() {
				(*request).Post()
				wg.Done()
			}()
		}
		wg.Wait()
		*rq = (*rq)[:0]
	}
	fmt.Println("# end")
}

// 递归扫描文件夹
func Scan(ctx *statement.Context, files *[]os.FileInfo, fs *[]*statement.Request, t *template.Template, wg *sync.WaitGroup, parentPath string) {
	for _, f := range *files {
		if f.IsDir() {
			p := utils.PathAssemble(parentPath, f.Name())
			files2, err := ioutil.ReadDir(p)
			if nil != err {
				fmt.Println(err)
				return
			}
			Scan(ctx, &files2, fs, t, wg, p)
		} else {
			ParseFile(ctx, &f, fs, t, wg, true, parentPath)
		}
	}
}

// 解析单个文件
func ParseFile(ctx *statement.Context, f *os.FileInfo, fs *[]*statement.Request, t *template.Template, wg *sync.WaitGroup, IsDir bool, parentPath string) {
	var fp string
	if IsDir {
		fp = utils.PathAssemble(parentPath, (*f).Name())
	} else {
		fp = parentPath
	}
	if !utils.IsContain(ExtSupported, path.Ext(fp)) {
		return
	}
	content, err := utils.ReadAll(fp)
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
			con, err := Parse(part[0])
			if nil != err {
				continue
			}
			var buf bytes.Buffer
			_ = t.Execute(&buf, con)
			func() {
				//开协程处理
				decoded, _ := base64.StdEncoding.DecodeString(cat)
				dCat := string(decoded)
				v := (*ctx).NewRequest(URL, key, token, dCat, con.Class, con.Title, buf.String())
				*fs = append(*fs, v)
				if len(*fs) >= speed {
					wg.Add(len(*fs))
					for _, request := range *fs {
						request := request
						go func() {
							(*request).Post()
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

// 将一条注释解析成Content
func Parse(annotation string) (*Content, error) {
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
