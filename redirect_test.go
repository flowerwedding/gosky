package gosky

import (
	"testing"
)

func TestRedirect_Render(t *testing.T) {
	r := Default()

	r.GET("/test", func(c *Context) {
		c.Redirect(302, "https://mail.qq.com")
	})

	_ = r.Run(":5523")
}