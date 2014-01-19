package mailer

import(
	"net/smtp"
)

type Email struct {
	From 	string
	To 		string
    Bcc     string
	Subject string
	Msg 	string
}


func (email *Email) Send() error {
	// Set up authentication information.
    // @TODO: move it to Config
    auth := smtp.PlainAuth(
        "",
        "andrey.evsyukov@gmail.com",
        "parasite",
        "smtp.gmail.com",
        /*"info@prflr.org",
        "eshukun",
        "smtp.spaceweb.ru",*/
    )

    // Connect to the server, authenticate, set the sender and recipient,
    // and send the email all in one step.
	return smtp.SendMail(
        "smtp.gmail.com:587",
        //"smtp.spaceweb.ru:25",
        auth,
        email.From,
        []string{email.To},
        []byte(email.Msg),
    )
}