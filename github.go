package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/v55/github"
)

type repository struct {
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
// TODO: Support both authenticated and non-authenticated clients
func newGhClient() *github.Client {
	return github.NewClient(nil)
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

// Filter open PRs which are not Draft.
// FIXME: The filter for open PRs, both Draft and non-Draft could probably be handle in the same function
func (r repository) sortPullRequests(client *github.Client, prState string, prDraft bool) strings.Builder {
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
	//Loop over the filtered list and return it as a string
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
