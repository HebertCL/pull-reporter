package main

import (
	"bytes"
	"fmt"
	"html/template"
)

type senderConfig struct {
	Server   string
	User     string
	Password string
	SmtpPort int
}

type emailData struct {
	Repository    repository
	OpenPulls     string
	ClosedPulls   string
	DraftPulls    string
	RecipientName string
}

// Creates a SMTP plain authentication
// TODO: Improve email authtentication
// func emailAuthentication(user string, password string, server string) smtp.Auth {
// 	auth := smtp.PlainAuth("", user, password, server)

// 	return auth
// }

func generateTemplateEmail(data emailData) (string, error) {
	message := `
		Subject: GitHub PR Report
		To: hebert.cuellar@gmail.com
		
		Greetings {{.RecipientName}}!,

		This is the Pull Request report digest for {{.Repository.Owner}}/{{.Repository.Name}} project's last week:

		Open Pull Requests:
		{{.OpenPulls}}

		Closed Pull Requests:
		{{.ClosedPulls}}

		Open Drafts:
		{{.DraftPulls}}
		
		Until next week.

		HebertCL
	`

	tmpl, err := template.New("email").Parse(message)
	if err != nil {
		return "", err
	}

	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, data); err != nil {
		return "", err
	}

	return emailBody.String(), nil
}

func (sc senderConfig) sendReport(recipient []string, data emailData) error {
	// emailAuth := emailAuthentication(sc.User, sc.Password, sc.Server)
	message, err := generateTemplateEmail(data)
	if err != nil {
		return err
	}

	fmt.Printf("Display email content:\n%s\n", message)

	// if err := smtp.SendMail(sc.Server+":"+fmt.Sprint(sc.SmtpPort),
	// 	emailAuth,
	// 	sc.User,
	// 	recipient,
	// 	[]byte(message),
	// ); err != nil {
	// 	return err
	// }

	return nil
}
