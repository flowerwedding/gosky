package gosky

import (
	"gopkg.in/gomail.v2"
)

func Email(to string,subject string,message string) {
	m := gomail.NewMessage()

	m.SetHeader("From", "2804991212@qq.com")
	m.SetHeader("To", to)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)
	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.qq.com", 587, "2804991212@qq.com", "xygdhlezsvirdebb")


	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}