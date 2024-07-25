package main

import (
	"log"
	"os"
	"strconv"

	"github.com/HebertCL/pull-reporter/mailer"
	"github.com/HebertCL/pull-reporter/repo"
	"github.com/joho/godotenv"
)

func main() {
	// Load values from .env file
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

	// Define email configuration
	smtpConfig := mailer.MailerConfig{
		Server:   mailServer,
		User:     mailUser,
		Password: mailPass,
		SmtpPort: port,
	}

	mailer := mailer.NewMailer(smtpConfig)

	// Create GitHub default client
	client := repo.NewGithubClient(mailer)

	client.ProcessRepo(repoOwner, repoName, recipientName, recipientEmail)

}
