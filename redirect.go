package gosky

import (
	"fmt"
	"net/http"
)

type Redirect struct {
	Code     int
	Request  *http.Request
	Location string
}

func (c *Context) Redirect(code int, location string) {
	c.Render = Redirect{
		Code:     code,
		Location: location,
		Request:  c.Req,
	}

	//因为一直报 http: superfluous response.WriteHeader call from 的错误，所以就没有 c.Writer.WriteHeader(code)，没找到原因
	if err := c.Render.Render(c.Writer); err != nil {
		panic(err)
	}
}

func (r Redirect) Render(w http.ResponseWriter) error {
	if (r.Code < http.StatusMultipleChoices || r.Code > http.StatusPermanentRedirect) && r.Code != http.StatusCreated {
		panic(fmt.Sprintf("Cannot redirect with status code %d", r.Code))
	}
	http.Redirect(w, r.Request, r.Location, r.Code)
	return nil
}