package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/HebertCL/pull-reporter/foundation/logger"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.New("PULL-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	sig := <-shutdown

	log.Infow("shutdown", "status", "shutdown started", "signal", sig)
	defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

	return nil
}

// TODO: Replace old logic and remove this commented main func.
// func main() {
// 	// Load values from .env file
// 	if err := godotenv.Load(); err != nil {
// 		log.Fatalf("Couldn't load values from .env: %v", err)
// 	}

// 	repoOwner := os.Getenv("REPO_OWNER")
// 	repoName := os.Getenv("REPO_NAME")
// 	recipientName := os.Getenv("RECIPIENT_NAME")
// 	recipientEmail := os.Getenv("RECIPIENT_EMAIL")
// 	mailServer := os.Getenv("MAIL_SERVER")
// 	mailUser := os.Getenv("MAIL_USER")
// 	mailPass := os.Getenv("MAIL_PASS")
// 	mailPort := os.Getenv("MAIL_PORT")

// 	port, err := strconv.Atoi(mailPort)
// 	if err != nil {
// 		log.Fatalf("Couldn't convert port to integer: %v", err)
// 	}

// 	// Create GitHub default client
// 	client := newGhClient()

// 	// Define recipient
// 	recipientList := []string{"hebert.cuellar@gmail.com"}

// 	// Define repository values
// 	repo := repository{
// 		Owner: repoOwner,
// 		Name:  repoName,
// 	}

// 	// Get all PR related values which will be used in the template email
// 	openPulls := repo.sortPullRequests(client, "open", false)
// 	closedPulls := repo.sortPullRequests(client, "closed", false)
// 	draftPulls := repo.sortPullRequests(client, "open", true)

// 	// Define email message body values for template
// 	emailBody := emailData{
// 		OpenPulls:      openPulls.String(),
// 		ClosedPulls:    closedPulls.String(),
// 		DraftPulls:     draftPulls.String(),
// 		Repository:     repo,
// 		RecipientName:  recipientName,
// 		RecipientEmail: recipientEmail,
// 	}

// 	// Define email configuration
// 	smtpConfig := senderConfig{
// 		Server:   mailServer,
// 		User:     mailUser,
// 		Password: mailPass,
// 		SmtpPort: port,
// 	}

// 	if err := smtpConfig.sendReport(recipientList, emailBody); err != nil {
// 		log.Fatal(err)
// 	}
// }
