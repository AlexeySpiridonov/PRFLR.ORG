package mailer

import(
	"net/smtp"
)

type Email struct {
	From 	string
	To 		string
	Subject string
	Msg 	string
}


func (email *Email) Send() error {
	// Set up authentication information.
    auth := smtp.PlainAuth(
        "",
        "andrey.evsyukov@gmail.com",
        "parasite",
        "smtp.gmail.com",
    )

    // Connect to the server, authenticate, set the sender and recipient,
    // and send the email all in one step.
	return smtp.SendMail(
        "smtp.gmail.com:587",
        auth,
        email.From,
        []string{email.To},
        []byte(email.Msg),
    )
}