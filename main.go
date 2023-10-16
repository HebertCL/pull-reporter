package main

import (
	"log"
	"os"
	"strconv"

	"github.com/HebertCL/pull-reporter/pkg/email"
	"github.com/HebertCL/pull-reporter/pkg/github"

	"github.com/joho/godotenv"
)

func main() {
	// Load values from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Couldn't load values from .env: %v", err)
	}

	repoOwner := os.Getenv("REPO_OWNER")
	repoName := os.Getenv("REPO_NAME")
	authToken := os.Getenv("AUTH_TOKEN")
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

	// Create GitHub client. Uses default client if authToken is not set.
	// Using default client for private repos or repos out of your orgs
	// will result in failed calls to GitHub API
	client := github.NewGhClient(authToken)

	// Define recipient
	recipientList := []string{"hebert.cuellar@gmail.com"}

	// Define repository values
	repo := github.Repository{
		Owner: repoOwner,
		Name:  repoName,
	}

	// Get all PR related values which will be used in the template email
	openPulls := repo.SortPullRequests(client, "open", false)
	closedPulls := repo.SortPullRequests(client, "closed", false)
	draftPulls := repo.SortPullRequests(client, "open", true)

	// Define email message body values for template
	emailBody := email.EmailData{
		OpenPulls:      openPulls.String(),
		ClosedPulls:    closedPulls.String(),
		DraftPulls:     draftPulls.String(),
		Repository:     repo,
		RecipientName:  recipientName,
		RecipientEmail: recipientEmail,
	}

	// Define email configuration
	smtpConfig := email.SenderConfig{
		Server:   mailServer,
		User:     mailUser,
		Password: mailPass,
		SmtpPort: port,
	}

	if err := smtpConfig.SendReport(recipientList, emailBody); err != nil {
		log.Fatal(err)
	}
}
