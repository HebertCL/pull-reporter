package main

import (
	"log"
	"os"
	"strconv"

	repo "github.com/HebertCL/pull-reporter/pkg/git"
	mailer "github.com/HebertCL/pull-reporter/pkg/notifications"
	"github.com/joho/godotenv"
)

func main() {
	// TODO: Replace third-party library to load env variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Couldn't load values from .env: %v", err)
	}
	repoOwner := os.Getenv("REPO_OWNER")
	repoName := os.Getenv("REPO_NAME")
	recipientName := os.Getenv("RECIPIENT_NAME")
	recipientEmail := os.Getenv("RECIPIENT_EMAIL")
	mailServer := os.Getenv("MAIL_SERVER")
	mailUser := os.Getenv("MAIL_USER")
	mailPass := os.Getenv("MAIL_PASS")
	mailPort := os.Getenv("MAIL_PORT")

	port, err := strconv.Atoi(mailPort)
	if err != nil {
		log.Fatalf("Couldn't convert port to integer: %v", err)
	}

	smtpConfig := mailer.MailerConfig{
		Server:   mailServer,
		User:     mailUser,
		Password: mailPass,
		SmtpPort: port,
	}

	// Define email configuration
	mailer := mailer.NewMailer(smtpConfig)

	// Create GitHub default client
	client := repo.NewGithubClient(mailer)

	client.ProcessRepo(repoOwner, repoName, recipientName, recipientEmail)
}
