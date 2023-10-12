package main

import (
	"log"
)

func main() {
	// Create GitHub default client
	client := newGhClient()

	// Define recipient
	recipientList := []string{"hebert.cuellar@gmail.com"}

	// Define repository values
	repo := repository{
		Owner: "opentofu",
		Name:  "opentofu",
	}

	// Get all PR related values which will be used in the template email
	openPulls := repo.sortPullRequests(client, "open", false)
	closedPulls := repo.sortPullRequests(client, "closed", false)
	draftPulls := repo.sortPullRequests(client, "open", true)

	// Define email message body values for template
	emailBody := emailData{
		OpenPulls:     openPulls.String(),
		ClosedPulls:   closedPulls.String(),
		DraftPulls:    draftPulls.String(),
		Repository:    repo,
		RecipientName: "HebertCL",
	}

	// Define email configuration
	smtpConfig := senderConfig{
		Server:   "smtp.gmail.com",
		User:     "hebert.andres.cl@gmail.com",
		Password: "LuvMyNut4!",
		SmtpPort: 587,
	}

	if err := smtpConfig.sendReport(recipientList, emailBody); err != nil {
		log.Fatal(err)
	}
}
