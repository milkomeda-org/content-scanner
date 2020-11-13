// Created by vinson on 2020/10/30.

package main

import (
	"annotation-parse/model"
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
	"time"
)

// 带有此标识的匹配注释才会进行文档解析
const CiFlag = "#doc"

// 支持的文件格式
var ExtSupported = []string{".go", ".js", ".java", ".php", ".py"}

func main() {
	t1 := time.Now()
	defer func() {
		t2 := time.Now()
		fmt.Println(t2.Sub(t1))
	}()
	var condition model.Condition
	flag.StringVar(&condition.T, "type", "showdoc", "Document system type, example for showdoc, yapi?")
	flag.StringVar(&condition.URL, "url", "http://172.16.2.101:4999/server/index.php?s=/api/item/updateByApi", "show doc api")
	flag.StringVar(&condition.Key, "key", "e9f0bdd396a768399c63ef86d70ccc322044412143", "show doc api_key")
	flag.StringVar(&condition.Token, "token", "834e06eb69e21565d997cf15a1159da21794468976", "show doc api_token")
	flag.StringVar(&condition.Cat, "cat", "", "doc cat, End with second /")
	flag.StringVar(&condition.SearchPath, "searchPath", "./controller", "search path")
	flag.StringVar(&condition.TemplatePath, "templatePath", "./template.txt", "customer template file path")
	flag.IntVar(&condition.Speed, "speed", 1, "for Concurrent requests")
	flag.BoolVar(&condition.Ask, "ask", true, "Ask first, then execute the program")
	flag.Parse()
	decoded, _ := base64.StdEncoding.DecodeString(condition.Cat)
	condition.Cat = string(decoded)
	var head = func() {
		fmt.Println(fmt.Sprintf("#\n"+
			"# start\n"+
			"# type: %s \n"+
			"# url: %s\n"+
			"# cat: %s\n"+
			"# search_path: %s\n"+
			"# template: %s\n"+
			"# multipost: %d\n"+
			"#",
			condition.T, condition.URL, condition.Cat, condition.SearchPath, condition.TemplatePath, condition.Speed))
	}
	var ps *statement.Context
	// get context
	switch condition.T {
	case "showdoc":
		ps = (&parsor.ShowDoc{}).New(&condition)
	case "yapi":
		ps = (&parsor.YApi{}).New(&condition)
	default:
		panic("missing document type")
	}
	// parse args
	var rq = (*ps).RequestQueue()
	var wg sync.WaitGroup
	if ok, fileInfo := utils.IsDir(condition.SearchPath); ok {
		// read the file list of searchPath and be sort they by modify time
		// wait for user select any options to continue
		files, err := utils.ReadDirOrderByModify(condition.SearchPath)
		if nil != err {
			fmt.Println(err)
			return
		}
		selection := -1
		if condition.Ask {
			selection = utils.Selection(&files)
		}
		head()
		if selection == -1 {
			Scan(ps, &files, rq, &condition, &wg, condition.SearchPath)
		} else {
			ParseFile(ps, &(files[selection]), rq, &condition, &wg, true, condition.SearchPath)
		}
	} else {
		// 单文件直接执行
		ParseFile(ps, fileInfo, rq, &condition, &wg, false, condition.SearchPath)
	}
	// redundant
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
func Scan(ctx *statement.Context, files *[]os.FileInfo, fs *[]*statement.Request, condition *model.Condition, wg *sync.WaitGroup, parentPath string) {
	for _, f := range *files {
		if f.IsDir() {
			p := utils.PathAssemble(parentPath, f.Name())
			files2, err := ioutil.ReadDir(p)
			if nil != err {
				fmt.Println(err)
				return
			}
			Scan(ctx, &files2, fs, condition, wg, p)
		} else {
			ParseFile(ctx, &f, fs, condition, wg, true, parentPath)
		}
	}
}

// 解析单个文件
func ParseFile(ctx *statement.Context, f *os.FileInfo, fs *[]*statement.Request, condition *model.Condition, wg *sync.WaitGroup, IsDir bool, parentPath string) {
	var fp string
	if IsDir {
		fp = utils.PathAssemble(parentPath, (*f).Name())
	} else {
		fp = parentPath
	}
	if !utils.IsContain(ExtSupported, path.Ext(fp)) {
		fmt.Println("file type is not supported")
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
			context, err := Parse(part[0])
			if nil != err {
				continue
			}
			v := (*ctx).NewRequest(condition, context)
			*fs = append(*fs, v)
			if len(*fs) >= condition.Speed {
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
		}
	} else {
		fmt.Println("reg initialize Failed")
	}
}

// the function can parse the annotation str to a Content obj
func Parse(annotation string) (*model.Content, error) {
	// 解析@
	if !strings.Contains(annotation, CiFlag) {
		return nil, errors.New("no doc annotation")
	}
	tagReg := regexp.MustCompile(`@[\w\W]*?\n`)
	tagLine := tagReg.FindAllStringSubmatch(annotation, -1)
	var content = &model.Content{}
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
