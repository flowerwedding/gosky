package gosky

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"testing"
)

func TestContext_FormFile(t *testing.T) {
	router := Default()

	router.POST("/upload", func(c *Context) {
		file, err := c.FormFile("file")
		log.Println(file.Filename)

		if err != nil{
			c.String(http.StatusOK, err.(interface{}).(string))
		}else{
			dst := path.Join("./", file.Filename)//上传文件时路径必须有
			_  = c.SaveUploadedFile(file,dst)
			c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
		}
	})
	_ = router.Run(":8080")
}