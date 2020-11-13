// Created by vinson on 2020/11/4.

package statement

import "annotation-parse/model"

type Context interface {
	New(condition *model.Condition) *Context
	NewRequest(condition *model.Condition, content *model.Content) *Request
	RequestQueue() *[]*Request
}
