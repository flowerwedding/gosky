package gosky

import "testing"

func TestEmail(t *testing.T) {
	to := "2965502421@qq.com"
	subject := "gosky"
	message := "Hello Flowerwedding"
	Email(to,subject,message)
}

func TestSendEmail(t *testing.T) {

}