package gosky

import (
	"fmt"
	"net/http"
)

func (c *Context)Redirect(code int, location string) {
	//因为一直报 http: superfluous response.WriteHeader call from 的错误，所以就没有 c.Writer.WriteHeader(code)，没找到原因
	if (code < http.StatusMultipleChoices || code > http.StatusPermanentRedirect) && code != http.StatusCreated {
		panic(fmt.Sprintf("Cannot redirect with status code %d", code))
	}
	http.Redirect(c.Writer, c.Req, location, code)
}