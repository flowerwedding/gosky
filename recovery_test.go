package gosky

import (
	"net/http"
	"testing"
)

func TestRecovery(t *testing.T) {
	r := Default()

	r.GET("/", func(c *Context) {
		c.String(http.StatusOK, "Hello Geektutu\n")
	})

	r.GET("/panic", func(c *Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	_ = r.Run(":9999")
}