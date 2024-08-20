package git

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	mailer "github.com/HebertCL/pull-reporter/pkg/notifications"
	"github.com/google/go-github/v55/github"
)

type pullRequest struct {
	CreationDate time.Time
	Title        string
	Url          string
	State        string
	Number       int
	Draft        bool
}

type Mailer interface {
	SendReport(recipientList []string, data mailer.EmailData) error
}

type githubClient struct {
	client *github.Client
	mailer Mailer
}

func NewGithubClient(mailer Mailer) *githubClient {
	return &githubClient{
		client: github.NewClient(nil),
		mailer: mailer,
	}
}

func (c *githubClient) ProcessRepo(owner, name, recipientName, recipientEmail string) {
	// Define recipient
	recipientList := []string{"hebert.cuellar@gmail.com"}

	// Get all PR related values which will be used in the template email
	openPulls := c.sortPullRequests(owner, name, "open", false)
	closedPulls := c.sortPullRequests(owner, name, "closed", false)
	draftPulls := c.sortPullRequests(owner, name, "open", true)

	// Define email message body values for template
	emailBody := mailer.EmailData{
		RepositoryOwner: owner,
		RepositoryName:  name,
		OpenPulls:       openPulls.String(),
		ClosedPulls:     closedPulls.String(),
		DraftPulls:      draftPulls.String(),
		RecipientName:   recipientName,
		RecipientEmail:  recipientEmail,
	}

	if err := c.mailer.SendReport(recipientList, emailBody); err != nil {
		log.Fatal(err)
	}
}

func (c *githubClient) listPullRequests(owner string, name string) []*github.PullRequest {
	var prList []*github.PullRequest
	ctx := context.Background()
	// Use defined PR options for the specified requirement
	options := &github.PullRequestListOptions{
		State: "all",
		// Sort:      "long-running",
		// Direction: "desc",
	}

	prs, resp, err := c.client.PullRequests.List(ctx, owner, name, options)
	if err != nil {
		log.Printf("Failed to get Pull Requests: %s\n Status Code: %v", err, resp.StatusCode)
		return []*github.PullRequest{}
	}
	for _, pr := range prs {
		hoursSinceCreation := time.Since(pr.CreatedAt.Time)
		if hoursSinceCreation.Hours() < 7*24 {
			prList = append(prList, pr)
		}
	}
	return prList
}

func (c *githubClient) sortPullRequests(owner, name, prState string, prDraft bool) strings.Builder {
	var unfilteredPrList []pullRequest
	prList := c.listPullRequests(owner, name)

	// Loop over GitHub PR object
	for _, pr := range prList {
		unfilteredPrList = append(unfilteredPrList, pullRequest{
			CreationDate: pr.CreatedAt.Time,
			Title:        *pr.Title,
			Url:          *pr.URL,
			State:        *pr.State,
			Number:       *pr.Number,
			Draft:        *pr.Draft,
		})
	}
	var pulls []pullRequest
	// Loop over unfiltered pull request custom object
	// and save open pull requests
	for _, pull := range unfilteredPrList {
		if pull.State == prState && pull.Draft == prDraft {
			pulls = append(pulls, pull)
		}
	}
	var pullstring strings.Builder
	// Loop over the filtered list and return it as a string
	for _, value := range pulls {
		if len(pulls) != 0 {
			fmt.Fprintf(&pullstring, "Number: %v\nTitle: %s\nCreated: %s\nState: %s\nDraft: %v \nURL: %s\n\n", value.Number, value.Title, value.CreationDate, value.State, value.Draft, value.Url)
		} else {
			fmt.Fprintf(&pullstring, "Hooray! There's no %s pull requests.", prState)
		}
	}
	fmt.Println(len(pulls))
	return pullstring
}
