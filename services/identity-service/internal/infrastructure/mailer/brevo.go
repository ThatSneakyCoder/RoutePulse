package mailer

import (
	"bytes"
	"html/template"

	"gopkg.in/gomail.v2"
)

type BrevoClient struct {
	host      string
	port      int
	username  string
	password  string
	fromEmail string
}

func NewBrevoClient(host string, port int, username, password, from string) *BrevoClient {
	return &BrevoClient{
		host:      host,
		port:      port,
		username:  username,
		password:  password,
		fromEmail: from,
	}
}

func (m *BrevoClient) Send(templateFile, to string, data any) error {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return err
	}

	body := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(body, "body", data); err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", m.fromEmail)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/html", body.String())

	d := gomail.NewDialer(
		m.host,
		m.port,
		m.username,
		m.password,
	)

	return d.DialAndSend(msg)
}
