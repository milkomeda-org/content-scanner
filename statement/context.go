// Created by vinson on 2020/11/4.

package statement

type Context interface {
	NewRequest(URL, key, token, cat, class, title, buf string) *Request
	RequestQueue() *[]*Request
}
