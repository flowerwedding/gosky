package gosky

import (
	"io"
	"mime/multipart"
	"os"
)

//文件上传
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {//单个文件
	if c.Req.MultipartForm == nil {
		if err := c.Req.ParseMultipartForm(c.MaxMultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := c.Req.FormFile(name)
	if err != nil {
		return nil, err
	}
	_ = f.Close()
	return fh, err
}

func (c *Context) MultipartForm() (*multipart.Form, error) {//多个文件
	err := c.Req.ParseMultipartForm(c.MaxMultipartMemory)
	return c.Req.MultipartForm, err
}

func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {//上传到指定路径
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}