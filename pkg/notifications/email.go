package notifications

import (
	"bytes"
	"fmt"
	"html/template"
)

type MailerConfig struct {
	Server   string
	User     string
	Password string
	SmtpPort int
}

type EmailData struct {
	RepositoryOwner string
	RepositoryName  string
	OpenPulls       string
	ClosedPulls     string
	DraftPulls      string
	RecipientName   string
	RecipientEmail  string
}

// Creates a SMTP plain authentication
// TODO: Improve email authtentication
// func emailAuthentication(user string, password string, server string) smtp.Auth {
// 	auth := smtp.PlainAuth("", user, password, server)

// 	return auth
// }

type mailer struct{}

func NewMailer(config MailerConfig) *mailer {
	return &mailer{}
}

func (m mailer) generateTemplateData(data EmailData) (string, error) {
	var emailBody bytes.Buffer

	message := `
Subject: GitHub PR Report
To: {{.RecipientEmail}}
Greetings {{.RecipientName}}!,

This is the Pull Request report digest for {{.RepositoryOwner}}/{{.RepositoryName}} project's last week:
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

	if err := tmpl.Execute(&emailBody, data); err != nil {
		return "", err
	}

	return emailBody.String(), nil
}

func (m mailer) SendReport(recipientList []string, data EmailData) error {
	// emailAuth := emailAuthentication(sc.User, sc.Password, sc.Server)
	message, err := m.generateTemplateData(data)
	if err != nil {
		return err
	}

	fmt.Printf("Display email content:\n%s\n", message)
	// TODO: Actually send the email
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
