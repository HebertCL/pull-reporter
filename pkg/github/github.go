package github

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/v55/github"
)

type Repository struct {
	Owner string
	Name  string
}

type pullRequest struct {
	CreationDate time.Time
	Title        string
	Url          string
	State        string
	Number       int
	Draft        bool
}

// Create a new GH client. Uses default, unauthenticated client
func NewGhClient(token string) *github.Client {
	// Return default token if no authentication token is provided
	if token == "" {
		return github.NewClient(nil)
	}

	client := github.NewClient(nil).WithAuthToken(token)

	return client
}

// Create a list with PRs with 7 or less days and return it
func listPullRequests(client *github.Client, name string, owner string) []*github.PullRequest {
	var prList []*github.PullRequest
	ctx := context.Background()
	// Use defined PR options for the specified requirement
	options := &github.PullRequestListOptions{
		State: "all",
		// Sort:      "long-running",
		// Direction: "desc",
	}

	prs, resp, err := client.PullRequests.List(ctx, owner, name, options)
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

// SortPullRequests creates a list of all PRs in the last 7 days and returns
// a formated string of values for those PRs taking into account the PR state
// and its draft status.
func (r Repository) SortPullRequests(client *github.Client, prState string, prDraft bool) strings.Builder {
	var unfilteredPrList []pullRequest
	prList := listPullRequests(client, r.Name, r.Owner)

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

	// Validate if the slice has content
	if len(pulls) == 0 {
		fmt.Fprintf(&pullstring, "\nHooray! Nothing to do here :D\n")

		return pullstring
	}

	//Loop over the filtered list and return it as a string
	for _, value := range pulls {
		if len(pulls) != 0 {
			fmt.Fprintf(&pullstring, "\nNumber: %v\nTitle: %s\nCreated: %s\nState: %s\nDraft: %v \nURL: %s\n\n", value.Number, value.Title, value.CreationDate, value.State, value.Draft, value.Url)
		}
	}

	return pullstring
}
