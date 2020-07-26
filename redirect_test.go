package gosky

import (
	"net/http"
	"testing"
)

func TestRedirect_Render(t *testing.T) {
	r := Default()

	r.GET("/test", func(c *Context) {
		c.Redirect(http.StatusMovedPermanently, "https://mail.qq.com")
	})

	_ = r.Run(":1234")
}