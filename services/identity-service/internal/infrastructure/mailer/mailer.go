package mailer

import "embed"

type SMTPMailer struct {
	host     string
	port     int
	username string
	password string
	from     string
}

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, email string, data any) error
}
