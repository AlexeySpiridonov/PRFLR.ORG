package mailer

import (
	"fmt"
	"net/mail"
	"net/smtp"
)

type Email struct {
	From    string
	To      string
	Bcc     string
	Subject string
	Msg     string
}

func (email *Email) Send() error {
	// Set up authentication information.
	// @TODO: move it to Config
	auth := smtp.PlainAuth(
		"",
		"no-reply@2hive.org",
		"robot06539010",
		"smtp.gmail.com",
	)

	from := mail.Address{"PRFLR Team", email.From}
	to := mail.Address{"", email.To}

	// setup a map for the headers
	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = email.Subject

	// setup the message
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + email.Msg

	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		email.From,
		[]string{email.To},
		[]byte(message),
	)
}
