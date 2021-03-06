// Created by vinson on 2020/11/13.

package model

type Content struct {
	Catalog     string      // catalog xx/xx/
	Class       string      // class string
	Title       string      // title title
	Description string      // description description
	Method      string      // method method
	URL         string      // url url
	Header      [][4]string // @header 名称 必选 类型 释义
	Query       [][4]string // @param 名称 必选 类型 释义
	Body        [][4]string // @body 名称 必选 类型 释义
	Return      string      // @return str
	ReturnParam [][3]string // @return_param 名称 类型 释义
	Remark      string      // @remark str
	Number      string      // number int
}

type Condition struct {
	T            string // doc type
	URL          string // post url
	Key          string // post auth key
	Token        string // post auth token
	Cat          string // pre path
	Speed        int    // post multi
	SearchPath   string // scanPath
	TemplatePath string // doc template engine
	Ask          bool   // not skip option options
}
