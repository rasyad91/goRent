package main

import (
	"fmt"
	"goRent/internal/model"
	"io/ioutil"
	"strings"
	"time"

	mail "github.com/xhit/go-simple-mail"
)

func listenForMail() {
	go func() {
		for {
			sendMsg(<-app.MailChan)
		}
	}()
}

func sendMsg(m model.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		app.Error.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template != "" {
		data, err := ioutil.ReadFile(fmt.Sprintf("../../templates/email/%s", m.Template))
		if err != nil {
			app.Error.Println("Error reading email template data")
		}
		mailTemplate := string(data)
		body := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		m.Content = body
	}

	email.SetBody(mail.TextHTML, m.Content)

	if err := email.Send(client); err != nil {
		app.Error.Println(err)
	}
}
