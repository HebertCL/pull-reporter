package main

import (
	"log"
	"os"
	"strconv"

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

	// Create GitHub default client
	client := newGhClient()

	// Define recipient
	recipientList := []string{"hebert.cuellar@gmail.com"}

	// Define repository values
	repo := repository{
		Owner: repoOwner,
		Name:  repoName,
	}

	// Get all PR related values which will be used in the template email
	openPulls := repo.sortPullRequests(client, "open", false)
	closedPulls := repo.sortPullRequests(client, "closed", false)
	draftPulls := repo.sortPullRequests(client, "open", true)

	// Define email message body values for template
	emailBody := emailData{
		OpenPulls:      openPulls.String(),
		ClosedPulls:    closedPulls.String(),
		DraftPulls:     draftPulls.String(),
		Repository:     repo,
		RecipientName:  recipientName,
		RecipientEmail: recipientEmail,
	}

	// Define email configuration
	smtpConfig := senderConfig{
		Server:   mailServer,
		User:     mailUser,
		Password: mailPass,
		SmtpPort: port,
	}

	if err := smtpConfig.sendReport(recipientList, emailBody); err != nil {
		log.Fatal(err)
	}
}
