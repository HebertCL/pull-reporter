package notifications

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/HebertCL/pull-reporter/pkg/github"
)

type SenderConfig struct {
	Server   string
	User     string
	Password string
	SmtpPort int
}

type EmailData struct {
	Repository     github.Repository
	OpenPulls      string
	ClosedPulls    string
	DraftPulls     string
	RecipientName  string
	RecipientEmail string
}

// Creates a SMTP plain authentication
// TODO: Improve email authtentication
// func emailAuthentication(user string, password string, server string) smtp.Auth {
// 	auth := smtp.PlainAuth("", user, password, server)

// 	return auth
// }

func generateTemplateEmail(data EmailData) (string, error) {
	message := `
Subject: GitHub PR Report
To: {{.RecipientEmail}}

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

func (sc SenderConfig) SendReport(recipient []string, data EmailData) error {
	// emailAuth := emailAuthentication(sc.User, sc.Password, sc.Server)
	message, err := generateTemplateEmail(data)
	if err != nil {
		return err
	}

	fmt.Printf("Display email content:\n%s\n", message)

	//TODO: Actually send the email
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
