// Created by vinson on 2020/11/13.

package utils

import (
	"annotation-parse/model"
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

// 带有此标识的匹配注释才会进行文档解析
const CiFlag = "#doc"

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
				q := [4]string{qs[0], qs[1], qs[2], qs[3]}
				content.Header = append(content.Header, q)
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
				q := [4]string{qs[0], qs[1], qs[2], qs[3]}
				content.Query = append(content.Query, q)
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
				q := [4]string{qs[0], qs[1], qs[2], qs[3]}
				content.Body = append(content.Body, q)
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
				q := [3]string{qs[0], qs[1], qs[2]}
				content.ReturnParam = append(content.ReturnParam, q)
			}()
			continue
		}
	}
	return content, nil
}
