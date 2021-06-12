package main

import (
	"fmt"
	"goRent/internal/model"
	"io/ioutil"
	"strings"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	go func() {
		for {
			sendMsg(<-app.MailChan)
		}
	}()
}

func sendMsg(m model.MailData) {

	client := mail.NewSMTPClient()
	client.Host = *mailhost
	client.Port = *mailport
	client.Username = *mailUsername
	client.Password = *mailPassword

	client.Encryption = mail.EncryptionSTARTTLS
	client.KeepAlive = false
	client.ConnectTimeout = 10 * time.Second
	client.SendTimeout = 10 * time.Second

	smtpClient, err := client.Connect()
	if err != nil {
		app.Error.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template != "" {
		data, err := ioutil.ReadFile(fmt.Sprintf("./templates/email/%s", m.Template))
		if err != nil {
			app.Error.Println("Error reading email template data: ", err)
		}
		mailTemplate := string(data)
		body := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		m.Content = body
	}

	email.SetBody(mail.TextHTML, m.Content)

	if err := email.Send(smtpClient); err != nil {
		app.Error.Println(err)
	}
	app.Info.Println("mail sent")
}
